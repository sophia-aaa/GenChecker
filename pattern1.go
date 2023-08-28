package main

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

/*
written by Jang, Hyang Gi on 28.August.2023

Desired codes:

from the function PatternObjectSlice
	func (h *Header[T]) ObjectSlice() []T {
		return h.Raw
	}

from the function PatternSet
	func (h *Header[T]) Set(i int, x T) {
		h.Raw[i] = x
	}

from the function PatternGet
	func (h *Header[T]) Get(i int) T {
		return h.Raw[i]
	}
*/

func patternObjectSlice(fnc basicStr) *ast.FuncDecl {
	flag0, flag1, flag2 := false, false, false
	receiverParam := []string{}
	funcName := ""
	typeLists := []string{}
	fieldName := ""
	var resultNode *ast.FuncDecl

	// name your function here
	if 'A' <= fnc.funcName[0] && fnc.funcName[0] <= 'Z' {
		funcName = "ObjectSlice"
	} else if 'a' <= fnc.funcName[0] && fnc.funcName[0] <= 'z' {
		funcName = "objectSlice"
	} else {
		funcName = fnc.funcName
	}

	for ind, elems := range fnc.value {
		// Check if it has a receiver
		if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") {
			if !flag0 {
				receiverParam = elems.value
				if len(receiverParam) == 1 {
					if ind+1 < len(fnc.value) {
						if strings.Contains(fnc.value[ind+1].path, "*ast.StarExpr") {
							receiverParam = append(receiverParam, fnc.value[ind+1].value[0])
						}
					}
				}
				flag0 = true
				continue
			}
		}
		// Check whether the return type is []TYPE
		if flag0 {
			if strings.Contains(elems.path, "*ast.Field -> *ast.ArrayType") {
				typeLists = elems.value
				if !flag1 {
					flag1 = true
					continue
				}
			}
		}
		// Check whether the return value looks like return (*(*[]bool)(unsafe.Pointer(&...)) ...
		if flag1 {
			if strings.Contains(elems.path, "*ast.ReturnStmt -> *ast.SliceExpr -> *ast.ParenExpr -> *ast.StarExpr -> *ast.CallExpr -> *ast.ParenExpr -> *ast.StarExpr -> *ast.ArrayType") &&
				reflect.DeepEqual(elems.value, typeLists) {
				if ind+3 < len(fnc.value) {
					//  Check unsafe Pointer Usage
					if strings.Contains(fnc.value[ind+1].path, "*ast.SelectorExpr") && isSameString(fnc.value[ind+1].value, "unsafe") && isSameString(fnc.value[ind+1].value, "Pointer") {
						//  Check &
						if strings.Contains(fnc.value[ind+2].path, "*ast.UnaryExpr") && len(fnc.value[ind+2].value) == 1 && strings.EqualFold(fnc.value[ind+2].value[0], "&") {
							// Check if ampersand of unsafe.Pointer matches the receiver
							if strings.Contains(fnc.value[ind+3].path, "*ast.SelectorExpr") &&
								len(fnc.value[ind+3].value) == 2 && contains(receiverParam, fnc.value[ind+3].value[0]) {
								flag2 = true
								fieldName = fnc.value[ind+3].value[1]
								break
							}
						}
					}
				}
			}
		}
	}

	/*
		Identifier names in the examples below may differ from actual code identifier names! This node looks like
		func (g *GenHeader[T]) Object() []T {
			return g.List
		}
	*/
	if flag2 && len(receiverParam) > 1 {
		resultNode = &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: receiverParam[0],
							},
						},
						Type: &ast.StarExpr{
							X: &ast.IndexExpr{
								X: &ast.Ident{
									Name: receiverParam[1],
								},
								Index: &ast.Ident{
									Name: "T",
								},
							},
						},
					},
				},
			},
			Name: &ast.Ident{
				Name: funcName,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Type: &ast.ArrayType{
								Elt: &ast.Ident{
									Name: "T",
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
									Name: fieldName,
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

func patternSet(fnc basicStr) *ast.FuncDecl {
	flag0, flag1, flag2, flag3 := false, false, false, false
	receiverParam := []string{}
	paramList := []string{}
	retValue := []string{}
	funcName := ""

	var resultNode *ast.FuncDecl

	// name your function here
	if 'A' <= fnc.funcName[0] && fnc.funcName[0] <= 'Z' && strings.Contains(fnc.funcName, "Set") {
		funcName = "Set"
	} else if 'a' <= fnc.funcName[0] && fnc.funcName[0] <= 'z' && strings.Contains(fnc.funcName, "set") {
		funcName = "set"
	} else {
		funcName = fnc.funcName
	}

	for ind, elems := range fnc.value {
		// Check if it has a receiver
		if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") {
			if !flag0 {
				receiverParam = elems.value
				if len(receiverParam) == 1 {
					if ind+1 < len(fnc.value) {
						if strings.Contains(fnc.value[ind+1].path, "*ast.StarExpr") {
							receiverParam = append(receiverParam, fnc.value[ind+1].value[0])
							flag0 = true
							ind = ind + 1
							continue
						}
					}
				}
			}
		}
		// Check if there exists int parameter
		if flag0 {
			if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") && contains(elems.value, "int") {
				if !flag1 {
					paramList = elems.value
					flag1 = true
					continue
				}
			}
		}
		// Check the return value
		if flag1 {
			if strings.EqualFold(elems.path, "*ast.Field") &&
				(len(elems.value) == 2 || (contains(elems.value, "unsafe") && contains(elems.value, "Pointer"))) {
				if !flag2 {
					retValue = elems.value
					flag2 = true
					continue
				}
			}
		}
		// Check Assignment-Statement
		if flag2 {
			if strings.Contains(elems.path, "*ast.AssignStmt") && strings.EqualFold(elems.value[0], "=") {
				if ind+1 < len(fnc.value) {
					if strings.Contains(fnc.value[ind+1].path, "*ast.IndexExpr -> *ast.CallExpr -> *ast.SelectorExpr") {
						if contains(fnc.value[ind+1].value, receiverParam[0]) && contains(fnc.value[ind+1].value, paramList[0]) && contains(fnc.value[ind+1].value, retValue[0]) {
							if !flag3 {
								flag3 = true
								break
							}
						}
					}
				}
			}
		}
	}

	/*
		Identifier names in the examples below may differ from actual code identifier names! This node looks like
		func (g *GenHeader[T]) SetGI(i int, x T) {
			g.List[i] = x
		}
	*/
	if flag3 && len(receiverParam) > 1 {
		resultNode = &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: receiverParam[0],
							},
						},
						Type: &ast.StarExpr{
							X: &ast.IndexExpr{
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
			},
			Name: &ast.Ident{
				Name: funcName,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: paramList[0],
								},
							},
							Type: &ast.Ident{
								Name: paramList[1],
							},
						},
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: retValue[0],
								},
							},
							Type: &ast.Ident{
								Name: typeName,
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.IndexExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: receiverParam[0],
									},
									Sel: &ast.Ident{
										Name: sliceName,
									},
								},
								Index: &ast.Ident{
									Name: paramList[0],
								},
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.Ident{
								Name: retValue[0],
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

func patternGet(fnc basicStr) *ast.FuncDecl {
	flag0, flag1, flag2, flag3 := false, false, false, false
	receiverParam := []string{}
	paramList := []string{}
	funcName := ""

	var resultNode *ast.FuncDecl

	// name your function here
	if 'A' <= fnc.funcName[0] && fnc.funcName[0] <= 'Z' && strings.Contains(fnc.funcName, "Get") {
		funcName = "Get"
	} else if 'a' <= fnc.funcName[0] && fnc.funcName[0] <= 'z' && strings.Contains(fnc.funcName, "get") {
		funcName = "get"
	} else {
		funcName = fnc.funcName
	}

	for ind, elems := range fnc.value {
		// Check if it has a receiver
		if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") {
			if !flag0 {
				receiverParam = elems.value
				if len(receiverParam) == 1 {
					if ind+1 < len(fnc.value) {
						if strings.Contains(fnc.value[ind+1].path, "*ast.StarExpr") {
							receiverParam = append(receiverParam, fnc.value[ind+1].value[0])
							flag0 = true
							continue
						}
					}
				}
			}
		}
		// Check if there exists int parameter
		if flag0 {
			if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") && contains(elems.value, "int") {
				if !flag1 {
					paramList = elems.value
					flag1 = true
					continue
				}
			}
		}
		// Check the return value
		if flag1 {
			if strings.EqualFold(elems.path, "*ast.FieldList -> *ast.Field") &&
				(len(elems.value) == 1 || reflect.DeepEqual(elems.value, []string{"unsafe", "Pointer"})) {
				if !flag2 {
					flag2 = true
					continue
				}
			}
		}
		// Check Assignment-Statement
		if flag2 {
			if strings.Contains(elems.path, "*ast.ReturnStmt -> *ast.IndexExpr -> *ast.CallExpr -> *ast.SelectorExpr") &&
				contains(elems.value, receiverParam[0]) && contains(elems.value, paramList[0]) {
				if !flag3 {
					flag3 = true
					continue
				}
			}

		}
	}

	/*
		Identifier names in the examples below may differ from actual code identifier names! This node looks like
		func (g *GenHeader[T]) GetGI(i int) T {
			return g.List[i]
		}
	*/
	if flag3 && len(receiverParam) > 1 {
		resultNode = &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: receiverParam[0],
							},
						},
						Type: &ast.StarExpr{
							X: &ast.IndexExpr{
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
			},
			Name: &ast.Ident{
				Name: funcName,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: paramList[0],
								},
							},
							Type: &ast.Ident{
								Name: paramList[1],
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Type: &ast.Ident{
								Name: typeName,
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.IndexExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: receiverParam[0],
									},
									Sel: &ast.Ident{
										Name: sliceName,
									},
								},
								Index: &ast.Ident{
									Name: paramList[0],
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
