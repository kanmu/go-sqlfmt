package sqlfmt

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"
)

// sqlfmt retrieves all strings from "Query" and "QueryRow" and "Exec" functions in .go file.
const (
	QUERY    = "Query"
	QUERYROW = "QueryRow"
	EXEC     = "Exec"
)

// replaceAst replace ast node with formatted SQL statement.
func replaceAst(f *ast.File, fset *token.FileSet, opts ...Option) {
	o := defaultOptions(opts...)

	ast.Inspect(f, func(n ast.Node) bool {
		x, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		fun, ok := x.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		funcName := fun.Sel.Name
		if funcName != QUERY && funcName != QUERYROW && funcName != EXEC {
			return true
		}

		if len(x.Args) == 0 {
			return true
		}

		// not for parsing url.Query
		arg, ok := x.Args[0].(*ast.BasicLit)
		if !ok {
			return true
		}

		sqlStmt := arg.Value
		if !strings.HasPrefix(sqlStmt, "`") {
			return true
		}

		src := strings.Trim(sqlStmt, "`")

		res, err := Format(src, opts...)
		if err != nil {
			log.Println(fmt.Sprintf("Format failed at %s: %v", fset.Position(arg.Pos()), err))

			return true
		}

		arg.Value = fmt.Sprintf("`%s%s`", res, strings.Repeat(group.WhiteSpace, o.Distance))

		return true
	})
}
