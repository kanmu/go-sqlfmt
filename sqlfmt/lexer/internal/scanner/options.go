package scanner

import "github.com/fredbi/go-sqlfmt/sqlfmt/lexer/internal/reader"

type (
	Option func(*scannerOptions)

	scannerOptions struct {
		readerOptions []reader.Option
	}
)

func defaultOptions(opts ...Option) *scannerOptions {
	o := &scannerOptions{
		readerOptions: []reader.Option{reader.WithLookAhead(4)},
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

// WithReaderOptions sets options for the underlying reader used by this scanner.
func WithReaderOptions(opts ...reader.Option) Option {
	return func(o *scannerOptions) {
		o.readerOptions = opts
	}
}
