package parserR

import (
	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

type SelectExpr struct {
	Values []interface{}
	Parent Expr
	SubQueryCnt int
}

func parseSelect(tokens []lexer.Token)(*SelectExpr, int, error){
	var (
		expr = &SelectExpr{}
		consumed = 0
		restTokens = tokens
	)

	// parseのそれぞれの関数がExprとconsumeしたcntだけを返すというインターフェースはそれで良さそう
	idx := 0
	for {
		t := restTokens[idx]

		if expr.endTType(t.Type) {
			return expr, idx, nil
		}


		// 一番最初のトークンはそのままアペンド
		// これはでも、fanctionの時しか必要ない？
		if idx == 0 {
			expr.Values = append(expr.Values, t)
			idx++
		} else {
			switch t.Type {
			case lexer.STARTPARENTHESIS:
				parseParenthesis(restTokens)
			case lexer.FUNCTION:
				cExpr, consumed, err := parseFunction(tokens[consumed:])
				if err != nil {
					// FIXME: エラーハンドリングする
					return nil, 0, err
				}

				cExpr.Parent = expr
				expr.Values = append(expr.Values, cExpr)
				idx += consumed
			default:
				expr.Values = append(expr.Values, t)
				idx++
			}
		}
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

func (f *SelectExpr) Build() string {
	return ""
}