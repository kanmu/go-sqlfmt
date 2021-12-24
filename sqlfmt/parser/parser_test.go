package parser

import (
	"testing"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				group.NewSelect(
					[]group.Reindenter{
						lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
						lexer.Token{Type: lexer.IDENT, Value: "name"},
						lexer.Token{Type: lexer.COMMA, Value: ","},
						lexer.Token{Type: lexer.IDENT, Value: "age"},
						lexer.Token{Type: lexer.COMMA, Value: ","},
						group.NewFunction(
							[]group.Reindenter{
								lexer.Token{Type: lexer.FUNCTION, Value: "SUM"},
								lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
								lexer.Token{Type: lexer.IDENT, Value: "xxx"},
								lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
							},
						),
						group.NewParenthesis(
							[]group.Reindenter{
								lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
								lexer.Token{Type: lexer.IDENT, Value: "xxx"},
								lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
							},
						),
						group.NewTypeCast(
							[]group.Reindenter{
								lexer.Token{Type: lexer.TYPE, Value: "TEXT"},
								lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
								lexer.Token{Type: lexer.IDENT, Value: "xxx"},
								lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
							},
						),
					},
				),
				group.NewFrom(
					[]group.Reindenter{
						lexer.Token{Type: lexer.FROM, Value: "FROM"},
						lexer.Token{Type: lexer.IDENT, Value: "user"},
					},
				),
				group.NewWhere(
					[]group.Reindenter{
						lexer.Token{Type: lexer.WHERE, Value: "WHERE"},
						lexer.Token{Type: lexer.IDENT, Value: "name"},
						lexer.Token{Type: lexer.IDENT, Value: "="},
						lexer.Token{Type: lexer.STRING, Value: "'xxx'"},
					},
				),
			},
		},
		{
			name:        "normal test case 2",
			tokenSource: testingData2,
			want: []group.Reindenter{
				group.NewSelect(
					[]group.Reindenter{
						lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				),
				group.NewFrom(
					[]group.Reindenter{
						lexer.Token{Type: lexer.FROM, Value: "FROM"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				),
				group.NewWhere(
					[]group.Reindenter{
						lexer.Token{Type: lexer.WHERE, Value: "WHERE"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
						lexer.Token{Type: lexer.IN, Value: "IN"},
						group.NewSubquery(
							[]group.Reindenter{
								lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
								group.NewSelect(
									[]group.Reindenter{
										lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
									},
									group.WithIndentLevel(1),
								),
								group.NewFrom(
									[]group.Reindenter{
										lexer.Token{Type: lexer.FROM, Value: "FROM"},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
									},
									group.WithIndentLevel(1),
								),
								group.NewJoin(
									[]group.Reindenter{
										lexer.Token{Type: lexer.JOIN, Value: "JOIN"},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
										lexer.Token{Type: lexer.ON, Value: "ON"},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
										lexer.Token{Type: lexer.IDENT, Value: "="},
										lexer.Token{Type: lexer.IDENT, Value: "xxx"},
									},
									group.WithIndentLevel(1),
								),
								lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
							},
							group.WithIndentLevel(1),
						),
					},
				),
				group.NewGroupBy(
					[]group.Reindenter{
						lexer.Token{Type: lexer.GROUP, Value: "GROUP"},
						lexer.Token{Type: lexer.BY, Value: "BY"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				),
				group.NewOrderBy(
					[]group.Reindenter{
						lexer.Token{Type: lexer.ORDER, Value: "ORDER"},
						lexer.Token{Type: lexer.BY, Value: "BY"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				),
				group.NewLimitClause(
					[]group.Reindenter{
						lexer.Token{Type: lexer.LIMIT, Value: "LIMIT"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				),
				group.NewTieClause(
					[]group.Reindenter{
						lexer.Token{Type: lexer.UNION, Value: "UNION"},
						lexer.Token{Type: lexer.ALL, Value: "ALL"},
					},
				),
				group.NewSelect(
					[]group.Reindenter{
						lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				),
				group.NewFrom(
					[]group.Reindenter{
						lexer.Token{Type: lexer.FROM, Value: "FROM"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				),
			},
		},
		{
			name:        "normal test case 3",
			tokenSource: testingData3,
			want: []group.Reindenter{
				group.NewUpdate(
					[]group.Reindenter{
						lexer.Token{Type: lexer.UPDATE, Value: "UPDATE"},
						lexer.Token{Type: lexer.IDENT, Value: "user"},
					},
				),
				group.NewSet(
					[]group.Reindenter{
						lexer.Token{Type: lexer.SET, Value: "SET"},
						lexer.Token{Type: lexer.IDENT, Value: "point"},
						lexer.Token{Type: lexer.IDENT, Value: "="},
						lexer.Token{Type: lexer.IDENT, Value: "0"},
					},
				),
			},
		},
	}
	for _, toPin := range tests {
		tt := toPin
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTokens(tt.tokenSource)
			require.NoError(t, err)

			for i, actual := range got {
				expected := tt.want[i]

				assert.EqualValuesf(t, expected, actual, "\ngroup[%d]\nwant %#v, \ngot %#v", i, expected, actual)
			}
		})
	}
}
