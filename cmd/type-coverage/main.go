// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"cmp"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

const loadMode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
	packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo

var directive = regexp.MustCompile(`^//\s*(nolint\b|lint:ignore\b|lint:file-ignore\b|go:linkname\b)`)

var unsafeReflect = map[string]bool{"NewAt": true, "UnsafeAddr": true, "UnsafePointer": true}

type hatch struct {
	kind    string
	message string
}

type audit struct {
	tracked      map[string]bool
	compiled     map[string]bool
	declarations map[token.Position]bool
	untyped      map[token.Position]string
	hatches      map[token.Position]hatch
}

func main() {
	platforms := flag.String("goos", runtime.GOOS, "comma-separated GOOS values to load the module under")
	flag.Parse()
	passed, err := run(os.Stdout, ".", strings.Split(*platforms, ","))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if !passed {
		os.Exit(1)
	}
}

func run(output io.Writer, directory string, platforms []string) (bool, error) {
	root, err := filepath.Abs(directory)
	if err != nil {
		return false, err
	}
	tracked, configs, err := trackedFiles(root)
	if err != nil {
		return false, err
	}
	result := audit{
		tracked:      tracked,
		compiled:     map[string]bool{},
		declarations: map[token.Position]bool{},
		untyped:      map[token.Position]string{},
		hatches:      map[token.Position]hatch{},
	}
	for _, platform := range platforms {
		if err := result.load(root, platform); err != nil {
			return false, err
		}
	}
	for file := range tracked {
		if !result.compiled[file] {
			result.hatches[token.Position{Filename: file}] = hatch{"skipped file", "no checked platform compiles this tracked file"}
		}
	}
	for _, config := range configs {
		result.hatches[token.Position{Filename: config}] = hatch{"staticcheck config", "staticcheck.conf can exempt code from the checks that -checks all enables"}
	}
	lines := map[token.Position]string{}
	maps.Copy(lines, result.untyped)
	kinds := map[string]int{}
	for position, h := range result.hatches {
		lines[position] = h.message
		kinds[h.kind]++
	}
	for _, position := range slices.SortedFunc(maps.Keys(lines), comparePositions) {
		message := lines[position]
		name, err := filepath.Rel(root, position.Filename)
		if err != nil {
			return false, err
		}
		position.Filename = name
		fmt.Fprintf(output, "%s: %s\n", position, message)
	}
	total := len(result.declarations)
	typed := total - len(result.untyped)
	share := 100.0
	if total > 0 {
		share = 100 * float64(typed) / float64(total)
	}
	fmt.Fprintf(output, "Type coverage %.2f%%, %d of %d declarations typed.\n", share, typed, total)
	fmt.Fprintf(output, "Escape hatches %d.", len(result.hatches))
	for _, kind := range slices.Sorted(maps.Keys(kinds)) {
		fmt.Fprintf(output, " %s %d.", kind, kinds[kind])
	}
	fmt.Fprintln(output)
	return len(lines) == 0, nil
}

func comparePositions(a, b token.Position) int {
	return cmp.Or(strings.Compare(a.Filename, b.Filename), cmp.Compare(a.Offset, b.Offset))
}

func trackedFiles(root string) (map[string]bool, []string, error) {
	command := exec.Command("git", "ls-files", "-z", "--", "*.go", ":(icase)*staticcheck.conf")
	command.Dir = root
	var stderr bytes.Buffer
	command.Stderr = &stderr
	listing, err := command.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("git ls-files failed in %s: %w %s", root, err, stderr.String())
	}
	files := map[string]bool{}
	var configs []string
	for name := range strings.SplitSeq(string(listing), "\x00") {
		switch {
		case strings.EqualFold(filepath.Base(name), "staticcheck.conf"):
			configs = append(configs, filepath.Join(root, name))
		case strings.HasSuffix(name, ".go"):
			files[filepath.Join(root, name)] = true
		}
	}
	return files, configs, nil
}

func (a *audit) load(root, platform string) error {
	config := &packages.Config{
		Mode:  loadMode,
		Dir:   root,
		Env:   append(os.Environ(), "GOOS="+platform),
		Tests: true,
	}
	loaded, err := packages.Load(config, "./...")
	if err != nil {
		return err
	}
	var problems []error
	for _, p := range loaded {
		for _, problem := range p.Errors {
			problems = append(problems, fmt.Errorf("GOOS=%s %s: %s", platform, p.ID, problem))
		}
	}
	if len(problems) > 0 {
		return errors.Join(problems...)
	}
	for _, p := range loaded {
		a.inspect(p)
	}
	return nil
}

func (a *audit) inspect(p *packages.Package) {
	at := func(node token.Pos) token.Position { return p.Fset.Position(node) }
	declare := func(node token.Pos, untyped bool, description string) {
		position := at(node)
		if !a.tracked[position.Filename] || a.declarations[position] {
			return
		}
		a.declarations[position] = true
		if untyped {
			a.untyped[position] = description + ", which holds an empty interface"
		}
	}
	for identifier, object := range p.TypesInfo.Defs {
		switch object := object.(type) {
		case *types.Var, *types.Const:
			declare(identifier.Pos(), holdsEmptyInterface(object.Type()),
				identifier.Name+" has type "+types.TypeString(object.Type(), types.RelativeTo(p.Types)))
		case *types.TypeName:
			if _, parameter := object.Type().(*types.TypeParam); !parameter {
				declare(identifier.Pos(), typeDeclarationHolds(object), "type "+identifier.Name)
			}
		}
	}
	for node, object := range p.TypesInfo.Implicits {
		if variable, ok := object.(*types.Var); ok {
			declare(node.Pos(), holdsEmptyInterface(variable.Type()),
				"implicit declaration has type "+types.TypeString(variable.Type(), types.RelativeTo(p.Types)))
		}
	}
	for _, file := range p.Syntax {
		name := p.Fset.File(file.Pos()).Name()
		if !a.tracked[name] {
			continue
		}
		a.compiled[name] = true
		a.inspectFile(file, p.TypesInfo, at)
	}
}

