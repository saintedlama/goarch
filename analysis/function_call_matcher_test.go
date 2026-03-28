package analysis

import (
	"strings"
	"testing"
)

func TestFunctionCalls_FindsExpectedFmtErrorfCalls(t *testing.T) {
	program := mustLoadFixtureProgramAST(t, "fixturemod")

	findings := program.FunctionCalls.Match(FunctionCallMatchFunc(func(call FunctionCallPointcut) (bool, string) {
		if call.Callee != "fmt.Errorf" {
			return false, ""
		}
		return true, "found fmt.Errorf"
	}))
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}

	var sawRoot, sawSub bool
	for _, f := range findings {
		if !strings.Contains(f.Message, "fmt.Errorf") {
			t.Fatalf("unexpected finding message: %q", f.Message)
		}
		if strings.HasSuffix(strings.ReplaceAll(f.Filename, "\\", "/"), "/main.go") {
			sawRoot = true
		}
		if strings.HasSuffix(strings.ReplaceAll(f.Filename, "\\", "/"), "/subpkg/sub.go") {
			sawSub = true
		}
	}

	if !sawRoot {
		t.Fatalf("did not find fmt.Errorf in fixture main.go")
	}
	if !sawSub {
		t.Fatalf("did not find fmt.Errorf in fixture subpkg/sub.go")
	}
}
