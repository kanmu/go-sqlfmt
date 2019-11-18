package parser

import (
	"reflect"
	"testing"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
	"github.com/kanmu/go-sqlfmt/sqlfmt/parser/group"
)

func TestParseTokens(t *testing.T) {
	testingData := []lexer.Token{
		{Type: lexer.SELECT, Value: "SELECT"},
		{Type: lexer.IDENT, Value: "name"},
		{Type: lexer.COMMA, Value: ","},
		{Type: lexer.IDENT, Value: "age"},
		{Type: lexer.COMMA, Value: ","},

		{Type: lexer.FUNCTION, Value: "SUM"},
		{Type: lexer.STARTPARENTHESIS, Value: "("},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.ENDPARENTHESIS, Value: ")"},

		{Type: lexer.STARTPARENTHESIS, Value: "("},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.ENDPARENTHESIS, Value: ")"},

		{Type: lexer.TYPE, Value: "TEXT"},
		{Type: lexer.STARTPARENTHESIS, Value: "("},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.ENDPARENTHESIS, Value: ")"},

		{Type: lexer.FROM, Value: "FROM"},
		{Type: lexer.IDENT, Value: "user"},
		{Type: lexer.WHERE, Value: "WHERE"},
		{Type: lexer.IDENT, Value: "name"},
		{Type: lexer.IDENT, Value: "="},
		{Type: lexer.STRING, Value: "'xxx'"},
		{Type: lexer.EOF, Value: "EOF"},
	}
	testingData2 := []lexer.Token{
		{Type: lexer.SELECT, Value: "SELECT"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.FROM, Value: "FROM"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.WHERE, Value: "WHERE"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.IN, Value: "IN"},
		{Type: lexer.STARTPARENTHESIS, Value: "("},
		{Type: lexer.SELECT, Value: "SELECT"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.FROM, Value: "FROM"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.JOIN, Value: "JOIN"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.ON, Value: "ON"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.IDENT, Value: "="},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.ENDPARENTHESIS, Value: ")"},
		{Type: lexer.GROUP, Value: "GROUP"},
		{Type: lexer.BY, Value: "BY"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.ORDER, Value: "ORDER"},
		{Type: lexer.BY, Value: "BY"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.LIMIT, Value: "LIMIT"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.UNION, Value: "UNION"},
		{Type: lexer.ALL, Value: "ALL"},
		{Type: lexer.SELECT, Value: "SELECT"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.FROM, Value: "FROM"},
		{Type: lexer.IDENT, Value: "xxx"},
		{Type: lexer.EOF, Value: "EOF"},
	}
	testingData3 := []lexer.Token{
		{Type: lexer.UPDATE, Value: "UPDATE"},
		{Type: lexer.IDENT, Value: "user"},
		{Type: lexer.SET, Value: "SET"},
		{Type: lexer.IDENT, Value: "point"},
		{Type: lexer.IDENT, Value: "="},
		{Type: lexer.IDENT, Value: "0"},
		{Type: lexer.EOF, Value: "EOF"},
	}

	tests := []struct {
		name        string
		tokenSource []lexer.Token
		want        []group.Reindenter
	}{
		{
			name:        "normal test case 1",
			tokenSource: testingData,
			want: []group.Reindenter{
				&group.Select{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
						lexer.Token{Type: lexer.IDENT, Value: "name"},
						lexer.Token{Type: lexer.COMMA, Value: ","},
						lexer.Token{Type: lexer.IDENT, Value: "age"},
						lexer.Token{Type: lexer.COMMA, Value: ","},
						&group.Function{
							Element: []group.Reindenter{
								lexer.Token{Type: lexer.FUNCTION, Value: "SUM"},
								lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
								lexer.Token{Type: lexer.IDENT, Value: "xxx"},
								lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
							},
						},
						&group.Parenthesis{
							Element: []group.Reindenter{
								lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
								lexer.Token{Type: lexer.IDENT, Value: "xxx"},
								lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
							},
						},
						&group.TypeCast{
							Element: []group.Reindenter{
								lexer.Token{Type: lexer.TYPE, Value: "TEXT"},
								lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
								lexer.Token{Type: lexer.IDENT, Value: "xxx"},
								lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
							},
						},
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
						lexer.Token{Type: lexer.IDENT, Value: "name"},
						lexer.Token{Type: lexer.IDENT, Value: "="},
						lexer.Token{Type: lexer.STRING, Value: "'xxx'"},
					},
				},
			},
		},
		{
			name:        "normal test case 2",
			tokenSource: testingData2,
			want: []group.Reindenter{
				&group.Select{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
				&group.From{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.FROM, Value: "FROM"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
				&group.Where{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.WHERE, Value: "WHERE"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
						lexer.Token{Type: lexer.IN, Value: "IN"},
						&group.Subquery{
							Element: []group.Reindenter{
								lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
								&group.Select{
									Element: []group.Reindenter{
										lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
									},
									IndentLevel: 1,
								},
								&group.From{
									Element: []group.Reindenter{
										lexer.Token{Type: lexer.FROM, Value: "FROM"},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
									},
									IndentLevel: 1,
								},
								&group.Join{
									Element: []group.Reindenter{
										lexer.Token{Type: lexer.JOIN, Value: "JOIN"},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
										lexer.Token{Type: lexer.ON, Value: "ON"},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
										lexer.Token{Type: lexer.IDENT, Value: "="},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
									},
									IndentLevel: 1,
								},
								lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
							},
							IndentLevel: 1,
						},
					},
				},
				&group.GroupBy{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.GROUP, Value: "GROUP"},
						lexer.Token{Type: lexer.BY, Value: "BY"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
				&group.OrderBy{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.ORDER, Value: "ORDER"},
						lexer.Token{Type: lexer.BY, Value: "BY"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
				&group.LimitClause{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.LIMIT, Value: "LIMIT"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
				&group.TieClause{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.UNION, Value: "UNION"},
						lexer.Token{Type: lexer.ALL, Value: "ALL"},
					},
				},
				&group.Select{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
				&group.From{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.FROM, Value: "FROM"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
			},
		},
		{
			name:        "normal test case 3",
			tokenSource: testingData3,
			want: []group.Reindenter{
				&group.Update{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.UPDATE, Value: "UPDATE"},
						lexer.Token{Type: lexer.IDENT, Value: "user"},
					},
				},
				&group.Set{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.SET, Value: "SET"},
						lexer.Token{Type: lexer.IDENT, Value: "point"},
						lexer.Token{Type: lexer.IDENT, Value: "="},
						lexer.Token{Type: lexer.IDENT, Value: "0"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got, err := ParseTokens(tt.tokenSource)
		if err != nil {
			t.Errorf("ERROR: %#v", err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("\nwant %#v, \ngot %#v", tt.want, got)
		}
	}
}
