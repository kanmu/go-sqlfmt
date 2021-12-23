package sqlfmt

type (
	Option func(*options)

	// options for go-sqlfmt
	options struct {
		Distance   int
		IsRawSQL   bool
		Colorized  bool
		LowerCased bool
	}
)

func defaultOptions(opts ...Option) *options {
	o := &options{
		Distance: 0,
	}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

func withOptions(o *options) Option {
	return func(opts *options) {
		*opts = *o
	}
}

// WithDistance sets the distance between formatted tokens
func WithDistance(distance int) Option {
	return func(opts *options) {
		opts.Distance = distance
	}
}

// WithRawSQL formats raw SQL files
func WithRawSQL(enabled bool) Option {
	return func(opts *options) {
		opts.IsRawSQL = enabled
	}
}

// WithColorized formats output with some colors (does not apply to go files)
func WithColorized(enabled bool) Option {
	return func(opts *options) {
		opts.Colorized = enabled
	}
}

// WithLowerCased formats output with SQL keywords lower cased
func WithLowerCased(enabled bool) Option {
	return func(opts *options) {
		opts.LowerCased = enabled
	}
}
