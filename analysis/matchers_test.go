package analysis

import "testing"

func TestMatchers_BuildFindingsFromPredicates(t *testing.T) {
	program := mustLoadFixtureProgramAST(t, "fixturemod")

	packageFindings := program.Packages.Match(PackageMatchFunc(func(pkg PackageAST) (bool, string) {
		return pkg.Name == "main", "package predicate matched"
	}))
	if len(packageFindings) == 0 {
		t.Fatalf("expected package findings")
	}

	typeFindings := program.Types.Match(TypeMatchFunc(func(typ TypePointcut) (bool, string) {
		return typ.Name == "widget", "type predicate matched"
	}))
	if len(typeFindings) != 1 {
		t.Fatalf("expected 1 type finding, got %d", len(typeFindings))
	}

	functionFindings := program.Functions.Match(FunctionMatchFunc(func(fn FunctionPointcut) (bool, string) {
		return fn.Name == "rootErr", "function predicate matched"
	}))
	if len(functionFindings) != 1 {
		t.Fatalf("expected 1 function finding, got %d", len(functionFindings))
	}

	variableFindings := program.Variables.Match(VariableMatchFunc(func(v VariablePointcut) (bool, string) {
		return v.Name == "globalCounter", "variable predicate matched"
	}))
	if len(variableFindings) != 1 {
		t.Fatalf("expected 1 variable finding, got %d", len(variableFindings))
	}

	callFindings := program.FunctionCalls.Match(FunctionCallMatchFunc(func(call FunctionCallPointcut) (bool, string) {
		return call.Callee == "fmt.Errorf", "call predicate matched"
	}))
	if len(callFindings) != 2 {
		t.Fatalf("expected 2 call findings, got %d", len(callFindings))
	}

	conditionFindings := program.Conditions.Match(ConditionMatchFunc(func(cond ConditionPointcut) (bool, string) {
		return cond.Kind == "if", "condition predicate matched"
	}))
	if len(conditionFindings) == 0 {
		t.Fatalf("expected at least 1 condition finding")
	}

	for _, findings := range [][]Finding{packageFindings, typeFindings, functionFindings, variableFindings, callFindings, conditionFindings} {
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
}
