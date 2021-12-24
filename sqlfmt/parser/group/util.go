package group

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
)

// separate elements by comma and the reserved word in select clause.
func separate(rs []Reindenter) []interface{} {
	var (
		result           []interface{}
		skipRange, count int
	)
	buf := &bytes.Buffer{}

	for _, r := range rs {
		switch token := r.(type) {
		case lexer.Token:
			switch {
			case skipRange > 0:
				skipRange--

			case token.IsKeyWordInSelect():
				// TODO: more elegant
				if buf.String() != "" {
					result = append(result, buf.String())
					buf.Reset()
				}
				result = append(result, token)

			case token.Type == lexer.COMMA:
				if buf.String() != "" {
					result = append(result, buf.String())
				}
				result = append(result, token)
				buf.Reset()
				count = 0

			case strings.HasPrefix(token.FormattedValue(), "::"):
				buf.WriteString(token.FormattedValue())

			default:
				if count == 0 {
					buf.WriteString(token.FormattedValue())
				} else {
					buf.WriteString(WhiteSpace + token.FormattedValue())
				}

				count++
			}

		default:
			if buf.String() != "" {
				result = append(result, buf.String())
				buf.Reset()
			}
			result = append(result, r)
		}
	}

	// append the last element in buf
	if buf.String() != "" {
		result = append(result, buf.String())
	}

	return result
}

// process bracket, singlequote and brace
// TODO: more elegant.
func processPunctuation(rs []Reindenter) ([]Reindenter, error) {
	var (
		result    []Reindenter
		skipRange int
	)

	for i, v := range rs {
		if token, ok := v.(lexer.Token); ok {
			switch {
			case skipRange > 0:
				skipRange--
			case token.Type == lexer.STARTBRACE || token.Type == lexer.STARTBRACKET:
				surrounding, sr, err := extractSurroundingArea(rs[i:])
				if err != nil {
					return nil, err
				}
				result = append(result, lexer.Token{
					Type:  lexer.SURROUNDING,
					Value: surrounding,
				})
				skipRange += sr
			default:
				result = append(result, token)
			}
		} else {
			result = append(result, v)
		}
	}

	return result, nil
}

// returns surrounding area including punctuation such as {xxx, xxx}.
func extractSurroundingArea(rs []Reindenter) (string, int, error) {
	var (
		countOfStart int
		countOfEnd   int
		result       string
		skipRange    int
	)

	for i, r := range rs {
		if token, ok := r.(lexer.Token); ok {
			switch {
			case token.Type == lexer.COMMA || token.Type == lexer.STARTBRACKET || token.Type == lexer.STARTBRACE || token.Type == lexer.ENDBRACKET || token.Type == lexer.ENDBRACE:
				result += fmt.Sprint(token.FormattedValue())
				// for next token of StartToken
			case i == 1:
				result += fmt.Sprint(token.FormattedValue())
			default:
				result += fmt.Sprint(WhiteSpace + token.FormattedValue())
			}

			if token.Type == lexer.STARTBRACKET || token.Type == lexer.STARTBRACE || token.Type == lexer.STARTPARENTHESIS {
				countOfStart++
			}
			if token.Type == lexer.ENDBRACKET || token.Type == lexer.ENDBRACE || token.Type == lexer.ENDPARENTHESIS {
				countOfEnd++
			}
			if countOfStart == countOfEnd {
				break
			}
			skipRange++
		} else {
			// TODO: should support group type in surrounding area?
			// I have not encountered any groups in surrounding area so far
			return "", -1, errors.New("group type is not supposed be here")
		}
	}

	return result, skipRange, nil
}
