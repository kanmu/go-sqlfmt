package sqlfmt

import (
	"fmt"
	"go/ast"
	"strings"
)

// sqlfmt retrieves all strings from "Query" and "QueryRow" and "Exec" functions in .go file
const (
	QUERY    = "Query"
	QUERYROW = "QueryRow"
	EXEC     = "Exec"
)

// replaceAstWithFormattedStmt replace ast node with formatted SQL statement
func replaceAstWithFormattedStmt(f *ast.File) error {
	var err error
	ast.Inspect(f, func(n ast.Node) bool {
		if x, ok := n.(*ast.CallExpr); ok {
			if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
				funcName := fun.Sel.Name
				if funcName == QUERY || funcName == QUERYROW || funcName == EXEC {
					if len(x.Args) > 0 {
						if arg, ok := x.Args[0].(*ast.BasicLit); ok {
							sqlStmt := arg.Value
							if !strings.HasPrefix(sqlStmt, "`") {
								return true
							}
							src := strings.Trim(sqlStmt, "`")
							formattedStmt, e := Format(src)
							if e != nil {
								err = &FormatError{
									Msg: e.Error(),
								}
								return false
							}
							arg.Value = "`" + formattedStmt + "`"
						}
					}
				}
			}
		}
		return true
	})
	return err
}

// FormatError is an error that occurs during Format
type FormatError struct {
	Msg string
}

// Error ...
func (f *FormatError) Error() string {
	return fmt.Sprintf("FORMAT ERROR :%#v\n", f.Msg)
}
