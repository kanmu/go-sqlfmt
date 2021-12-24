//nolint:dupl
package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentValuesGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.VALUES, Value: "VALUES"},
				lexer.Token{Type: lexer.IDENT, Value: "xxxxx"},
				lexer.Token{Type: lexer.ON, Value: "ON"},
				lexer.Token{Type: lexer.IDENT, Value: "xxxxx"},
				lexer.Token{Type: lexer.DO, Value: "DO"},
			},
			want: "\nVALUES xxxxx\nON xxxxx\nDO ",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		valuesGroup := NewValues(tt.tokenSource)

		if err := valuesGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
