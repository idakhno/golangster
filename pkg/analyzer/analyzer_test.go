package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/idakhno/golangster/pkg/analyzer"
)

func TestAnalyzer(t *testing.T) {
	// resolve testdata path relative to the package directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// pkg/analyzer is two levels below the project root
	testdata := filepath.Join(wd, "..", "..", "testdata")

	analysistest.Run(t, testdata, analyzer.Analyzer,
		"lowercase",
		"english",
		"special_chars",
		"sensitive",
	)
}
