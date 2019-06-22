package sqlfmt

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
	"github.com/kanmu/go-sqlfmt/sqlfmt/parser"
	"github.com/kanmu/go-sqlfmt/sqlfmt/parser/group"
	"github.com/pkg/errors"
)

// Format formats src in 3 steps
// 1: tokenize src
// 2: parse tokens by SQL clause group
// 3: for each clause group (Reindenter), add indentation or new line in the correct position
func Format(src string, distance int) (string, error) {
	t := lexer.NewTokenizer(src)
	tokens, err := t.GetTokens()
	if err != nil {
		return src, errors.Wrap(err, "Tokenize failed")
	}

	rs, err := parser.ParseTokens(tokens)
	if err != nil {
		return src, errors.Wrap(err, "ParseTokens failed")
	}

	res, err := getFormattedStmt(rs, distance)
	if err != nil {
		return src, errors.Wrap(err, "getFormattedStmt failed")
	}

	if !compare(src, res) {
		return src, fmt.Errorf("the formatted statement has diffed from the source")
	}
	return res, nil
}

func getFormattedStmt(rs []group.Reindenter, distance int) (string, error) {
	var buf bytes.Buffer

	for _, r := range rs {
		if err := r.Reindent(&buf); err != nil {
			return "", errors.Wrap(err, "Reindent failed")
		}
	}

	res := strings.NewReader(buf.String())
	s := bufio.NewScanner(res)
	var result string
	for s.Scan() {
		r := fmt.Sprintf("%s%s", strings.Repeat(" ", distance), s.Text())
		result += r + "\n"
	}

	return result, nil
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
