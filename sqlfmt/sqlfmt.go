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

	f, err := parser.ParseFile(fset, filename, src, parserMode)
	if err != nil {
		return nil, err
	}

	if err := replaceAstWithFormattedStmt(f); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = printer.Fprint(&buf, fset, f)
	if err != nil {
		return nil, err
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return out, nil
}
