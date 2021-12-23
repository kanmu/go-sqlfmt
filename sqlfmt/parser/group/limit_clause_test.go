package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentLimitGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normalcase",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.LIMIT, Value: "LIMIT"},
				lexer.Token{Type: lexer.IDENT, Value: "123"},
			},
			want: "\nLIMIT 123",
		},
		{
			name: "normalcase",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.OFFSET, Value: "OFFSET"},
			},
			want: "\nOFFSET",
		},
		{
			name: "normalcase",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.FETCH, Value: "FETCH"},
				lexer.Token{Type: lexer.FIRST, Value: "FIRST"},
			},
			want: "\nFETCH FIRST",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		limitGroup := &LimitClause{Element: tt.tokenSource}

		if err := limitGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
