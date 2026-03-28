package internaltest

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/saintedlama/goarch"
)

func LoadFixtureWorkspace(t testing.TB, fixtureName string) *goarch.Workspace {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	dir := filepath.Join(filepath.Dir(filename), "..", "testdata", fixtureName)
	workspace, err := goarch.LoadWorkspace(context.Background(), dir)
	if err != nil {
		t.Fatalf("LoadWorkspace(%q) failed: %v", dir, err)
	}

	return workspace
}
