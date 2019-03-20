package group

import (
	"bytes"
	"testing"

	"github.com/kanmu/go-sqlfmt/lexer"
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
		returningGroup := &Returning{Element: tt.tokenSource}

		returningGroup.Reindent(buf)
		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
