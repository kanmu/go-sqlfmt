package sqlfmt

import (
	"bytes"

	"github.com/kanmu/go-sqlfmt/parser/group"
)

// writer represents writer
type writer struct {
	buf *bytes.Buffer
	rs  []group.Reindenter
}

// NewWriter returns a pointer of writer
func NewWriter(rs []group.Reindenter) *writer {
	return &writer{
		buf: &bytes.Buffer{},
		rs:  rs,
	}
}

// Write writes formatted statement to buffer
func (w *writer) Write() (string, error) {
	for _, v := range w.rs {
		if err := v.Reindent(w.buf); err != nil {
			return "", err
		}
	}
	return w.buf.String(), nil
}
