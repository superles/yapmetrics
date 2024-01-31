package exitcheck_test

import (
	"testing"

	"github.com/superles/yapmetrics/pkg/exitcheck"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzerNegative(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, exitcheck.ExitCheckAnalyzer, "a")
}

func TestAnalyzerPositive(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, exitcheck.ExitCheckAnalyzer, "b")
}
