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

		}

        restTokens = restTokens[consumed:]
		exprs = append(exprs, expr)
	}

	return exprs, nil
}
