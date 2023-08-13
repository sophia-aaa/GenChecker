package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"os"
	"reflect"
	"strings"
)

func checkDataPattern(listFunctions []basicStr) []string {
	dataCheck1 := false   // check reflect.SliceHeader
	var dataVar1 string   // save the variable name for reflect.SliceHeader
	var dataAssign string // a variable which assigns to dataVar1
	dataCheck2 := false   // check return statement has an Interface Type
	var dataVar2 string   // save the variable name for return value with Interface Type
	dataCheck3 := false   // check if there exists relation between dataVar1 and dataVar2

	var checkData []string
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			for idx, val := range listFunctions[s].value {
				if reflect.DeepEqual(val.value, []string{"reflect", "SliceHeader"}) {
					dataCheck1 = true
				}
				if dataCheck1 {
					for j := idx; j >= 0; j-- {
						if strings.Contains(listFunctions[s].value[j].path, "*ast.AssignStmt") {
							dataVar1 = listFunctions[s].value[j].value[0]
							dataCheck1 = false
							break
						}
					}
				}
				if strings.Contains(val.path, "*ast.UnaryExpr") && isSameString(val.value, dataVar1) {
					for j := idx; j >= 0; j-- {
						if strings.Contains(listFunctions[s].value[j].path, "*ast.AssignStmt") {
							dataAssign = listFunctions[s].value[j].value[0]
							dataCheck1 = false
							break
						}
					}
				}
				if strings.Contains(val.path, "*ast.ReturnStmt") && isSameString(val.value, "Interface") {
					dataCheck2 = true
					dataVar2 = val.value[0]
				}
				if dataCheck2 {
					for j := idx; j >= 0; j-- {
						if dataCheck3 && strings.Contains(listFunctions[s].value[j].path, "*ast.AssignStmt") && isSameString(listFunctions[s].value[j].value, dataVar2) {
							checkData = append(checkData, listFunctions[s].funcName)
							dataCheck3 = false
							dataCheck2 = false
							break
						}
						if strings.Contains(listFunctions[s].value[j].path, "*ast.CallExpr") && contains(listFunctions[s].value[j].value, dataAssign) {
							// this if statement steps into first, if the process that a dataVar1 assigns to another variable is going on

							dataCheck3 = true
						}
					}
				}
			}
			dataCheck1 = false
			dataCheck2 = false
			dataCheck3 = false
		}
	}
	return checkData
}

// example replacement
func replacePattern(filename string, pattern int) {
	switch pattern {
	case 2:
		fset := token.NewFileSet()
		filename := "test.go"
		f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		//spew.Dump(f)
		//ast.Print(fset, f)
		if err != nil {
			log.Fatal(err)
		}
		once := true
		n := astutil.Apply(f, func(c *astutil.Cursor) bool {
			//ast.Print(fset, c.Node().(*ast.File).Decls[0].(*ast.GenDecl).Specs[0].End())

			if _, ok := c.Node().(*ast.GenDecl); ok {
				if once {
					// Generics Getter
					c.InsertAfter(&ast.FuncDecl{
						Recv: &ast.FieldList{
							List: []*ast.Field{
								&ast.Field{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "g",
										},
									},
									Type: &ast.StarExpr{
										X: &ast.IndexExpr{
											X: &ast.Ident{
												Name: "GenHeader",
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
							Name: "GetGI",
						},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "i",
											},
										},
										Type: &ast.Ident{
											Name: "int",
										},
									},
								},
							},
							Results: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Type: &ast.Ident{
											Name: "T",
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
													Name: "g",
												},
												Sel: &ast.Ident{
													Name: "lists",
												},
											},
											Index: &ast.Ident{
												Name: "i",
											},
										},
									},
								},
							},
						},
					})

					// Generics struct
					c.InsertAfter(&ast.GenDecl{
						Tok: token.TYPE,
						Specs: []ast.Spec{
							&ast.TypeSpec{
								Name: &ast.Ident{
									Name: "GenHeader",
								},
								TypeParams: &ast.FieldList{
									List: []*ast.Field{
										&ast.Field{
											Names: []*ast.Ident{
												&ast.Ident{
													Name: "T",
												},
											},
											Type: &ast.Ident{
												Name: "any",
											},
										},
									},
								},
								Type: &ast.StructType{
									Fields: &ast.FieldList{
										List: []*ast.Field{
											&ast.Field{
												Names: []*ast.Ident{
													&ast.Ident{
														Name: "lists",
													},
												},
												Type: &ast.ArrayType{
													Elt: &ast.Ident{
														Name: "T",
													},
												},
											},
										},
									},
								},
							},
						},
					})
				}
				once = false
			}
			return true
		}, nil)

		newFile, err := os.Create(filename[0:len(filename)-3] + "_test.go")
		if err != nil {
			log.Fatal(err)
		}
		defer func(new *os.File) {
			err := new.Close()
			if err != nil {
			}
		}(newFile)

		if err := printer.Fprint(newFile, fset, n); err != nil {
			log.Fatal(err)
		}
	case 3:
	case 4:
	default:
		break
	}
}
