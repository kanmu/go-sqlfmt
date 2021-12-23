package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Update clause
type Update struct {
	Element     []Reindenter
	IndentLevel int
	baseReindenter
}

// Reindent reindents its elements
func (u *Update) Reindent(buf *bytes.Buffer) error {
	u.start = 0

	src, err := processPunctuation(u.Element)
	if err != nil {
		return err
	}

	for _, el := range separate(src) {
		switch v := el.(type) {
		case lexer.Token, string:
			if erw := writeWithComma(buf, v, &u.start, u.IndentLevel); erw != nil {
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
func (u *Update) IncrementIndentLevel(lev int) {
	u.IndentLevel += lev
}
