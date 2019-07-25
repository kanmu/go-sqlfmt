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

func replace(node *ast.BasicLit, options *Options) error {
	stmt := node.Value

	if !strings.HasPrefix(stmt, "`") {
		return nil
	}

	src := strings.Trim(stmt, "`")
	// optionはここでは渡さずに、resに対して
	res, err := format(src, options)
	if err != nil {
		return errors.Wrap(err, "format failed")
	}

	if options.Distance != 0 {
		res = getStmtWithDistance(res, options.Distance)
	}

	// FIXME: more elegant
	// this is for the backquote appearing after SQL statements
	node.Value = "`" + res + strings.Repeat(group.WhiteSpace, options.Distance) + "`"

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
