package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentDeleteGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.DELETE, Value: "DELETE"},
				lexer.Token{Type: lexer.FROM, Value: "FROM"},
				lexer.Token{Type: lexer.IDENT, Value: "xxxxxx"},
			},
			want: "\nDELETE\nFROM xxxxxx",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		deleteGroup := &Delete{Element: tt.tokenSource}

		if err := deleteGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
