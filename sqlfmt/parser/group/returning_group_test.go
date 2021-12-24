//nolint:dupl
package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentReturningGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.RETURNING, Value: "RETURNING"},
				lexer.Token{Type: lexer.IDENT, Value: "something1"},
				lexer.Token{Type: lexer.COMMA, Value: ","},
				lexer.Token{Type: lexer.IDENT, Value: "something1"},
			},
			want: "\nRETURNING\n  something1\n  , something1",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		returningGroup := NewReturning(tt.tokenSource)

		if err := returningGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
