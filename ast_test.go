package sqlfmt

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var testingGoFile = `
package main 
import "fmt"
func main(){
	fmt.Println("select * from sometable")
}`

func TestNewSQLFormatter(t *testing.T) {
	test := struct {
		src io.Reader
	}{
		src: strings.NewReader(testingGoFile),
	}
	t.Run("test that ast file is made correctly", func(t *testing.T) {
		_, err := NewSQLFormatter(test.src)
		if err != nil {
			t.Fatalf("ERROR %#v", err)
		}
	})
}

var testFiles, _ = filepath.Glob("./testdata/*.go")

func TestFormatInAst(t *testing.T) {
	for _, file := range testFiles {
		t.Run(file, func(t *testing.T) {
			f, err := os.Open(file)
			if err != nil {
				t.Fatalf("ERROR: %#v", err)
			}

			sfmt, err := NewSQLFormatter(f)
			if err != nil {
				t.Fatalf("ERROR %#v", err)
			}
			if err := sfmt.Format(); err != nil {
				t.Errorf("ERROR:%#v", err)
			}
		})
	}
}
