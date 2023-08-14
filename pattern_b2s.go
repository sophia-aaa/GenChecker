package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"reflect"
	"strings"
)

// Check for type conversions from a slice of byte to string using unsafe.Pointer
func patternB2S(fileName string, lists []basicStr, unsafeList []basicFunc) {
	flag0, flag1, flag2, flag3, flag4, flag5, flag6 := false, false, false, false, false, false, false
	byteName := ""
	newName := ""
	for _, val := range lists {
		if checkDuplicateInFunc(unsafeList, val.funcName, val.funcToken) {
			for _, elem := range val.value {
				// Checking if parameter contains []bytes
				if strings.Contains(elem.path, "*ast.ArrayType") && isSameString(elem.value, "byte") {
					flag0 = true
					continue
				}

				// Checking if the type of return value is string
				if flag0 {
					if strings.Contains(elem.path, "*ast.FieldList -> *ast.Field") && isSameString(elem.value, "string") {
						flag1 = true
						continue
					}
				}

				// Checking if the return value looks like "return *(*string)...."
				if flag1 {
					if strings.Contains(elem.path, "*ast.ReturnStmt -> *ast.StarExpr -> *ast.CallExpr -> *ast.ParenExpr -> *ast.StarExpr") &&
						isSameString(elem.value, "string") {
						flag2 = true
						continue
					}
				}

				// Detecting "unsafe.Pointer" usages
				if flag2 {
					if strings.Contains(elem.path, "*ast.SelectorExpr") &&
						isSameString(elem.value, "unsafe") && isSameString(elem.value, "Pointer") {
						flag3 = true
						continue
					}
				}

				// Detecting if "&" is used and saving the name of the []byte
				if flag3 {
					if strings.Contains(elem.path, "*ast.UnaryExpr") {
						for i, ident := range elem.value {
							if strings.EqualFold(ident, "&") {
								if i+1 < len(elem.value) {
									flag4 = true
									continue
								} else {
									break
								}
							}
							if flag4 {
								byteName = ident
							}
						}
					}
				}
			}

			// Replacement
			if flag4 {
				fset := token.NewFileSet()
				node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
				if err != nil {
					log.Fatal(err)
				}

				ast.Inspect(node, func(n ast.Node) bool {
					switch x := n.(type) {
					case *ast.FuncDecl:
						if strings.EqualFold(x.Name.String(), val.funcName) && x.Pos() == val.funcToken {
							flag5 = true
						}
					case *ast.ReturnStmt:
						if flag5 {
							x.Results = []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "string",
									},
									Args: []ast.Expr{
										&ast.SliceExpr{
											X: &ast.Ident{
												Name: byteName,
											},
											Slice3: false,
										},
									},
								},
							}
							flag6 = true
						}
					}
					return true
				})

				newName = fileName[0:len(fileName)-3] + "_replaced.go"

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

				flag0, flag1, flag2, flag3, flag4, flag5 = false, false, false, false, false, false
			}

		}
	}

	// remove import "unsafe" if there does not exist unsafe usage
	if flag6 {
		newAST := buildAstDataStr(newName)
		checkUnsafe := buildUnsafeList(newAST)
		if len(checkUnsafe) == 0 {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, newName, nil, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}

			ast.Inspect(node, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.GenDecl:
					fmt.Println("genDecl")
					for i := range x.Specs {
						fmt.Println(reflect.TypeOf(x.Specs[i]))
						fmt.Println(x.Specs[i].(*ast.ImportSpec).Path.Value)
						if strings.EqualFold(x.Specs[i].(*ast.ImportSpec).Path.Value, "\"unsafe\"") {
							fmt.Println("true")
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
			})

			newName = fileName[0:len(fileName)-3] + "_unsafe_removed.go"
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
		flag6 = false
	}
}
