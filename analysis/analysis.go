package analysis

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"

	pccommon "goarch/analysis/pointcuts/common"
	pcconditions "goarch/analysis/pointcuts/conditions"
	pcfunctioncalls "goarch/analysis/pointcuts/functioncalls"
	pcfunctions "goarch/analysis/pointcuts/functions"
	pcpackages "goarch/analysis/pointcuts/packages"
	pctypes "goarch/analysis/pointcuts/types"
	pcvariables "goarch/analysis/pointcuts/variables"

	toolspackages "golang.org/x/tools/go/packages"
)

// ProgramAST is the loaded AST structure for all discovered packages.
type ProgramAST struct {
	Packages      pcpackages.Collection
	Types         pctypes.Collection
	Functions     pcfunctions.Collection
	Variables     pcvariables.Collection
	FunctionCalls pcfunctioncalls.Collection
	Conditions    pcconditions.Collection
}

// Compatibility aliases for consumers using the previous analysis package surface.
type Finding = pccommon.Finding
type PointcutRef = pccommon.Ref

type PackagePointcuts = pcpackages.Collection
type PackageAST = pcpackages.Item
type FileAST = pcpackages.File

type TypePointcuts = pctypes.Collection
type TypePointcut = pctypes.Item

type FunctionPointcuts = pcfunctions.Collection
type FunctionPointcut = pcfunctions.Item

type VariablePointcuts = pcvariables.Collection
type VariablePointcut = pcvariables.Item

type FunctionCallPointcuts = pcfunctioncalls.Collection
type FunctionCallPointcut = pcfunctioncalls.Item

type ConditionPointcuts = pcconditions.Collection
type ConditionPointcut = pcconditions.Item

// LoadProgramAST loads all packages in dir and returns a full AST structure.
// progress is called with status messages; pass nil to suppress them.
func LoadProgramAST(ctx context.Context, dir string, progress func(string)) (*ProgramAST, error) {
	cfg := &toolspackages.Config{
		Dir: dir,
		Mode: toolspackages.NeedName | toolspackages.NeedFiles | toolspackages.NeedSyntax |
			toolspackages.NeedCompiledGoFiles | toolspackages.NeedImports,
		Context: ctx,
	}

	report := func(msg string) {
		if progress != nil {
			progress(msg)
		}
	}

	report("Loading packages (./...)...")
	pkgs, err := toolspackages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("loading packages: %w", err)
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found in %q", dir)
	}
	report(fmt.Sprintf("Loaded %d package(s)", len(pkgs)))

	program := &ProgramAST{}
	for _, pkg := range pkgs {
		report(fmt.Sprintf("Analyzing %s...", pkg.ID))

		p := pcpackages.Item{
			ID:     pkg.ID,
			Name:   pkg.Name,
			Fset:   pkg.Fset,
			Errors: pkg.Errors,
		}

		for i, file := range pkg.Syntax {
			if file == nil {
				continue
			}

			filename := ""
			if i < len(pkg.CompiledGoFiles) {
				filename = pkg.CompiledGoFiles[i]
			} else if i < len(pkg.GoFiles) {
				filename = pkg.GoFiles[i]
			}

			p.Files = append(p.Files, pcpackages.File{
				Filename: filename,
				Node:     file,
			})

			indexFilePointcuts(&program.Types, &program.Functions, &program.Variables, &program.FunctionCalls, &program.Conditions, p, filename, file)
		}

		program.Packages.Add(p)
	}

	return program, nil
}

func indexFilePointcuts(
	types *pctypes.Collection,
	functions *pcfunctions.Collection,
	variables *pcvariables.Collection,
	functionCalls *pcfunctioncalls.Collection,
	conditions *pcconditions.Collection,
	pkg pcpackages.Item,
	filename string,
	file *ast.File,
) {
	if file == nil {
		return
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.TypeSpec:
			types.Add(pctypes.Item{
				Ref:  newPointcutRef(pkg, filename, node),
				Name: node.Name.Name,
				Kind: exprKind(node.Type),
				Node: node,
			})

		case *ast.FuncDecl:
			receiver := ""
			if node.Recv != nil && len(node.Recv.List) > 0 {
				receiver = exprText(node.Recv.List[0].Type)
			}
			functions.Add(pcfunctions.Item{
				Ref:      newPointcutRef(pkg, filename, node),
				Name:     node.Name.Name,
				Receiver: receiver,
				Node:     node,
			})

		case *ast.ValueSpec:
			kind := "var"
			if genDecl, ok := enclosingGenDecl(file, node); ok && genDecl.Tok == token.CONST {
				kind = "const"
			}
			for _, name := range node.Names {
				variables.Add(pcvariables.Item{
					Ref:  newPointcutRef(pkg, filename, name),
					Name: name.Name,
					Kind: kind,
					Node: name,
				})
			}

		case *ast.CallExpr:
			functionCalls.Add(pcfunctioncalls.Item{
				Ref:    newPointcutRef(pkg, filename, node),
				Callee: calleeName(node.Fun),
				Node:   node,
			})

		case *ast.IfStmt:
			conditions.Add(pcconditions.Item{Ref: newPointcutRef(pkg, filename, node), Kind: "if", Node: node})
		case *ast.SwitchStmt:
			conditions.Add(pcconditions.Item{Ref: newPointcutRef(pkg, filename, node), Kind: "switch", Node: node})
		case *ast.TypeSwitchStmt:
			conditions.Add(pcconditions.Item{Ref: newPointcutRef(pkg, filename, node), Kind: "type-switch", Node: node})
		case *ast.CaseClause:
			conditions.Add(pcconditions.Item{Ref: newPointcutRef(pkg, filename, node), Kind: "case", Node: node})
		case *ast.SelectStmt:
			conditions.Add(pcconditions.Item{Ref: newPointcutRef(pkg, filename, node), Kind: "select", Node: node})
		case *ast.CommClause:
			conditions.Add(pcconditions.Item{Ref: newPointcutRef(pkg, filename, node), Kind: "comm", Node: node})
		}

		return true
	})
}

func newPointcutRef(pkg pcpackages.Item, fallbackFilename string, n ast.Node) pccommon.Ref {
	pos := pkg.Fset.PositionFor(n.Pos(), true)
	filename := fallbackFilename
	if pos.Filename != "" {
		filename = pos.Filename
	}

	return pccommon.Ref{
		PackageID:   pkg.ID,
		PackageName: pkg.Name,
		Filename:    filename,
		Line:        pos.Line,
		Column:      pos.Column,
	}
}

func calleeName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		prefix := exprText(e.X)
		if prefix == "" {
			return e.Sel.Name
		}
		return prefix + "." + e.Sel.Name
	default:
		return exprText(expr)
	}
}

func exprKind(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.StructType:
		return "struct"
	case *ast.InterfaceType:
		return "interface"
	case *ast.ArrayType:
		return "array"
	case *ast.MapType:
		return "map"
	case *ast.FuncType:
		return "func"
	case *ast.ChanType:
		return "chan"
	default:
		return "type"
	}
}

func exprText(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		inner := exprText(e.X)
		if inner == "" {
			return "*"
		}
		return "*" + inner
	case *ast.SelectorExpr:
		left := exprText(e.X)
		if left == "" {
			return e.Sel.Name
		}
		return left + "." + e.Sel.Name
	default:
		return ""
	}
}

func enclosingGenDecl(file *ast.File, target *ast.ValueSpec) (*ast.GenDecl, bool) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			if spec == target {
				return genDecl, true
			}
		}
	}
	return nil, false
}
