package lexer

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer/internal/reader"
	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer/internal/scanner"
	"github.com/pkg/errors"
)

type runeWriter interface {
	WriteRune(rune) (int, error)
	String() string
	Reset()
}

// Tokenizer tokenizes SQL statements.
type Tokenizer struct {
	*scanner.RuneScanner
	w      runeWriter // w  writes token value. It resets its value when the end of token appears
	result []Token
	*options
}

// NewTokenizer creates Tokenizer.
func NewTokenizer(src string, opts ...Option) *Tokenizer {
	return &Tokenizer{
		RuneScanner: scanner.NewRuneScanner(src,
			scanner.WithReaderOptions(reader.WithLookAhead(maxOperatorLength+1)),
		),
		w:       &bytes.Buffer{},
		options: defaultOptions(opts...),
	}
}

// GetTokens returns tokens for parsing.
func (t *Tokenizer) GetTokens() ([]Token, error) {
	tokens, err := t.tokenize()
	if err != nil {
		return nil, errors.Wrap(err, "Tokenize failed")
	}

	result := make([]Token, 0, len(tokens))

	// replace all tokens without whitespaces and new lines
	// if "AND" or "OR" appears after new line, token value will be ANDGROUP, ORGROUP
	for i, tok := range tokens {
		switch {
		case tok.Type == WS || tok.Type == NEWLINE:
			continue

		case tok.Type == AND && tokens[i-1].Type == NEWLINE:
			tok = Token{Type: ANDGROUP, Value: tok.Value, options: t.options}

		case tok.Type == OR && tokens[i-1].Type == NEWLINE:
			tok = Token{Type: ORGROUP, Value: tok.Value, options: t.options}

		case tok.Type == LEFT && i < len(tokens)-1 && tokens[i+1].Type == STARTPARENTHESIS:
			// LEFT depends on context: may be keyword or function
			tok = Token{Type: FUNCTION, Value: tok.Value, options: t.options}
		}

		result = append(result, tok)
	}

	return result, nil
}