func (a *audit) inspectFile(file *ast.File, info *types.Info, at func(token.Pos) token.Position) {
	if ast.IsGenerated(file) {
		a.hatches[at(file.Package)] = hatch{"generated header", "generated-code header makes linters skip this file"}
	}
	for _, group := range file.Comments {
		for _, comment := range group.List {
			if match := directive.FindStringSubmatch(comment.Text); match != nil {
				a.hatches[at(comment.Pos())] = hatch{"directive", "//" + match[1] + " directive escapes a check"}
			}
		}
	}
	for _, spec := range file.Imports {
		if spec.Path.Value == `"unsafe"` {
			a.hatches[at(spec.Pos())] = hatch{"unsafe", "import of package unsafe"}
		}
	}
	constraints := map[ast.Expr]bool{}
	ast.Inspect(file, func(node ast.Node) bool {
		if expression, ok := node.(ast.Expr); ok && constraints[expression] {
			return false
		}
		switch node := node.(type) {
		case *ast.FuncType:
			markConstraints(node.TypeParams, constraints)
		case *ast.TypeSpec:
			markConstraints(node.TypeParams, constraints)
		case *ast.TypeAssertExpr:
			if _, checked := info.Types[node].Type.(*types.Tuple); node.Type != nil && !checked {
				a.hatches[at(node.Lparen)] = hatch{"unchecked assertion", "type assertion without the comma-ok form"}
			}
		case *ast.AssignStmt:
			a.discardedAssertion(node.Lhs, node.Rhs, at)
		case *ast.ValueSpec:
			names := make([]ast.Expr, len(node.Names))
			for i, name := range node.Names {
				names[i] = name
			}
			a.discardedAssertion(names, node.Values, at)
		case *ast.SelectorExpr:
			if function, ok := info.Uses[node.Sel].(*types.Func); ok && function.Pkg() != nil &&
				function.Pkg().Path() == "reflect" && unsafeReflect[function.Name()] {
				a.hatches[at(node.Sel.Pos())] = hatch{"unsafe", "reflect " + function.Name() + " bypasses type safety"}
			}
		case *ast.Ident:
			if object, ok := info.Uses[node].(*types.TypeName); ok && object.Pkg() == nil && object.Name() == "any" {
				a.hatches[at(node.Pos())] = hatch{"empty interface", "any written outside a type parameter constraint"}
			}
		case *ast.InterfaceType:
			if len(node.Methods.List) == 0 {
				a.hatches[at(node.Pos())] = hatch{"empty interface", "interface{} written outside a type parameter constraint"}
			}
		}
		return true
	})
}

func (a *audit) discardedAssertion(targets, values []ast.Expr, at func(token.Pos) token.Position) {
	if len(targets) != 2 || len(values) != 1 {
		return
	}
	assertion, isAssertion := ast.Unparen(values[0]).(*ast.TypeAssertExpr)
	if blank, isIdent := targets[1].(*ast.Ident); isAssertion && isIdent && blank.Name == "_" {
		a.hatches[at(assertion.Lparen)] = hatch{"unchecked assertion", "type assertion discards its ok result"}
	}
}

func markConstraints(parameters *ast.FieldList, constraints map[ast.Expr]bool) {
	if parameters == nil {
		return
	}
	for _, field := range parameters.List {
		constraints[field.Type] = true
	}
}

func typeDeclarationHolds(object *types.TypeName) bool {
	if object.IsAlias() {
		return holdsEmptyInterface(object.Type())
	}
	switch underlying := object.Type().Underlying().(type) {
	case *types.Interface:
		return underlying.Empty()
	case *types.Struct:
		return false
	default:
		return holdsEmptyInterface(underlying)
	}
}

func holdsEmptyInterface(t types.Type) bool {
	return holds(t, map[types.Type]bool{})
}

func holds(t types.Type, seen map[types.Type]bool) bool {
	if seen[t] {
		return false
	}
	seen[t] = true
	switch t := t.(type) {
	case *types.Alias:
		return holds(types.Unalias(t), seen)
	case *types.Named:
		if underlying, ok := t.Underlying().(*types.Interface); ok && underlying.Empty() {
			return true
		}
		return slices.ContainsFunc(slices.Collect(t.TypeArgs().Types()), func(argument types.Type) bool {
			return holds(argument, seen)
		})
	case *types.Interface:
		return t.Empty()
	case *types.Pointer:
		return holds(t.Elem(), seen)
	case *types.Slice:
		return holds(t.Elem(), seen)
	case *types.Array:
		return holds(t.Elem(), seen)
	case *types.Chan:
		return holds(t.Elem(), seen)
	case *types.Map:
		return holds(t.Key(), seen) || holds(t.Elem(), seen)
	case *types.Signature:
		return holds(t.Params(), seen) || holds(t.Results(), seen)
	case *types.Tuple:
		return slices.ContainsFunc(slices.Collect(t.Variables()), func(variable *types.Var) bool {
			return holds(variable.Type(), seen)
		})
	case *types.Struct:
		return slices.ContainsFunc(slices.Collect(t.Fields()), func(field *types.Var) bool {
			return holds(field.Type(), seen)
		})
	}
	return false
}
