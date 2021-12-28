package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTokens(t *testing.T) {
	t.Parallel()

	options := defaultOptions()

	var testingSQLStatement = `select name, age,'age',sum,z+d^2, sum(case when x = xxx then false else true end), "old"::double Precision
	,"new"::bit varying(30), test::character varying(2)[]
		from user where name + xxx = 0 and 'age' = 'xxx' limit 100 except 100`

	want := []Token{
		{Type: SELECT, Value: "SELECT"},
		{Type: IDENT, Value: "name"},
		{Type: COMMA, Value: ","},
		{Type: IDENT, Value: "age"},
		{Type: COMMA, Value: ","},
		{Type: STRING, Value: "'age'"},
		{Type: COMMA, Value: ","},
		{Type: IDENT, Value: "sum"}, // this token is not considered a function
		{Type: COMMA, Value: ","},
		{Type: IDENT, Value: "z"},
		{Type: OPERATOR, Value: "+"},
		{Type: IDENT, Value: "d"},
		{Type: OPERATOR, Value: "^"},
		{Type: IDENT, Value: "2"},
		{Type: COMMA, Value: ","},
		{Type: FUNCTION, Value: "SUM"},
		{Type: STARTPARENTHESIS, Value: "("},
		{Type: CASE, Value: "CASE"},
		{Type: WHEN, Value: "WHEN"},
		{Type: IDENT, Value: "x"},
		{Type: OPERATOR, Value: "="},
		{Type: IDENT, Value: "xxx"},
		{Type: THEN, Value: "THEN"},
		{Type: RESERVEDVALUE, Value: "FALSE"},
		{Type: ELSE, Value: "ELSE"},
		{Type: RESERVEDVALUE, Value: "TRUE"},
		{Type: END, Value: "END"},
		{Type: ENDPARENTHESIS, Value: ")"},
		{Type: COMMA, Value: ","},
		{Type: STRING, Value: `"old"`},
		{Type: OPERATOR, Value: "::"},
		{Type: TYPE, Value: "DOUBLE PRECISION"},
		{Type: COMMA, Value: ","},
		{Type: STRING, Value: `"new"`},
		{Type: OPERATOR, Value: "::"},
		{Type: TYPE, Value: "BIT VARYING"},
		{Type: STARTPARENTHESIS, Value: "("},
		{Type: IDENT, Value: "30"},
		{Type: ENDPARENTHESIS, Value: ")"},
		{Type: COMMA, Value: ","},
		{Type: IDENT, Value: "test"},
		{Type: OPERATOR, Value: "::"},
		{Type: TYPE, Value: "CHARACTER VARYING"},
		{Type: STARTPARENTHESIS, Value: "("},
		{Type: IDENT, Value: "2"},
		{Type: ENDPARENTHESIS, Value: ")"},
		{Type: STARTBRACKET, Value: "["},
		{Type: ENDBRACKET, Value: "]"},

		{Type: FROM, Value: "FROM"},
		{Type: IDENT, Value: "user"},
		{Type: WHERE, Value: "WHERE"},
		{Type: IDENT, Value: "name"},
		{Type: OPERATOR, Value: "+"},
		{Type: IDENT, Value: "xxx"},
		{Type: OPERATOR, Value: "="},
		{Type: IDENT, Value: "0"},
		{Type: AND, Value: "AND"},
		{Type: STRING, Value: "'age'"},
		{Type: OPERATOR, Value: "="},
		{Type: STRING, Value: "'xxx'"},
		{Type: LIMIT, Value: "LIMIT"},
		{Type: IDENT, Value: "100"},
		{Type: EXCEPT, Value: "EXCEPT"},
		{Type: IDENT, Value: "100"},

		{Type: EOF, Value: "EOF"},
	}

	for i := range want {
		want[i].options = options
	}

	tnz := NewTokenizer(testingSQLStatement)
	tnz.options = options
	got, err := tnz.GetTokens()
	require.NoError(t, err)

	if assert.EqualValues(t, want, got) {
		return
	}

	// assert detailed diff
	assert.Lenf(t, got, len(want), "expected %d tokens, got %d", len(want), len(got))
	for i, token := range got {
		if i >= len(want) {
			break
		}

		assert.EqualValuesf(t, want[i], token, "unexpected token N°%d", i)
	}
}

func TestIsWhiteSpace(t *testing.T) {
	t.Parallel()

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
	for _, toPin := range tests {
		tt := toPin

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, isWhiteSpace(tt.src))
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

	for _, toPin := range tests {
		tt := toPin

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tnz := NewTokenizer(tt.src)

			got, err := tnz.scan()
			require.NoError(t, err)

			require.Equal(t, tt.want, got)
		})
	}
}

