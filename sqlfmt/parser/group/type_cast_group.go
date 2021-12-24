package group

import (
	"bytes"
	"fmt"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// TypeCast group.
type TypeCast struct {
	elementReindenter
}

func NewTypeCast(element []Reindenter, opts ...Option) *TypeCast {
	return &TypeCast{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements.
func (t *TypeCast) Reindent(buf *bytes.Buffer) error {
	elements, err := t.processPunctuation()
	if err != nil {
		return err
	}

	return t.elementsTokenApply(elements, buf, t.writeTypeCast)
}

func (t *TypeCast) writeTypeCast(buf *bytes.Buffer, token lexer.Token, _ int) error {
	switch token.Type {
	case lexer.TYPE:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	case lexer.COMMA:
		buf.WriteString(fmt.Sprintf(
			"%s%s",
			token.FormattedValue(),
			WhiteSpace,
		))
	default:
		buf.WriteString(token.FormattedValue())
	}

	return nil
}
