package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type tokenizer struct {
	r *bufio.Reader
}

// Tokenize tokenize src and returns slice of Token
// It ignores Token of white-space, new-line and tab
func Tokenize(src string) ([]Token, error) {
	t := &tokenizer{
		r: bufio.NewReader(strings.NewReader(src)),
	}

	tokens, err := t.tokenize()
	if err != nil {
		return nil, errors.Wrap(err, "failed to tokenize")
	}

	return tokens, nil
}

// scan until END OF FILE
func (t *tokenizer) tokenize() ([]Token, error) {
	var tokens []Token
	for {
		token, err := t.scan()
		if err != nil {
			return nil, err
		}

		// ignorig space (white-space, new-line and tab)
		// go-sqlfmt formats src consistent with any space forcibly so far, but should I make a option to choose whether to ignore space..?
		if !(token.Type == SPACE) {
			tokens = append(tokens, token)
		}
		if token.Type == EOF {
			return tokens, nil
		}
	}
}

// unread undoes t.r.readRune method to get last character
func (t *tokenizer) unread() error {
	if err := t.r.UnreadRune(); err != nil {
		return err
	}
	return nil
}

// firstCharactor returns the first charactor of t.r without reading t.r
func (t *tokenizer) firstCharactor() (rune, error) {
	ch, _, err := t.r.ReadRune()
	if err != nil {
		return ch, err
	}

	// unread one charactor consumed already
	t.unread()
	return ch, nil
}

// scan reads the first charactor of t.r and creates Token
func (t *tokenizer) scan() (Token, error) {
	ch, err := t.firstCharactor()

	// create EOF Token if END OF FILE
	if err != nil {
		if err.Error() == "EOF" {
			return Token{Type: EOF, Value: "EOF"}, nil
		}
		return Token{}, err
	}

	var buf bytes.Buffer
	switch {
	case isSpace(ch):
		token, err := t.scanSpace(&buf)
		if err != nil {
			return Token{}, err
		}
		return token, nil
	case isPunctuation(ch):
		token, err := t.scanPunctuation(&buf)
		if err != nil {
			return Token{}, err
		}
		return token, nil
	// scan string surrounded by single quote such as 'xxxxxxxx'
	case isSingleQuote(ch):
		token, err := t.scanString(&buf)
		if err != nil {
			return Token{}, err
		}
		return token, nil
	default:
		token, err := t.scanIdent(&buf)
		if err != nil {
			return Token{}, err
		}
		return token, err
	}
}

// create token of space
func (t *tokenizer) scanSpace(buf *bytes.Buffer) (Token, error) {
	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return Token{}, err
			}
		}
		if !isSpace(ch) {
			t.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Token{Type: SPACE, Value: buf.String()}, nil
}

// create token of punctuation
func (t *tokenizer) scanPunctuation(buf *bytes.Buffer) (Token, error) {
	// token of punctuation is consisted of one charactor, so it reads t.r once except DOUBLECOLON token
	ch, _, err := t.r.ReadRune()
	if err != nil {
		return Token{}, err
	}
	buf.WriteRune(ch)

	// create token of colon or double-colon
	// TODO: more elegant
	if isColon(ch) {
		nextCh, _, err := t.r.ReadRune()
		if err != nil {
			return Token{}, err
		}
		// double-colon
		if isColon(nextCh) {
			return Token{Type: DOUBLECOLON, Value: fmt.Sprintf("%s%s", string(ch), string(nextCh))}, nil
		} else {
			// it already read the charactor of next token when colon does not appear twice
			// t.unread() makes it possible for caller function to scan next charactor that consumed already
			t.unread()
			return Token{Type: COLON, Value: string(ch)}, nil
		}
	}

	if ttype, ok := punctuationMap[buf.String()]; ok {
		return Token{Type: ttype, Value: buf.String()}, nil
	}

	return Token{}, fmt.Errorf("unexpected value: %v", buf.String())
}

// create token of string
// scan value surrounded with single-quote and return STRING token
func (t *tokenizer) scanString(buf *bytes.Buffer) (Token, error) {
	// read and write the first charactor before scanning so that it can ignore the first single quote and read until the last single-quote appears
	// TODO: more elegant way to scan string in the SQL
	sq, _, err := t.r.ReadRune()
	if err != nil {
		return Token{}, err
	}
	buf.WriteRune(sq)

	// read until next single-quote appears
	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return Token{}, err
			}
		}

		buf.WriteRune(ch)
		if isSingleQuote(ch) {
			break
		}
	}

	return Token{Type: STRING, Value: buf.String()}, nil
}

// create token of iden
// append all ch to result until ch is a white-space, new-line or punctuation
// if ident is SQL keyword, it returns Token of the keyword
func (t *tokenizer) scanIdent(buf *bytes.Buffer) (Token, error) {
	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				return Token{}, err
			}
		}
		if isPunctuation(ch) || isSpace(ch) || isSingleQuote(ch) {
			t.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	upperValue := strings.ToUpper(buf.String())
	if ttype, ok := keywordMap[upperValue]; ok {
		return Token{Type: ttype, Value: upperValue}, nil
	}

	return Token{Type: IDENT, Value: buf.String()}, nil
}

