package sqlfmt

import (
	"bytes"
	fmt "go/format"
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

	astFile, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, formatErr(errors.Wrap(err, "parser.ParseFile failed"))
	}

	inspectAndReplace(fset, astFile, options)

	var buf bytes.Buffer
	if err = printer.Fprint(&buf, fset, astFile); err != nil {
		return nil, formatErr(errors.Wrap(err, "printer.Fprint failed"))
	}

	out, err := fmt.Source(buf.Bytes())
	if err != nil {
		return nil, formatErr(errors.Wrap(err, "format.Source failed"))
	}
	return out, nil
}

func formatErr(err error) error {
	return &FormatError{msg: err.Error()}
}
