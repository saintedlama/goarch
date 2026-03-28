package conditions_test

import (
	"testing"

	"github.com/saintedlama/goarch/conditions"
	"github.com/saintedlama/goarch/internaltest"
)

func TestConditions_MatchBuildsRefsFromPredicates(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	refs := workspace.Conditions.Match(func(cond conditions.Item) bool {
		return cond.Kind == "if"
	})
	if len(refs) == 0 {
		t.Fatalf("expected at least 1 condition ref")
	}

	for _, f := range refs {
		if f.PackageName == "" {
			t.Fatalf("ref package should not be empty")
		}
		if f.Line <= 0 {
			t.Fatalf("ref line should be > 0")
		}
	}
}
