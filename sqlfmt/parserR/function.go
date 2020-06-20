package parserR

import "github.com/kanmu/go-sqlfmt/sqlfmt/lexer"

type FunctionExpr struct {
	Values []interface{}
	Parent Expr
	SubQueryCnt int
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

func (f FunctionExpr) Build() string {
	return ""
}
