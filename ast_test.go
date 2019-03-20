package sqlfmt

import (
	"io"
	"os"
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

func TestFormatInAst(t *testing.T) {
	f, err := os.Open("./testdata/testing_gofile.go")
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
}
