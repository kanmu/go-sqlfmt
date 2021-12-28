package lexer

import "unicode"

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '　' || unicode.IsSpace(ch)
}

func isComma(ch rune) bool {
	return ch == ','
}

func isSemiColon(ch rune) bool {
	return ch == ';'
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

func isDoubleQuote(ch rune) bool {
	return ch == '"'
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

func isEOF(ch rune) bool {
	return ch == eof
}

func isOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' ||
		ch == '<' || ch == '>' || ch == '=' || ch == '~' ||
		ch == '!' || ch == '@' || ch == '#' || ch == '%' ||
		ch == '^' || ch == '&' || ch == '|' || ch == '`' ||
		ch == '?' || ch == ':'
}

func isSeparator(ch rune) bool {
	return isWhiteSpace(ch) ||
		isComma(ch) ||
		isStartParenthesis(ch) ||
		isEndParenthesis(ch) ||
		isSingleQuote(ch) ||
		isStartBracket(ch) ||
		isEndBracket(ch) ||
		isStartBrace(ch) ||
		isEndBrace(ch)
}