func TestScanWhiteSpace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want Token
		next rune
	}{
		{
			name: "whitespace test case 1",
			src:  ` `,
			want: Token{Type: WS, Value: " "},
			next: eof,
		},
		{
			name: "whitespace test case 2",
			src:  "\n",
			want: Token{Type: NEWLINE, Value: "\n"},
			next: eof,
		},
		{
			name: "whitespace test case 3",
			src:  "    \n    \r   x",
			want: Token{Type: NEWLINE, Value: "\n"},
			next: 'x',
		},
	}
	options := defaultOptions()
	for i := range tests {
		tests[i].want.options = options
	}

	for _, toPin := range tests {
		tt := toPin

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tnz := NewTokenizer(tt.src)
			tnz.options = options

			ch, err := tnz.Read()
			require.NoError(t, err)
			require.True(t, isWhiteSpace(ch))

			require.NoError(t, tnz.scanWhiteSpace(ch))

			require.EqualValues(t, tt.want, tnz.result[0])

			ch, err = tnz.Read()
			require.NoError(t, err)
			require.Equal(t, tt.next, ch)
		})
	}
}

func TestScanIdent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		src         string
		want        Token
		next        rune
		expectEmpty bool
	}{
		{
			name: "ident test case 1",
			src:  `select`,
			want: Token{Type: SELECT, Value: "SELECT"},
			next: eof,
		},
		{
			name: "ident test case 2",
			src:  "table",
			want: Token{Type: TABLE, Value: "TABLE"},
			next: eof,
		},
		{
			name: "ident test case 3",
			src:  "end",
			want: Token{Type: END, Value: "END"},
			next: eof,
		},
		{
			name: "ident test case 4",
			src:  "end(other)",
			want: Token{Type: END, Value: "END"},
			next: '(',
		},
		{
			name: "ident test case 5",
			src:  "end other",
			want: Token{Type: END, Value: "END"},
			next: ' ',
		},
		{
			name: "ident test case 6",
			src:  "end+other",
			want: Token{Type: END, Value: "END"},
			next: '+',
		},
		{
			name: "ident test case 7",
			src:  "foo->>other",
			want: Token{Type: IDENT, Value: "foo"},
			next: '-',
		},
		{
			name:        "ident test case 8",
			src:         "'foo'::other",
			expectEmpty: true, // not an ident
			next:        'f',
		},
	}
	options := defaultOptions()

	for _, toPin := range tests {
		tt := toPin

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.want.options = options
			tnz := NewTokenizer(tt.src)
			tnz.options = options

			ch, err := tnz.Read()
			require.NoError(t, err)

			require.NoError(t, tnz.scanIdent(ch))
			if tt.expectEmpty {
				require.Empty(t, tnz.result)
			} else {
				require.EqualValues(t, tt.want, tnz.result[0])
				require.NotEmpty(t, tnz.result)
			}

			ch, err = tnz.Read()
			require.NoError(t, err)

			require.Equalf(t, tt.next, ch, "expected next rune to be %q, but got %q", string(tt.next), string(ch))
		})
	}
}

func TestScanOperator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		src    string
		expect bool
		want   Token
		next   rune
	}{
		{
			name:   "operator +",
			src:    `+`,
			expect: true,
			want:   Token{Type: OPERATOR, Value: "+"},
			next:   eof,
		},
		{
			name:   "operator *",
			src:    `*y`,
			expect: true,
			want:   Token{Type: OPERATOR, Value: "*"},
			next:   'y',
		},
		{
			name:   "operator ::",
			src:    "::",
			expect: true,
			want:   Token{Type: OPERATOR, Value: "::"},
			next:   eof,
		},
		{
			name:   "operator ::(type)",
			src:    "::bit(3)",
			expect: true,
			want:   Token{Type: OPERATOR, Value: "::"},
			next:   'b',
		},
		{
			name:   "non operator x",
			src:    "x",
			expect: false,
			next:   'x',
		},
		{
			name:   "operator -",
			src:    "-",
			expect: true,
			want:   Token{Type: OPERATOR, Value: "-"},
			next:   eof,
		},
		{
			name:   "operator ->>",
			src:    "->>x",
			expect: true,
			want:   Token{Type: OPERATOR, Value: "->>"},
			next:   'x',
		},
		{
			name:   "operator !=",
			src:    "!= x",
			expect: true,
			want:   Token{Type: OPERATOR, Value: "!="},
			next:   ' ',
		},
		{
			name:   "non operator !@",
			src:    "!@ x",
			expect: false,
			next:   '!',
		},
	}
	options := defaultOptions()
	for i := range tests {
		tests[i].want.options = options
	}

	for _, toPin := range tests {
		tt := toPin

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tnz := NewTokenizer(tt.src)
			tnz.options = options

			ch, err := tnz.Read()
			require.NoError(t, err)

			if tt.expect {
				require.True(t, isOperator(ch))
			}

			ok, err := tnz.scanOperator(ch)
			require.NoError(t, err)
			if !tt.expect {
				require.False(t, ok)

				return
			}

			require.True(t, ok)
			require.EqualValues(t, tt.want, tnz.result[0])

			ch, err = tnz.Read()
			require.NoError(t, err)

			require.Equalf(t, tt.next, ch, "expected next rune to be %q, got %q", string(tt.next), string(ch))
		})
	}
}
