package main

import (
	"go/ast"
	"strings"
)

/*
written by Jang, Hyang Gi on 28.August.2023

Desired codes:

from function patternGenSlice
		func (g GenHeader[T]) Object() []T {
			return g.List
		}

*/

func patternGenSlice(fnc basicStr) *ast.FuncDecl {
	flag0, flag1, flag2, flag3, flag4, flag5, flag6 := false, false, false, false, false, false, false
	receiverParam := []string{}
	var sliceHeaderIdent string
	var sliceOfIdent string
	var indirectIdent string
	var unsafePointerIdent string
	var resultNode *ast.FuncDecl

	for ind, elems := range fnc.value {
		// Check if it has a receiver
		if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") {
			if !flag0 {
				receiverParam = elems.value
				if len(receiverParam) == 1 {
					if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.StarExpr") {
						receiverParam = append(receiverParam, fnc.value[ind+1].value[0])
					}
				}
				flag0 = true
				continue
			}

		}

		// Check if the return type is Interface{}
		if flag0 {
			if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field -> *ast.InterfaceType") && len(elems.value) == 0 {
				if !flag1 {
					flag1 = true
					continue
				}
			}
		}
		// Check reflect SliceHeader
		if flag1 {
			if strings.Contains(elems.path, "*ast.CompositeLit -> *ast.SelectorExpr") &&
				contains(elems.value, "reflect") && contains(elems.value, "SliceHeader") {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") {
					if len(fnc.value[ind-1].value) > 1 {
						sliceHeaderIdent = fnc.value[ind-1].value[1]
						if !flag2 {
							flag2 = true
							continue
						}
					}
				}
			}
		}
		// Check reflect SliceOf
		if flag2 {
			if strings.Contains(elems.path, "*ast.CallExpr -> *ast.SelectorExpr") &&
				contains(elems.value, "reflect") && contains(elems.value, "SliceOf") {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") {
					if len(fnc.value[ind-1].value) > 1 {
						sliceOfIdent = fnc.value[ind-1].value[1]
						if ind+1 < len(fnc.value) && contains(fnc.value[ind+1].value, receiverParam[0]) {
							if !flag3 {
								flag3 = true
								continue
							}
						}
					}
				}
			}
		}
		// Check Unsafe Pointer
		if flag3 {
			if strings.Contains(elems.path, "*ast.CallExpr -> *ast.SelectorExpr") &&
				(contains(elems.value, "unsafe") && contains(elems.value, "Pointer")) {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") {
					if len(fnc.value[ind-1].value) > 1 {
						unsafePointerIdent = fnc.value[ind-1].value[1]
						if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.UnaryExpr") &&
							strings.Contains(fnc.value[ind+1].value[1], sliceHeaderIdent) {
							if !flag4 {
								flag4 = true
								continue
							}
						}
					}
				}
			}
		}

		// Check Indirect Check
		if flag4 {
			if strings.Contains(elems.path, "*ast.CallExpr -> *ast.SelectorExpr") &&
				(contains(elems.value, "reflect") && contains(elems.value, "Indirect")) {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") &&
					len(fnc.value[ind-1].value) > 1 {
					indirectIdent = fnc.value[ind-1].value[1]
					if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.CallExpr -> *ast.SelectorExpr") &&
						contains(fnc.value[ind+1].value, "reflect") && contains(fnc.value[ind+1].value, "NewAt") &&
						contains(fnc.value[ind+1].value, sliceOfIdent) && contains(fnc.value[ind+1].value, unsafePointerIdent) {
						if !flag5 {
							flag5 = true
							continue
						}
					}
				}
			}
		}

		// Check return value
		if flag5 {
			if strings.Contains(elems.path, "*ast.ReturnStmt -> *ast.CallExpr -> *ast.SelectorExpr") &&
				contains(elems.value, indirectIdent) && contains(elems.value, "Interface") {
				if !flag6 {
					flag6 = true
					continue
				}
			}
		}
	}

	/*
		Identifier names in the examples below may differ from actual code identifier names! This node looks like
		func (g GenHeader[T]) Object() []T {
			return g.List
		}
	*/
	if flag6 && len(receiverParam) > 1 {
		resultNode = &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: receiverParam[0],
							},
						},
						Type: &ast.IndexExpr{
							X: &ast.Ident{
								Name: receiverParam[1],
							},
							Index: &ast.Ident{
								Name: typeName,
							},
						},
					},
				},
			},
			Name: &ast.Ident{
				Name: fnc.funcName,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Type: &ast.ArrayType{
								Elt: &ast.Ident{
									Name: typeName,
								},
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: receiverParam[0],
								},
								Sel: &ast.Ident{
									Name: sliceName,
								},
							},
						},
					},
				},
			},
		}

		return resultNode
	} else {
		return nil
	}
}
