package lexer

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/pkg/errors"
)

// Tokenizer tokenizes SQL statements.
type Tokenizer struct {
	r      *bufio.Reader
	w      *bytes.Buffer // w  writes token value. It resets its value when the end of token appears
	result []Token
	*options
}

const (
	ErrEOF = "EOF"
)

// rune that can't be contained in SQL statement
// TODO: I have to make better solution of making rune of eof in stead of using '∂'.
var eof rune

func init() {
	eof = '∂'
}

// value of literal.
const (
	Comma            = ","
	StartParenthesis = "("
	EndParenthesis   = ")"
	StartBracket     = "["
	EndBracket       = "]"
	StartBrace       = "{"
	EndBrace         = "}"
	SingleQuote      = "'"
	NewLine          = "\n"
)

// NewTokenizer creates Tokenizer.
func NewTokenizer(src string, opts ...Option) *Tokenizer {
	return &Tokenizer{
		r:       bufio.NewReader(strings.NewReader(src)),
		w:       &bytes.Buffer{},
		options: defaultOptions(opts...),
	}
}

// GetTokens returns tokens for parsing.
func (t *Tokenizer) GetTokens() ([]Token, error) {
	tokens, err := t.Tokenize()
	if err != nil {
		return nil, errors.Wrap(err, "Tokenize failed")
	}

	result := make([]Token, 0, len(tokens))

	// replace all tokens without whitespaces and new lines
	// if "AND" or "OR" appears after new line, token value will be ANDGROUP, ORGROUP
	for i, tok := range tokens {
		if tok.Type == AND && tokens[i-1].Type == NEWLINE {
			andGroupToken := Token{Type: ANDGROUP, Value: tok.Value, options: t.options}
			result = append(result, andGroupToken)

			continue
		}

		if tok.Type == OR && tokens[i-1].Type == NEWLINE {
			orGroupToken := Token{Type: ORGROUP, Value: tok.Value, options: t.options}
			result = append(result, orGroupToken)

			continue
		}

		if tok.Type == WS || tok.Type == NEWLINE {
			continue
		}

		result = append(result, tok)
	}

	return result, nil
}

// Tokenize analyses every rune in SQL statement
// every token is identified when whitespace appears.
func (t *Tokenizer) Tokenize() ([]Token, error) {
	for {
		isEOF, err := t.scan()

		if isEOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	return t.result, nil
}

// unread undoes t.r.readRune method to get last character.
func (t *Tokenizer) unread() { _ = t.r.UnreadRune() }

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '　'
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

// scan scans each character and appends to result until "eof" appears
// when it finishes scanning all characters, it returns true.
func (t *Tokenizer) scan() (bool, error) {
	ch, _, err := t.r.ReadRune()
	if err != nil {
		if err.Error() != ErrEOF {
			return false, errors.Wrap(err, "read rune failed")
		}

		ch = eof
	}

	switch {
	case ch == eof:
		tok := Token{Type: EOF, Value: "EOF", options: t.options}
		t.result = append(t.result, tok)

		return true, nil
	case isWhiteSpace(ch):
		if err := t.scanWhiteSpace(); err != nil {
			return false, err
		}

		return false, nil
	// extract string
	case isSingleQuote(ch):
		if err := t.scanString(); err != nil {
			return false, err
		}

		return false, nil
	case isComma(ch):
		token := Token{Type: COMMA, Value: Comma, options: t.options}
		t.result = append(t.result, token)

		return false, nil
	case isStartParenthesis(ch):
		token := Token{Type: STARTPARENTHESIS, Value: StartParenthesis, options: t.options}
		t.result = append(t.result, token)

		return false, nil
	case isEndParenthesis(ch):
		token := Token{Type: ENDPARENTHESIS, Value: EndParenthesis, options: t.options}
		t.result = append(t.result, token)

		return false, nil
	case isStartBracket(ch):
		token := Token{Type: STARTBRACKET, Value: StartBracket, options: t.options}
		t.result = append(t.result, token)

		return false, nil
	case isEndBracket(ch):
		token := Token{Type: ENDBRACKET, Value: EndBracket, options: t.options}
		t.result = append(t.result, token)

		return false, nil
	case isStartBrace(ch):
		token := Token{Type: STARTBRACE, Value: StartBrace, options: t.options}
		t.result = append(t.result, token)

		return false, nil
	case isEndBrace(ch):
		token := Token{Type: ENDBRACE, Value: EndBrace, options: t.options}
		t.result = append(t.result, token)

		return false, nil
	default:
		if err := t.scanIdent(); err != nil {
			return false, err
		}

		return false, nil
	}
}

func (t *Tokenizer) scanWhiteSpace() error {
	t.unread()

	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == ErrEOF {
				break
			} else {
				return err
			}
		}

		if !isWhiteSpace(ch) {
			t.unread()

			break
		} else {
			t.w.WriteRune(ch)
		}
	}

	if strings.Contains(t.w.String(), "\n") {
		tok := Token{Type: NEWLINE, Value: "\n", options: t.options}
		t.result = append(t.result, tok)
	} else {
		tok := Token{Type: WS, Value: t.w.String(), options: t.options}
		t.result = append(t.result, tok)
	}

	t.w.Reset()

	return nil
}

