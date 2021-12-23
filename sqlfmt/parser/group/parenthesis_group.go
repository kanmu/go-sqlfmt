package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Parenthesis clause
type Parenthesis struct {
	Element      []Reindenter
	IndentLevel  int
	InColumnArea bool
	ColumnCount  int
	baseReindenter
}

// Reindent reindents its elements
func (p *Parenthesis) Reindent(buf *bytes.Buffer) error {
	var hasStartBefore bool

	elements, err := processPunctuation(p.Element)
	if err != nil {
		return err
	}
	for i, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			hasStartBefore = (i == 1)
			writeParenthesis(buf, token, p.IndentLevel, p.ColumnCount, p.InColumnArea, hasStartBefore)
		} else {
			if eri := el.Reindent(buf); eri != nil {
				return eri
			}
		}
	}

	return nil
}

// IncrementIndentLevel indents by its specified indent level
func (p *Parenthesis) IncrementIndentLevel(lev int) {
	p.IndentLevel += lev
}
