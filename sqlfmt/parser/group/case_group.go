package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

const (
	DoubleColumn = "::"
)

// Case Clause.
type Case struct {
	elementReindenter
	hasCommaBefore bool
}

func NewCase(element []Reindenter, opts ...Option) *Case {
	return &Case{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements.
func (c *Case) Reindent(buf *bytes.Buffer) error {
	elements, err := c.processPunctuation()
	if err != nil {
		return err
	}

	return c.elementsTokenApply(elements, buf, c.writeCase)
}

func (c *Case) writeCase(buf *bytes.Buffer, token lexer.Token, indent int) error {
	if c.hasCommaBefore {
		c.writeCaseWithCommaBefore(buf, token, indent)
	} else {
		c.writeCaseWithoutCommaBefore(buf, token, indent)
	}

	return nil
}

func (c *Case) writeCaseWithCommaBefore(buf *bytes.Buffer, token lexer.Token, indent int) {
	switch token.Type {
	case lexer.CASE:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	case lexer.WHEN, lexer.ELSE:
		buf.WriteString(fmt.Sprintf("%s%s%s%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, WhiteSpace, WhiteSpace, DoubleWhiteSpace, token.FormattedValue()))
	case lexer.END:
		buf.WriteString(fmt.Sprintf("%s%s%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, WhiteSpace, WhiteSpace, token.FormattedValue()))
	case lexer.COMMA:
		buf.WriteString(token.FormattedValue())
	default:
		if strings.HasPrefix(token.FormattedValue(), DoubleColumn) {
			buf.WriteString(token.FormattedValue())
		} else {
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
		}
	}
}

func (c *Case) writeCaseWithoutCommaBefore(buf *bytes.Buffer, token lexer.Token, indent int) {
	switch token.Type {
	case lexer.CASE, lexer.END:
		buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, token.FormattedValue()))
	case lexer.WHEN, lexer.ELSE:
		buf.WriteString(fmt.Sprintf("%s%s%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, WhiteSpace, DoubleWhiteSpace, token.FormattedValue()))
	case lexer.COMMA:
		buf.WriteString(token.FormattedValue())
	default:
		if strings.HasPrefix(token.FormattedValue(), DoubleColumn) {
			buf.WriteString(token.FormattedValue())
		} else {
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
		}
	}
}
