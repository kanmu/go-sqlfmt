package sqlfmt

import (
	"testing"

	"github.com/kanmu/go-sqlfmt/lexer"
	"github.com/kanmu/go-sqlfmt/parser/group"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name string
		src  []group.Reindenter
		want string
	}{
		{
			name: "normal test case 1",
			src: []group.Reindenter{
				&group.Select{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
						lexer.Token{Type: lexer.IDENT, Value: "name"},
						lexer.Token{Type: lexer.COMMA, Value: ","},
						lexer.Token{Type: lexer.IDENT, Value: "age"},
					},
				},
				&group.From{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.FROM, Value: "FROM"},
						lexer.Token{Type: lexer.IDENT, Value: "user"},
					},
				},
				&group.Where{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.WHERE, Value: "WHERE"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
			},
			want: "\nSELECT\n  name\n  , age\nFROM user\nWHERE xxx",
		},
	}

	for _, tt := range tests {
		w := NewWriter(tt.src)
		got, err := w.Write()
		if err != nil {
			t.Errorf("ERROR:%#v", err)
		}
		if got != tt.want {
			t.Errorf("\nwant %#v, \ngot %#v", tt.want, got)
		}
	}

}
