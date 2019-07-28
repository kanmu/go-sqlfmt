package parser

import (
	"github.com/kanmu/go-sqlfmt/sqlfmt/lexer"
)

// SQL ...
type SQL struct {
	Clauses []Clause
}

// Clause ...
type Clause struct {
	Name        string
	IndentLevel int
}

// Value represents all of the constructs in the SQL and their subsequent rules.
type Value struct {
	name          string
	values        []interface{}
	prevTokenType lexer.TokenType
	hasParent     bool
}

type parser struct {
	values []Value
}

// name of value
const (
	CRUD         = "crud"
	BRACKET      = "bracket"
	BRACE        = "brace"
	PARENTHESIS  = "parenthesis"
	FUNCEXPR     = "func_expr"
	TYPECASTEXPR = "type_cast_expr"
	CASEEXPR     = "case_expr"
	IDENT        = "ident"
)

func (p *parser) parseValues(tokens []lexer.Token) []*Value {
	var (
		values []*Value
		idx    int
	)

LOOP:
	for {
		token := tokens[idx]
		if idx == 0 {
			// first value must be SELECT, UPDATE, INSERT, DELETE and LOCK
			values = append(values, &Value{
				name:   CRUD,
				values: []interface{}{token},
			})
		}

		switch token.Type {
		case lexer.EOF:
			break LOOP
		case lexer.STARTBRACKET:
			v, len := p.parseBracket(tokens[idx:], tokens[idx-1])
			values = append(values, v)
			idx += len
		case lexer.STARTBRACE:
			v, len := p.parseBrace(tokens[idx:], tokens[idx-1])
			values = append(values, v)
			idx += len
		case lexer.FUNCTION:
			v, len := p.parseFunction(tokens[idx:], tokens[idx-1])
			values = append(values, v)
			idx += len
		case lexer.TYPE:
			v, len := p.parseTypeCast(tokens[idx:], tokens[idx-1])
			values = append(values, v)
			idx += len
		case lexer.CASE:
			v, len := p.parseBracket(tokens[idx:], tokens[idx-1])
			values = append(values, v)
			idx += len
		case lexer.STARTPARENTHESIS:
			v, len := p.parseBracket(tokens[idx:], tokens[idx-1])
			values = append(values, v)
			idx += len
		default:
			v, len := p.parseIdent(tokens[idx:], tokens[idx-1])
			values = append(values, v)
			idx += len
		}
	}
	return values
}

// ここで、distinct from　や、その他 句のパースで便利なものはここで分けておく
func (p *parser) parseIdent(tokens []lexer.Token, prevToken lexer.Token) (*Value, int) {
	var result []interface{}

	// ひとまず初めのトークンだけをidentとして返す
	// 後々、識別子ごとに処理を分ける
	result = append(result, tokens[0])

	return &Value{
		name:          IDENT,
		values:        result,
		prevTokenType: prevToken.Type,
	}, len(result)
}

func (p *parser) parseBracket(tokens []lexer.Token, prevToken lexer.Token) (*Value, int) {
	var (
		result   []interface{}
		startCnt int
		endCnt   int
	)

	for _, token := range tokens {
		switch {
		case startCnt == endCnt:
			break
		case token.Type == lexer.STARTBRACKET:
			startCnt++
			result = append(result, token)
		case token.Type == lexer.ENDBRACKET:
			endCnt++
			result = append(result, token)
		default:
			result = append(result, token)
		}
	}
	return &Value{
		name:          BRACKET,
		values:        result,
		prevTokenType: prevToken.Type,
	}, len(result)
}

func (p *parser) parseBrace(tokens []lexer.Token, prevToken lexer.Token) (*Value, int) {
	var (
		result   []interface{}
		startCnt int
		endCnt   int
	)
LOOP:
	for _, token := range tokens {
		switch {
		case startCnt == endCnt:
			break LOOP
		case token.Type == lexer.STARTBRACE:
			startCnt++
			result = append(result, token)
		case token.Type == lexer.ENDBRACE:
			endCnt++
			result = append(result, token)
		default:
			result = append(result, token)
		}
	}

	return &Value{
		name:          BRACE,
		values:        result,
		prevTokenType: prevToken.Type,
	}, len(result)
}

// 初めのバリューは構造体に持たせる
// switch のところを共通化して、それぞれのパースファンクションで呼ぶようにする
func (p *parser) parseFunction(tokens []lexer.Token, prevToken lexer.Token) (*Value, int) {
	var (
		result      []interface{}
		idx         int
		hasParent   bool
		endFunction = lexer.ENDPARENTHESIS
	)

LOOP:
	for {
		token := tokens[idx]

		// in case of FUNCTION expression, the first value must be a name of functions and the second value must be start of parenthesis
		// That's why values when idx == 0 and 1 are automatically appended to the result
		if idx == 0 || idx == 1 {
			result = append(result, token)
			idx++
			continue
		}

		switch token.Type {
		case lexer.EOF:
			// error処理
		case endFunction:
			result = append(result, token)
			break LOOP
		case lexer.FUNCTION, lexer.TYPE, lexer.STARTPARENTHESIS, lexer.CASE:
			v, len := p.parseNestedValue(token, tokens[idx:], tokens[idx-1])
			result = append(result, v)
			idx += len
		default:
			result = append(result, token)
			idx++
		}
	}
	return &Value{
		name:          FUNCEXPR,
		values:        result,
		prevTokenType: prevToken.Type,
		hasParent:     hasParent,
		// (idx + 1) is the length of the result which will be used as the index of next token
		// I chose not to use len(result) because result may have another nested value , which makes counting tokens more complicated
	}, (idx + 1)
}

