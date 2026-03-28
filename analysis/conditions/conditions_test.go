package conditions_test

import (
	"testing"

	"github.com/saintedlama/goarch/analysis/conditions"
	"github.com/saintedlama/goarch/analysis/internaltest"
)

func TestConditions_MatchBuildsFindingsFromPredicates(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	findings := workspace.Conditions.Match(conditions.MatchFunc(func(cond conditions.Item) (bool, string) {
		return cond.Kind == "if", "condition predicate matched"
	}))
	if len(findings) == 0 {
		t.Fatalf("expected at least 1 condition finding")
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
