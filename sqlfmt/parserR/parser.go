package parserR

import (
	"fmt"
	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

type Expr interface {
	Build()string
}

func ParseTokens(tokens []lexer.Token) ([]Expr, error) {
	var (
		err error
		expr Expr
		exprs []Expr
		consumed int
	)

	restTokens := tokens
	for t := restTokens[0]; t.Type == lexer.EOF; {
		switch t.Type {
		case lexer.FUNCTION:
			expr, consumed, err = parseFunction(restTokens)
			if err != nil{
				fmt.Println(err)
			}
		case lexer.IDENT:
			// ...
		}

        restTokens = restTokens[consumed:]
		exprs = append(exprs, expr)
	}

	return exprs, nil
}

// parentがない場合はnilを入れる
type FunctionExpr struct {
	Values []interface{}
	Parent Expr
	SubQueryCnt int
}

func (f FunctionExpr) Build() string {
	return ""
}

func parseFunction(tokens []lexer.Token)(*FunctionExpr, int, error){
	var (
		expr = &FunctionExpr{}
	    consumed = 0
	    restTokens = tokens
	)

	// parseのそれぞれの関数がExprとconsumeしたcntだけを返すというインターフェースはそれで良さそう
	for t := restTokens[0]; expr.endTType(t.Type); {
		switch t.Type {
		case restTokens[0].Type:
			// 一発目は自分自身をパースしてまうので、そのままTokenを入れておく
			expr.Values = append(expr.Values, t)
			consumed ++
		case lexer.FUNCTION:
			cExpr, cConsumed, err := parseFunction(tokens[consumed:])
			if err != nil {
				// FIXME: エラーハンドリングする
				return nil, 0, err
			}

			cExpr.Parent = expr
			expr.Values = append(expr.Values, cExpr)
			consumed += cConsumed
		default:
			expr.Values = append(expr.Values, t)
			consumed ++
		}
		restTokens = restTokens[consumed:]
	}

	return expr, consumed, nil
}

func (expr *FunctionExpr) endTType(ttype lexer.TokenType) bool{
	for _, end := range lexer.EndOfFunction{
		if ttype == end {
			return true
		}
	}
	return false
}