package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

type (
	// Reindenter interface
	// specific values of Reindenter would be clause group or token.
	Reindenter interface {
		Reindent(*bytes.Buffer) error
		IncrementIndentLevel(int)
		GetStart() int
	}

	baseReindenter struct {
		start int
		*options
	}

	elementReindenter struct {
		Element []Reindenter
		baseReindenter
	}
)

func (g baseReindenter) GetStart() int {
	return g.start
}

// GetOpts retrieves options. This is useful for comparing types when testing.
func (g *baseReindenter) GetOpts() []Option {
	if g.options == nil {
		return nil
	}

	return []Option{
		WithIndentLevel(g.options.IndentLevel),
		WithCommaStyle(g.options.commaStyle),
	}
}

// SetOpts sets options. This is useful for comparing types when testing.
func (g *baseReindenter) SetOpts(opts ...Option) {
	g.options = defaultOptions(opts...)
}

// IncrementIndentLevel increments by its specified indent level.
func (g *baseReindenter) IncrementIndentLevel(lev int) {
	g.IndentLevel += lev
}

func (g *baseReindenter) writeComma(buf *bytes.Buffer, token lexer.Token, indent int) {
	switch g.commaStyle {
	case CommaStyleRight:
		buf.WriteString(fmt.Sprintf(
			"%s%s%s%s",
			token.FormattedValue(),
			NewLine,
			strings.Repeat(DoubleWhiteSpace, indent),
			WhiteSpace,
		))
	default:
		buf.WriteString(fmt.Sprintf(
			"%s%s%s%s",
			NewLine,
			strings.Repeat(DoubleWhiteSpace, indent),
			DoubleWhiteSpace,
			token.FormattedValue(),
		))
	}
}

func newElementReindenter(element []Reindenter, opts ...Option) elementReindenter {
	o := defaultOptions(opts...)

	return elementReindenter{
		Element: element,
		baseReindenter: baseReindenter{
			options: o,
		},
	}
}

func (e *elementReindenter) processPunctuation() ([]Reindenter, error) {
	elements, err := processPunctuation(e.Element)
	if err != nil {
		return nil, err
	}

	return elements, nil
}

func (e *elementReindenter) elementsTokenApply(
	elements []Reindenter,
	buf *bytes.Buffer,
	apply func(*bytes.Buffer, lexer.Token, int) error,
) error {
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			if err := apply(buf, token, e.IndentLevel); err != nil {
				return err
			}
		} else {
			if err := el.Reindent(buf); err != nil {
				return err
			}
		}
	}

	return nil
}

// Reindent reindents its elements.
func (e *elementReindenter) Reindent(buf *bytes.Buffer) error {
	elements, err := e.processPunctuation()
	if err != nil {
		return err
	}

	return e.elementsTokenApply(elements, buf, write)
}

func (e *elementReindenter) writeWithComma(buf *bytes.Buffer, v interface{}, indent int) error {
	columnCount := e.start
	defer func() {
		e.start = columnCount
	}()

	if token, ok := v.(lexer.Token); ok {
		switch {
		case token.IsNeedNewLineBefore():
			buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.FormattedValue()))
		case token.Type == lexer.BY:
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
		case token.Type == lexer.COMMA:
			e.writeComma(buf, token, indent)
		default:
			return fmt.Errorf("can not reindent %#v", token.FormattedValue())
		}
	} else if str, ok := v.(string); ok {
		str = strings.TrimRight(str, " ")

		switch {
		case columnCount == 0:
			buf.WriteString(fmt.Sprintf("%s%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), DoubleWhiteSpace, str))
		case strings.HasPrefix(token.FormattedValue(), "::"):
			buf.WriteString(str)
		default:
			buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, str))
		}

		columnCount++
	}

	return nil
}

// to reindent.
const (
	NewLine          = "\n"
	WhiteSpace       = " "
	DoubleWhiteSpace = "  "
)

var (
	_ Reindenter = &elementReindenter{}
	_ Reindenter = &AndGroup{}
	_ Reindenter = &Case{}
	_ Reindenter = &Delete{}
	_ Reindenter = &From{}
	_ Reindenter = &Function{}
	_ Reindenter = &GroupBy{}
	_ Reindenter = &Having{}
	_ Reindenter = &Insert{}
	_ Reindenter = &Join{}
	_ Reindenter = &LimitClause{}
	_ Reindenter = &Lock{}
	_ Reindenter = &OrderBy{}
	_ Reindenter = &OrGroup{}
	_ Reindenter = &Parenthesis{}
	_ Reindenter = &Returning{}
	_ Reindenter = &Select{}
	_ Reindenter = &Set{}
	_ Reindenter = &Subquery{}
	_ Reindenter = &TieClause{}
	_ Reindenter = &TypeCast{}
	_ Reindenter = &Update{}
	_ Reindenter = &Values{}
	_ Reindenter = &Where{}
	_ Reindenter = &With{}
)

func write(buf *bytes.Buffer, token lexer.Token, indent int) error {
	switch {
	case token.IsNeedNewLineBefore():
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, strings.Repeat(DoubleWhiteSpace, indent), token.FormattedValue()))
	case token.Type == lexer.COMMA:
		buf.WriteString(token.FormattedValue())
	case token.Type == lexer.DO:
		buf.WriteString(fmt.Sprintf("%s%s%s", NewLine, token.FormattedValue(), WhiteSpace))
	case strings.HasPrefix(token.FormattedValue(), "::"):
		buf.WriteString(token.FormattedValue())
	case token.Type == lexer.WITH:
		buf.WriteString(fmt.Sprintf("%s%s", NewLine, token.FormattedValue()))
	default:
		buf.WriteString(fmt.Sprintf("%s%s", WhiteSpace, token.FormattedValue()))
	}

	return nil
}
