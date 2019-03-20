package main

import (
	"flag"
	"go/printer"
	"log"
	"os"

	"github.com/kanmu/go-sqlfmt"
)

var (
	srcFile    = flag.String("s", "", "the source file")
	outputFile = flag.String("o", "", "the output file")
)

const (
	tabWidth    = 8
	printerMode = printer.UseSpaces | printer.TabIndent
)

func main() {
	flag.Parse()

	if *srcFile == "" {
		log.Fatal("-s is required")
	}

	f, err := os.Open(*srcFile)
	if err != nil {
		log.Fatal(err)
	}

	sfmt, err := sqlfmt.NewSQLFormatter(f)
	if err != nil {
		log.Fatal(err)
	}

	if err := sfmt.Format(); err != nil {
		log.Println(err)
	}

	if *outputFile == "" {
		cfg := printer.Config{Mode: printerMode, Tabwidth: tabWidth}
		cfg.Fprint(os.Stdout, sfmt.Fset, sfmt.AstNode)
	} else {
		dst, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal(err)
		}
		cfg := printer.Config{Mode: printerMode, Tabwidth: tabWidth}
		cfg.Fprint(dst, sfmt.Fset, sfmt.AstNode)
	}
}
