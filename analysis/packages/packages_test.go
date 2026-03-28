package packages_test

import (
	"testing"

	"github.com/saintedlama/goarch/analysis/internaltest"
	"github.com/saintedlama/goarch/analysis/packages"
)

func TestPackages_MatchBuildsFindingsFromPredicates(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	findings := workspace.Packages.Match(packages.MatchFunc(func(pkg packages.Item) (bool, string) {
		return pkg.Name == "main", "package predicate matched"
	}))
	if len(findings) == 0 {
		t.Fatalf("expected package findings")
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
