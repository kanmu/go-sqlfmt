package sqlfmt

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

var testFiles, _ = filepath.Glob("./testdata/*.go")

func TestReplaceAstWithFormattedStmt(t *testing.T) {
	for _, file := range testFiles {
		t.Run(file, func(t *testing.T) {
			parserMode := parser.ParseComments
			f, err := parser.ParseFile(token.NewFileSet(), file, nil, parserMode)
			if err != nil {
				t.Fatalf("parser.ParseFile failed: %#v", err)
			}
			if err := replaceAstWithFormattedStmt(f); err != nil {
				t.Errorf("ERROR:%#v", err)
			}
		})
	}
}
