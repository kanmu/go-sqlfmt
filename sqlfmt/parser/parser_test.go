package parser

import (
	"reflect"
	"testing"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
	"github.com/kanmu/go-sqlfmt/sqlfmt/parser/group"
)

func TestParseTokens(t *testing.T) {
	testingData := []lexer.Token{
		lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
		lexer.Token{Type: lexer.IDENT, Value: "name"},
		lexer.Token{Type: lexer.COMMA, Value: ","},
		lexer.Token{Type: lexer.IDENT, Value: "age"},
		lexer.Token{Type: lexer.COMMA, Value: ","},

		lexer.Token{Type: lexer.FUNCTION, Value: "SUM"},
		lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},

		lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},

		lexer.Token{Type: lexer.TYPE, Value: "TEXT"},
		lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},

		lexer.Token{Type: lexer.FROM, Value: "FROM"},
		lexer.Token{Type: lexer.IDENT, Value: "user"},
		lexer.Token{Type: lexer.WHERE, Value: "WHERE"},
		lexer.Token{Type: lexer.IDENT, Value: "name"},
		lexer.Token{Type: lexer.IDENT, Value: "="},
		lexer.Token{Type: lexer.STRING, Value: "'xxx'"},
		lexer.Token{Type: lexer.EOF, Value: "EOF"},
	}
	testingData2 := []lexer.Token{
		lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.FROM, Value: "FROM"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.WHERE, Value: "WHERE"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.IN, Value: "IN"},
		lexer.Token{Type: lexer.STARTPARENTHESIS, Value: "("},
		lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.FROM, Value: "FROM"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.JOIN, Value: "JOIN"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.ON, Value: "ON"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.IDENT, Value: "="},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.ENDPARENTHESIS, Value: ")"},
		lexer.Token{Type: lexer.GROUP, Value: "GROUP"},
		lexer.Token{Type: lexer.BY, Value: "BY"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.ORDER, Value: "ORDER"},
		lexer.Token{Type: lexer.BY, Value: "BY"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.LIMIT, Value: "LIMIT"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.UNION, Value: "UNION"},
		lexer.Token{Type: lexer.ALL, Value: "ALL"},
		lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.FROM, Value: "FROM"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.EOF, Value: "EOF"},
	}
	testingData3 := []lexer.Token{
		lexer.Token{Type: lexer.UPDATE, Value: "UPDATE"},
		lexer.Token{Type: lexer.IDENT, Value: "user"},
		lexer.Token{Type: lexer.SET, Value: "SET"},
		lexer.Token{Type: lexer.IDENT, Value: "point"},
		lexer.Token{Type: lexer.IDENT, Value: "="},
		lexer.Token{Type: lexer.IDENT, Value: "0"},
		lexer.Token{Type: lexer.EOF, Value: "EOF"},
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

func TestParseSelectStmt(t *testing.T) {
	testingData := []lexer.Token{
		lexer.Token{Type: lexer.SELECT, Value: "SELECT"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.FROM, Value: "FROM"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.WHERE, Value: "WHERE"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.EOF, Value: "EOF"},
	}
	tests := []struct {
		tokenSource []lexer.Token
		want        []group.Reindenter
	}{
		{
			tokenSource: testingData,
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
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got, err := ParseTokens(tt.tokenSource)
		if err != nil {
			t.Fatalf("ERROR: %#v", err)
		}
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("want %#v got %#v", tt.want, got)
		}
	}
}

func TestParseUpdateStmt(t *testing.T) {
	testingData := []lexer.Token{
		lexer.Token{Type: lexer.UPDATE, Value: "UPDATE"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.SET, Value: "SET"},
		lexer.Token{Type: lexer.IDENT, Value: "status"},
		lexer.Token{Type: lexer.IDENT, Value: "="},
		lexer.Token{Type: lexer.IDENT, Value: "$1"},
		lexer.Token{Type: lexer.WHERE, Value: "WHERE"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.IDENT, Value: "="},
		lexer.Token{Type: lexer.IDENT, Value: "'xxx'"},
		lexer.Token{Type: lexer.RETURNING, Value: "RETURNING"},
		lexer.Token{Type: lexer.IDENT, Value: "'xxx'"},
		lexer.Token{Type: lexer.EOF, Value: "EOF"},
	}
	tests := []struct {
		tokenSource []lexer.Token
		want        []group.Reindenter
	}{
		{
			tokenSource: testingData,
			want: []group.Reindenter{
				&group.Update{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.UPDATE, Value: "UPDATE"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
				&group.Set{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.SET, Value: "SET"},
						lexer.Token{Type: lexer.IDENT, Value: "status"},
						lexer.Token{Type: lexer.IDENT, Value: "="},
						lexer.Token{Type: lexer.IDENT, Value: "$1"},
					},
				},
				&group.Where{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.WHERE, Value: "WHERE"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
						lexer.Token{Type: lexer.IDENT, Value: "="},
						lexer.Token{Type: lexer.IDENT, Value: "'xxx'"},
					},
				},
				&group.Returning{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.RETURNING, Value: "RETURNING"},
						lexer.Token{Type: lexer.IDENT, Value: "'xxx'"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got, err := ParseTokens(tt.tokenSource)
		if err != nil {
			t.Fatalf("ERROR: %#v", err)
		}
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("want %#v got %#v", tt.want, got)
		}
	}
}

func TestParseDeleteStmt(t *testing.T) {
	testingData := []lexer.Token{
		lexer.Token{Type: lexer.DELETE, Value: "DELETE"},
		lexer.Token{Type: lexer.FROM, Value: "FROM"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.EOF, Value: "EOF"},
	}
	tests := []struct {
		tokenSource []lexer.Token
		want        []group.Reindenter
	}{
		{
			tokenSource: testingData,
			want: []group.Reindenter{
				&group.Delete{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.DELETE, Value: "DELETE"},
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
	}
	for _, tt := range tests {
		got, err := ParseTokens(tt.tokenSource)
		if err != nil {
			t.Fatalf("ERROR: %#v", err)
		}
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("want %#v got %#v", tt.want, got)
		}
	}
}

func TestParseInsertStmt(t *testing.T) {
	testingData := []lexer.Token{
		lexer.Token{Type: lexer.INSERT, Value: "INSERT"},
		lexer.Token{Type: lexer.INTO, Value: "INTO"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.VALUES, Value: "VALUES"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.ON, Value: "ON"},
		lexer.Token{Type: lexer.IDENT, Value: "conflict"},
		lexer.Token{Type: lexer.IDENT, Value: "xxx"},
		lexer.Token{Type: lexer.EOF, Value: "EOF"},
	}
	tests := []struct {
		tokenSource []lexer.Token
		want        []group.Reindenter
	}{
		{
			tokenSource: testingData,
			want: []group.Reindenter{
				&group.Insert{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.INSERT, Value: "INSERT"},
						lexer.Token{Type: lexer.INTO, Value: "INTO"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
				&group.Values{
					Element: []group.Reindenter{
						lexer.Token{Type: lexer.VALUES, Value: "VALUES"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
						lexer.Token{Type: lexer.ON, Value: "ON"},
						lexer.Token{Type: lexer.IDENT, Value: "conflict"},
						lexer.Token{Type: lexer.IDENT, Value: "xxx"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got, err := ParseTokens(tt.tokenSource)
		if err != nil {
			t.Fatalf("ERROR: %#v", err)
		}
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("want %#v got %#v", tt.want, got)
		}
	}
}
