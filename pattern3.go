package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
)

/*
written by Jang, Hyang Gi on 28.August.2023

Desired codes:

from function patternArraySet
	func (a *main.array[T]) SetGen(i int, x T) {
		a.SetG(i, x)
	}

from function patternArrayGet
	func (a *main.array[T]) SetGen(i int, x T) {
			a.SetG(i, x)
	}

from function patternMemset
	func (a *main.array[T]) SetGen(i int, x T) {
		a.SetG(i, x)
	}

from function patternMemsetIter
	func (a *main.array[T]) SetGen(i int, x T) {
		a.SetG(i, x)
	}

from function patternEq
it will be printed out No parameterised methods
*/

func patternArraySet(fnc basicStr) *ast.FuncDecl {
	flag0, flag1, flag2, flag3, flag4, flag5, flag6 := false, false, false, false, false, false, false
	receiverParam := []string{}
	paramList := []string{}
	var interfaceVar string
	var caseVar string
	var typeVar string
	var resultNode *ast.FuncDecl
	beforeBlockStm := false

	for ind, elems := range fnc.value {
		if strings.Contains(elems.path, "*ast.BlockStmt") {
			beforeBlockStm = true
		}
		// Check if it has a receiver
		if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") {
			if !flag0 {
				receiverParam = elems.value
				if len(receiverParam) == 1 {
					if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.StarExpr") {
						for _, val := range fnc.value[ind+1].value {
							receiverParam = append(receiverParam, val)
						}
					}
				}
				flag0 = true
				continue
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

		// Check if there exists a parameter which type is Interface{}
		if flag1 {
			if !beforeBlockStm && strings.Contains(elems.path, "*ast.Field") {
				if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.InterfaceType") && len(elems.value) > 0 {
					if !flag2 {
						interfaceVar = elems.value[0]
						flag2 = true
						continue
					}
				}

			}
		}

		// Switch statement begin and check whether the switch statement handles a type
		if flag2 {
			if strings.Contains(elems.path, "*ast.SwitchStmt -> *ast.CallExpr -> *ast.SelectorExpr") &&
				contains(elems.value, receiverParam[0]) && contains(elems.value, "Kind") {
				if !flag3 {
					flag3 = true
					continue
				}
			}
		}
		// Check Case Clause
		if flag3 {
			if strings.Contains(elems.path, "*ast.CaseClause -> *ast.SelectorExpr") &&
				len(elems.value) > 1 && contains(elems.value, "reflect") {
				if !flag4 {
					caseVar = elems.value[1]
					flag4 = true
					continue
				}
			}
		}
		// Check a variable for setting
		if flag4 {
			if strings.Contains(elems.path, "*ast.TypeAssertExpr") &&
				(contains(elems.value, interfaceVar) && contains(elems.value, caseVar)) {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") {
					if len(fnc.value[ind-1].value) > 1 {
						if !flag5 {
							typeVar = fnc.value[ind-1].value[1]
							flag5 = true
							continue
						}
					}
				}
			}
		}

		// Check Set function
		if flag5 {
			if strings.Contains(elems.path, "*ast.ExprStmt -> *ast.CallExpr -> *ast.SelectorExpr") &&
				contains(elems.value, receiverParam[0]) && istValueInArr(elems.value, fnc.funcName) &&
				contains(elems.value, paramList[0]) && contains(elems.value, typeVar) {
				if !flag6 {
					flag6 = true
					break
				}
			}
		}
	}

	/*
		Identifier names in the examples below may differ from actual code identifier names! This node looks like
		func (a *main.array[T]) SetGen(i int, x T) {
			a.SetG(i, x)
		}
	*/
	if flag6 && len(receiverParam) > 2 && len(paramList) > 1 {
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
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: receiverParam[1],
									},
									Sel: &ast.Ident{
										Name: receiverParam[2],
									},
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
				Name: fnc.funcName,
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
									Name: interfaceVar,
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
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: receiverParam[0],
								},
								Sel: &ast.Ident{
									Name: fnc.funcName,
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: paramList[0],
								},
								&ast.Ident{
									Name: interfaceVar,
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

func patternArrayGet(fnc basicStr) *ast.FuncDecl {
	flag0, flag1, flag2, flag3, flag4, flag5 := false, false, false, false, false, false
	receiverParam := []string{}
	paramList := []string{}

	var resultNode *ast.FuncDecl
	beforeBlockStm := false

	for ind, elems := range fnc.value {
		if strings.Contains(elems.path, "*ast.BlockStmt") {
			beforeBlockStm = true
		}
		// Check if it has a receiver
		if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") {
			if !flag0 {
				receiverParam = elems.value
				if len(receiverParam) == 1 {
					if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.StarExpr") {
						for _, val := range fnc.value[ind+1].value {
							receiverParam = append(receiverParam, val)
						}
					}
				}
				flag0 = true
				continue
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

		// Check if the return value is Interface{} type
		if flag1 {
			if !beforeBlockStm && strings.Contains(elems.path, "ast.FieldList -> *ast.Field -> *ast.InterfaceType") && len(elems.value) == 0 {
				if !flag2 {
					flag2 = true
					continue
				}
			}
		}

		// Switch statement begin and check whether the switch statement handles a type
		if flag2 {
			if strings.Contains(elems.path, "*ast.SwitchStmt -> *ast.CallExpr -> *ast.SelectorExpr") &&
				contains(elems.value, receiverParam[0]) && contains(elems.value, "Kind") {
				if !flag3 {
					flag3 = true
					continue
				}
			}
		}
		// Check Case Clause
		if flag3 {
			if strings.Contains(elems.path, "*ast.CaseClause -> *ast.SelectorExpr") &&
				len(elems.value) > 1 && contains(elems.value, "reflect") {
				if !flag4 {
					flag4 = true
					continue
				}
			}
		}

		// Check Set function
		if flag4 {
			if strings.Contains(elems.path, "*ast.ReturnStmt -> *ast.CallExpr -> *ast.SelectorExpr") &&
				contains(elems.value, receiverParam[0]) && istValueInArr(elems.value, fnc.funcName) &&
				contains(elems.value, paramList[0]) {
				if !flag5 {
					flag5 = true
					break
				}
			}
		}
	}

	/*
		Identifier names in the examples below may differ from actual code identifier names! This node looks like
		func (a *main.array[T]) SetGen(i int, x T) {
			a.SetG(i, x)
		}
	*/
	if flag5 && len(receiverParam) > 2 && len(paramList) > 1 {
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
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: receiverParam[1],
									},
									Sel: &ast.Ident{
										Name: receiverParam[2],
									},
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
				Name: fnc.funcName,
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
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: receiverParam[0],
									},
									Sel: &ast.Ident{
										Name: fnc.funcName,
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: paramList[0],
									},
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

func patternMemset(fnc basicStr) *ast.FuncDecl {
	flag0, flag1, flag2, flag3, flag4, flag5, flag6, flag7, flag8, flag9 := false, false, false, false, false, false, false, false, false, false
	receiverParam := []string{}
	var interfaceVar string
	var caseVar string
	var setVar string
	var arrayVar string
	var rangeVar string
	var returnType string
	var returnValue string
	var resultNode *ast.FuncDecl
	beforeBlockStm := true

	for ind, elems := range fnc.value {

		if strings.Contains(elems.path, "*ast.BlockStmt") {
			beforeBlockStm = false
		}
		// Check if it has a receiver
		if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") {
			if !flag0 {
				receiverParam = elems.value
				if len(receiverParam) == 1 {
					if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.StarExpr") {
						for _, val := range fnc.value[ind+1].value {
							receiverParam = append(receiverParam, val)
						}
					}
				}
				flag0 = true
				continue
			}

		}

		// Check if there exists Interface parameter
		if flag0 {
			if beforeBlockStm && strings.Contains(elems.path, "*ast.InterfaceType") && len(elems.value) == 0 {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.FieldList -> *ast.Field") &&
					len(fnc.value[ind-1].value) > 0 {
					if !flag1 {
						interfaceVar = fnc.value[ind-1].value[0]
						flag1 = true
						continue
					}
				}
			}
		}

		// Save the return value
		if flag1 {
			if beforeBlockStm && strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") && len(elems.value) > 0 {
				if !flag2 {
					returnType = elems.value[0]
					flag2 = true
					continue
				}
			}
		}

		// Switch statement begin and check whether the switch statement handles receiver parameter
		if flag2 {
			if strings.Contains(elems.path, "*ast.SwitchStmt -> *ast.SelectorExpr") &&
				contains(elems.value, receiverParam[0]) {
				if !flag3 {
					flag3 = true
					continue
				}
			}
		}
		// Check Case Clause
		if flag3 {
			if strings.Contains(elems.path, "*ast.CaseClause") && len(elems.value) == 1 {
				if !flag4 {
					caseVar = elems.value[0]
					flag4 = true
					continue
				}
			}
		}

		// If there is a interface variable in the astValue?
		if flag4 {
			if strings.Contains(elems.path, "*ast.TypeAssertExpr") && contains(elems.value, interfaceVar) && contains(elems.value, caseVar) {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") &&
					len(fnc.value[ind-1].value) > 1 {
					if !flag5 {
						setVar = fnc.value[ind-1].value[1]
						flag5 = true
						continue
					}
				}
			}
		}

		// If there is a receiver parameter in the astValue?
		if flag5 {
			if strings.Contains(elems.path, "*ast.CallExpr -> *ast.SelectorExpr") && contains(elems.value, receiverParam[0]) && istValueInArr(elems.value, caseVar) {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") &&
					len(fnc.value[ind-1].value) > 1 {
					if !flag6 {
						arrayVar = fnc.value[ind-1].value[1]
						flag6 = true
						continue
					}
				}
			}
		}

		// If there is a Range-Statement?
		if flag6 {
			if strings.Contains(elems.path, "*ast.RangeStmt") && contains(elems.value, arrayVar) && len(elems.value) > 1 {
				if !flag7 {
					rangeVar = elems.value[0]
					flag7 = true
					continue
				}
			}
		}

		// If there is Index Expression ?
		if flag7 {
			if strings.Contains(elems.path, "*ast.IndexExpr") && contains(elems.value, arrayVar) && contains(elems.value, rangeVar) && contains(elems.value, setVar) {
				if !flag8 {
					flag8 = true
					continue
				}
			}
		}

		// Check return value
		if flag8 {
			if strings.EqualFold(elems.path, "*ast.ReturnStmt") && len(elems.value) == 1 && strings.EqualFold(elems.value[0], "nil") {
				if !flag9 {
					returnValue = elems.value[0]
					flag9 = true
					break
				}
			}
		}
	}
	/*
		Identifier names in the examples below may differ from actual code identifier names! This node looks like
		func (a *main.array[T]) SetGen(i int, x T) {
			a.SetG(i, x)
		}
	*/
	if flag9 && len(receiverParam) > 2 {
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
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: receiverParam[1],
									},
									Sel: &ast.Ident{
										Name: receiverParam[2],
									},
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
				Name: fnc.funcName,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: interfaceVar,
								},
							},
							Type: &ast.Ident{
								Name: typeName,
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Type: &ast.Ident{
								Name: returnType,
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: arrayVar,
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
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
					&ast.RangeStmt{
						Key: &ast.Ident{
							Name: rangeVar,
						},
						Tok: token.DEFINE,
						X: &ast.Ident{
							Name: arrayVar,
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.AssignStmt{
									Lhs: []ast.Expr{
										&ast.IndexExpr{
											X: &ast.Ident{
												Name: arrayVar,
											},
											Index: &ast.Ident{
												Name: rangeVar,
											},
										},
									},
									Tok: token.ASSIGN,
									Rhs: []ast.Expr{
										&ast.Ident{
											Name: interfaceVar,
										},
									},
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: returnValue,
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

func patternMemsetIter(fnc basicStr) *ast.FuncDecl {
	flag0, flag1, flag2, flag3, flag4, flag5, flag6, flag7, flag8, flag9, flag10, flag11 := false, false, false, false, false, false, false, false, false, false, false, false
	receiverParam := []string{}
	var interfaceVar string
	var iteratorVar string
	var intVar string
	var caseVar string
	var setVar string
	var arrayVar string
	var returnValue string
	var errorFunc string
	var resultNode *ast.FuncDecl
	beforeBlockStm := true

	for ind, elems := range fnc.value {

		if strings.Contains(elems.path, "*ast.BlockStmt") {
			beforeBlockStm = false
		}
		// Check if it has a receiver
		if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") {
			if !flag0 {
				receiverParam = elems.value
				if len(receiverParam) == 1 {
					if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.StarExpr") {
						for _, val := range fnc.value[ind+1].value {
							receiverParam = append(receiverParam, val)
						}
					}
				}
				flag0 = true
				continue
			}

		}

		// Check if there exists Interface parameter
		if flag0 {
			if beforeBlockStm && strings.Contains(elems.path, "*ast.InterfaceType") && len(elems.value) == 0 {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.FieldList -> *ast.Field") &&
					len(fnc.value[ind-1].value) > 0 {
					if !flag1 {
						interfaceVar = fnc.value[ind-1].value[0]
						flag1 = true
						continue
					}
				}
			}
		}

		// Check if there exists Iterator parameter
		if flag1 {
			if beforeBlockStm && strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") &&
				len(elems.value) > 1 && contains(elems.value, "Iterator") {
				if !flag2 {
					iteratorVar = elems.value[0]
					flag2 = true
					continue
				}
			}
		}

		// Save the return value
		if flag2 {
			if beforeBlockStm && strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") &&
				len(elems.value) > 1 && contains(elems.value, "error") {
				if !flag3 {
					returnValue = elems.value[0]
					flag3 = true
					continue
				}
			}
		}

		// Save the int value for for-Loop
		if flag3 {
			if strings.Contains(elems.path, "ast.GenDecl -> *ast.ValueSpec") &&
				len(elems.value) > 1 && strings.EqualFold(elems.value[1], "int") {
				if !flag4 {
					intVar = elems.value[0]
					flag4 = true
					continue
				}
			}
		}

		// Switch statement begin and check whether the switch statement handles receiver parameter
		if flag4 {
			if strings.Contains(elems.path, "*ast.SwitchStmt -> *ast.SelectorExpr") &&
				contains(elems.value, receiverParam[0]) {
				if !flag5 {
					flag5 = true
					continue
				}
			}
		}

		// Check Case Clause
		if flag5 {
			if strings.Contains(elems.path, "*ast.CaseClause") && len(elems.value) == 1 {
				if !flag6 {
					caseVar = elems.value[0]
					flag6 = true
					continue
				}
			}
		}

		// If there is a interface variable in the astValue?
		if flag6 {
			if strings.Contains(elems.path, "*ast.TypeAssertExpr") && contains(elems.value, interfaceVar) && contains(elems.value, caseVar) {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") &&
					len(fnc.value[ind-1].value) > 1 {
					if !flag7 {
						setVar = fnc.value[ind-1].value[1]
						flag7 = true
						continue
					}
				}
			}
		}

		// If there is a receiver parameter in the astValue?
		if flag7 {
			if strings.Contains(elems.path, "*ast.CallExpr -> *ast.SelectorExpr") && contains(elems.value, receiverParam[0]) && istValueInArr(elems.value, caseVar) {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") &&
					len(fnc.value[ind-1].value) > 1 {
					if !flag8 {
						arrayVar = fnc.value[ind-1].value[1]
						flag8 = true
						continue
					}
				}
			}
		}

		// if there is a for-loop ?
		if flag8 {
			if strings.Contains(elems.path, "ForStmt -> *ast.AssignStmt") && contains(elems.value, returnValue) && contains(elems.value, intVar) {
				if ind+4 < len(fnc.value) && reflect.DeepEqual(fnc.value[ind+1].value, []string{iteratorVar, "Next"}) &&
					reflect.DeepEqual(fnc.value[ind+2].value, []string{"==", returnValue, "nil"}) &&
					reflect.DeepEqual(fnc.value[ind+3].value, []string{"=", intVar, returnValue}) &&
					reflect.DeepEqual(fnc.value[ind+4].value, []string{iteratorVar, "Next"}) {
					if !flag9 {
						flag9 = true
						continue
					}
				}

			}
		}

		// If there is Index Expression ?
		if flag9 {
			if strings.Contains(elems.path, "*ast.IndexExpr") && contains(elems.value, arrayVar) && contains(elems.value, intVar) && contains(elems.value, setVar) {
				if !flag10 {
					flag10 = true
					continue
				}
			}
		}

		// If there is a return value in Astvalue
		if flag10 {
			if strings.EqualFold(elems.path, "*ast.AssignStmt") && contains(elems.value, returnValue) {
				if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.CallExpr") && contains(fnc.value[ind+1].value, returnValue) &&
					len(fnc.value[ind+1].value) > 1 {
					if !flag11 {
						errorFunc = fnc.value[ind+1].value[0]
						flag11 = true
						break
					}
				}

			}
		}
	}
	/*
		Identifier names in the examples below may differ from actual code identifier names! This node looks like
		func (a *main.array[T]) SetGen(i int, x T) {
			a.SetG(i, x)
		}
	*/
	if flag11 && len(receiverParam) > 2 {
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
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: receiverParam[1],
									},
									Sel: &ast.Ident{
										Name: receiverParam[2],
									},
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
				Name: fnc.funcName,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: interfaceVar,
								},
							},
							Type: &ast.Ident{
								Name: typeName,
							},
						},
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: iteratorVar,
								},
							},
							Type: &ast.Ident{
								Name: "Iterator",
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: returnValue,
								},
							},
							Type: &ast.Ident{
								Name: "error",
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.DeclStmt{
						Decl: &ast.GenDecl{
							Tok: token.VAR,
							Specs: []ast.Spec{
								&ast.ValueSpec{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: intVar,
										},
									},
									Type: &ast.Ident{
										Name: "int",
									},
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: arrayVar,
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
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
					&ast.ForStmt{
						Init: &ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: intVar,
								},
								&ast.Ident{
									Name: returnValue,
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: iteratorVar,
										},
										Sel: &ast.Ident{
											Name: "Next",
										},
									},
								},
							},
						},
						Cond: &ast.BinaryExpr{
							X: &ast.Ident{
								Name: returnValue,
							},
							Op: token.EQL,
							Y: &ast.Ident{
								Name: "nil",
							},
						},
						Post: &ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: intVar,
								},
								&ast.Ident{
									Name: returnValue,
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: iteratorVar,
										},
										Sel: &ast.Ident{
											Name: "Next",
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
											X: &ast.Ident{
												Name: arrayVar,
											},
											Index: &ast.Ident{
												Name: intVar,
											},
										},
									},
									Tok: token.ASSIGN,
									Rhs: []ast.Expr{
										&ast.Ident{
											Name: interfaceVar,
										},
									},
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: returnValue,
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.Ident{
									Name: errorFunc,
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: returnValue,
									},
								},
							},
						},
					},
					&ast.ReturnStmt{},
				},
			},
		}

		return resultNode
	} else {
		return nil
	}
}

