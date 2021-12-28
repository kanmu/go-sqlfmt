package scanner

import (
	"errors"
	"io"
	"unicode/utf8"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer/internal/reader"
)

// EOFRune symbolizes the end of input.
const EOFRune = '∂'

// ErrInvalidRune indicates that an invalid UTF-8 sequence has been read from the input.
var ErrInvalidRune = errors.New("invalid rune")

// RuneScanner knows how to read and rewind across runes.
type RuneScanner struct {
	r io.RuneScanner // r is a custom reader which may rewind runes
}

func NewRuneScanner(src string, opts ...Option) *RuneScanner {
	o := defaultOptions(opts...)

	return &RuneScanner{
		r: reader.NewRewindRuneReader(src, o.readerOptions...),
	}
}

// Unread undoes t.r.UnreadRune method to get the previous character.
//
// It panics if the caller attempts to unread beyond the permitted look-ahead.
func (s *RuneScanner) Unread() {
	if err := s.r.UnreadRune(); err != nil {
		// dev error
		panic(err)
	}
}

func (s *RuneScanner) Rewind(back int) {
	for i := 0; i < back; i++ {
		s.Unread()
	}
}

// Read a rune.
//
// Returns the special rune EOFRune whenever the input is scanned, rather than error.
//
// Error cases are therefore only hit when runes are invalid UTF-8.
func (s *RuneScanner) Read() (rune, error) {
	ch, _, err := s.r.ReadRune()
	if err != nil && errors.Is(err, io.EOF) {
		return EOFRune, nil
	}

	if ch == utf8.RuneError {
		return ch, ErrInvalidRune
	}

	return ch, err
}
