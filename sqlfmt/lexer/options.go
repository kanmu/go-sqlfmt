package lexer

import (
	"strings"

	"github.com/fatih/color"
)

type (
	// TokenFormatter knows how to format a token
	TokenFormatter func(string) string

	// Option for token formatting
	Option func(*options)

	options struct {
		colorizer TokenFormatter
		recaser   TokenFormatter
	}
)

func defaultOptions(opts ...Option) *options {
	o := &options{}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

// WithColorizer sets an arbitrary token formatter as colorizer
func WithColorizer(colorizer TokenFormatter) Option {
	return func(opt *options) {
		opt.colorizer = colorizer
	}
}

// WithRecaser sets an arbitrary token formatter as recaser
func WithRecaser(recaser TokenFormatter) Option {
	return func(opt *options) {
		opt.recaser = recaser
	}
}

// Colorized is the default colorizer, with SQL keywords in yellow
func Colorized() Option {
	return func(opt *options) {
		opt.colorizer = func(in string) string { return color.YellowString(in) }
	}
}

// LowerCased normalizes all SQL keywords as lower case
func LowerCased() Option {
	return func(opt *options) {
		opt.recaser = strings.ToLower
	}
}
