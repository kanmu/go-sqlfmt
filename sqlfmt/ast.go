package sqlfmt

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"

	"github.com/kanmu/go-sqlfmt/sqlfmt/parser/group"
	"github.com/pkg/errors"
)

// sqlfmt retrieves all strings from "Query" and "QueryRow" and "Exec" functions in .go file
const (
	QUERY    = "Query"
	QUERYROW = "QueryRow"
	EXEC     = "Exec"
)

// inspectAndReplace inspect ast node and replace it with formatted SQL
func inspectAndReplace(fset *token.FileSet, f *ast.File, options *Options) {
	ast.Inspect(f, func(n ast.Node) bool {
		if x, ok := n.(*ast.CallExpr); ok {
			if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
				funcName := fun.Sel.Name
				if funcName == QUERY || funcName == QUERYROW || funcName == EXEC {
					// not for parsing url.Query
					if len(x.Args) > 0 {
						if node, ok := x.Args[0].(*ast.BasicLit); ok {
							if err := replace(node, options); err != nil {
								log.Println(fmt.Sprintf("Format failed at %s: %v", fset.Position(node.Pos()), err))
							}
						}
					}
				}
			}
		}
		return true
	})
}

// replace replace the node with formatted statement
func replace(node *ast.BasicLit, options *Options) error {
	// node.Value should have SQL statement
	stmt := node.Value

	// ignore SQL without backquote
	if !strings.HasPrefix(stmt, "`") {
		return nil
	}

	src := strings.Trim(stmt, "`")
	res, err := format(src)
	if err != nil {
		return errors.Wrap(err, "format failed")
	}

	// TODO: i will delete this getStmtWithDistance later, so that i don't have to do things below
	// i put new line before the last backquote because format.Source ("go/format") does not work without new line somehow ...
	if options.Distance == 0 {
		node.Value = "`" + res + "\n`"
		return nil
	}

	result := getStmtWithDistance(res, options.Distance)

	// FIXME: more elegant
	// this is for the backquote appearing after SQL statements
	node.Value = "`" + result + strings.Repeat(group.WhiteSpace, options.Distance) + "`"
	fmt.Println(node.Value)

	return nil
}

func getStmtWithDistance(src string, distance int) string {
	scanner := bufio.NewScanner(strings.NewReader(src))

	var result string
	for scanner.Scan() {
		// FIXME: more elegant
		// this is for putting the newline before SQL statements
		if scanner.Text() == "" {
			result += group.NewLine
			continue
		}
		result += fmt.Sprintf("%s%s%s", strings.Repeat(group.WhiteSpace, distance), scanner.Text(), "\n")
	}
	return result
}
