package lexer

import (
	"bytes"
	"fmt"
)

// Token is a token struct.
type Token struct {
	Type  TokenType
	Value string
	*options
}

// MakeToken builds an immutable token.
func MakeToken(ttype TokenType, value string, opts ...Option) Token {
	o := defaultOptions(opts...)

	return Token{
		Type:    ttype,
		Value:   value,
		options: o,
	}
}

// Reindent is a placeholder for implementing Reindenter interface.
func (t Token) Reindent(buf *bytes.Buffer) error { return nil }

// GetStart is a placeholder for implementing Reindenter interface.
func (t Token) GetStart() int { return 0 }

// IncrementIndentLevel is a placeholder implementing Reindenter interface.
func (t Token) IncrementIndentLevel(lev int) {}

func (t Token) formatKeyword() string {
	if t.options == nil {
		return t.Value
	}
	in := t.Value

	switch t.Type {
	case STRING, IDENT:
		// no op
	case RESERVEDVALUE:
		switch t.Value {
		case "NAN":
			in = "NaN"
		case "INFINITY":
			in = "Infinity"
		case "-INFINITY":
			in = "-Infinity"
		}
	case FUNCTION:
		recased, ok := casedFunctions.Get([]byte(t.Value))
		if ok {
			in = recased.(string)
		}
	}

	if t.recaser != nil {
		in = t.recaser(in)
	}

	if t.colorizer != nil {
		in = t.colorizer(t.Type, in)
	}

	return in
}

func (t Token) formatPunctuation() string {
	if t.Type == SEMICOLON {
		return fmt.Sprintf("%s%s", NewLine, t.Value)
	}

	return t.Value
}

// FormattedValue returns the token with some formatting options.
func (t Token) FormattedValue() string {
	switch t.Type {
	case EOF,
		WS,
		NEWLINE,
		COMMA,
		SEMICOLON,
		STARTPARENTHESIS,
		ENDPARENTHESIS,
		STARTBRACKET,
		ENDBRACKET,
		STARTBRACE,
		ENDBRACE:
		// ANDGROUP,
		// ORGROUP:
		return t.formatPunctuation()
	default:
		return t.formatKeyword()
	}
}

// end keywords of each clause.
var (
	EndOfSelect      []TokenType
	EndOfCase        []TokenType
	EndOfFrom        []TokenType
	EndOfJoin        []TokenType
	EndOfWhere       []TokenType
	EndOfAndGroup    []TokenType
	EndOfOrGroup     []TokenType
	EndOfGroupBy     []TokenType
	EndOfHaving      []TokenType
	EndOfOrderBy     []TokenType
	EndOfLimitClause []TokenType
	EndOfParenthesis []TokenType
	EndOfTieClause   []TokenType
	EndOfUpdate      []TokenType
	EndOfSet         []TokenType
	EndOfReturning   []TokenType
	EndOfDelete      []TokenType
	EndOfInsert      []TokenType
	EndOfValues      []TokenType
	EndOfFunction    []TokenType
	EndOfTypeCast    []TokenType
	EndOfLock        []TokenType
	EndOfWith        []TokenType
)

