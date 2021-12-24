package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentFromGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.FROM, Value: "FROM"},
				lexer.Token{Type: lexer.IDENT, Value: "sometable"},
			},
			want: "\nFROM sometable",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		fromGroup := NewFrom(tt.tokenSource)

		if err := fromGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
