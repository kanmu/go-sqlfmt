package group

import (
	"bytes"
	"testing"

	"github.com/kanmu/go-sqlfmt/lexer"
)

func TestReindentOrGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normalcase",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.ORGROUP, Value: "OR"},
				lexer.Token{Type: lexer.IDENT, Value: "something1"},
				lexer.Token{Type: lexer.IDENT, Value: "something2"},
			},
			want: "\nOR something1 something2",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		orGroup := &OrGroup{Element: tt.tokenSource}

		orGroup.Reindent(buf)
		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