// scan string token including single quotes.
func (t *Tokenizer) scanString() error {
	var counter int
	t.unread()

	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == ErrEOF {
				break
			} else {
				return err
			}
		}

		// ignore the first single quote
		if counter != 0 && isSingleQuote(ch) {
			t.w.WriteRune(ch)

			break
		} else {
			t.w.WriteRune(ch)
		}
		counter++
	}

	tok := Token{Type: STRING, Value: t.w.String(), options: t.options}
	t.result = append(t.result, tok)
	t.w.Reset()

	return nil
}

// append all ch to result until ch is a white space
// if ident is keyword, Type will be the keyword and value will be the uppercase keyword.
func (t *Tokenizer) scanIdent() error {
	t.unread()

LOOP:
	for {
		ch, _, err := t.r.ReadRune()
		if err != nil {
			if err.Error() == ErrEOF {
				break
			} else {
				return err
			}
		}
		switch {
		case isWhiteSpace(ch):
			t.unread()

			break LOOP
		case isComma(ch):
			t.unread()

			break LOOP
		case isStartParenthesis(ch):
			t.unread()

			break LOOP
		case isEndParenthesis(ch):
			t.unread()

			break LOOP
		case isSingleQuote(ch):
			t.unread()

			break LOOP
		case isStartBracket(ch):
			t.unread()

			break LOOP
		case isEndBracket(ch):
			t.unread()

			break LOOP
		case isStartBrace(ch):
			t.unread()

			break LOOP
		case isEndBrace(ch):
			t.unread()

			break LOOP
		default:
			t.w.WriteRune(ch)
		}
	}

	t.append(t.w.String())

	return nil
}

func (t *Tokenizer) append(v string) {
	upperValue := strings.ToUpper(v)

	if ttype, ok := t.isSQLKeyWord(upperValue); ok {
		t.result = append(t.result, Token{
			Type:    ttype,
			Value:   upperValue,
			options: t.options,
		})
	} else {
		t.result = append(t.result, Token{
			Type:    ttype,
			Value:   v,
			options: t.options,
		})
	}

	t.w.Reset()
}

func (t *Tokenizer) isSQLKeyWord(v string) (TokenType, bool) {
	if ttype, ok := sqlKeywordMap[v]; ok {
		return ttype, ok
	} else if ttype, ok := typeWithParenMap[v]; ok {
		if r, _, err := t.r.ReadRune(); err == nil && string(r) == StartParenthesis {
			t.unread()

			return ttype, ok
		}
		t.unread()

		return IDENT, ok
	}

	return IDENT, false
}

var (
	sqlKeywordMap    map[string]TokenType
	typeWithParenMap map[string]TokenType
)

