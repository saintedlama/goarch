# ![archscout logo](./assets/logo.png) ArchScout

`archscout` helps you keep architecture visible and enforceable in Go codebases.

Use it to:

- Explore code structure quickly (packages, files, types, calls, dependencies)
- Write architecture tests as code
- Validate dependency boundaries continuously in CI

## Why archscout

Architecture often lives in docs, not in tests. `archscout` lets you move those rules into executable checks.

Examples:

- "domain must not depend on infrastructure"
- "library code must not call panic or os.Exit"
- "application layer may only depend on domain"

When a rule is violated, you get source refs you can print in test failures.

## Install

```bash
go get github.com/saintedlama/archscout
```

## Quick Start

```go
package architecture_test

import (
  "context"
  "testing"

  "github.com/saintedlama/archscout"
)

func TestDomainDoesNotDependOnInfrastructure(t *testing.T) {
  workspace, err := archscout.LoadWorkspace(context.Background(), ".")
  if err != nil {
    t.Fatalf("LoadWorkspace failed: %v", err)
  }

  rule := archscout.Rule("domain must not depend on infrastructure").
    Dependencies().
    InPackage("github.com/your-project/domain/...").
    DependOn("github.com/your-project/infrastructure/...")

  rule.Test(t, workspace)
}
```

## Core Workflows

### 1. Explore a codebase

```go
refs := workspace.FunctionCalls.
  InPackage("github.com/your-project/...").
  IsNotTest().
  Match(func(call archscout.FunctionCall) bool {
    return call.Callee == "fmt.Errorf"
  })
```

### 2. Validate architecture with reusable rules

```go
forbidden := map[string]bool{"panic": true, "os.Exit": true}

rule := archscout.Rule("panic and os.Exit forbidden in library code").
  FunctionCalls().
  InPackage("github.com/your-project/...").
  NotInPackage("github.com/your-project/internal/...").
  IsNotTest().
  Match(func(fc archscout.FunctionCall) bool {
    return forbidden[fc.Callee]
  })

rule.Test(t, workspace)
```

### 3. Assert existence

Use `ShouldExist()` on `Packages`, `Types`, and `Functions` rules to assert that at
least one entry survives the filter chain. Combine with `Match` to pin a specific item:

```go
archscout.Rule("domain package must exist").
  Packages().
  InPackage("github.com/your-project/domain").
  ShouldExist().
  Test(t, workspace)

archscout.Rule("Repository interface must be defined in domain").
  Types().
  InPackage("github.com/your-project/domain").
  ShouldExist().
  Match(func(t archscout.Type) bool { return t.Name == "Repository" }).
  Test(t, workspace)
```

### 4. Reason about dependencies

Dependency checks can be done directly or through files/packages.

```go
rule := archscout.Rule("files with no stdlib deps").
  Files().
  Match(func(file archscout.File) bool {
    return file.Dependencies().IsStandardLibrary().Len() == 0
  })

rule.Test(t, workspace)
```

For hierarchy-style reporting, use `workspace.Dependencies.Tree()`.

### 5. Explore dependencies in large codebases

Three aggregation helpers make it easy to answer high-level questions without
counting raw import statements:

```go
mod := archscout.Module("github.com/your-project")

// What does the UI layer reach (workspace-internal, non-test)?
targets := workspace.Dependencies.
  InPackage(mod.Pkg("ui/...")).
  IsNotTest().
  IsWithinWorkspace().
  UniqueTargets()
// → ["github.com/your-project/audio", "github.com/your-project/domain", ...]

// Who imports the domain layer?
importers := workspace.Dependencies.
  DependOn(mod.Pkg("domain/...")).
  IsNotTest().
  UniqueSourcePackages()
// → ["github.com/your-project/application", "github.com/your-project/ui/tracker", ...]

// Full per-package breakdown
for pkg, deps := range workspace.Dependencies.IsNotTest().IsWithinWorkspace().GroupBySourcePackage() {
  fmt.Printf("%s → %v\n", pkg, deps.UniqueTargets())
}
```

### 6. Reduce repetition with Module

Use `Module` to avoid repeating the module path across patterns:

```go
mod := archscout.Module("github.com/your-project")

archscout.Rule("ui/common must not depend on other internal packages").
  Dependencies().
  InPackage(mod.Pkg("ui/common/...")).
  IsNotTest().
  DependOn(mod.Pkgs(
    "audio/...",
    "persistence/...",
    "player/...",
  )...).
  Test(t, workspace)
```

