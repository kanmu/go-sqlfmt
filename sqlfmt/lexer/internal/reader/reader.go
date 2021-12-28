package reader

import (
	"container/list"
	"fmt"
	"io"
	"strings"
)

type (
	// RewindRuneReader is derived from the strings.Reader.
	//
	// It knows how to rewind across the runes it reads.
	RewindRuneReader struct {
		reader *strings.Reader
		stack  *runeStack
		*readerOptions
	}

	runeStack struct {
		stack *list.List
		limit int
	}

	runeWithOffset struct {
		ch   rune
		size int
	}
)

var _ io.RuneScanner = &RewindRuneReader{}

func NewRewindRuneReader(src string, opts ...Option) *RewindRuneReader {
	o := defaultOptions(opts...)

	return &RewindRuneReader{
		reader: strings.NewReader(src),
		stack: &runeStack{
			stack: list.New(),
			limit: o.lookAhead,
		},
		readerOptions: o,
	}
}

func (r *RewindRuneReader) ReadRune() (rune, int, error) {
	ch, size, err := r.reader.ReadRune()
	if err != nil {
		return 0, 0, err
	}

	r.stack.Push(runeWithOffset{
		ch:   ch,
		size: size,
	})

	return ch, size, err
}

func (r *RewindRuneReader) UnreadRune() error {
	ro, err := r.stack.Pop()
	if err != nil {
		return err
	}

	_, err = r.reader.Seek(int64(-1*ro.size), io.SeekCurrent)

	return err
}

func (c *runeStack) Push(value runeWithOffset) {
	c.stack.PushFront(value)

	if c.limit > 0 && c.stack.Len() > c.limit {
		_ = c.stack.Remove(c.stack.Back())
	}
}

func (c *runeStack) Pop() (runeWithOffset, error) {
	if c.stack.Len() > 0 {
		ele := c.stack.Front()
		val := c.stack.Remove(ele)

		return val.(runeWithOffset), nil
	}

	return runeWithOffset{}, fmt.Errorf("Pop Error: Stack is empty")
}
