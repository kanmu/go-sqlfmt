package sqlfmt

import "github.com/fredbi/go-sqlfmt/sqlfmt/lexer"

type (
	Option func(*options)

	// options for go-sqlfmt
	options struct {
		Distance     int
		IsRawSQL     bool
		Colorized    bool
		LowerCased   bool
		lexerOptions []lexer.Option
	}
)

func (o *options) ToLexerOptions() []lexer.Option {
	lexerOptions := o.lexerOptions
	if o.Colorized {
		lexerOptions = append(lexerOptions, lexer.Colorized())
	}
	if o.LowerCased {
		lexerOptions = append(lexerOptions, lexer.LowerCased())
	}

	return lexerOptions
}

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

// WithLexerOptions sets some lexer options (e.g. formatting, ...)
func WithLexerOptions(lexerOptions ...lexer.Option) Option {
	return func(opts *options) {
		lexerOptions = lexerOptions
	}
}
