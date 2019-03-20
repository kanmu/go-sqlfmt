package sqlfmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
)

// sqlfmt retrieves all strings from "Query" and "QueryRow" and "Exec" functions in .go file
const (
	QUERY    = "Query"
	QUERYROW = "QueryRow"
	EXEC     = "Exec"
)

// SQLFormatter represents SQLformatter
type SQLFormatter struct {
	AstNode   *ast.File
	Formatter *Formatter
	Fset      *token.FileSet
}

// NewSQLFormatter creates SQLFormatter
func NewSQLFormatter(src io.Reader) (*SQLFormatter, error) {
	var (
		node *ast.File
		err  error
	)

	fset := token.NewFileSet()
	parserMode := parser.ParseComments

	if file, ok := src.(*os.File); ok {
		node, err = parser.ParseFile(fset, file.Name(), nil, parserMode)
		if err != nil {
			return nil, err
		}
	} else {
		node, err = parser.ParseFile(fset, "file.go", src, parserMode)
		if err != nil {
			return nil, err
		}
	}
	return &SQLFormatter{
		AstNode:   node,
		Formatter: &Formatter{},
		Fset:      fset,
	}, nil
}

// Format formats SQL statements after retrieving from AstNode
func (s *SQLFormatter) Format() error {
	var err error
	ast.Inspect(s.AstNode, func(n ast.Node) bool {
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
							formattedStmt, e := s.Formatter.Format(src)
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
