package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// GroupBy clause
// nolint:revive
type GroupBy struct {
	elementReindenter
}

func NewGroupBy(element []Reindenter, opts ...Option) *GroupBy {
	return &GroupBy{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements.
func (g *GroupBy) Reindent(buf *bytes.Buffer) error {
	g.start = 0

	elements, err := g.processPunctuation()
	if err != nil {
		return err
	}

	for _, el := range separate(elements) {
		switch v := el.(type) {
		case lexer.Token, string:
			if erw := g.writeWithComma(buf, v, g.IndentLevel); erw != nil {
				return erw
			}
		case Reindenter:
			if eri := v.Reindent(buf); eri != nil {
				return eri
			}
		}
	}

	return nil
}
