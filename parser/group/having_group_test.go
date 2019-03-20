package group

import (
	"bytes"
	"testing"

	"github.com/kanmu/go-sqlfmt/lexer"
)

func TestReindentHavingGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.HAVING, Value: "HAVING"},
				lexer.Token{Type: lexer.IDENT, Value: "xxxxxxxx"},
			},
			want: "\nHAVING xxxxxxxx",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		havingGroup := &Having{Element: tt.tokenSource}

		havingGroup.Reindent(buf)
		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
