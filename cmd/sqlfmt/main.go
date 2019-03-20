package main

import (
	"flag"
	"go/printer"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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

	cfg := printer.Config{Mode: printerMode, Tabwidth: tabWidth}
	if *outputFile == "" {
		cfg.Fprint(os.Stdout, sfmt.Fset, sfmt.AstNode)
	} else {
		if err = writeFile(*outputFile, cfg, sfmt); err != nil {
			log.Fatal(err)
		}
	}
}

// atomic write
func writeFile(filename string, cfg printer.Config, sfmt *sqlfmt.SQLFormatter) error {
	tmpFile, err := filepath.Abs(filename + ".")
	if err != nil {
		return err
	}
	f, err := ioutil.TempFile(filepath.Dir(tmpFile), filepath.Base(tmpFile))
	if err != nil {
		return err
	}
	if err := cfg.Fprint(f, sfmt.Fset, sfmt.AstNode); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err = os.Rename(f.Name(), filename); err != nil {
		return err
	}
	return nil
}
