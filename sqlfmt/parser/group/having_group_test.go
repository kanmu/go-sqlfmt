package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
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
		havingGroup := NewHaving(tt.tokenSource)

		if err := havingGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
