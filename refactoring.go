package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"os"
	"strings"
)

func refactoring(filename string) {
	var patternObjectNode *ast.FuncDecl
	var patternSetNode *ast.FuncDecl
	var patternGetNode *ast.FuncDecl
	var patternArraySetNode *ast.FuncDecl
	var patternArrayGetNode *ast.FuncDecl
	var patternMemsetNode *ast.FuncDecl
	var patternMemsetIterNode *ast.FuncDecl

	var toReplace []pattern3Result
	listFunctions := buildAstDataStr(filename)

	var funcList []string
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			funcList = append(funcList, listFunctions[s].funcName)
		}
	}
	fmt.Println(len(funcList), " Function list: ")
	for ind, val := range funcList {
		if len(funcList) == 1 {
			fmt.Print("{ ", val, " }")
		} else if ind == len(funcList)-1 {
			fmt.Println(val, "}")
		} else if ind == 0 {
			fmt.Print("{ ", val, ", ")
		} else {
			fmt.Print(val, ", ")
		}
	}
	fmt.Println()
	fmt.Println()

	var methodList []string
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			if len(listFunctions[s].value) > 0 &&
				strings.Contains(listFunctions[s].value[0].path, "*ast.FieldList -> *ast.Field") {
				methodList = append(methodList, listFunctions[s].funcName)
			}
		}
	}

	fmt.Println(len(methodList), " Method list: ")
	for ind, val := range methodList {
		if len(methodList) == 1 {
			fmt.Print("{ ", val, " }")
		} else if ind == len(methodList)-1 {
			fmt.Println(val, "}")
		} else if ind == 0 {
			fmt.Print("{ ", val, ", ")
		} else {
			fmt.Print(val, ", ")
		}
	}
	fmt.Println()
	fmt.Println()

	// check leaf of SelectorExpr and unsafe Pointer
	modListFunctions := checkSelectorExpr(listFunctions)
	createTextFile(filename, modListFunctions)
	unsafeList := buildUnsafePointerList(modListFunctions)

	if len(unsafeList) > 0 {
		fmt.Println(len(unsafeList), " This function contains unsafe.Pointer:")
		for s := range unsafeList {
			if unsafeList[s].funcName != "" {
				if len(unsafeList) == 1 {
					fmt.Print("{ ", unsafeList[s].funcName, " }")
				} else if s == len(unsafeList)-1 {
					fmt.Println(unsafeList[s].funcName, "}")
				} else if s == 0 {
					fmt.Print("{ ", unsafeList[s].funcName, ", ")
				} else {
					fmt.Print(unsafeList[s].funcName, ", ")
				}
			}
		}
		fmt.Println()
		fmt.Println()
	}

	//patternB2S(filename, modListFunctions, unsafeList)
	// Check Generic Replacement
	genCheck := checkGenerics(modListFunctions, funcList, typeList)
	for s := range genCheck {
		if len(genCheck[s]) > 1 {
			//pattern1 = true
			fmt.Println()
			fmt.Print(len(genCheck[s]), " These functions have a same structure and the code are reused:\n")
			for ind, val := range genCheck[s] {
				if len(genCheck[s]) == 1 {
					fmt.Print(" { ", val.funcName, " & ", val.funcPos, " }")
				} else if ind == len(genCheck[s])-1 {
					fmt.Println(val.funcName, " & ", val.funcPos, "}")
				} else if ind == 0 {
					fmt.Print(" { ", val.funcName, " & ", val.funcPos, ", ")
				} else {
					fmt.Print(val.funcName, " & ", val.funcPos, ", ")
				}
			}
			fmt.Println()
		}
	}

	for _, fnc := range genCheck {
		if len(fnc) > 1 {
			for _, val := range modListFunctions {
				if strings.EqualFold(val.funcName, fnc[0].funcName) && val.funcToken == fnc[0].funcPos {
					patternObjectNode = patternObjectSlice(val)
					if patternObjectNode != nil {
						if !checkReplaceFunc(toReplace, fnc) {
							toReplace = append(toReplace, pattern3Result{patternObjectNode, fnc})
						}
						patternObjectNode = nil
					}
					patternSetNode = patternSet(val)
					if patternSetNode != nil {
						if !checkReplaceFunc(toReplace, fnc) {
							toReplace = append(toReplace, pattern3Result{patternSetNode, fnc})
						}
						patternSetNode = nil
					}
					patternGetNode = patternGet(val)
					if patternGetNode != nil {
						if !checkReplaceFunc(toReplace, fnc) {
							toReplace = append(toReplace, pattern3Result{patternGetNode, fnc})
						}
						patternGetNode = nil
					}
				}
			}
		}
	}

	// This variable is for checking switch statement
	existsSwitch, caseList := checkSwitchStatement(filename, modListFunctions)
	if existsSwitch {
		if len(caseList) > 0 {
			for _, val := range modListFunctions {
				for _, cases := range caseList {
					if strings.EqualFold(val.funcName, cases.funcName) {
						patternArraySetNode = patternArraySet(val)
						funcInfoList := []funcNamePos{funcNamePos{val.funcName, val.funcToken}}
						if patternArraySetNode != nil {
							fmt.Println("This function ", val.funcName, " has a Set pattern.")
							if !checkReplaceFunc(toReplace, funcInfoList) {
								toReplace = append(toReplace, pattern3Result{patternArraySetNode, funcInfoList})
							}
						}
						patternArrayGetNode = patternArrayGet(val)
						if patternArrayGetNode != nil {
							fmt.Println("This function ", val.funcName, " has a Get pattern.")

							if !checkReplaceFunc(toReplace, funcInfoList) {
								toReplace = append(toReplace, pattern3Result{patternArrayGetNode, funcInfoList})
							}
						}
						patternMemsetNode = patternMemset(val)
						if patternMemsetNode != nil {
							fmt.Println("This function ", val.funcName, " has a Memset pattern.")
							if !checkReplaceFunc(toReplace, funcInfoList) {
								toReplace = append(toReplace, pattern3Result{patternMemsetNode, funcInfoList})
							}
						}
						patternMemsetIterNode = patternMemsetIter(val)
						if patternMemsetIterNode != nil {
							fmt.Println("This function ", val.funcName, " has a memsetIter pattern.")
							if !checkReplaceFunc(toReplace, funcInfoList) {
								toReplace = append(toReplace, pattern3Result{patternMemsetIterNode, funcInfoList})
							}
						}
						patternEq(val)
						patternReduce(val)
					}
				}
			}
			fmt.Println("\nThis function has a switch statement with reused cases: ")
			for ind, val := range caseList {
				if len(caseList) == 1 {
					fmt.Print("{ ", val.funcName, " }")
				} else if ind == len(caseList)-1 {
					fmt.Println(val.funcName, "}")
				} else if ind == 0 {
					fmt.Print("{ ", val.funcName, ", ")
				} else {
					fmt.Print(val.funcName, ", ")
				}
			}
		}
	}

	var nodeList []*ast.FuncDecl
	var removeNameList []funcNamePos
	var toResult patternReplace
	for _, val := range toReplace {
		nodeList = append(nodeList, val.nodes)
		for _, value := range val.funcList {
			removeNameList = append(removeNameList, value)
		}

	}
	toReplace = []pattern3Result{}
	toResult = patternReplace{nodeList, removeNameList}
	if len(toResult.nodes) > 0 {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		//count := 0
		//var tmp *ast.Decl
		count := 0

		astutil.Apply(node, func(c *astutil.Cursor) bool {
			n := c.Node()

			if d, ok := n.(*ast.FuncDecl); ok {
				if checkDuplicateInFuncGen(toResult.funcRemove, d.Name.String(), d.Pos()) {
					if count == 0 {
						for _, val := range toResult.nodes {
							c.InsertBefore(val)
						}
					}
					c.Delete()
					count++
				}
			}
			/*if _, ok := n.(*ast.GenDecl); ok {
				c.InsertBefore(&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "GenHeader",
								Obj: &ast.Object{
									Kind: ast.Typ,
									Name: "GenHeader",
								},
							},

							Type: &ast.ArrayType{
								Len: &ast.Ident{
									Name: "T",
								},
							},
						},
					},
				},
				)
			}*/
			return true
		}, nil)
		/*
			panic: runtime error: index out of range [19] with length 18
			length cannot be changed! except usage unsafe pointer. But this tool avoid using unsafe.Pointer, so that I should accept the limit of ast Package and astutil
				ast.Inspect(node, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.File:
					for i := range x.Decls {
						for _, val := range toResult.nodes {
							if _, ok := x.Decls[i].(*ast.FuncDecl); ok && strings.EqualFold(x.Decls[i].(*ast.FuncDecl).Name.String(), val.Name.String()) {
								x.Decls[i] = val
								//x.Specs[i].(*ast.ImportSpec).Path = &ast.BasicLit{}
							}
						}
						for _, val := range toResult.funcRemove {
							if _, ok := x.Decls[i].(*ast.FuncDecl); ok && strings.EqualFold(x.Decls[i].(*ast.FuncDecl).Name.String(), val.funcName) {
								if len(x.Decls) == 1 {
									x.Decls = []ast.Decl{}
								} else {
									length := len(x.Decls)
									j := i + 1
									tmp := x.Decls
									x.Decls = x.Decls[:i]
									for ; j < length; j++ {
										x.Decls = append(x.Decls, tmp[j])
									}

								}
								//x.Specs[i].(*ast.ImportSpec).Path = &ast.BasicLit{}
							}

						}

					}

				}
				return true
			})

			ast.Inspect(node, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.GenDecl:
					for i := range x.Specs {
						if _, ok := x.Specs[i].(*ast.ImportSpec); ok && strings.EqualFold(x.Specs[i].(*ast.ImportSpec).Name.String(), "\"unsafe\"") {
							if len(x.Specs) == 1 {
								x.Specs = []ast.Spec{}
							} else {
								x.Specs = append(x.Specs[:i], x.Specs[i+1:]...)
							}
							//x.Specs[i].(*ast.ImportSpec).Path = &ast.BasicLit{}
						}
					}

				}
				return true
			})*/

		newName := filename[0:len(filename)-3] + "_replaced.go"

		newFile, err := os.Create(newName)
		if err != nil {
			log.Fatal(err)
		}
		defer func(new *os.File) {
			err := new.Close()
			if err != nil {
			}
		}(newFile)

		if err := printer.Fprint(newFile, fset, node); err != nil {
			log.Fatal(err)
		}
	}

	for _, val := range modListFunctions {
		if patternGenSlice(val) != nil {
			fmt.Println()
			fmt.Println("This function ", val.funcName, " has (a) function(s) with reflect.SliceHeader and Interface of return value.")
			fmt.Println("It recommends to use Generics Slice.")
			fmt.Println("No replacement because of Generics Replacement Suggestion.")
			fmt.Println()
			fmt.Println()
		}
	}
}
