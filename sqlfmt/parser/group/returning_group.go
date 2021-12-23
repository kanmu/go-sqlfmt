package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Returning clause
type Returning struct {
	Element     []Reindenter
	IndentLevel int
	baseReindenter
}

// Reindent reindents its elements
func (r *Returning) Reindent(buf *bytes.Buffer) error {
	src, err := processPunctuation(r.Element)
	if err != nil {
		return err
	}

	for _, el := range separate(src) {
		switch v := el.(type) {
		case lexer.Token, string:
			if erw := writeWithComma(buf, v, &r.start, r.IndentLevel); erw != nil {
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
func (r *Returning) IncrementIndentLevel(lev int) {
	r.IndentLevel += lev
}
