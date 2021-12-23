//nolint:dupl
package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentInsertGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normalcase",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.INSERT, Value: "INSERT"},
				lexer.Token{Type: lexer.INTO, Value: "INTO"},
				lexer.Token{Type: lexer.IDENT, Value: "xxxxxx"},
				lexer.Token{Type: lexer.IDENT, Value: "xxxxxx"},
			},
			want: "\nINSERT INTO xxxxxx xxxxxx",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		insertGroup := &Insert{Element: tt.tokenSource}

		if err := insertGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
