package lexer

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetTokens(t *testing.T) {
	var testingSQLStatement = strings.Trim(`select name, age,sum, sum(case xxx) from user where name xxx and age = 'xxx' limit 100 except 100`, "`")
	want := []Token{
		{Type: SELECT, Value: "SELECT"},
		{Type: IDENT, Value: "name"},
		{Type: COMMA, Value: ","},
		{Type: IDENT, Value: "age"},
		{Type: COMMA, Value: ","},
		{Type: IDENT, Value: "SUM"},
		{Type: COMMA, Value: ","},
		{Type: FUNCTION, Value: "SUM"},
		{Type: STARTPARENTHESIS, Value: "("},
		{Type: CASE, Value: "CASE"},
		{Type: IDENT, Value: "xxx"},
		{Type: ENDPARENTHESIS, Value: ")"},

		{Type: FROM, Value: "FROM"},
		{Type: IDENT, Value: "user"},
		{Type: WHERE, Value: "WHERE"},
		{Type: IDENT, Value: "name"},
		{Type: IDENT, Value: "xxx"},
		{Type: AND, Value: "AND"},
		{Type: IDENT, Value: "age"},
		{Type: IDENT, Value: "="},
		{Type: STRING, Value: "'xxx'"},
		{Type: LIMIT, Value: "LIMIT"},
		{Type: IDENT, Value: "100"},
		{Type: EXCEPT, Value: "EXCEPT"},
		{Type: IDENT, Value: "100"},

		{Type: EOF, Value: "EOF"},
	}
	tnz := NewTokenizer(testingSQLStatement)
	got, err := tnz.GetTokens()
	if err != nil {
		t.Fatalf("\nERROR: %#v", err)
	} else if !reflect.DeepEqual(want, got) {
		t.Errorf("\nwant %#v, \ngot %#v", want, got)
	}
}

func TestIsWhiteSpace(t *testing.T) {
	tests := []struct {
		name string
		src  rune
		want bool
	}{
		{
			name: "normal test case 1",
			src:  '\n',
			want: true,
		},
		{
			name: "normal test case 2",
			src:  '\t',
			want: true,
		},
		{
			name: "normal test case 3",
			src:  ' ',
			want: true,
		},
		{
			name: "abnormal case",
			src:  'a',
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isWhiteSpace(tt.src); got != tt.want {
				t.Errorf("\nwant %v, \ngot %v", tt.want, got)
			}
		})
	}
}

func TestScan(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want bool
	}{
		{
			name: "normal test case 1",
			src:  `select`,
			want: false,
		},
		{
			name: "normal test case 2",
			src:  `table`,
			want: false,
		},
		{
			name: "normal test case 3",
			src:  ` `,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tnz := NewTokenizer(tt.src)

			got, err := tnz.scan()
			if err != nil {
				t.Errorf("\nERROR: %#v", err)
			}
			if got != tt.want {
				t.Errorf("\nwant %v, \ngot %v", tt.want, got)
			}
		})
	}
}

func TestScanWhiteSpace(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want Token
	}{
		{
			name: "normal test case 1",
			src:  ` `,
			want: Token{Type: WS, Value: " "},
		},
		{
			name: "normal test case 2",
			src:  "\n",
			want: Token{Type: NEWLINE, Value: "\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tnz := NewTokenizer(tt.src)
			tnz.scanWhiteSpace()

			if got := tnz.result[0]; got != tt.want {
				t.Errorf("\nwant %v, \ngot %v", tt.want, got)
			}
		})
	}
}

func TestScanIdent(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want Token
	}{
		{
			name: "normal test case 1",
			src:  `select`,
			want: Token{Type: SELECT, Value: "SELECT"},
		},
		{
			name: "normal test case 2",
			src:  "table",
			want: Token{Type: IDENT, Value: "table"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tnz := NewTokenizer(tt.src)
			tnz.scanIdent()

			if got := tnz.result[0]; got != tt.want {
				t.Errorf("\nwant %v, \ngot %v", tt.want, got)
			}
		})
	}
}
