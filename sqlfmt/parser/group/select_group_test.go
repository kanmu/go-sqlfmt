package group

import (
	"bytes"
	"testing"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

func TestReindentSelectGroup(t *testing.T) {
	tests := []struct {
		name        string
		tokenSource []Reindenter
		want        string
	}{
		{
			name: "normal case",
			tokenSource: []Reindenter{
				lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
				lexer.Token{Type: lexer.IDENT, Value: "name"},
				lexer.Token{Type: lexer.COMMA, Value: ","},
				lexer.Token{Type: lexer.IDENT, Value: "age"},
			},
			want: "\nSELECT\n  name\n  , age",
		},
	}
	for _, tt := range tests {
		buf := &bytes.Buffer{}
		selectGroup := &Select{Element: tt.tokenSource}

		selectGroup.Reindent(buf)
		got := buf.String()
		if tt.want != got {
			t.Errorf("want%#v, got %#v", tt.want, got)
		}
	}
}

func TestIncrementIndentLevel(t *testing.T) {
	s := &Select{}
	s.IncrementIndentLevel(1)
	got := s.IndentLevel
	want := 1
	if got != want {
		t.Errorf("want %#v got %#v", want, got)
	}
}
