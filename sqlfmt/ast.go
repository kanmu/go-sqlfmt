package sqlfmt

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"
)

// sqlfmt retrieves all strings from "Query" and "QueryRow" and "Exec" functions in .go file
const (
	QUERY    = "Query"
	QUERYROW = "QueryRow"
	EXEC     = "Exec"
)

// replaceAstWithFormattedStmt replace ast node with formatted SQL statement
func replaceAstWithFormattedStmt(f *ast.File, fset *token.FileSet) {
	ast.Inspect(f, func(n ast.Node) bool {
		if x, ok := n.(*ast.CallExpr); ok {
			if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
				funcName := fun.Sel.Name
				if funcName == QUERY || funcName == QUERYROW || funcName == EXEC {
					// not for parsing url.Query
					if len(x.Args) > 0 {
						if arg, ok := x.Args[0].(*ast.BasicLit); ok {
							sqlStmt := arg.Value
							if !strings.HasPrefix(sqlStmt, "`") {
								return true
							}
							src := strings.Trim(sqlStmt, "`")
							formattedStmt, err := Format(src)
							if err != nil {
								log.Println(fmt.Sprintf("format failed at %s\n", fset.Position(arg.Pos())))
								return true
							}
							arg.Value = "`" + formattedStmt + "`"
						}
					}
				}
			}
		}
		return true
	})
}
