package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Update clause
type Update struct {
	elementReindenter
}

func NewUpdate(element []Reindenter, opts ...Option) *Update {
	return &Update{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements
func (u *Update) Reindent(buf *bytes.Buffer) error {
	u.start = 0

	elements, err := processPunctuation(u.Element)
	if err != nil {
		return err
	}

	for _, el := range separate(elements) {
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
