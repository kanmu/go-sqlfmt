package parserR

import (
	"fmt"
	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

type SelectExpr struct {
	Values []interface{}
	SubQueryCnt int
}

func parseSelect(tokens []lexer.Token)(*SelectExpr, int, error){
	expr := &SelectExpr{}

	var (
		idx int
		value interface{}
		consumed int
		err error
	)
	for {
		token := tokens[idx]
		if expr.endTType(token.Type) {
			return expr, idx, nil
		}

		// if any expr appears from the second token, it should be parsed as one expr and consumed will be the count of tokens in the expr
		// in other cases, value will be the token and consumed will be 1
		value = token
		consumed = 1
		if idx > 0 {
			switch token.Type {
			case lexer.STARTPARENTHESIS:
				// TODO
			case lexer.FUNCTION:
				// TODO
			}
		}

		fmt.Println(err)
		expr.append(value)
		idx += consumed
	}
}

func (expr *SelectExpr) endTType(ttype lexer.TokenType) bool{
	for _, end := range lexer.EndOfSelect{
		if ttype == end {
			return true
		}
	}
	return false
}

func (expr *SelectExpr) append(elm interface{}) {
	expr.Values = append(expr.Values, elm)
}

func (expr *SelectExpr) Build() string {
	return ""
}