package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Returning clause
type Returning struct {
	elementReindenter
}

func NewReturning(element []Reindenter, opts ...Option) *Returning {
	return &Returning{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements
func (r *Returning) Reindent(buf *bytes.Buffer) error {
	elements, err := r.processPunctuation()
	if err != nil {
		return err
	}

	for _, el := range separate(elements) {
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
