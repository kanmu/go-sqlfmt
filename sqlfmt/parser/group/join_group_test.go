package group

import (
	"bytes"
	"testing"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentJoinGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normalcase",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.LEFT, Value: "LEFT"},
				lexer.Token{Type: lexer.OUTER, Value: "OUTER"},
				lexer.Token{Type: lexer.JOIN, Value: "JOIN"},
				lexer.Token{Type: lexer.IDENT, Value: "sometable"},
				lexer.Token{Type: lexer.ON, Value: "ON"},
				lexer.Token{Type: lexer.IDENT, Value: "status1"},
				lexer.Token{Type: lexer.IDENT, Value: "="},
				lexer.Token{Type: lexer.IDENT, Value: "status2"},
			},

			want: "\nLEFT OUTER JOIN sometable\nON status1 = status2",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		joinGroup := &Join{Element: tt.tokenSource}

		joinGroup.Reindent(buf)
		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
