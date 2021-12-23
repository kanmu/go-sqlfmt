package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// OrderBy clause
type OrderBy struct {
	Element     []Reindenter
	IndentLevel int
	baseReindenter
}

// Reindent reindents its elements
func (o *OrderBy) Reindent(buf *bytes.Buffer) error {
	o.start = 0

	src, err := processPunctuation(o.Element)
	if err != nil {
		return err
	}

	for _, el := range separate(src) {
		switch v := el.(type) {
		case lexer.Token, string:
			if erw := writeWithComma(buf, v, &o.start, o.IndentLevel); erw != nil {
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
func (o *OrderBy) IncrementIndentLevel(lev int) {
	o.IndentLevel += lev
}
