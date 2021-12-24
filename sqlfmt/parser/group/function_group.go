package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Function clause
type Function struct {
	elementReindenter
	InColumnArea bool
	ColumnCount  int
}

func NewFunction(element []Reindenter, opts ...Option) *Function {
	return &Function{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements
func (f *Function) Reindent(buf *bytes.Buffer) error {
	elements, err := f.processPunctuation()
	if err != nil {
		return err
	}

	for i, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			var prev lexer.Token

			if i > 0 {
				if preToken, ok := elements[i-1].(lexer.Token); ok {
					prev = preToken
				}
			}

			f.writeFunction(buf, token, prev, f.IndentLevel)
		} else {
			if eri := el.Reindent(buf); eri != nil {
				return eri
			}
		}
	}

	return nil
}

func (f *Function) writeFunction(buf *bytes.Buffer, token, prev lexer.Token, indent int) {
	switch {
	case prev.Type == lexer.STARTPARENTHESIS || token.Type == lexer.STARTPARENTHESIS || token.Type == lexer.ENDPARENTHESIS:
		buf.WriteString(token.FormattedValue())
	case token.Type == lexer.FUNCTION && f.ColumnCount == 0 && f.InColumnArea:
		buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, token.FormattedValue()))
	case token.Type == lexer.FUNCTION:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	case token.Type == lexer.COMMA:
		buf.WriteString(token.FormattedValue())
	case strings.HasPrefix(token.FormattedValue(), "::"):
		buf.WriteString(token.FormattedValue())
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	}
}
