package lexer

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/pkg/errors"
)

type tokenizer struct {
	r *bufio.Reader
}

// Tokenize tokenize src to Token, ignorig white-space, new-line and tab
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

// tokenize tokenize the source
func (t *tokenizer) tokenize() ([]Token, error) {
	var tokens []Token
	for {
		token, err := t.scanIgnoreSpace()
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
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

// scan reads the first charactor of t.r and creates Token
func (t *tokenizer) scanIgnoreSpace() (Token, error) {
	ch, _, err := t.r.ReadRune()
	// create EOF Token if end of file
	if err != nil {
		if err.Error() == "EOF" {
			return Token{Type: EOF, Value: "EOF"}, nil
		}
		return Token{}, err
	}

	// read until end of space
	if isSpace(ch) {
		if err := t.readUntilIdent(); err != nil {
			return Token{}, err
		}
	}

	// create punctuation Token if ch represents some punctuation
	if isPunctuation(ch) {
		token := createPunctuationToken(ch)
		return token, nil
	}

	var buf bytes.Buffer
	switch {
	// scan string surrounded by single quote
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

func (t *tokenizer) readUntilIdent() error {
	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			}
			return err
		}
		if !isSpace(ch) {
			// to scan the first charactor of next token
			t.unread()
			return nil
		}
	}
}

// scanString scans values surrounded by singleQuote such as 'xxxxxxxx'
// scanString writes rune of singleQuote to buf until the last single quote appears
func (t *tokenizer) scanString(buf *bytes.Buffer) (Token, error) {
	// read and write the first charactor before scanning so that it can ignore the first single quote and read until the last single-quote appears
	// TODO: more elegant way to scan string in the SQL
	sq, _, err := t.r.ReadRune()
	if err != nil {
		return Token{}, err
	}
	buf.WriteRune(sq)

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
		if isPunctuation(ch) || isSpace(ch) {
			t.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	upperValue := strings.ToUpper(buf.String())
	if ttype, ok := sqlKeywordMap[upperValue]; ok {
		return Token{Type: ttype, Value: upperValue}, nil
	}
	return Token{Type: IDENT, Value: buf.String()}, nil
}

var sqlKeywordMap = map[string]TokenType{
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

func isPunctuation(ch rune) bool {
	return isStartParenthesis(ch) || isEndParenthesis(ch) || isStartBracket(ch) || isEndBracket(ch) || isStartBrace(ch) || isEndBrace(ch) || isComma(ch) || isSingleQuote(ch)
}

func createPunctuationToken(ch rune) Token {
	switch {
	case isComma(ch):
		return Token{Type: COMMA, Value: string(ch)}
	case isStartParenthesis(ch):
		return Token{Type: STARTPARENTHESIS, Value: string(ch)}
	case isEndParenthesis(ch):
		return Token{Type: ENDPARENTHESIS, Value: string(ch)}
	case isStartBracket(ch):
		return Token{Type: STARTBRACKET, Value: string(ch)}
	case isEndBracket(ch):
		return Token{Type: ENDBRACKET, Value: string(ch)}
	case isStartBrace(ch):
		return Token{Type: STARTBRACE, Value: string(ch)}
	case isEndBrace(ch):
		return Token{Type: ENDBRACE, Value: string(ch)}
	}
	return Token{}
}
