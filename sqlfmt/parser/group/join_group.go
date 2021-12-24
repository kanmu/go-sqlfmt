package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Join clause.
type Join struct {
	elementReindenter
}

func NewJoin(element []Reindenter, opts ...Option) *Join {
	return &Join{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindent its elements.
func (j *Join) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(j.Element)
	if err != nil {
		return err
	}

	for i, v := range elements {
		if token, ok := v.(lexer.Token); ok {
			j.writeJoin(buf, token, j.IndentLevel, i == 0)
		} else {
			if eri := v.Reindent(buf); eri != nil {
				return eri
			}
		}
	}

	return nil
}

func (j *Join) writeJoin(buf *bytes.Buffer, token lexer.Token, indent int, isFirst bool) {
	switch {
	case isFirst && token.IsJoinStart():
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.FormattedValue()))
	case token.Type == lexer.ON || token.Type == lexer.USING:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.FormattedValue()))
	case strings.HasPrefix(token.FormattedValue(), "::"):
		buf.WriteString(token.FormattedValue())
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	}
}
