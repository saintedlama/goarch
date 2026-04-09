package variables_test

import (
	"testing"

	"github.com/saintedlama/archscout/internaltest"
	"github.com/saintedlama/archscout/variables"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariables_MatchBuildsRefsFromPredicates(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	refs := workspace.Variables.Match(func(v variables.Item) bool {
		return v.Name == "ErrNotFound"
	})
	require.Len(t, refs, 1, "expected 1 variable ref")

	for _, f := range refs {
		assert.NotEmpty(t, f.PackageName, "ref package should not be empty")
		assert.Greater(t, f.Line, 0, "ref line should be > 0")
	}
}

func TestVariables_IsExported_ReturnsOnlyExportedVariables(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	exported := workspace.Variables.IsExported().All()

	require.NotEmpty(t, exported)
	for _, v := range exported {
		assert.True(t, v.Name[0] >= 'A' && v.Name[0] <= 'Z', "expected exported name, got %q", v.Name)
	}
}

func TestVariables_IsUnexported_ReturnsOnlyUnexportedVariables(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	all := workspace.Variables.All()
	exported := workspace.Variables.IsExported().All()
	unexported := workspace.Variables.IsUnexported().All()

	assert.Equal(t, len(all), len(exported)+len(unexported), "exported + unexported should cover all variables")
}

func TestVariables_NameMatchesRegex_FiltersOnPattern(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	items := workspace.Variables.NameMatchesRegex(`^Err`).All()

	require.NotEmpty(t, items)
	for _, v := range items {
		assert.Contains(t, v.Name, "Err")
	}
}

func TestVariables_NameMatches_FiltersOnPredicate(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	items := workspace.Variables.NameMatches(func(name string) bool {
		return name == "ErrNotFound"
	}).All()

	require.Len(t, items, 1)
	assert.Equal(t, "ErrNotFound", items[0].Name)
}
