package group

import (
	"bytes"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// Lock clause
type Lock struct {
	Element     []Reindenter
	IndentLevel int
	baseReindenter
}

// Reindent reindent its elements
func (l *Lock) Reindent(buf *bytes.Buffer) error {
	for _, v := range l.Element {
		if token, ok := v.(lexer.Token); ok {
			writeLock(buf, token)
		} else {
			if err := v.Reindent(buf); err != nil {
				return err
			}
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified increment level
func (l *Lock) IncrementIndentLevel(lev int) {
	l.IndentLevel += lev
}
