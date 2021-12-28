package sqlfmt

import (
	"sync"

	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer"
	"github.com/fredbi/go-sqlfmt/sqlfmt/lexer/postgis"
	"github.com/fredbi/go-sqlfmt/sqlfmt/parser"
	"github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"
)

type (
	Option func(*options)

	// options for go-sqlfmt.
	options struct {
		Distance    int
		IsRawSQL    bool
		Colorized   bool
		LowerCased  bool
		WithPostgis bool
		CommaStyle  group.CommaStyle

		lexerOptions  []lexer.Option
		parserOptions []parser.Option
	}
)

var onceRegisterPostgis sync.Once

func (o *options) ToLexerOptions() []lexer.Option {
	var lexerOptions []lexer.Option
	if o.Colorized {
		lexerOptions = append(lexerOptions, lexer.Colorized())
	}
	if o.LowerCased {
		lexerOptions = append(lexerOptions, lexer.LowerCased())
	}

	if o.WithPostgis {
		onceRegisterPostgis.Do(func() {
			lexer.Register(postgis.Registry{})
		})
	}

	lexerOptions = append(lexerOptions, o.lexerOptions...)

	return lexerOptions
}

func (o *options) ToParserOptions() []parser.Option {
	var parserOptions []parser.Option

	parserOptions = append(parserOptions, parser.WithGroupOptions(group.WithCommaStyle(o.CommaStyle)))

	parserOptions = append(parserOptions, o.parserOptions...)

	return parserOptions
}

func defaultOptions(opts ...Option) *options {
	o := &options{
		Distance:   0,
		CommaStyle: group.CommaStyleLeft,
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

// WithDistance sets the distance between formatted tokens.
func WithDistance(distance int) Option {
	return func(opts *options) {
		opts.Distance = distance
	}
}

// WithRawSQL formats raw SQL files.
func WithRawSQL(enabled bool) Option {
	return func(opts *options) {
		opts.IsRawSQL = enabled
	}
}

// WithColorized formats output with some colors (does not apply to go files).
func WithColorized(enabled bool) Option {
	return func(opts *options) {
		opts.Colorized = enabled
	}
}

// WithLowerCased formats output with SQL keywords lower cased.
func WithLowerCased(enabled bool) Option {
	return func(opts *options) {
		opts.LowerCased = enabled
	}
}

// WithLexerOptions sets some lexer options (e.g. formatting, ...).
func WithLexerOptions(lexerOptions ...lexer.Option) Option {
	return func(opts *options) {
		opts.lexerOptions = lexerOptions
	}
}

// WithParserOptions sets some parser options (e.g. grouping, ...).
func WithParserOptions(parserOptions ...parser.Option) Option {
	return func(opts *options) {
		opts.parserOptions = parserOptions
	}
}

// WithCommaStyles defines the comma justification style (left comma, right comma).
func WithCommaStyle(style group.CommaStyle) Option {
	return func(opts *options) {
		opts.CommaStyle = style
	}
}

// WithPostgis enables postgis dictionary of functions,types and operators.
func WithPostgis(enabled bool) Option {
	return func(opts *options) {
		opts.WithPostgis = enabled
	}
}
