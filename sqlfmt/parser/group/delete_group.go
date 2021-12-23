package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Delete clause
type Delete struct {
	Element     []Reindenter
	IndentLevel int
	baseReindenter
}

// Reindent reindents its elements
func (d *Delete) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(d.Element)
	if err != nil {
		return err
	}
	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, d.IndentLevel)
		} else {
			if eri := el.Reindent(buf); eri != nil {
				return eri
			}
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (d *Delete) IncrementIndentLevel(lev int) {
	d.IndentLevel += lev
}
