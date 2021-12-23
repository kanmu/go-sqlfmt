package group

import (
	"bytes"
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentSubqueryGroup(t *testing.T) {
	tests := []struct {
		name string
		src  []Reindenter
		want string
	}{
		{
			name: "normalcase",
			src: []Reindenter{
				lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
				&Select{
					Element: []Reindenter{
						lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
						lexer.Token{Type: lexer.IDENT, Value: "xxxxxx"},
					},
					IndentLevel: 1,
				},
				&From{
					Element: []Reindenter{
						lexer.Token{Type: lexer.FROM, Value: "FROM"},
						lexer.Token{Type: lexer.IDENT, Value: "xxxxxx"},
					},
					IndentLevel: 1,
				},
				lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
			},
			want: " (\n  SELECT\n    xxxxxx\n  FROM xxxxxx)",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		parenGroup := &Parenthesis{Element: tt.src, IndentLevel: 1}

		if err := parenGroup.Reindent(buf); err != nil {
			t.Errorf("unexpected error: %v", err)

			return
		}

		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}
