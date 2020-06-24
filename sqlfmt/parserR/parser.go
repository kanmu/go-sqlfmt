package parserR

import (
	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

type Expr interface {
	Build()string
}

type Result struct {
	Values []Expr
}

func (pr *Result) Build() string {
	return ""
}

func ParseTokens(tokens []lexer.Token) ([]Expr, error) {
	rslt := &Result{}
	var (
		idx int
		expr Expr
		consumed int
		err error
	)

	for {
		t := tokens[idx]
		if rslt.endTType(t.Type){
			return rslt.Values, nil
		}

		switch t.Type {
		case lexer.SELECT:
			expr, consumed, err = parseSelect(tokens[idx:])
		case lexer.FROM:
			expr, consumed, err = parseFrom(tokens[idx:])
		case lexer.FUNCTION:
		case lexer.EOF:
			return rslt.Values, nil
		}
		if err != nil{
			return nil, err
		}

		rslt.append(expr)
		idx += consumed
	}
}


func (rslt *Result) append(elm Expr){
	rslt.Values = append(rslt.Values, elm)
}

func (rslt *Result) endTType(ttype lexer.TokenType) bool {
	if ttype == lexer.EOF{
		return true
	}
	return false
}