func init() {
	sqlKeywordMap = map[string]TokenType{
		"ALL":         ALL,
		"AND":         AND,
		"AS":          AS,
		"ASC":         ASC,
		"AT":          AT,
		"BETWEEN":     BETWEEN,
		"BY":          BY,
		"CASE":        CASE,
		"COLLATE":     COLLATE,
		"CROSS":       CROSS,
		"DELETE":      DELETE,
		"DESC":        DESC,
		"DISTINCT":    DISTINCT,
		"DISTINCTROW": DISTINCTROW,
		"DO":          DO,
		"ELSE":        ELSE,
		"END":         END,
		"EXCEPT":      EXCEPT,
		"EXISTS":      EXISTS,
		"FETCH":       FETCH,
		"FILTER":      FILTER,
		"FIRST":       FIRST,
		"FOR":         FOR,
		"FROM":        FROM,
		"GROUP":       GROUP,
		"HAVING":      HAVING,
		"IN":          IN,
		"INNER":       INNER,
		"INSERT":      INSERT,
		"INTERSECT":   INTERSECT,
		"INTO":        INTO,
		"IS":          IS,
		"JOIN":        JOIN,
		"LAST":        LAST,
		"LEFT":        LEFT,
		"LIKE":        LIKE,
		"LIMIT":       LIMIT,
		"LOCK":        LOCK,
		"NATURAL":     NATURAL,
		"NOT":         NOT,
		"NULL":        NULL,
		"NULLS":       NULLS,
		"OFFSET":      OFFSET,
		"ON":          ON,
		"OR":          OR,
		"ORDER":       ORDER,
		"OUTER":       OUTER,
		"OVERLAPS":    OVERLAPS,
		"RETURNING":   RETURNING,
		"RIGHT":       RIGHT,
		"ROWS":        ROWS,
		"SELECT":      SELECT,
		"SET":         SET,
		"THEN":        THEN,
		"UNION":       UNION,
		"UPDATE":      UPDATE,
		"USING":       USING,
		"VALUES":      VALUES,
		"WHEN":        WHEN,
		"WHERE":       WHERE,
		"WITH":        WITH,
		"WITHIN":      WITHIN,
		"ZONE":        ZONE,
	}

	typeWithParenMap = map[string]TokenType{
		"ARRAY_AGG":       FUNCTION,
		"AVG":             FUNCTION,
		"BIG":             TYPE,
		"BIGSERIAL":       TYPE,
		"BIT":             TYPE,
		"BOOLEAN":         TYPE,
		"CAST":            FUNCTION,
		"CHAR":            TYPE,
		"COALESCE":        FUNCTION,
		"COUNT":           FUNCTION,
		"CUSTOMTYPE":      TYPE,
		"DATE_PART":       FUNCTION,
		"DATE_TRUNC":      FUNCTION,
		"DEC":             TYPE,
		"DECIMAL":         TYPE,
		"EXTRACT":         FUNCTION,
		"FLOAT":           TYPE,
		"GREATEST":        FUNCTION,
		"INTEGER":         TYPE,
		"INTERVAL":        TYPE,
		"LEAST":           FUNCTION,
		"MAX":             FUNCTION,
		"MIN":             FUNCTION,
		"NUMERIC":         TYPE,
		"OVER":            FUNCTION,
		"OVERLAY":         FUNCTION,
		"PERCENTILE_DISC": FUNCTION,
		"POSITION":        FUNCTION,
		"RANDOM":          FUNCTION,
		"ROW_NUMBER":      FUNCTION,
		"SECOND":          TYPE,
		"SUBSTRING":       FUNCTION,
		"SUM":             FUNCTION,
		"TEXT":            TYPE,
		"TIME":            TYPE,
		"TIMESTAMP":       TYPE,
		"TRIM":            FUNCTION,
		"VARBIT":          TYPE,
		"VARCHAR":         TYPE,
		"XMLCONCAT":       FUNCTION,
		"XMLELEMENT":      FUNCTION,
		"XMLFOREST":       FUNCTION,
	}
}
