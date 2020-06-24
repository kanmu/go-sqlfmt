package parserR

import (
	"fmt"
	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

type FromExpr struct {
	Values []interface{}
	Parent Expr
	SubQueryCnt int
}

func parseFrom(tokens []lexer.Token)(*FromExpr, int, error){
	expr := &FromExpr{}
	var (
		idx int
		value interface{}
		consumed int
		err error
	)
	// parseのそれぞれの関数がExprとconsumeしたcntだけを返すというインターフェースはそれで良さそう
	for {
		t := tokens[idx]
		if expr.endTType(t.Type) {
			return expr, idx, nil
		}

		value = t
		consumed = 1
		if idx > 0{
			switch t.Type {
			case lexer.STARTPARENTHESIS:
			case lexer.FUNCTION:
			}
		}

		fmt.Println(err)
		expr.append(value)
		idx += consumed
	}
}

func (expr *FromExpr) endTType(ttype lexer.TokenType) bool{
	for _, end := range lexer.EndOfFrom{
		if ttype == end {
			return true
		}
	}
	return false
}

func (f *FromExpr) Build() string {
	return ""
}

func (f *FromExpr) append(elm interface{}) {
	f.Values = append(f.Values, elm)
}