package group

type (
	CommaStyle uint8

	Option func(*options)

	options struct {
		IndentLevel int
		commaStyle  CommaStyle
	}
)

const (
	CommaStyleLeft CommaStyle = iota
	CommaStyleRight
)

func defaultOptions(opts ...Option) *options {
	o := &options{
		commaStyle: CommaStyleLeft,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

func WithIndentLevel(level int) Option {
	return func(opts *options) {
		opts.IndentLevel = level
	}
}

func WithCommaStyle(style CommaStyle) Option {
	return func(opts *options) {
		opts.commaStyle = style
	}
}