// Tokenize analyses every rune in SQL statement
// every token is identified when whitespace appears.
func (t *Tokenizer) tokenize() ([]Token, error) {
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

// scan scans each character and appends to result until "eof" appears.
//
// When it has finished scanning all characters, it returns true.
func (t *Tokenizer) scan() (bool, error) {
	ch, err := t.Read()
	if err != nil {
		return false, errors.Wrap(err, "read rune failed")
	}

	switch {
	case isEOF(ch):
		tok := Token{Type: EOF, Value: "EOF", options: t.options}
		t.result = append(t.result, tok)

		return true, nil
	case isWhiteSpace(ch):
		// skip white space
		if err := t.scanWhiteSpace(ch); err != nil {
			return false, err
		}

		return false, nil

	case isSingleQuote(ch):
		// extract quoted string
		// TODO: use start rune
		if err := t.scanString(ch); err != nil {
			return false, err
		}

		return false, nil

	//case isDoubleQuote(ch): // TODO
	case isComma(ch):
		token := Token{Type: COMMA, Value: Comma, options: t.options}
		t.result = append(t.result, token)

		return false, nil
	case isSemiColon(ch):
		token := Token{Type: SEMICOLON, Value: SemiColon, options: t.options}
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
		if isOperator(ch) {
			// operators may not be separated by whitespace
			// extract operator starting with this rune
			ok, err := t.scanOperator(ch)
			if err != nil {
				return false, err
			}

			if ok {
				return false, nil
			}
		}

		if err := t.scanIdent(ch); err != nil {
			return false, err
		}

		return false, nil
	}
}

func (t *Tokenizer) scanWhiteSpace(start rune) error {
	var err error
	ch := start

LOOP:
	for {
		switch {
		case isEOF(ch):
			break LOOP
		case !isWhiteSpace(ch):
			t.Unread()

			break LOOP
		default:
			_, _ = t.w.WriteRune(ch)
		}

		ch, err = t.Read()
		if err != nil {
			return err
		}
	}

	if strings.Contains(t.w.String(), "\n") {
		// TODO: when does this occur?
		tok := Token{Type: NEWLINE, Value: "\n", options: t.options}
		t.result = append(t.result, tok)
	} else {
		// TODO: why do we add whitespace token?
		tok := Token{Type: WS, Value: t.w.String(), options: t.options}
		t.result = append(t.result, tok)
	}

	t.w.Reset()

	return nil
}

// scan string token including single quotes.
func (t *Tokenizer) scanString(start rune) error {
	var err error
	ch := start

	for {
		_, _ = t.w.WriteRune(ch)

		ch, err = t.Read()
		if err != nil {
			return err
		}

		if isSingleQuote(ch) {
			_, _ = t.w.WriteRune(ch)

			break
		}

		if isEOF(ch) {
			break
		}
	}

	tok := Token{Type: STRING, Value: t.w.String(), options: t.options}
	t.result = append(t.result, tok)
	t.w.Reset()

	return nil
}

// isOperatorToken returns an operator token if it finds one, starting from some valid operator rune.
//
// The index returned indicates the number of extra consumed runes from the reader:
// this allows the caller to rewind.
func (t *Tokenizer) isOperatorToken(start rune) (Token, bool, int, error) {
	var (
		counter int
		token   string
		err     error
	)

	w := bytes.NewBuffer(make([]byte, 0, maxOperatorBytes+1))
	ch := start

	for {
		_, _ = w.WriteRune(ch)

		if _, ok := operatorsIndex.Root().Get(w.Bytes()); ok {
			// There is a legit operator corresponding to that sequence.
			// Keep it, and find out if we have a longer match.
			token = w.String()
		}

		if !existsOperatorWithPrefix(w.Bytes()) {
			break
		}

		ch, err = t.Read()
		if err != nil {
			return Token{}, false, counter, err
		}

		if isEOF(ch) {
			break // do not increment count whenever EOF is reached
		}

		counter++
	}

	if token != "" {
		tok := Token{Type: OPERATOR, Value: token, options: t.options}

		return tok, true, counter, nil
	}

	return Token{}, false, counter, nil
}

// scanOperator extracts an operator
func (t *Tokenizer) scanOperator(ch rune) (bool, error) {
	token, ok, counter, err := t.isOperatorToken(ch)
	if ok {
		t.result = append(t.result, token)
		counter = counter - utf8.RuneCountInString(token.Value) + 1
	}

	// rewind to the next rune after the search
	t.Rewind(counter)

	return ok, err
}

func existsOperatorWithPrefix(key []byte) bool {
	iterator := operatorsIndex.Root().Iterator()
	iterator.SeekPrefix(key)
	_, _, ok := iterator.Next()

	return ok
}

// scanIdent appends all runes to result until a separator or a recognized operator is found.
func (t *Tokenizer) scanIdent(start rune) error {
	var (
		counter int
		err     error
	)

	ch := start
LOOP:
	for {
		switch {
		case isEOF(ch):
			if ident := t.w.String(); len(ident) > 0 {
				t.append(t.w.String())
			}

			return nil
		// TODO: double quoted identifiers too, $$ quoting as well
		case isSeparator(ch):
			break LOOP
		default:
			if isOperator(ch) {
				_, isOperator, consumed, err := t.isOperatorToken(ch)
				if err != nil {
					return err
				}

				// rewind looked-ahead runes
				t.Rewind(consumed)
				if isOperator {
					break LOOP
				}
			}

			counter = 0
			_, _ = t.w.WriteRune(ch)

			ch, err = t.Read()
			if err != nil {
				return err
			}

			counter++
		}
	}

	if ident := t.w.String(); len(ident) > 0 {
		t.append(t.w.String())
	}

	t.Rewind(counter)

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
	}

	ttype, ok := typeWithParenMap[v]
	if ok {
		if ttype == TYPE ||
			ttype == OPERATOR ||
			ttype == RESERVEDVALUE {
			return ttype, ok
		}

		if ttype == FUNCTION {
			t.Unread()
			if r, err := t.Read(); err == nil {
				// TODO: some functions may be called without parenthesis --> consider as RESERVED_VALUES ??
				if isStartParenthesis(r) {
					return ttype, true
				}
			}
		}
	}

	return IDENT, false
}
