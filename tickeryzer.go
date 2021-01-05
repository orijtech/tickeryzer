// Copyright 2021 Orijtech, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tickeryzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `check that (*time.Ticker).Stop is called

The (*time.Ticker).Stop method must be called to release
associated resources hold by the Ticker.
`

var Analyzer = &analysis.Analyzer{
	Name: "tickeryzer",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		ctrlflow.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Fast path: if the package doesn't import time,
	// skip the traversal.
	if !imports(pass.Pkg, "time") {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}
		call := n.(*ast.CallExpr)
		if !isTimeNewTicker(pass.TypesInfo, call) {
			return true // the function call is not related to this check.
		}

		// Find the innermost containing block, and get the list
		// of statements starting with the one containing call.
		stmts := restOfBlock(stack)
		if len(stmts) == 0 {
			return true
		}

		asg, ok := stmts[0].(*ast.AssignStmt)
		if !ok {
			return true // the first statement is not assignment.
		}
		ticker := rootIdent(asg.Lhs[0])
		if ticker == nil {
			return true // could not find the *time.Ticker in the assignment.
		}

		for _, stmt := range stmts[1:] {
			switch stmt := stmt.(type) {
			case *ast.DeferStmt:
				if root := rootIdent(stmt.Call.Fun); root.Obj == ticker.Obj {
					return true
				}
			case *ast.ExprStmt:
				if root := rootIdent(stmt.X); root.Obj == ticker.Obj {
					return true
				}
			}
		}

		pass.ReportRangef(ticker, "missing %s.Stop() call", ticker.Name)
		return true
	})
	return nil, nil
}

func isTimeNewTicker(info *types.Info, expr *ast.CallExpr) bool {
	fun, _ := expr.Fun.(*ast.SelectorExpr)
	sig, _ := info.Types[fun].Type.(*types.Signature)
	if sig == nil {
		return false
	}

	res := sig.Results()
	if res.Len() != 1 {
		return false // the function called does not return one value.
	}
	if ptr, ok := res.At(0).Type().(*types.Pointer); !ok || !isNamedType(ptr.Elem(), "time", "Ticker") {
		return false // the first return type is not *time.Ticker.
	}
	return true
}

// restOfBlock, given a traversal stack, finds the innermost containing
// block and returns the suffix of its statements starting with the
// current node (the last element of stack).
func restOfBlock(stack []ast.Node) []ast.Stmt {
	for i := len(stack) - 1; i >= 0; i-- {
		if b, ok := stack[i].(*ast.BlockStmt); ok {
			for j, v := range b.List {
				if v == stack[i+1] {
					return b.List[j:]
				}
			}
			break
		}
	}
	return nil
}

// rootIdent finds the root identifier x in a chain of selections x.y.z, or nil if not found.
func rootIdent(n ast.Node) *ast.Ident {
	switch n := n.(type) {
	case *ast.SelectorExpr:
		return rootIdent(n.X)
	case *ast.Ident:
		return n
	default:
		return nil
	}
}

// isNamedType reports whether t is the named type path.name.
func isNamedType(t types.Type, path, name string) bool {
	n, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := n.Obj()
	return obj.Name() == name && obj.Pkg() != nil && obj.Pkg().Path() == path
}

// Imports returns true if path is imported by pkg.
func imports(pkg *types.Package, path string) bool {
	for _, imp := range pkg.Imports() {
		if imp.Path() == path {
			return true
		}
	}
	return false
}
