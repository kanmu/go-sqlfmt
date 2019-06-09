package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/scanner"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kanmu/go-sqlfmt/sqlfmt"
)

var (
	exitCode = 0
	// main operation modes
	list    = flag.Bool("l", false, "list files whose formatting differs from goreturns's")
	write   = flag.Bool("w", false, "write result to (source) file instead of stdout")
	doDiff  = flag.Bool("d", false, "display diffs instead of rewriting files")
	options = &sqlfmt.Options{}
)

func init() {
	flag.BoolVar(&options.PrintErrors, "p", false, "print non-fatal typechecking errors to stderr")
	flag.BoolVar(&options.AllErrors, "e", false, "report all errors (not just the first 10 on different lines)")
	flag.BoolVar(&options.RemoveBareReturns, "b", false, "remove bare returns")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: sqlfmt [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func isGoFile(info os.FileInfo) bool {
	name := info.Name()
	return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func visitFile(path string, info os.FileInfo, err error) error {
	if err == nil && isGoFile(info) {
		err = processFile(path, nil, os.Stdout, false)
	}
	if err != nil {
		report(err)
	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func processFile(filename string, in io.Reader, out io.Writer, stdin bool) error {
	opt := options
	if stdin {
		nopt := *options
		nopt.Fragment = true
		opt = &nopt
	}

	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	res, err := sqlfmt.Process(filename, src, opt)
	if err != nil {
		return err
	}

	if !bytes.Equal(src, res) {
		if *list {
			fmt.Fprintln(out, filename)
		}
		if *write {
			if err = ioutil.WriteFile(filename, res, 0); err != nil {
				return err
			}
		}
		if *doDiff {
			data, err := diff(src, res)
			if err != nil {
				return fmt.Errorf("computing diff: %s", err)
			}
			fmt.Printf("diff %s gofmt/%s\n", filename, filename)
			out.Write(data)
		}
		if !*list && !*write && !*doDiff {
			_, err = out.Write(res)
		}
	}

	return nil
}

func sqlfmtMain() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		if err := processFile("", os.Stdin, os.Stdout, true); err != nil {
			report(err)
		}
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			walkDir(path)
		default:
			info, err := os.Stat(path)
			if err != nil {
				report(err)
			}
			if isGoFile(info) {
				err = processFile(path, nil, os.Stdout, false)
				if err != nil {
					report(err)
				}
			}
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	sqlfmtMain()
	os.Exit(exitCode)
}

func diff(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "sqlfmt")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "sqlfmt")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return
}
