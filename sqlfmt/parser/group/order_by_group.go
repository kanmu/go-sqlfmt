package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// OrderBy clause
type OrderBy struct {
	elementReindenter
}

func NewOrderBy(element []Reindenter, opts ...Option) *OrderBy {
	return &OrderBy{
		elementReindenter: newElementReindenter(element, opts...),
	}
}

// Reindent reindents its elements
func (o *OrderBy) Reindent(buf *bytes.Buffer) error {
	o.start = 0

	element, err := o.processPunctuation()
	if err != nil {
		return err
	}

	for _, el := range separate(element) {
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