`mod.Pkg("sub/path")` returns a single fully-qualified pattern.
`mod.Pkgs("a/...", "b/...")` returns a `[]string` of fully-qualified patterns.

## What You Can Query

`archscout` exposes seven collections on `Workspace`:

| Field           | Item type      | Notable fields                                                                      |
| --------------- | -------------- | ----------------------------------------------------------------------------------- |
| `Packages`      | `Package`      | `ID`, `Name`, `Files`, `Dependencies()`                                             |
| `Files`         | `File`         | `Filename`, `Dependencies()`                                                        |
| `Types`         | `Type`         | `Name`, `Kind`                                                                      |
| `Functions`     | `Function`     | `Name`, `Receiver`                                                                  |
| `Variables`     | `Variable`     | `Name`, `Kind`                                                                      |
| `FunctionCalls` | `FunctionCall` | `Callee`                                                                            |
| `Dependencies`  | `Dependency`   | `ImportPath`, `WithinWorkspace`, `External`, `StandardLibrary`, `TargetPackageName` |

All collections support:

| Method                      | Description                                             |
| --------------------------- | ------------------------------------------------------- |
| `All()`                     | Returns a snapshot slice of all items                   |
| `Len()`                     | Number of items                                         |
| `Match(func)`               | Applies a predicate; returns matching `Refs`            |
| `InPackage(patterns...)`    | Keeps items whose source package matches any pattern    |
| `NotInPackage(patterns...)` | Excludes items whose source package matches any pattern |
| `IsTest()`                  | Keeps items from `_test.go` files                       |
| `IsNotTest()`               | Excludes items from `_test.go` files                    |

Dependencies additionally support:

| Method                       | Description                                               |
| ---------------------------- | --------------------------------------------------------- |
| `DependOn(patterns...)`      | Keeps items whose import path matches any pattern         |
| `DependsOn(pattern)`         | Keeps items whose import path matches a single pattern    |
| `DoNotDependOn(patterns...)` | Excludes items whose import path matches any pattern      |
| `IsWithinWorkspace()`        | Keeps imports that resolve to workspace packages          |
| `IsExternal()`               | Keeps imports that resolve outside the workspace          |
| `IsStandardLibrary()`        | Keeps standard library imports                            |
| `IsThirdParty()`             | Keeps external, non-stdlib imports                        |
| `UniqueTargets()`            | Sorted, deduplicated import paths in the collection       |
| `UniqueSourcePackages()`     | Sorted, deduplicated source package IDs in the collection |
| `GroupBySourcePackage()`     | Partitions into one sub-collection per source package     |
| `Tree()`                     | Builds a hierarchical `TreeNode` from import paths        |

## Refs and Formatting

Rule violations are returned as `Refs` — each `Ref` identifies a source location:

```go
refs, err := rule.Evaluate(workspace)
fmt.Println(archscout.FormatRefs(refs))

// Customise output
fmt.Println(archscout.FormatRefs(refs,
  archscout.WithRefPackage(),
  archscout.WithRefKind(),
  archscout.WithoutRefColumn(),
))
```

Available format options: `WithRefPackage()`, `WithRefKind()`, `WithoutRefFile()`,
`WithoutRefLine()`, `WithoutRefColumn()`, `WithoutRefMatch()`, `WithRefSeparator(sep)`,
`WithoutSeparator()`.

## Public API

- `LoadWorkspace(ctx, dir, opts...) (*Workspace, error)`
- `WithReporter(func(string)) LoadWorkspaceOption` — progress callback
- `WithInMemoryCache() LoadWorkspaceOption` — reuse a loaded workspace within the process
- `Module(path)` — helper for building fully-qualified package patterns
- `Rule(name)` — entry point for all rule construction

Rule types expose:

- fluent filters (package/test and kind-specific filters)
- `ShouldExist()` — assert at least one match exists (`Packages`, `Types`, `Functions`)
- `Match(func)`
- `Evaluate(workspace) (Refs, error)`
- `Test(t, workspace)` — fails the test if any refs are returned (or none when `ShouldExist`)

## Development

```bash
make fmt
make vet
make build
make test-verbose
```

## Notes

- `LoadWorkspace` expects a Go module directory with `go.mod`.
- `WithReporter(...)` is optional and useful for progress output.
- `WithInMemoryCache()` is optional and reuses a loaded workspace by path.
- Pattern matching: a pattern ending in `/...` matches the base path and all sub-paths.

## License

This project is licensed under the MIT License. See `LICENSE` for details.
