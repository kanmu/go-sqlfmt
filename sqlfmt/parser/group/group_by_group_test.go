package group

import (
	"bytes"
	"testing"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentGroupByGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.GROUP, Value: "GROUP"},
				lexer.Token{Type: lexer.BY, Value: "BY"},
				lexer.Token{Type: lexer.IDENT, Value: "xxxxxx"},
			},
			want: "\nGROUP BY\n  xxxxxx",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		groupByGroup := &GroupBy{Element: tt.tokenSource}

		groupByGroup.Reindent(buf)
		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