func init() {
	EndOfSelect = []TokenType{FROM, UNION, EOF}
	EndOfCase = []TokenType{END}
	EndOfFrom = []TokenType{WHERE, INNER, OUTER, LEFT, RIGHT, JOIN, NATURAL, CROSS, ORDER, GROUP, UNION, OFFSET, LIMIT, FETCH, EXCEPT, INTERSECT, EOF, ENDPARENTHESIS}
	EndOfJoin = []TokenType{WHERE, ORDER, GROUP, LIMIT, OFFSET, FETCH, ANDGROUP, ORGROUP, LEFT, RIGHT, INNER, OUTER, NATURAL, CROSS, UNION, EXCEPT, INTERSECT, EOF, ENDPARENTHESIS}
	EndOfWhere = []TokenType{GROUP, ORDER, LIMIT, OFFSET, FETCH, ANDGROUP, OR, UNION, EXCEPT, INTERSECT, RETURNING, EOF, ENDPARENTHESIS}
	EndOfAndGroup = []TokenType{GROUP, ORDER, LIMIT, OFFSET, FETCH, UNION, EXCEPT, INTERSECT, ANDGROUP, ORGROUP, EOF, ENDPARENTHESIS}
	EndOfOrGroup = []TokenType{GROUP, ORDER, LIMIT, OFFSET, FETCH, UNION, EXCEPT, INTERSECT, ANDGROUP, ORGROUP, EOF, ENDPARENTHESIS}
	EndOfGroupBy = []TokenType{ORDER, LIMIT, FETCH, OFFSET, UNION, EXCEPT, INTERSECT, HAVING, EOF, ENDPARENTHESIS}
	EndOfHaving = []TokenType{LIMIT, OFFSET, FETCH, ORDER, UNION, EXCEPT, INTERSECT, EOF, ENDPARENTHESIS}
	EndOfOrderBy = []TokenType{LIMIT, FETCH, OFFSET, UNION, EXCEPT, INTERSECT, EOF, ENDPARENTHESIS}
	EndOfLimitClause = []TokenType{UNION, EXCEPT, INTERSECT, EOF, ENDPARENTHESIS}
	EndOfParenthesis = []TokenType{ENDPARENTHESIS}
	EndOfTieClause = []TokenType{SELECT}
	EndOfUpdate = []TokenType{WHERE, SET, RETURNING, EOF}
	EndOfSet = []TokenType{WHERE, RETURNING, EOF}
	EndOfReturning = []TokenType{EOF}
	EndOfDelete = []TokenType{WHERE, FROM, EOF}
	EndOfInsert = []TokenType{VALUES, EOF}
	EndOfValues = []TokenType{UPDATE, RETURNING, EOF}
	EndOfFunction = []TokenType{ENDPARENTHESIS}
	EndOfTypeCast = []TokenType{ENDPARENTHESIS}
	EndOfLock = []TokenType{EOF}
	EndOfWith = []TokenType{EOF}
}

// token types that contain the keyword to make subGroup.
var (
	TokenTypesOfGroupMaker []TokenType
	TokenTypesOfJoinMaker  []TokenType
	TokenTypeOfTieClause   []TokenType
	TokenTypeOfLimitClause []TokenType
)

func init() {
	TokenTypesOfGroupMaker = []TokenType{
		SELECT, CASE, FROM, WHERE, ORDER, GROUP, LIMIT,
		ANDGROUP, ORGROUP, HAVING,
		UNION, EXCEPT, INTERSECT,
		FUNCTION,
		STARTPARENTHESIS,
		TYPE,
	}
	TokenTypesOfJoinMaker = []TokenType{
		JOIN, INNER, OUTER, LEFT, RIGHT, NATURAL, CROSS, LATERAL,
	}
	TokenTypeOfTieClause = []TokenType{UNION, INTERSECT, EXCEPT}
	TokenTypeOfLimitClause = []TokenType{LIMIT, FETCH, OFFSET}
}

// IsJoinStart determines if ttype is included in TokenTypesOfJoinMaker.
func (t Token) IsJoinStart() bool {
	for _, v := range TokenTypesOfJoinMaker {
		if t.Type == v {
			return true
		}
	}

	return false
}

// IsTieClauseStart determines if ttype is included in TokenTypesOfTieClause.
func (t Token) IsTieClauseStart() bool {
	for _, v := range TokenTypeOfTieClause {
		if t.Type == v {
			return true
		}
	}

	return false
}

// IsLimitClauseStart determines ttype is included in TokenTypesOfLimitClause.
func (t Token) IsLimitClauseStart() bool {
	for _, v := range TokenTypeOfLimitClause {
		if t.Type == v {
			return true
		}
	}

	return false
}

// IsNeedNewLineBefore returns true if token needs new line before written in buffer.
func (t Token) IsNeedNewLineBefore() bool {
	var ttypes = []TokenType{
		SELECT, UPDATE, INSERT, DELETE,
		ANDGROUP,
		FROM, GROUP, ORGROUP,
		ORDER, HAVING, LIMIT, OFFSET, FETCH, RETURNING,
		SET, UNION, INTERSECT, EXCEPT, VALUES,
		WHERE, ON, USING, UNION, EXCEPT, INTERSECT,
	}
	for _, v := range ttypes {
		if t.Type == v {
			return true
		}
	}

	return false
}

// IsKeyWordInSelect returns true if token is a keyword in select group.
func (t Token) IsKeyWordInSelect() bool {
	return t.Type == SELECT ||
		t.Type == EXISTS ||
		t.Type == DISTINCT ||
		t.Type == DISTINCTROW ||
		t.Type == INTO ||
		t.Type == AS ||
		t.Type == GROUP ||
		t.Type == ORDER ||
		t.Type == BY ||
		t.Type == ON ||
		t.Type == RETURNING ||
		t.Type == SET ||
		t.Type == UPDATE
}
