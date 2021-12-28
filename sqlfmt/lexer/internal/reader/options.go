package reader

type (
	Option func(*readerOptions)

	readerOptions struct {
		lookAhead int
	}
)

func defaultOptions(opts ...Option) *readerOptions {
	o := &readerOptions{
		lookAhead: 3,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

// WithLookAhead sets the max number of runes that can be unread.
func WithLookAhead(lookAhead int) Option {
	return func(o *readerOptions) {
		o.lookAhead = lookAhead
	}
}
