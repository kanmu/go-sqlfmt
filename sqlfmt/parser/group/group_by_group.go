package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// GroupBy clause
// nolint:revive
type GroupBy struct {
	Element     []Reindenter
	IndentLevel int
	baseReindenter
}

// Reindent reindents its elements
func (g *GroupBy) Reindent(buf *bytes.Buffer) error {
	g.start = 0

	elements, err := processPunctuation(g.Element)
	if err != nil {
		return err
	}

	for _, el := range separate(elements) {
		switch v := el.(type) {
		case lexer.Token, string:
			if erw := writeWithComma(buf, v, &g.start, g.IndentLevel); erw != nil {
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

// IncrementIndentLevel increments by its specified indent level
func (g *GroupBy) IncrementIndentLevel(lev int) {
	g.IndentLevel += lev
}
