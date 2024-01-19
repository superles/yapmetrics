package exitcheck_test

import (
	"testing"

	"github.com/superles/yapmetrics/internal/exitcheck"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, exitcheck.ExitCheckAnalyzer, "a")
}
