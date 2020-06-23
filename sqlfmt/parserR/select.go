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

	var(
		idx int
		consumed int
		value interface{}
		err error
	)
	for {
		t := tokens[idx]
		if expr.endTType(t.Type) {
			return expr, idx, nil
		}

		value = t
		consumed = 1
		if idx > 0 {
			switch t.Type {
			case lexer.STARTPARENTHESIS:
				// TODO
			case lexer.FUNCTION:
				// TODO
			}
		}

		fmt.Println(err)
		expr.append(value)
		idx = nextIDX(idx, consumed)
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