func patternEq(fnc basicStr) {
	flag0, flag1, flag2, flag3, flag4 := false, false, false, false, false
	receiverParam := []string{}
	var interfaceVar string
	var interfaceVar1 string
	var caseVar string
	var receiverRangeVar string
	beforeBlockStm := true
	noParametrisationMethod := false

	for ind, elems := range fnc.value {

		if strings.Contains(elems.path, "*ast.BlockStmt") {
			beforeBlockStm = false
		}
		// Check if it has a receiver
		if strings.Contains(elems.path, "*ast.FieldList -> *ast.Field") {
			if !flag0 {
				receiverParam = elems.value
				if len(receiverParam) == 1 {
					if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.StarExpr") {
						for _, val := range fnc.value[ind+1].value {
							receiverParam = append(receiverParam, val)
						}
					}
				}
				flag0 = true
				continue
			}

		}

		// Check if there exists Interface parameter
		if flag0 {
			if beforeBlockStm && strings.Contains(elems.path, "*ast.InterfaceType") && len(elems.value) == 0 {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.FieldList -> *ast.Field") &&
					len(fnc.value[ind-1].value) > 0 {
					if !flag1 {
						interfaceVar = fnc.value[ind-1].value[0]
						flag1 = true
						continue
					}
				}
			}
		}

		// If there is a interface variable in the astValue?
		if flag1 {
			if strings.Contains(elems.path, "*ast.TypeAssertExpr") && contains(elems.value, interfaceVar) {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.AssignStmt") &&
					len(fnc.value[ind-1].value) > 1 {
					if !flag2 {
						interfaceVar1 = fnc.value[ind-1].value[1]
						flag2 = true
						continue
					}
				}
			}
		}

		// Check Case Clause
		if flag2 {
			if strings.Contains(elems.path, "*ast.CaseClause -> *ast.SelectorExpr") && len(elems.value) > 1 &&
				contains(elems.value, "reflect") {
				if !flag3 {
					caseVar = elems.value[1]
					flag3 = true
					continue
				}
			}
		}

		// If there is a receiver parameter in the astValue and for-Loop?
		if flag3 {
			if strings.Contains(elems.path, "*ast.CallExpr -> *ast.SelectorExpr") && contains(elems.value, receiverParam[0]) && istValueInArr(elems.value, caseVar) {
				if ind-1 > 0 && strings.Contains(fnc.value[ind-1].path, "*ast.RangeStmt") &&
					len(fnc.value[ind-1].value) > 1 {
					if !flag4 {
						receiverRangeVar = fnc.value[ind-1].value[1]
						flag4 = true
						continue
					}
				}
			}
		}

		// If there is a compare statement
		if flag4 {
			if strings.Contains(elems.path, "*ast.BinaryExpr") &&
				len(elems.value) > 0 && (strings.EqualFold(elems.value[0], "==") || strings.EqualFold(elems.value[0], "!=")) {
				if ind+1 < len(fnc.value) && strings.Contains(fnc.value[ind+1].path, "*ast.CallExpr -> *ast.SelectorExpr") &&
					len(fnc.value[ind+1].value) > 1 {
					for i, value := range fnc.value[ind+1].value {
						if strings.EqualFold(value, interfaceVar1) {
							if i+1 < len(fnc.value[ind+1].value) && !strings.EqualFold(fnc.value[ind+1].value[i+1], receiverRangeVar) {
								noParametrisationMethod = true
								break
							}
						}
					}
				}
				if noParametrisationMethod {
					break
				}
			}
		}

	}
	if noParametrisationMethod {
		fmt.Println("This function ", fnc.funcName, " has reused cases in the switch statement, but it cannot be refactored by Generics.")
		fmt.Println("The reason is: No Parameterized Method")
		fmt.Println()
	}
}
