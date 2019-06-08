package lexer

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetTokens(t *testing.T) {
	var testingSQLStatement = strings.Trim(`select name, age,sum, sum(case xxx) from user where name xxx and age = 'xxx' limit 100 except 100`, "`")
	want := []Token{
		Token{Type: SELECT, Value: "SELECT"},
		Token{Type: IDENT, Value: "name"},
		Token{Type: COMMA, Value: ","},
		Token{Type: IDENT, Value: "age"},
		Token{Type: COMMA, Value: ","},
		Token{Type: IDENT, Value: "SUM"},
		Token{Type: COMMA, Value: ","},
		Token{Type: FUNCTION, Value: "SUM"},
		Token{Type: STARTPARENTHESIS, Value: "("},
		Token{Type: CASE, Value: "CASE"},
		Token{Type: IDENT, Value: "xxx"},
		Token{Type: ENDPARENTHESIS, Value: ")"},

		Token{Type: FROM, Value: "FROM"},
		Token{Type: IDENT, Value: "user"},
		Token{Type: WHERE, Value: "WHERE"},
		Token{Type: IDENT, Value: "name"},
		Token{Type: IDENT, Value: "xxx"},
		Token{Type: AND, Value: "AND"},
		Token{Type: IDENT, Value: "age"},
		Token{Type: IDENT, Value: "="},
		Token{Type: STRING, Value: "'xxx'"},
		Token{Type: LIMIT, Value: "LIMIT"},
		Token{Type: IDENT, Value: "100"},
		Token{Type: EXCEPT, Value: "EXCEPT"},
		Token{Type: IDENT, Value: "100"},

		Token{Type: EOF, Value: "EOF"},
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
