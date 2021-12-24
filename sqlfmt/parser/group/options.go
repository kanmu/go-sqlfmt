package group

type (
	Option func(*options)

	options struct {
		indentLevel int
	}
)

func defaultOptions(opts ...Option) *options {
	o := &options{}

	for _, apply := range opts {
		apply(o)
	}

	return o
}
func WithIndentLevel(level int) Option {
	return func(opts *options) {
		opts.indentLevel = level
	}
}
