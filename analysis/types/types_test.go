package types_test

import (
	"testing"

	"github.com/saintedlama/goarch/analysis/internaltest"
	"github.com/saintedlama/goarch/analysis/types"
)

func TestTypes_MatchBuildsFindingsFromPredicates(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	findings := workspace.Types.Match(types.MatchFunc(func(typ types.Item) (bool, string) {
		return typ.Name == "widget", "type predicate matched"
	}))
	if len(findings) != 1 {
		t.Fatalf("expected 1 type finding, got %d", len(findings))
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
