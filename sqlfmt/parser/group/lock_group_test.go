//nolint:dupl
package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentLockGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normalcase",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.LOCK, Value: "LOCK"},
				lexer.Token{Type: lexer.IDENT, Value: "table"},
				lexer.Token{Type: lexer.IN, Value: "IN"},
				lexer.Token{Type: lexer.IDENT, Value: "xxx"},
			},
			want: "\nLOCK table\nIN xxx",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		lock := NewLock(tt.tokenSource)

		if err := lock.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
