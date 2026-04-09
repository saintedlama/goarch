package functions_test

import (
	"testing"

	"github.com/saintedlama/archscout/functions"
	"github.com/saintedlama/archscout/internaltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctions_MatchBuildsRefsFromPredicates(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	refs := workspace.Functions.Match(func(fn functions.Item) bool {
		return fn.Name == "NewOrder"
	})
	require.Len(t, refs, 1, "expected 1 function ref")

	for _, f := range refs {
		assert.NotEmpty(t, f.PackageName, "ref package should not be empty")
		assert.Greater(t, f.Line, 0, "ref line should be > 0")
	}
}

func TestFunctions_IsMethod_ReturnsOnlyMethods(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	methods := workspace.Functions.IsMethod().All()

	require.NotEmpty(t, methods)
	for _, fn := range methods {
		assert.NotEmpty(t, fn.Receiver, "method should have a non-empty receiver, got %q", fn.Name)
	}
}

func TestFunctions_IsFunction_ReturnsOnlyFreeFunctions(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	functions := workspace.Functions.IsFunction().All()

	require.NotEmpty(t, functions)
	for _, fn := range functions {
		assert.Empty(t, fn.Receiver, "free function should have an empty receiver, got %q", fn.Name)
	}
}

func TestFunctions_IsMethodAndIsFunction_PartitionAllFunctions(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	total := workspace.Functions.Len()
	methods := workspace.Functions.IsMethod().Len()
	freeFuncs := workspace.Functions.IsFunction().Len()

	assert.Equal(t, total, methods+freeFuncs, "methods + free functions should equal total")
}

func TestFunctions_HasReceiver_FiltersToSpecificReceiver(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	items := workspace.Functions.HasReceiver("OrderService").All()

	require.NotEmpty(t, items)
	for _, fn := range items {
		assert.Contains(t, fn.Receiver, "OrderService")
	}
}

func TestFunctions_IsExported_ReturnsOnlyExportedFunctions(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	exported := workspace.Functions.IsExported().All()

	require.NotEmpty(t, exported)
	for _, fn := range exported {
		assert.True(t, fn.Name[0] >= 'A' && fn.Name[0] <= 'Z', "expected exported name, got %q", fn.Name)
	}
}

func TestFunctions_NameMatchesRegex_FiltersOnPattern(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	items := workspace.Functions.NameMatchesRegex(`^New`).All()

	require.NotEmpty(t, items)
	for _, fn := range items {
		assert.Contains(t, fn.Name, "New")
	}
}

func TestFunctions_NameMatches_FiltersOnPredicate(t *testing.T) {
	workspace := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	items := workspace.Functions.NameMatches(func(name string) bool {
		return name == "NewOrder"
	}).All()

	require.Len(t, items, 1)
	assert.Equal(t, "NewOrder", items[0].Name)
}
