package lexer

import "bytes"

// runeWriter is used to write sequences of runes as tokens.
type runeWriter interface {
	WriteRune(rune) (int, error)
	String() string
	Reset()
}

func newRuneWriter() runeWriter {
	return &bytes.Buffer{}
}
