package main

import (
	"bytes"
	stderrors "errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fredbi/go-sqlfmt/sqlfmt"
	"github.com/pkg/errors"
)

func isGoFile(info os.FileInfo) bool {
	name := info.Name()

	return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func visitFile(opts *cliOptions) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() && (opts.IsRawSQL || isGoFile(info)) {
				if isGoFile(info) {
					opts.IsRawSQL = false
				}

				err = processFile(path, nil, os.Stdout, opts)
			}
		}

		if err != nil {
			processError(errors.Wrap(err, fmt.Sprintf("visit file failed: %s", path)))
		}

		return err
	}
}

func walkDir(path string, opts *cliOptions) error {
	return filepath.Walk(path, visitFile(opts))
}

func processFile(filename string, in io.Reader, out io.Writer, opts *cliOptions) error {
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

	res, err := sqlfmt.Process(filename, src, opts.ToOptions()...)
	if err != nil {
		return errors.Wrap(err, "sqlfmt.Process failed")
	}

	if bytes.Equal(src, res) {
		return nil
	}

	if options.List {
		fmt.Fprintln(out, filename)
	}

	if options.Write {
		if err = ioutil.WriteFile(filename, res, 0); err != nil {
			return errors.Wrap(err, "ioutil.WriteFile failed")
		}
	}

	if options.Diff {
		data, erd := diff(src, res)
		if erd != nil {
			return errors.Wrap(erd, "diff failed")
		}

		fmt.Fprintf(os.Stdout, "diff %s gofmt/%s\n", filename, filename)

		_, err = out.Write(data)
		if err != nil {
			return errors.Wrap(err, "out.Write failed")
		}
	}

	if !options.List && !options.Write && !options.Diff {
		_, err = out.Write(res)
		if err != nil {
			return errors.Wrap(err, "out.Write failed")
		}
	}

	return nil
}

func sqlfmtMain() {
	// the user is piping their source into go-sqlfmt
	if flag.NArg() == 0 {
		if options.Write {
			log.Fatal("can not use -w while using pipeline")
		}

		if err := processFile("<standard input>", os.Stdin, os.Stdout, options); err != nil {
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
			erw := walkDir(path, options)
			if erw != nil {
				processError(erw)
			}
		default:
			info, ers := os.Stat(path)
			if ers != nil {
				processError(ers)
			}

			if options.IsRawSQL || isGoFile(info) {
				err = processFile(path, nil, os.Stdout, options)
				if err != nil {
					processError(err)
				}
			}
		}
	}
}

func diff(b1, b2 []byte) ([]byte, error) {
	f1, err := ioutil.TempFile("", "sqlfmt")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Remove(f1.Name())
	}()
	defer func() {
		_ = f1.Close()
	}()

	f2, err := ioutil.TempFile("", "sqlfmt")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Remove(f2.Name())
	}()
	defer func() {
		_ = f2.Close()
	}()

	_, err = f1.Write(b1)
	if err != nil {
		processError(err)

		return nil, err
	}
	_, err = f2.Write(b2)
	if err != nil {
		processError(err)

		return nil, err
	}

	// TODO: use lib instead of shell
	data, err := exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput() //#nosec
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	if err != nil {
		processError(err)

		return nil, err
	}

	return data, nil
}

func processError(err error) {
	if stderrors.Is(err, &sqlfmt.FormatError{}) {
		log.Println(err)
	} else {
		log.Fatal(err)
	}
}