var keywordMap = map[string]TokenType{
	"SELECT":          SELECT,
	"FROM":            FROM,
	"WHERE":           WHERE,
	"CASE":            CASE,
	"ORDER":           ORDER,
	"BY":              BY,
	"AS":              AS,
	"JOIN":            JOIN,
	"LEFT":            LEFT,
	"RIGHT":           RIGHT,
	"INNER":           INNER,
	"OUTER":           OUTER,
	"ON":              ON,
	"WHEN":            WHEN,
	"END":             END,
	"GROUP":           GROUP,
	"DESC":            DESC,
	"ASC":             ASC,
	"LIMIT":           LIMIT,
	"AND":             AND,
	"OR":              OR,
	"IN":              IN,
	"IS":              IS,
	"NOT":             NOT,
	"NULL":            NULL,
	"DISTINCT":        DISTINCT,
	"LIKE":            LIKE,
	"BETWEEN":         BETWEEN,
	"UNION":           UNION,
	"ALL":             ALL,
	"HAVING":          HAVING,
	"EXISTS":          EXISTS,
	"UPDATE":          UPDATE,
	"SET":             SET,
	"RETURNING":       RETURNING,
	"DELETE":          DELETE,
	"INSERT":          INSERT,
	"INTO":            INTO,
	"DO":              DO,
	"VALUES":          VALUES,
	"FOR":             FOR,
	"THEN":            THEN,
	"ELSE":            ELSE,
	"DISTINCTROW":     DISTINCTROW,
	"FILTER":          FILTER,
	"WITHIN":          WITHIN,
	"COLLATE":         COLLATE,
	"INTERSECT":       INTERSECT,
	"EXCEPT":          EXCEPT,
	"OFFSET":          OFFSET,
	"FETCH":           FETCH,
	"FIRST":           FIRST,
	"ROWS":            ROWS,
	"USING":           USING,
	"OVERLAPS":        OVERLAPS,
	"NATURAL":         NATURAL,
	"CROSS":           CROSS,
	"ZONE":            ZONE,
	"NULLS":           NULLS,
	"LAST":            LAST,
	"AT":              AT,
	"LOCK":            LOCK,
	"WITH":            WITH,
	"SUM":             FUNCTION,
	"AVG":             FUNCTION,
	"MAX":             FUNCTION,
	"MIN":             FUNCTION,
	"COUNT":           FUNCTION,
	"COALESCE":        FUNCTION,
	"EXTRACT":         FUNCTION,
	"OVERLAY":         FUNCTION,
	"POSITION":        FUNCTION,
	"CAST":            FUNCTION,
	"SUBSTRING":       FUNCTION,
	"TRIM":            FUNCTION,
	"XMLELEMENT":      FUNCTION,
	"XMLFOREST":       FUNCTION,
	"XMLCONCAT":       FUNCTION,
	"RANDOM":          FUNCTION,
	"DATE_PART":       FUNCTION,
	"DATE_TRUNC":      FUNCTION,
	"ARRAY_AGG":       FUNCTION,
	"PERCENTILE_DISC": FUNCTION,
	"GREATEST":        FUNCTION,
	"LEAST":           FUNCTION,
	"OVER":            FUNCTION,
	"ROW_NUMBER":      FUNCTION,
	"BIG":             TYPE,
	"BIGSERIAL":       TYPE,
	"BOOLEAN":         TYPE,
	"CHAR":            TYPE,
	"BIT":             TYPE,
	"TEXT":            TYPE,
	"INTEGER":         TYPE,
	"NUMERIC":         TYPE,
	"DECIMAL":         TYPE,
	"DEC":             TYPE,
	"FLOAT":           TYPE,
	"CUSTOMTYPE":      TYPE,
	"VARCHAR":         TYPE,
	"VARBIT":          TYPE,
	"TIMESTAMP":       TYPE,
	"TIME":            TYPE,
	"SECOND":          TYPE,
	"INTERVAL":        TYPE,
}

var punctuationMap = map[string]TokenType{
	"(": STARTPARENTHESIS,
	")": ENDPARENTHESIS,
	"[": STARTBRACKET,
	"]": ENDBRACKET,
	"{": STARTBRACE,
	"}": ENDBRACKET,
	",": COMMA,
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '　'
}

func isTab(ch rune) bool {
	return ch == '\t'
}

func isNewLine(ch rune) bool {
	return ch == '\n'
}

func isSpace(ch rune) bool {
	return isWhiteSpace(ch) || isNewLine(ch) || isTab(ch)
}

func isComma(ch rune) bool {
	return ch == ','
}

func isStartParenthesis(ch rune) bool {
	return ch == '('
}

func isEndParenthesis(ch rune) bool {
	return ch == ')'
}

func isSingleQuote(ch rune) bool {
	return ch == '\''
}

func isStartBracket(ch rune) bool {
	return ch == '['
}

func isEndBracket(ch rune) bool {
	return ch == ']'
}

func isStartBrace(ch rune) bool {
	return ch == '{'
}

func isEndBrace(ch rune) bool {
	return ch == '}'
}

func isColon(ch rune) bool {
	return ch == ':'
}

func isPunctuation(ch rune) bool {
	return isStartParenthesis(ch) || isEndParenthesis(ch) || isStartBracket(ch) || isEndBracket(ch) || isStartBrace(ch) || isEndBrace(ch) || isComma(ch) || isColon(ch)
}