package variables_test

import (
	"testing"

	"github.com/saintedlama/goarch/analysis/internaltest"
	"github.com/saintedlama/goarch/analysis/variables"
)

func TestVariables_MatchBuildsFindingsFromPredicates(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	findings := workspace.Variables.Match(variables.MatchFunc(func(v variables.Item) bool {
		return v.Name == "GlobalCounter"
	}))
	if len(findings) != 1 {
		t.Fatalf("expected 1 variable finding, got %d", len(findings))
	}

	for _, f := range findings {
		if f.PackageName == "" {
			t.Fatalf("finding package should not be empty")
		}
		if f.Line <= 0 {
			t.Fatalf("finding line should be > 0")
		}
	}
}
