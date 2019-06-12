package sqlfmt

import (
	"fmt"
	"strings"

	"bytes"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
	"github.com/kanmu/go-sqlfmt/sqlfmt/parser"
	"github.com/pkg/errors"
)

var (
	buf bytes.Buffer
	res string
)

// Format formats src in 3 steps
// 1: tokenize src
// 2: parse tokens by SQL clause group
// 3: for each clause group (Reindenter), add indentation or new line in the correct position
func Format(src string) (string, error) {
	t := lexer.NewTokenizer(src)
	tokens, err := t.GetTokens()
	if err != nil {
		return src, errors.Wrap(err, "Tokenize failed")
	}

	rs, err := parser.ParseTokens(tokens)
	if err != nil {
		return src, errors.Wrap(err, "ParseTokens failed")
	}

	for _, r := range rs {
		if err := r.Reindent(&buf); err != nil {
			return src, errors.Wrap(err, "Reindent failed")
		}
	}
	res = buf.String()

	if !compare(src, res) {
		return src, fmt.Errorf("value of formatted statement has diffed from source statement")
	}
	return res, nil
}

// returns false if the value of formatted statement  (without any space) differs from source statement
func compare(src string, res string) bool {
	before := removeSpace(src)
	after := removeSpace(res)

	if v := strings.Compare(before, after); v != 0 {
		return false
	}
	return true
}

// removes whitespaces and new lines from src
func removeSpace(src string) string {
	var result []rune
	for _, r := range src {
		if string(r) == "\n" || string(r) == " " || string(r) == "\t" || string(r) == "　" {
			continue
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
