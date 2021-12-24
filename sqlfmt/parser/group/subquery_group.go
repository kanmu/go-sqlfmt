package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Subquery group.
type Subquery struct {
	elementReindenter
	InColumnArea bool
	ColumnCount  int
}

func NewSubquery(element []Reindenter, opts ...Option) *Subquery {
	return &Subquery{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements.
func (s *Subquery) Reindent(buf *bytes.Buffer) error {
	elements, err := s.processPunctuation()
	if err != nil {
		return err
	}

	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			s.writeSubquery(buf, token, s.IndentLevel)
		} else {
			if s.InColumnArea {
				el.IncrementIndentLevel(1)
			}

			if eri := el.Reindent(buf); eri != nil {
				return eri
			}
		}
	}

	return nil
}

func (s *Subquery) writeSubquery(buf *bytes.Buffer, token lexer.Token, indent int) {
	switch {
	case token.Type == lexer.STARTPARENTHESIS && s.ColumnCount == 0 && s.InColumnArea:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.FormattedValue()))
	case token.Type == lexer.STARTPARENTHESIS:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	case token.Type == lexer.ENDPARENTHESIS && s.ColumnCount > 0:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.FormattedValue()))
	case token.Type == lexer.ENDPARENTHESIS:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent-1), token.FormattedValue()))
	case strings.HasPrefix(token.FormattedValue(), "::"):
		buf.WriteString(token.FormattedValue())
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	}
}
