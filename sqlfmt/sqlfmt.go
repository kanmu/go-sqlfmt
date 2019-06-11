package sqlfmt

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
)

// Options specifies options for processing files.
type Options struct {
	Fragment          bool
	PrintErrors       bool
	AllErrors         bool
	RemoveBareReturns bool
}

// Process formats SQL statement in .go file
func Process(filename string, src []byte, opt *Options) ([]byte, error) {
	fset := token.NewFileSet()
	parserMode := parser.ParseComments
	if opt.AllErrors {
		parserMode |= parser.AllErrors
	}

	astFile, err := parser.ParseFile(fset, filename, src, parserMode)
	if err != nil {
		return nil, err
	}

	replaceAstWithFormattedStmt(astFile, fset)

	var buf bytes.Buffer

	if err = printer.Fprint(&buf, fset, astFile); err != nil {
		return nil, err
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return out, nil
}

// func formatErr(err error) *FormatError {
// 	return &FormatError{err: err}
// }
