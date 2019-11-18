package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/kanmu/go-sqlfmt/sqlfmt"
)

var (
	// main operation modes
	list    = flag.Bool("l", false, "list files whose formatting differs from goreturns's")
	write   = flag.Bool("w", false, "write result to (source) file instead of stdout")
	doDiff  = flag.Bool("d", false, "display diffs instead of rewriting files")
	options = &sqlfmt.Options{}
)

func init() {
	flag.IntVar(&options.Distance, "distance", 0, "write the distane from the edge to the begin of SQL statements")
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: sqlfmt [flags] [path ...]\n")
	flag.PrintDefaults()
}

func isGoFile(info os.FileInfo) bool {
	name := info.Name()
	return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func visitFile(path string, info os.FileInfo, err error) error {
	if err == nil && isGoFile(info) {
		err = processFile(path, nil, os.Stdout)
	}
	if err != nil {
		processError(errors.Wrap(err, "visit file failed"))

	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func processFile(filename string, in io.Reader, out io.Writer) error {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return errors.Wrap(err, "os.Open failed")
		}
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll failed")
	}

	res, err := sqlfmt.Process(filename, src, options)
	if err != nil {
		return errors.Wrap(err, "sqlfmt.Process failed")
	}

	if !bytes.Equal(src, res) {
		if *list {
			fmt.Fprintln(out, filename)
		}
		if *write {
			if err = ioutil.WriteFile(filename, res, 0); err != nil {
				return errors.Wrap(err, "ioutil.WriteFile failed")
			}
		}
		if *doDiff {
			data, err := diff(src, res)
			if err != nil {
				return errors.Wrap(err, "diff failed")
			}
			fmt.Printf("diff %s gofmt/%s\n", filename, filename)
			out.Write(data)
		}
		if !*list && !*write && !*doDiff {
			_, err = out.Write(res)
			if err != nil {
				return errors.Wrap(err, "out.Write failed")
			}
		}
	}
	return nil
}

func sqlfmtMain() {
	flag.Usage = usage
	flag.Parse()

	// the user is piping their source into go-sqlfmt
	if flag.NArg() == 0 {
		if *write {
			log.Fatal("can not use -w while using pipeline")
		}
		if err := processFile("<standard input>", os.Stdin, os.Stdout); err != nil {
			processError(errors.Wrap(err, "processFile failed"))
		}
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			processError(err)
		case dir.IsDir():
			walkDir(path)
		default:
			info, err := os.Stat(path)
			if err != nil {
				processError(err)
			}
			if isGoFile(info) {
				err = processFile(path, nil, os.Stdout)
				if err != nil {
					processError(err)
				}
			}
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	sqlfmtMain()
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

func processError(err error) {
	switch err.(type) {
	case *sqlfmt.FormatError:
		log.Println(err)
	default:
		log.Fatal(err)
	}
}
