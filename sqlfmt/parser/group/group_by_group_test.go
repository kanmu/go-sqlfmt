package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
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
		groupByGroup := NewGroupBy(tt.tokenSource)

		if err := groupByGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
