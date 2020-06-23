package parserR

import (
	"fmt"
	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
	"testing"
)


func TestParseTokens(t *testing.T) {
	testTokens := []lexer.Token{
		{Type: lexer.SELECT, Value: "SELECT"},
		{Type: lexer.IDENT, Value: "name"},
		{Type: lexer.COMMA, Value: ","},
		{Type: lexer.IDENT, Value: "age"},
		{Type: lexer.FROM, Value: "FROM"},
		{Type: lexer.IDENT, Value: "user"},
		{Type: lexer.EOF, Value: "EOF"},
	}

	res, err := ParseTokens(testTokens)
	if err != nil {
		fmt.Println("Error")
		t.Fatal(err)
	}
	fmt.Println(res)
}