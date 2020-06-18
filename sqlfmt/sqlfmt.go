package sqlfmt

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/pkg/errors"
)

// Options for go-sqlfmt
type Options struct {
	Distance int
}

// Process formats SQL statement in .go file
func Process(filename string, src []byte, options *Options) ([]byte, error) {
	fset := token.NewFileSet()
	parserMode := parser.ParseComments

	f, err := parser.ParseFile(fset, filename, src, parserMode)
	if err != nil {
		return nil, errors.Wrap(err, "parser.ParseFile failed")
	}

	Replace(f, options)

	var buf bytes.Buffer
	if err = printer.Fprint(&buf, fset, f); err != nil {
		return nil, errors.Wrap(err, "printer.Fprint failed")
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "format.Source failed")
	}
	return out, nil
}
