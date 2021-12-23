package sqlfmt

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/pkg/errors"
)

// Process formats SQL statement in .go file or sql files
func Process(filename string, src []byte, opts ...Option) ([]byte, error) {
	o := defaultOptions(opts...)

	if o.IsRawSQL {
		sql, err := Format(string(src), opts...)
		if err != nil {
			return nil, err
		}

		return []byte(sql), nil
	}

	o.Colorized = false // colors do not apply to go code formatting

	fset := token.NewFileSet()
	parserMode := parser.ParseComments
	astFile, err := parser.ParseFile(fset, filename, src, parserMode)
	if err != nil {
		return nil, formatErr(errors.Wrap(err, "parser.ParseFile failed"))
	}

	replaceAst(astFile, fset, withOptions(o))

	var buf bytes.Buffer
	if err = printer.Fprint(&buf, fset, astFile); err != nil {
		return nil, formatErr(errors.Wrap(err, "printer.Fprint failed"))
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, formatErr(errors.Wrap(err, "format.Source failed"))
	}
	return out, nil
}

func formatErr(err error) error {
	return &FormatError{msg: err.Error()}
}
