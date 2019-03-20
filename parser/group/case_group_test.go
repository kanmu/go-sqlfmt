package group

import (
	"bytes"
	"testing"

	"github.com/kanmu/go-sqlfmt/lexer"
)

func TestReindentCaseGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.CASE, Value: "CASE"},
				lexer.Token{Type: lexer.WHEN, Value: "WHEN"},
				lexer.Token{Type: lexer.IDENT, Value: "something"},
				lexer.Token{Type: lexer.IDENT, Value: "something"},
				lexer.Token{Type: lexer.END, Value: "END"},
			},
			want: "\n  CASE\n     WHEN something something\n  END",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		caseGroup := &Case{Element: tt.tokenSource}

		caseGroup.Reindent(buf)
		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
