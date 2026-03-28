package analysis

import (
	"context"
	"path/filepath"
	"testing"
)

func mustLoadFixtureProgramAST(t *testing.T, fixtureName string) *ProgramAST {
	t.Helper()

	dir := fixtureDir(t, fixtureName)
	program, err := LoadProgramAST(context.Background(), dir, nil)
	if err != nil {
		t.Fatalf("LoadProgramAST(%q) failed: %v", dir, err)
	}

	return program
}

func fixtureDir(t *testing.T, fixtureName string) string {
	t.Helper()
	return filepath.Join("testdata", fixtureName)
}
