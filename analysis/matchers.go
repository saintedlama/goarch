package analysis

import (
	pcconditions "goarch/analysis/pointcuts/conditions"
	pcfunctioncalls "goarch/analysis/pointcuts/functioncalls"
	pcfunctions "goarch/analysis/pointcuts/functions"
	pcpackages "goarch/analysis/pointcuts/packages"
	pctypes "goarch/analysis/pointcuts/types"
	pcvariables "goarch/analysis/pointcuts/variables"
)

// Compatibility aliases for top-level matcher APIs.
type PackageMatcher = pcpackages.Matcher
type PackageMatchFunc = pcpackages.MatchFunc

type TypeMatcher = pctypes.Matcher
type TypeMatchFunc = pctypes.MatchFunc

type FunctionMatcher = pcfunctions.Matcher
type FunctionMatchFunc = pcfunctions.MatchFunc

type VariableMatcher = pcvariables.Matcher
type VariableMatchFunc = pcvariables.MatchFunc

type FunctionCallMatcher = pcfunctioncalls.Matcher
type FunctionCallMatchFunc = pcfunctioncalls.MatchFunc

type ConditionMatcher = pcconditions.Matcher
type ConditionMatchFunc = pcconditions.MatchFunc

// MatchPackages runs a matcher over all packages and returns generated findings.
func MatchPackages(program *ProgramAST, matcher PackageMatcher) []Finding {
	if program == nil || matcher == nil {
		return nil
	}
	return program.Packages.Match(matcher)
}

// MatchTypes runs a matcher over all type pointcuts and returns generated findings.
func MatchTypes(program *ProgramAST, matcher TypeMatcher) []Finding {
	if program == nil || matcher == nil {
		return nil
	}
	return program.Types.Match(matcher)
}

// MatchFunctions runs a matcher over all function pointcuts and returns generated findings.
func MatchFunctions(program *ProgramAST, matcher FunctionMatcher) []Finding {
	if program == nil || matcher == nil {
		return nil
	}
	return program.Functions.Match(matcher)
}

// MatchVariables runs a matcher over all variable pointcuts and returns generated findings.
func MatchVariables(program *ProgramAST, matcher VariableMatcher) []Finding {
	if program == nil || matcher == nil {
		return nil
	}
	return program.Variables.Match(matcher)
}

// MatchFunctionCalls runs a matcher over all call pointcuts and returns generated findings.
func MatchFunctionCalls(program *ProgramAST, matcher FunctionCallMatcher) []Finding {
	if program == nil || matcher == nil {
		return nil
	}
	return program.FunctionCalls.Match(matcher)
}

// MatchConditions runs a matcher over all condition pointcuts and returns generated findings.
func MatchConditions(program *ProgramAST, matcher ConditionMatcher) []Finding {
	if program == nil || matcher == nil {
		return nil
	}
	return program.Conditions.Match(matcher)
}
