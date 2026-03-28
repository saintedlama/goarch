package functions_test

import (
	"testing"

	"github.com/saintedlama/goarch/analysis/functions"
	"github.com/saintedlama/goarch/analysis/internaltest"
)

func TestFunctions_MatchBuildsFindingsFromPredicates(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	findings := workspace.Functions.Match(functions.MatchFunc(func(fn functions.Item) (bool, string) {
		return fn.Name == "rootErr", "function predicate matched"
	}))
	if len(findings) != 1 {
		t.Fatalf("expected 1 function finding, got %d", len(findings))
	}

	for _, f := range findings {
		if f.Message == "" {
			t.Fatalf("finding message should not be empty")
		}
		if f.Package == "" {
			t.Fatalf("finding package should not be empty")
		}
		if f.Line <= 0 {
			t.Fatalf("finding line should be > 0")
		}
	}
}
