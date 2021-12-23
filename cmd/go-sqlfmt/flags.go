package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fredbi/go-sqlfmt/sqlfmt"
)

type cliOptions struct {
	List  bool
	Write bool
	Diff  bool

	Distance   int
	IsRawSQL   bool
	Colorized  bool
	LowerCased bool
}

func (o cliOptions) ToOptions() []sqlfmt.Option {
	opts := make([]sqlfmt.Option, 0, 4)

	opts = append(opts, sqlfmt.WithDistance(o.Distance))
	opts = append(opts, sqlfmt.WithRawSQL(o.IsRawSQL))
	opts = append(opts, sqlfmt.WithColorized(o.Colorized))
	opts = append(opts, sqlfmt.WithLowerCased(o.LowerCased))

	return opts
}

func defaultCliOptions() *cliOptions {
	return &cliOptions{}
}

var (
	options = defaultCliOptions()
)

func init() {
	// main operation modes
	flag.BoolVar(&options.List, "l", defaultCliOptions().List, "list files whose formatting differs from goreturns's")
	flag.BoolVar(&options.Write, "w", defaultCliOptions().Write, "write result to (source) file instead of stdout")
	flag.BoolVar(&options.Diff, "d", defaultCliOptions().Diff, "display diffs instead of rewriting files")

	// formatting options
	flag.IntVar(&options.Distance, "distance", defaultCliOptions().Distance, "write the distance from the edge to the begin of SQL statements")
	flag.BoolVar(&options.IsRawSQL, "raw", defaultCliOptions().IsRawSQL, "parse raw SQL file")
	flag.BoolVar(&options.Colorized, "colorized", defaultCliOptions().Colorized, "colorize output")
	flag.BoolVar(&options.LowerCased, "lower", defaultCliOptions().LowerCased, "SQL keywords are lower-cased")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: sqlfmt [flags] [path ...]\n")
	flag.PrintDefaults()
}
