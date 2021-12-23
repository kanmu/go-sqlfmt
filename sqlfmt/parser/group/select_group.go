package group

import (
	"bytes"
	"fmt"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/pkg/errors"
)

// Select clause
type Select struct {
	Element     []Reindenter
	IndentLevel int
	baseReindenter
}

// Reindent reindens its elements
func (s *Select) Reindent(buf *bytes.Buffer) error {
	s.start = 0

	src, err := processPunctuation(s.Element)
	if err != nil {
		return err
	}
	elements := separate(src)

	for i, element := range elements {
		switch v := element.(type) {
		case lexer.Token, string:
			if err := writeSelect(buf, element, &s.start, s.IndentLevel); err != nil {
				return errors.Wrap(err, "writeSelect failed")
			}
		case *Case:
			if tok, ok := elements[i-1].(lexer.Token); ok {
				if tok.Type == lexer.COMMA {
					v.hasCommaBefore = true
				}
			}
			v.Reindent(buf)
			// Case group in Select clause must be in column area
			s.start++
		case *Parenthesis:
			v.InColumnArea = true
			v.ColumnCount = s.start
			v.Reindent(buf)
			s.start++
		case *Subquery:
			if token, ok := elements[i-1].(lexer.Token); ok {
				if token.Type == lexer.EXISTS {
					v.Reindent(buf)
					continue
				}
			}
			v.InColumnArea = true
			v.ColumnCount = s.start
			v.Reindent(buf)
		case *Function:
			v.InColumnArea = true
			v.ColumnCount = s.start
			v.Reindent(buf)
			s.start++
		case Reindenter:
			v.Reindent(buf)
			s.start++
		default:
			return fmt.Errorf("can not reindent %#v", v)
		}
	}
	return nil
}

// IncrementIndentLevel increments by its specified indent level
func (s *Select) IncrementIndentLevel(lev int) {
	s.IndentLevel += lev
}
