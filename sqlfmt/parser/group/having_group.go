package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Having clause
type Having struct {
	Element     []Reindenter
	IndentLevel int
	baseReindenter
}

// Reindent reindents its elements
func (h *Having) Reindent(buf *bytes.Buffer) error {
	elements, err := processPunctuation(h.Element)
	if err != nil {
		return err
	}

	for _, el := range elements {
		if token, ok := el.(lexer.Token); ok {
			write(buf, token, h.IndentLevel)
		} else {
			if eri := el.Reindent(buf); eri != nil {
				return eri
			}
		}
	}

	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (h *Having) IncrementIndentLevel(lev int) {
	h.IndentLevel += lev
}
