package conditions_test

import (
	"testing"

	"github.com/saintedlama/goarch/analysis/conditions"
	"github.com/saintedlama/goarch/analysis/internaltest"
)

func TestConditions_MatchBuildsFindingsFromPredicates(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	findings := workspace.Conditions.Match(conditions.MatchFunc(func(cond conditions.Item) bool {
		return cond.Kind == "if"
	}))
	if len(findings) == 0 {
		t.Fatalf("expected at least 1 condition finding")
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
