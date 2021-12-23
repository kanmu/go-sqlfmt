package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentUpdateGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.UPDATE, Value: "UPDATE"},
				lexer.Token{Type: lexer.IDENT, Value: "something1"},
			},
			want: "\nUPDATE\n  something1",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		updateGroup := &Update{Element: tt.tokenSource}

		if err := updateGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
