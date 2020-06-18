package sqlfmt

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
)

// Replace replace ast node with formatted SQL statement
func Replace(f *ast.File, options *Options) {
	ast.Inspect(f, func(n ast.Node) bool {
		sql, found := findSQL(n)
		if found {
			res, err := Format(sql, options)
			if err != nil {
				log.Println(err)

				// XXX for debugging
				log.Println(sql)
			} else {
				replace(n, res)
			}
		}
		return true
	})
}

func replace(n ast.Node, sql string) {
	replaceFunc := func(cr *astutil.Cursor) bool {
		cr.Replace(&ast.BasicLit{
			Kind:  token.STRING,
			Value: sql,
		})
		return true
	}
	astutil.Apply(n, replaceFunc, nil)
}

func findSQL(n ast.Node) (string, bool) {
	ce, ok := n.(*ast.CallExpr)
	if !ok {
		return "", false
	}
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	// check func name
	ok = validateFuncName(se.Sel.Name)
	if !ok {
		return "", false
	}

	// check length of the parameter
	// this is not for parsing "url.Query()"
	// FIXME: very adhoc
	if len(ce.Args) == 0 {
		return "", false
	}

	// SQL statement should appear in the first parameter
	arg, ok := ce.Args[0].(*ast.BasicLit)
	if !ok {
		return "", false
	}
	return arg.Value, true
}

// go-sqlfmt only formats the value passed as the parameter of "Exec(string, ... any type)", "Query(string, ... any type)" and "QueryRow(string, ... any type)"
func validateFuncName(name string) bool {
	switch name {
	case "Exec", "Query", "QueryRow":
		return true
	}
	return false
}