func (p *parser) parseCase(tokens []lexer.Token, prevToken lexer.Token) (*Value, int) {
	var (
		result    []interface{}
		idx       int
		hasParent bool
		endCase   = lexer.END
	)
LOOP:
	for {
		token := tokens[idx]

		if idx == 0 {
			result = append(result, token)
			idx++
			continue
		}

		switch token.Type {
		case lexer.EOF:
			// error処理
		case endCase:
			result = append(result, token)
			break LOOP
		case lexer.FUNCTION, lexer.TYPE, lexer.STARTPARENTHESIS, lexer.CASE:
			v, len := p.parseNestedValue(token, tokens[idx:], tokens[idx-1])
			result = append(result, v)
			idx += len
		default:
			result = append(result, token)
			idx++
		}
	}
	return &Value{
		name:          CASEEXPR,
		values:        result,
		prevTokenType: prevToken.Type,
		hasParent:     hasParent,
		// (idx + 1) is the length of the result which will be used as the index of next token
		// I chose not to use len(result) because result may have another nested value , which makes counting tokens more complicated
	}, (idx + 1)
}

func (p *parser) parseTypeCast(tokens []lexer.Token, prevToken lexer.Token) (*Value, int) {
	var (
		result      []interface{}
		idx         int
		hasParent   bool
		endTypeCast = lexer.ENDPARENTHESIS
	)
LOOP:
	for {
		token := tokens[idx]

		// in case of type cast expression, the first value must be a name of type and the second value must be start of parenthesis
		// That's why values when idx == 0 and 1 are automatically appended to the result
		if idx == 0 || idx == 1 {
			result = append(result, token)
			idx++
			continue
		}

		switch token.Type {
		case lexer.EOF:
			// error処理
		case endTypeCast:
			result = append(result, token)
			break LOOP
		case lexer.FUNCTION, lexer.TYPE, lexer.STARTPARENTHESIS, lexer.CASE:
			v, len := p.parseNestedValue(token, tokens[idx:], tokens[idx-1])
			result = append(result, v)
			idx += len
		default:
			result = append(result, token)
			idx++
		}
	}
	return &Value{
		name:          TYPECASTEXPR,
		values:        result,
		prevTokenType: prevToken.Type,
		hasParent:     hasParent,
		// (idx + 1) is the length of the result which will be used as the index of next token
		// I chose not to use len(result) because result may have another nested value , which makes counting tokens more complicated
	}, (idx + 1)
}

func (p *parser) parseParenthesis(tokens []lexer.Token, prevToken lexer.Token) (*Value, int) {
	var (
		result         []interface{}
		idx            int
		hasParent      bool
		endParenthesis = lexer.ENDPARENTHESIS
	)
LOOP:
	for {
		token := tokens[idx]

		if idx == 0 {
			result = append(result, token)
			idx++
			continue
		}

		switch token.Type {
		case lexer.EOF:
			// error処理
		case endParenthesis:
			result = append(result, token)
			break LOOP
		case lexer.FUNCTION, lexer.TYPE, lexer.STARTPARENTHESIS, lexer.CASE:
			v, len := p.parseNestedValue(token, tokens[idx:], tokens[idx-1])
			result = append(result, v)
			idx += len
		default:
			result = append(result, token)
			idx++
		}
	}
	return &Value{
		name:          PARENTHESIS,
		values:        result,
		prevTokenType: prevToken.Type,
		hasParent:     hasParent,
		// (idx + 1) is the length of the result which will be used as the index of next token
		// I chose not to use len(result) because result may have another nested value , which makes counting tokens more complicated
	}, (idx + 1)
}

// スイッチケースの共通化を行うため
// 他のケーすをどう扱うかは今からっか
func (p *parser) parseNestedValue(token lexer.Token, tokens []lexer.Token, prevToken lexer.Token) (*Value, int) {
	var (
		v   *Value
		len int
	)

	switch token.Type {
	case lexer.FUNCTION:
		v, len = p.parseFunction(tokens, prevToken)
	case lexer.TYPE:
		v, len = p.parseTypeCast(tokens, prevToken)
	case lexer.STARTPARENTHESIS:
		v, len = p.parseParenthesis(tokens, prevToken)
	case lexer.CASE:
		v, len = p.parseCase(tokens, prevToken)
	}
	return v, len
}
