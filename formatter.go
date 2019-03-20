package sqlfmt

import (
	"fmt"
	"strings"

	"github.com/kanmu/go-sqlfmt/lexer"
	"github.com/kanmu/go-sqlfmt/parser"
	"github.com/pkg/errors"
)

// Formatter formats SQL statements
type Formatter struct{}

// Format formats src in 3 steps
// 1: tokenize src
// 2: parse tokens by SQL clause group
// 3: for each clause group (Reindenter), add indentation or new line in the correct position
func (f *Formatter) Format(src string) (string, error) {
	t := lexer.NewTokenizer(src)
	tokens, err := t.GetTokens()
	if err != nil {
		return src, errors.Wrapf(err, "Tokenize failed at:%#v", src)
	}

	rs, err := parser.ParseTokens(tokens)
	if err != nil {
		return src, errors.Wrapf(err, "ParseTokens failed at:%s", src)
	}

	w := NewWriter(rs)
	formattedStmt, err := w.Write()
	if err != nil {
		return "", errors.Wrapf(err, "Write failed at:%s", src)
	}

	if !compareStmtValue(src, formattedStmt) {
		return src, fmt.Errorf("Format failed at:%s", src)
	}
	return formattedStmt, nil
}

// returns false if the value of formatted statement  (without any space) differs from source statement
func compareStmtValue(stmt string, formattedStmt string) bool {
	before := removeSpace(stmt)
	after := removeSpace(formattedStmt)

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
