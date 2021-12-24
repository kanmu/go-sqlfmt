package parser

import "github.com/fredbi/go-sqlfmt/sqlfmt/parser/group"

type (
	// Option for the parser.
	Option func(*options)

	options struct {
		groupOptions []group.Option
	}
)

func defaultOptions(opts ...Option) *options {
	o := &options{}

	for _, apply := range opts {
		apply(o)
	}

	return o
}

// WithGroupOptions specifies some grouping options.
func WithGroupOptions(groupOptions ...group.Option) Option {
	return func(opts *options) {
		opts.groupOptions = groupOptions
	}
}

func withOptions(o *options) Option {
	return func(opts *options) {
		*opts = *o
	}
}
