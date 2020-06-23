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
	idx := 0
	for {
		t := restTokens[idx]

		switch t.Type {
		case lexer.SELECT:
			expr, consumed, err = parseSelect(restTokens[idx:])
		case lexer.FROM:
			expr, consumed, err = parseFrom(restTokens[idx:])
		case lexer.FUNCTION:
			expr, consumed, err = parseFunction(restTokens)
		case lexer.EOF:
			return exprs, nil
		}

		if err != nil{
			fmt.Println(err)
		}

		idx += consumed
		exprs = append(exprs, expr)
	}
}
