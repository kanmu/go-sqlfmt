package group

import (
	"bytes"
	"testing"

	"github.com/kanmu/go-sqlfmt/lexer"
)

func TestReindentAndGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal test",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.ANDGROUP, Value: "AND"},
				lexer.Token{Type: lexer.IDENT, Value: "something1"},
				lexer.Token{Type: lexer.IDENT, Value: "something2"},
			},
			want: "\nAND something1 something2",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		andGroup := &AndGroup{Element: tt.tokenSource}

		if err := andGroup.Reindent(buf); err != nil {
			t.Errorf("error %#v", err)
		}
		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
