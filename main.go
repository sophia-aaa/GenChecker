package main

import (
	"container/list"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	_ "golang.org/x/exp/slices"
	"log"
	"os"
	"reflect"
	"strings"
)

type elem struct {
	path  string
	value []string
}

type function1 struct {
	funcName string
	value    *list.List
}

// Basic structure for every function in the input file
type function2 struct {
	funcName string
	value    []elem
}

type funcNameList struct {
	lists []string
}

type casesCheck struct {
	funcName string
	cases    []function2
}

func resetValues() funcNameList {
	return funcNameList{lists: []string{}}

}

func contains(strArr []string, str string) bool {
	for _, val := range strArr {
		if val == str {
			return true
		}
	}
	return false
}

func contains2D(strArr [][]string, str string) bool {
	for i := 0; i < len(strArr); i++ {
		for j := 0; j < len(strArr[i]); j++ {
			if strArr[i][j] == str {
				return true
			}
		}
	}
	return false
}

func checkUnsafeUsages(str string) bool {
	return contains([]string{"unsafe", "Pointer"}, str)
}
func main() {
	// filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run gen.go - test.go
	filename := os.Args[2]
	//fmt.Println(filename)
	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	typeList := []string{
		"bool", "bType", "int", "iType", "int8", "i8Type", "int16", "i16Type", "int32", "i32Type", "int64", "i64Type", "uint", "uType",
		"uint8", "u8Type", "uint16", "u16Type", "uint32", "u32Type", "uint64", "u64Type", "uintptr", "uintptrType",
		"float32", "f32Type", "float64", "f64Type", "complex64", "c64Type", "complex128", "c128Type", "string", "strType",
		"unsafe", "Pointer", "unsafePointerType",
	}
	/*test := values{lists: []string{"test1", "test2"}}
	fmt.Println(test)
	test = resetValues()
	fmt.Println("is reset?", test)*/

	//ast.Print(fset, astTree)

	node_list := list.New()
	var elem_list []elem
	nameFunction := ""
	//nameCheck := 0
	astNode := ""
	var listFunctions2 []function2
	var astValue []string

	ast.Inspect(astTree, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.FuncDecl:
			//	fmt.Println(x.Name)
			//	fmt.Println(x, "\t\t", reflect.TypeOf(x).String())
			//	fmt.Println(fset.Position(x.Pos()), fset.Position(x.End()))
			if len(astValue) != 0 { // if a new root meets
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				func2 := function2{nameFunction, elem_list}
				listFunctions2 = append(listFunctions2, func2)
				elem_list = []elem{}
				astNode = ""
				astValue = []string{}

				//nameCheck = 0
			}
			nameFunction = x.Name.String()
			node_list.PushBack(elem{reflect.TypeOf(x).String(), []string{nameFunction}})

		case *ast.Ident:
			//fmt.Println(fset.Position(x.Pos()), reflect.TypeOf(x).String(), "\t", x.Name)
			astValue = append(astValue, x.Name)

		case *ast.GenDecl, *ast.ImportSpec, *ast.BasicLit, *ast.CommentGroup, *ast.Comment:
			fmt.Print(" ")
		case *ast.CaseClause, *ast.SwitchStmt:

			// *ast.CaseClause -> *ast.SelectorExpr
			elem_list = append(elem_list, elem{astNode, astValue})
			node_list.PushBack(elem{astNode, astValue})
			astNode = ""
			astValue = []string{}

			if astNode == "" {
				astNode += reflect.TypeOf(x).String()
			} else {
				astNode += " -> " + reflect.TypeOf(x).String()
			}

		default:
			if x != nil {
				//fmt.Println(fset.Position(x.Pos()), reflect.TypeOf(x).String())
				//			fmt.Println(fset.Position(x.Pos()))
				if len(astValue) != 0 {
					elem_list = append(elem_list, elem{astNode, astValue})
					node_list.PushBack(elem{astNode, astValue})
					astNode = ""
					astValue = []string{}
				}
				if astNode == "" {
					astNode += reflect.TypeOf(x).String()
				} else {
					astNode += " -> " + reflect.TypeOf(x).String()
				}
			}

			/*	if x.Name == nameFunction {
						nameCheck++
					}
					if x.Name != nameFunction && nameCheck != 0 {
						astValue = append(astValue, x.Name)
					}
				}

					if node_list != nil {
						fmt.Println(nameFunction, " added! node_list")
						listFunctions1 = append(listFunctions1, function1{nameFunction, node_list})
						node_list.Init()
					}
					if elem_list != nil {
						fmt.Println(nameFunction, " added! elem_list")
						listFunctions2 = append(listFunctions2, function2{nameFunction, elem_list})
						elem_list = []elem{}
					}*/
		}
		return true
	})

	if len(astValue) != 0 { // if a new root meets
		elem_list = append(elem_list, elem{astNode, astValue})
		node_list.PushBack(elem{astNode, astValue})
		func2 := function2{nameFunction, elem_list}
		listFunctions2 = append(listFunctions2, func2)
		elem_list = []elem{}
		astNode = ""
		astValue = []string{}
	}
	/*for item := node_list.Front(); item != nil; item = item.Next() {
		fmt.Println(item.Value)
	}*/

	var funcList []string
	for s := range listFunctions2 {
		if listFunctions2[s].funcName != "" {
			funcList = append(funcList, listFunctions2[s].funcName)
		}
	}

	// compare functions
	var genCheck [][]string
	var genFunc []string
	flag := false
	for i := 0; i < len(listFunctions2); i++ {
		if listFunctions2[i].funcName != "" {
			if len(genCheck) != 0 {
				if !contains2D(genCheck, listFunctions2[i].funcName) {
					genFunc = []string{listFunctions2[i].funcName}
				} else {
					continue
				}
			} else {
				genFunc = []string{listFunctions2[i].funcName}
			}
			for j := i + 1; j < len(listFunctions2); j++ {
				if listFunctions2[j].funcName != "" {
					if len(listFunctions2[i].value) == len(listFunctions2[j].value) {
						// Compare details between functions
						for idx := range listFunctions2[i].value {
							if strings.Compare(listFunctions2[i].value[idx].path, listFunctions2[j].value[idx].path) == 0 {
								if len(listFunctions2[i].value[idx].value) == len(listFunctions2[j].value[idx].value) {
									for idxValue := range listFunctions2[i].value[idx].value {
										if strings.Compare(listFunctions2[i].value[idx].value[idxValue], listFunctions2[j].value[idx].value[idxValue]) == 0 {
											flag = true
										} else {
											if (contains(funcList, listFunctions2[i].value[idx].value[idxValue]) && contains(funcList, listFunctions2[j].value[idx].value[idxValue])) ||
												(contains(typeList, listFunctions2[i].value[idx].value[idxValue]) && contains(typeList, listFunctions2[j].value[idx].value[idxValue])) {
												flag = true
											} else {
												flag = false
												break
											}
										}
									}
								} else {
									flag = false
									break
								}
							} else if strings.Contains(listFunctions2[i].value[idx].path, "*ast.SelectorExpr") {
								// listFunctions2[i].value[idx] looks like "... -> *ast.SelectorExpr [unsafe Pointer]"
								// and listFunctions2[j].value[idx] looks like " ... [TYPE]"
								if !contains(typeList, listFunctions2[j].value[idx].value[0]) {
									flag = false
									break
								}
								for _, val := range listFunctions2[i].value[idx].value {
									if !checkUnsafeUsages(val) {
										flag = false
										break
									}
								}
								flag = true
							} else if strings.Contains(listFunctions2[j].value[idx].path, "*ast.SelectorExpr") {
								// listFunctions2[i].value[idx] looks like " ... [TYPE]"
								// and listFunctions2[j].value[idx] looks like "... -> *ast.SelectorExpr [unsafe Pointer]"
								if !contains(typeList, listFunctions2[i].value[idx].value[0]) {
									flag = false
									break
								}
								for _, val := range listFunctions2[j].value[idx].value {
									if !checkUnsafeUsages(val) {
										flag = false
										break
									}
								}
								flag = true
							} else {
								flag = false
								break
							}
						}
					} else {
						flag = false
						continue // continue the progress comparing a next function to the compared one

						/*
							// TODO: SetUnsafePointer
							// how to handle the *ast.SelectorExpr ?
							else if len(listFunctions2[i].value) > len(listFunctions2[j].value) {
								longLen := len(listFunctions2[i].value)
								puffer := 0
								for k := 0; k < len(listFunctions2[i].value); k++ {
									if listFunctions2[i].value[k].path == "*ast.SelectorExpr" {

									}
								}
							} else if len(listFunctions2[i].value) < len(listFunctions2[j].value) {

							}*/
					}
				}
				if flag == true {
					genFunc = append(genFunc, listFunctions2[j].funcName)
					flag = false
				}
			}
			genCheck = append(genCheck, genFunc)
		}
	}

	// build data structure for comparing switch statement
	var switchCheck []string
	var caseCheck []function2
	var caseName string
	var caseList []elem
	var caseWFunc []casesCheck
	caseBool := false
	for i := range listFunctions2 {
		for j := range listFunctions2[i].value {
			if strings.Contains(listFunctions2[i].value[j].path, "*ast.SwitchStmt") {
				if !contains(switchCheck, listFunctions2[i].funcName) {
					switchCheck = append(switchCheck, listFunctions2[i].funcName)
				}
			}
			if strings.Contains(listFunctions2[i].value[j].path, "*ast.CaseClause -> *ast.SelectorExpr") {
				if len(caseList) > 0 {
					caseCheck = append(caseCheck, function2{caseName, caseList})
					caseList = []elem{}
				}
				caseName = listFunctions2[i].value[j].value[1]
				caseList = append(caseList, elem{listFunctions2[i].value[j].path, listFunctions2[i].value[j].value})
				caseBool = true
				continue
			} else if strings.Contains(listFunctions2[i].value[j].path, "*ast.CaseClause") {
				if len(caseList) > 0 {
					caseBool = false
					caseCheck = append(caseCheck, function2{caseName, caseList})
					caseList = []elem{}
				}
				continue
			}
			if caseBool {
				caseList = append(caseList, elem{listFunctions2[i].value[j].path, listFunctions2[i].value[j].value})
			}
		}
		if len(caseCheck) > 0 {
			caseWFunc = append(caseWFunc, casesCheck{listFunctions2[i].funcName, caseCheck})
		}
	}
	if len(switchCheck) > 0 {
		fmt.Print("\nThese(This) function(s) contain(s) switch statement: ")
		for _, val := range switchCheck {
			fmt.Print(val, " ")
		}
		fmt.Println()
		/*
			for s := range caseCheck {
				for _, value := range caseCheck[s].value {
					fmt.Printf("%s\n", value)
				}
				fmt.Println()

			}*/
	}

	// TODO for multiple functions with switch statement
	var caseReplacement []string
	caseFlag := false
	for i := 0; i < len(caseCheck); i++ {
		if caseCheck[i].funcName != "" {
			for j := i + 1; j < len(caseCheck); j++ {
				if caseCheck[j].funcName != "" {
					if len(caseCheck[i].value) == len(caseCheck[j].value) {
						// Compare details between cases
						for idx := range caseCheck[i].value {
							if strings.Compare(caseCheck[i].value[idx].path, caseCheck[j].value[idx].path) == 0 {
								if len(caseCheck[i].value[idx].value) == len(caseCheck[j].value[idx].value) {
									for idxValue := range caseCheck[i].value[idx].value {
										if strings.Compare(caseCheck[i].value[idx].value[idxValue], caseCheck[j].value[idx].value[idxValue]) == 0 {
											caseFlag = true
										} else {
											if (contains(funcList, caseCheck[i].value[idx].value[idxValue]) && contains(funcList, caseCheck[j].value[idx].value[idxValue])) ||
												(contains(typeList, caseCheck[i].value[idx].value[idxValue]) && contains(typeList, caseCheck[j].value[idx].value[idxValue])) ||
												(strings.Contains(caseCheck[i].value[idx].value[idxValue], "Set") && strings.Contains(caseCheck[j].value[idx].value[idxValue], "Set")) ||
												(strings.Contains(caseCheck[i].value[idx].value[idxValue], "Get") && strings.Contains(caseCheck[j].value[idx].value[idxValue], "Get")) {
												caseFlag = true
											} else {
												caseFlag = false
												break
											}
										}
									}
								} else {
									caseFlag = false
									break
								}
							} else if strings.Contains(caseCheck[i].value[idx].path, "*ast.SelectorExpr") {
								// listFunctions2[i].value[idx] looks like "... -> *ast.SelectorExpr [unsafe Pointer]"
								// and listFunctions2[j].value[idx] looks like " ... [TYPE]"
								if !contains(typeList, caseCheck[j].value[idx].value[0]) {
									caseFlag = false
									break
								}
								for _, val := range caseCheck[i].value[idx].value {
									if !checkUnsafeUsages(val) {
										caseFlag = false
										break
									}
								}
								caseFlag = true
							} else if strings.Contains(caseCheck[j].value[idx].path, "*ast.SelectorExpr") {
								// listFunctions2[i].value[idx] looks like " ... [TYPE]"
								// and listFunctions2[j].value[idx] looks like "... -> *ast.SelectorExpr [unsafe Pointer]"
								if !contains(typeList, caseCheck[i].value[idx].value[0]) {
									caseFlag = false
									break
								}
								for _, val := range caseCheck[j].value[idx].value {
									if !checkUnsafeUsages(val) {
										caseFlag = false
										break
									}
								}
								caseFlag = true
							} else {
								caseFlag = false
								break
							}
						}
					} else {
						caseFlag = false
						continue // continue the progress comparing a next function to the compared one

					}
				}
				if caseFlag == true {
					if !contains(caseReplacement, caseCheck[j].funcName) {
						caseReplacement = append(caseReplacement, caseCheck[j].funcName)
					}
					caseFlag = false
				}
			}
		}

	}

	if len(caseReplacement) > 0 {
		fmt.Print("These cases have same structure: ")
		for _, val := range caseReplacement {
			fmt.Print(val, " ")
		}
	}

	// create a text file
	f, err := os.Create(filename + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	var str string
	nextLine := "\n\n"
	for s := range listFunctions2 {
		if listFunctions2[s].funcName != "" {
			str = "function name is " + listFunctions2[s].funcName + " \n"
			_, err2 := f.WriteString(str)
			if err2 != nil {
				log.Fatal(err2)
			}
			for _, value := range listFunctions2[s].value {
				str = fmt.Sprintln("\t", value)
				_, err2 := f.WriteString(str)
				if err2 != nil {
					log.Fatal(err2)
				}
			}
			_, err2 = f.WriteString(nextLine)
			if err2 != nil {
				log.Fatal(err2)
			}

		} else {
			for _, value := range listFunctions2[s].value {
				str = fmt.Sprintln(value)
				_, err2 := f.WriteString(str)
				if err2 != nil {
					log.Fatal(err2)
				}
			}
			_, err2 := f.WriteString(nextLine)
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	}
	/*for s := range listFunctions2 {
		fmt.Print("function name is ", listFunctions2[s].funcName, " ")
		fmt.Printf("%d\n", len(listFunctions2[s].value))
		for _, value := range listFunctions2[s].value {
			fmt.Printf("%s\n", value)
		}
		fmt.Println()
		fmt.Println()

	}*/

	for s := range genCheck {
		if len(genCheck[s]) > 1 {
			fmt.Print("These functions have a same structure and the code are reused: ")
			for _, value := range genCheck[s] {
				fmt.Print(value, " ")
			}
			fmt.Println()
		}
	}
}
