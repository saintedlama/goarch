package goarch_test

import (
	"testing"

	"github.com/saintedlama/goarch/internaltest"
)

func TestLoadWorkspace_LoadsPackagesAndFiles(t *testing.T) {
	program := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	if program.Packages.Len() < 2 {
		t.Fatalf("expected at least 2 packages, got %d", program.Packages.Len())
	}

	for _, pkg := range program.Packages.All() {
		if pkg.ID == "" {
			t.Fatalf("package ID should not be empty")
		}
		if pkg.FileSet == nil {
			t.Fatalf("package %q has nil file set", pkg.ID)
		}
		if len(pkg.Files) == 0 {
			t.Fatalf("package %q should contain at least one file", pkg.ID)
		}
		for _, file := range pkg.Files {
			if file.Node == nil {
				t.Fatalf("package %q has file with nil AST node", pkg.ID)
			}
		}
	}
}

func TestLoadWorkspace_BuildsTopLevelCollections(t *testing.T) {
	program := internaltest.LoadFixtureWorkspace(t, "fixturemod")

	if program.Types.Len() == 0 {
		t.Fatalf("expected at least one type entry")
	}
	if program.Functions.Len() == 0 {
		t.Fatalf("expected at least one function entry")
	}
	if program.Variables.Len() == 0 {
		t.Fatalf("expected at least one variable entry")
	}
	if program.FunctionCalls.Len() == 0 {
		t.Fatalf("expected at least one function call entry")
	}
}
