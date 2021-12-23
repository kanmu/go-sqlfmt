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

func WithColorizer(colorizer TokenFormatter) Option {
	return func(opt *options) {
		opt.colorizer = colorizer
	}
}

func WithRecaser(recaser TokenFormatter) Option {
	return func(opt *options) {
		opt.recaser = recaser
	}
}

func Colorized() Option {
	return func(opt *options) {
		opt.colorizer = func(in string) string { return color.YellowString(in) }
	}
}

func LowerCased() Option {
	return func(opt *options) {
		opt.recaser = strings.ToLower
	}
}
