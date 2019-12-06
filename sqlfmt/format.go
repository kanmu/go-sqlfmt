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

// Format format SQL
func Format(src string, options *Options) (string, error) {
	tokens, err := lexer.Tokenize(src)
	if err != nil {
		return src, errors.Wrap(err, "Tokenize failed")
	}

	rs, err := parser.ParseTokens(tokens)
	if err != nil {
		return src, errors.Wrap(err, "ParseTokens failed")
	}
	// ここでプリンタを用意して、そこに読み込みを行う

	res, err := getFormattedStmt(rs, options.Distance)
	if err != nil {
		return src, errors.Wrap(err, "getFormattedStmt failed")
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

	if distance != 0 {
		return putDistance(buf.String(), distance), nil
	}
	return buf.String(), nil
}

func putDistance(src string, distance int) string {
	scanner := bufio.NewScanner(strings.NewReader(src))

	var result string
	for scanner.Scan() {
		result += fmt.Sprintf("%s%s%s", strings.Repeat(group.WhiteSpace, distance), scanner.Text(), "\n")
	}
	return result
}
