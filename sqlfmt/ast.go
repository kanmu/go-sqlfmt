package sqlfmt

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// sqlfmt retrieves string value which is surrounded by backquote from "Query" and "QueryRow" and "Exec" functions in .go file
const (
	Backquote        = "`"
	FuncNameQUERY    = "Query"
	FuncNameQUERYROW = "QueryRow"
	FuncNameEXEC     = "Exec"
)

// replaceAstWithFormattedSQL traverse ast node until it finds a node with SQL
// sqlfmt format values passed to "Query" and "QueryRow" and "Exec" functions
func replaceAstWithFormattedSQL(fset *token.FileSet, f *ast.File, ops *Options) {
	ast.Inspect(f, func(n ast.Node) bool {
		if x, ok := n.(*ast.CallExpr); ok {
			if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
				funcName := fun.Sel.Name
				if funcName == FuncNameQUERY || funcName == FuncNameQUERYROW || funcName == FuncNameEXEC {
					// not for parsing url.Query
					if len(x.Args) > 0 {
						if lit, ok := x.Args[0].(*ast.BasicLit); ok {
							// lit.Value has SQL value
							// parseAndCheck parse, format and check the value
							// if ok to replace, go-sqlfmt replace the node with res
							res, replaceable, err := parseAndCheck(lit.Value, ops)
							if err == nil && replaceable {
								lit.Value = res
							}
							// should I print error ???
						}
					}
				}
			}
		}
		return true
	})
}

func parseAndCheck(v string, ops *Options) (string, bool, error) {
	// ignore SQL without backquote
	// go-sqlfmt only support SQL with back-quote like `select * from table`
	if !isSQLWithBackQuote(v) {
		return v, false, nil
	}

	// pass value without back-quote
	src := strings.Trim(v, Backquote)

	res, err := Format(src, ops)
	if err != nil {
		return res, false, err
	}
	// if src has changed destructively, return error in order not to replace ast with res
	// TODO: consider to  add a option to replace node anyway ?
	if !checkSum(src, res) {
		return res, false, fmt.Errorf("failed to format: %s", src)
	}

	// return value with back-quote
	return setBackquote(res), true, nil
}

// TODO: ！！
func checkSum(a, b string) bool {
	return false
}

func isSQLWithBackQuote(src string) bool {
	if strings.HasPrefix(src, Backquote) && strings.HasSuffix(src, Backquote) {
		return true
	}
	return false
}

func setBackquote(src string) string {
	return fmt.Sprintf("%s%s%s", Backquote, src, Backquote)
}
