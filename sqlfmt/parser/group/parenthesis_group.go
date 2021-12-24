package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Parenthesis clause
type Parenthesis struct {
	elementReindenter
	InColumnArea bool
	ColumnCount  int
}

func NewParenthesis(element []Reindenter, opts ...Option) *Parenthesis {
	return &Parenthesis{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements
func (p *Parenthesis) Reindent(buf *bytes.Buffer) error {
	var hasStartBefore bool

	elements, err := p.processPunctuation()
	if err != nil {
		return err
	}

	for i, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			hasStartBefore = (i == 1)
			p.writeParenthesis(buf, token, p.IndentLevel, hasStartBefore)
		} else {
			if eri := el.Reindent(buf); eri != nil {
				return eri
			}
		}
	}

	return nil
}

func (p *Parenthesis) writeParenthesis(buf *bytes.Buffer, token lexer.Token, indent int, hasStartBefore bool) {
	switch {
	case token.Type == lexer.STARTPARENTHESIS && p.ColumnCount == 0 && p.InColumnArea:
		buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, token.FormattedValue()))
	case token.Type == lexer.STARTPARENTHESIS:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	case token.Type == lexer.ENDPARENTHESIS:
		buf.WriteString(token.FormattedValue())
	case token.Type == lexer.COMMA:
		buf.WriteString(token.FormattedValue())
	case hasStartBefore:
		buf.WriteString(token.FormattedValue())
	case strings.HasPrefix(token.FormattedValue(), "::"):
		buf.WriteString(token.FormattedValue())
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	}
}
