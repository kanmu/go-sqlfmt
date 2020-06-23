package parserR

import "github.com/kanmu/go-sqlfmt/sqlfmt/lexer"

type ParenthesisExpr struct {
	Values []interface{}
	Parent Expr
	SubQueryCnt int
}

func parseParenthesis(tokens []lexer.Token)(*ParenthesisExpr, int, error){
	var (
		expr = &ParenthesisExpr{}
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
		case lexer.SELECT:
			// ParseSubquery的な関数を読んだら良さそう
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

func (expr *ParenthesisExpr) endTType(ttype lexer.TokenType) bool{
	for _, end := range lexer.EndOfParenthesis{
		if ttype == end {
			return true
		}
	}
	return false
}

func (f *ParenthesisExpr) Build() string {
	return ""
}


