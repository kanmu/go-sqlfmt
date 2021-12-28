package lexer

import (
	"strings"

	"github.com/fatih/color"
)

type (
	// TokenFormatter knows how to format a token.
	TokenFormatter func(string) string

	// TokenTypeFormatter knows how to format a token by its type.
	TokenTypeFormatter func(TokenType, string) string

	// Option for token formatting.
	Option func(*options)

	options struct {
		colorizer TokenTypeFormatter
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

// WithColorizer sets an arbitrary token formatter as colorizer.
func WithColorizer(colorizer TokenTypeFormatter) Option {
	return func(opt *options) {
		opt.colorizer = colorizer
	}
}

// WithRecaser sets an arbitrary token formatter as recaser.
func WithRecaser(recaser TokenFormatter) Option {
	return func(opt *options) {
		opt.recaser = recaser
	}
}

// Colorized is the default colorizer, with SQL keywords in yellow.
func Colorized() Option {
	return func(opt *options) {
		opt.colorizer = func(tokenType TokenType, in string) string {
			switch tokenType {
			case IDENT, // field or table name
				STRING: // values surrounded with single quotes
				return color.HiGreenString(in)
			case FUNCTION:
				return color.CyanString(in)
			case TYPE:
				return color.MagentaString(in)
			case RESERVEDVALUE:
				return color.HiBlueString(in)
			default:
				return color.YellowString(in)
			}
		}
	}
}

// LowerCased normalizes all SQL keywords as lower case.
func LowerCased() Option {
	return func(opt *options) {
		opt.recaser = strings.ToLower
	}
}

// WithOptionsFrom replicates some options from another token.
func WithOptionsFrom(like Token) Option {
	return func(opt *options) {
		var o *options
		if like.options == nil {
			o = defaultOptions()
		} else {
			o = like.options
		}

		*opt = *o
	}
}
