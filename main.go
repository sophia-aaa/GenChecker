package main

import (
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

// Basic structure for every function in the input file
type basicStr struct {
	funcName string
	value    []elem
}

type funcNameList struct {
	lists []string
}

type checkCases struct {
	funcName string
	cases    []basicStr
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
func buildAstDataStr(filename string) []basicStr {
	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var elem_list []elem
	nameFunction := ""
	astNode := ""
	var listFunctions2 []basicStr
	var astValue []string

	ast.Inspect(astTree, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.FuncDecl:
			//	fmt.Println(x, "\t\t", reflect.TypeOf(x).String())
			//	fmt.Println(fset.Position(x.Pos()), fset.Position(x.End()))
			if nameFunction != "" { // if a new root meets
				elem_list = append(elem_list, elem{astNode, astValue})
				listFunctions2 = append(listFunctions2, basicStr{nameFunction, elem_list})
				elem_list = []elem{}
				astNode = ""
				astValue = []string{}
			}
			nameFunction = x.Name.String()
		case *ast.Ident:
			//fmt.Println(fset.Position(x.Pos()), reflect.TypeOf(x).String(), "\t", x.Name)
			astValue = append(astValue, x.Name)
		case *ast.CaseClause, *ast.SwitchStmt:
			// *ast.CaseClause -> *ast.SelectorExpr
			elem_list = append(elem_list, elem{astNode, astValue})
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
					astNode = ""
					astValue = []string{}
				}
				if astNode == "" {
					astNode += reflect.TypeOf(x).String()
				} else {
					astNode += " -> " + reflect.TypeOf(x).String()
				}
			}
		}
		return true
	})
	if nameFunction != "" { // add last function
		elem_list = append(elem_list, elem{astNode, astValue})
		listFunctions2 = append(listFunctions2, basicStr{nameFunction, elem_list})
		elem_list = []elem{}
		astNode = ""
		astValue = []string{}
	}
	return listFunctions2
}

func checkUnsafeUsages(str string) bool {
	return contains([]string{"unsafe", "Pointer"}, str)
}

func checkGenerics(listFunctions2 []basicStr, funcList []string, typeList []string) [][]string {
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
							} else if listFunctions2[i].value[idx].path+" -> *ast.SelectorExpr" == listFunctions2[j].value[idx].path {
								// This is a special case for modified elem list
								// listFunctions2[j].value[idx].path is listFunctions2[i].value[idx].path " -> *ast.SelectorExpr"
								// listFunctions2[j].value[idx].value must look like [... unsafe Pointer] and
								// listFunctions2[i].value[idx].value must have type variable in its elem value so like [... type]
								len1 := len(listFunctions2[i].value[idx].value)
								len2 := len(listFunctions2[j].value[idx].value)
								if len2 > 1 && !contains(typeList, listFunctions2[i].value[idx].value[len1-1]) ||
									(!checkUnsafeUsages(listFunctions2[j].value[idx].value[len2-1]) &&
										!checkUnsafeUsages(listFunctions2[j].value[idx].value[len2-2])) {
									flag = false
									break
								}
								flag = true
							} else if listFunctions2[j].value[idx].path+" -> *ast.SelectorExpr" == listFunctions2[i].value[idx].path {
								// This is a special case for modified elem list
								// listFunctions2[i].value[idx].path is listFunctions2[j].value[idx].path " -> *ast.SelectorExpr"
								// listFunctions2[i].value[idx].value must look like [... unsafe Pointer] and
								// listFunctions2[j].value[idx].value must have type variable in its elem value so like [... type]
								len1 := len(listFunctions2[j].value[idx].value)
								len2 := len(listFunctions2[i].value[idx].value)
								if len2 > 2 && !contains(typeList, listFunctions2[j].value[idx].value[len1-1]) ||
									(!checkUnsafeUsages(listFunctions2[i].value[idx].value[len2-1]) &&
										!checkUnsafeUsages(listFunctions2[i].value[idx].value[len2-2])) {
									flag = false
									break
								}
								flag = true
							} else if strings.Contains(listFunctions2[i].value[idx].path, "*ast.SelectorExpr") {
								// listFunctions2[i].value[idx] looks like "... -> *ast.SelectorExpr [... unsafe Pointer]"
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
								// and listFunctions2[j].value[idx] looks like "... -> *ast.SelectorExpr [... unsafe Pointer]"
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
	return genCheck
}

func checkReusedCases(caseWFunc []checkCases, funcList []string, typeList []string) []elem {
	var caseListCheck []elem
	var caseReplacement []string
	caseFlag := false

	if len(caseWFunc) > 0 {
		for k := range caseWFunc {
			for i := 0; i < len(caseWFunc[k].cases); i++ {
				if caseWFunc[k].cases[i].funcName != "" {
					for j := i + 1; j < len(caseWFunc[k].cases); j++ {
						if caseWFunc[k].cases[j].funcName != "" {
							if len(caseWFunc[k].cases[i].value) == len(caseWFunc[k].cases[j].value) {
								// Compare details between cases
								for idx := range caseWFunc[k].cases[i].value {
									if strings.Compare(caseWFunc[k].cases[i].value[idx].path, caseWFunc[k].cases[j].value[idx].path) == 0 {
										if len(caseWFunc[k].cases[i].value[idx].value) == len(caseWFunc[k].cases[j].value[idx].value) {
											for idxValue := range caseWFunc[k].cases[i].value[idx].value {
												if strings.Compare(caseWFunc[k].cases[i].value[idx].value[idxValue], caseWFunc[k].cases[j].value[idx].value[idxValue]) == 0 {
													caseFlag = true
												} else {
													if (contains(funcList, caseWFunc[k].cases[i].value[idx].value[idxValue]) && contains(funcList, caseWFunc[k].cases[j].value[idx].value[idxValue])) ||
														(contains(typeList, caseWFunc[k].cases[i].value[idx].value[idxValue]) && contains(typeList, caseWFunc[k].cases[j].value[idx].value[idxValue])) ||
														(strings.Contains(caseWFunc[k].cases[i].value[idx].value[idxValue], caseWFunc[k].funcName) && strings.Contains(caseWFunc[k].cases[j].value[idx].value[idxValue], caseWFunc[k].funcName)) {
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
									} else if caseWFunc[k].cases[i].value[idx].path+" -> *ast.SelectorExpr" == caseWFunc[k].cases[j].value[idx].path {
										// This is a special case for modified elem list
										// listFunctions2[j].value[idx].path is listFunctions2[i].value[idx].path " -> *ast.SelectorExpr"
										// listFunctions2[j].value[idx].value must look like [... unsafe Pointer] and
										// listFunctions2[i].value[idx].value must have type variable in its elem value so like [... type]
										len1 := len(caseWFunc[k].cases[i].value[idx].value)
										len2 := len(caseWFunc[k].cases[j].value[idx].value)
										if len2 > 1 && !contains(typeList, caseWFunc[k].cases[i].value[idx].value[len1-1]) ||
											(!checkUnsafeUsages(caseWFunc[k].cases[j].value[idx].value[len2-1]) &&
												!checkUnsafeUsages(caseWFunc[k].cases[j].value[idx].value[len2-2])) {
											caseFlag = false
											break
										}
										caseFlag = true
									} else if caseWFunc[k].cases[j].value[idx].path+" -> *ast.SelectorExpr" == caseWFunc[k].cases[i].value[idx].path {
										// This is a special case for modified elem list
										// listFunctions2[i].value[idx].path is listFunctions2[j].value[idx].path " -> *ast.SelectorExpr"
										// listFunctions2[i].value[idx].value must look like [... unsafe Pointer] and
										// listFunctions2[j].value[idx].value must have type variable in its elem value so like [... type]
										len1 := len(caseWFunc[k].cases[j].value[idx].value)
										len2 := len(caseWFunc[k].cases[i].value[idx].value)
										if len2 > 2 && !contains(typeList, caseWFunc[k].cases[j].value[idx].value[len1-1]) ||
											(!checkUnsafeUsages(caseWFunc[k].cases[i].value[idx].value[len2-1]) &&
												!checkUnsafeUsages(caseWFunc[k].cases[i].value[idx].value[len2-2])) {
											caseFlag = false
											break
										}
										caseFlag = true
									} else if strings.Contains(caseWFunc[k].cases[i].value[idx].path, "*ast.SelectorExpr") {
										// listFunctions2[i].value[idx] looks like "... -> *ast.SelectorExpr [unsafe Pointer]"
										// and listFunctions2[j].value[idx] looks like " ... [TYPE]"
										if !contains(typeList, caseWFunc[k].cases[j].value[idx].value[0]) {
											caseFlag = false
											break
										}
										for _, val := range caseWFunc[k].cases[i].value[idx].value {
											if !checkUnsafeUsages(val) {
												caseFlag = false
												break
											}
										}
										caseFlag = true
									} else if strings.Contains(caseWFunc[k].cases[j].value[idx].path, "*ast.SelectorExpr") {
										// listFunctions2[i].value[idx] looks like " ... [TYPE]"
										// and listFunctions2[j].value[idx] looks like "... -> *ast.SelectorExpr [unsafe Pointer]"
										if !contains(typeList, caseWFunc[k].cases[i].value[idx].value[0]) {
											caseFlag = false
											break
										}
										for _, val := range caseWFunc[k].cases[j].value[idx].value {
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
							if !contains(caseReplacement, caseWFunc[k].cases[j].funcName) {
								caseReplacement = append(caseReplacement, caseWFunc[k].cases[j].funcName)
							}
							caseFlag = false
						}
					}
				}
			}
			if len(caseReplacement) > 0 {
				caseListCheck = append(caseListCheck, elem{caseWFunc[k].funcName, caseReplacement})
				caseReplacement = []string{}
			}
		}
	}
	return caseListCheck
}

func checkSelectorExpr(listFunctions2 []basicStr) []basicStr {
	var modElemList []elem
	var modListFunctions2 []basicStr
	var modFuncName string
	var modPath string
	var modAstValue []string
	if len(listFunctions2) > 0 {
		for l := range listFunctions2 {
			modFuncName = listFunctions2[l].funcName
			for i := 0; i < len(listFunctions2[l].value); i++ {
				modPath = listFunctions2[l].value[i].path
				for _, val := range listFunctions2[l].value[i].value {
					modAstValue = append(modAstValue, val)
				}
				if i < len(listFunctions2[l].value)-1 && listFunctions2[l].value[i+1].path == "*ast.SelectorExpr" {
					if listFunctions2[l].value[i+1].value[0] == "unsafe" && listFunctions2[l].value[i+1].value[1] == "Pointer" {
						modPath += " -> *ast.SelectorExpr"
						modAstValue = append(modAstValue, "unsafe")
						modAstValue = append(modAstValue, "Pointer")
						i++
					}
				}
				modElemList = append(modElemList, elem{modPath, modAstValue})
				modPath = ""
				modAstValue = []string{}
			}
			modListFunctions2 = append(modListFunctions2, basicStr{modFuncName, modElemList})
			modFuncName = ""
			modElemList = []elem{}
		}
	}

	return modListFunctions2
}

func createTextFile(filename string, listFunctions2 []basicStr) {
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
}

func checkSwitchStatement(listFunctions2 []basicStr) []checkCases {
	var switchCheck []string
	var caseListWVar []basicStr
	var caseName string
	var caseList []elem
	var caseWFunc []checkCases
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
					caseListWVar = append(caseListWVar, basicStr{caseName, caseList})
					caseList = []elem{}
				}
				caseName = listFunctions2[i].value[j].value[1]
				caseList = append(caseList, elem{listFunctions2[i].value[j].path, listFunctions2[i].value[j].value})
				caseBool = true
				continue
			} else if strings.Contains(listFunctions2[i].value[j].path, "*ast.CaseClause") {
				if len(caseList) > 0 {
					caseBool = false
					caseListWVar = append(caseListWVar, basicStr{caseName, caseList})
					caseList = []elem{}
				}
				continue
			}
			if caseBool {
				caseList = append(caseList, elem{listFunctions2[i].value[j].path, listFunctions2[i].value[j].value})
			}
		}
		if len(caseListWVar) > 0 {
			caseWFunc = append(caseWFunc, checkCases{listFunctions2[i].funcName, caseListWVar})
			caseListWVar = []basicStr{}

		}
	}

	if len(switchCheck) > 2 {
		fmt.Println()
		fmt.Println("These functions have (a) switch statement(s):")
		for _, val := range switchCheck {
			fmt.Print(val, " ")
		}
	} else if len(switchCheck) == 1 {
		fmt.Println()
		fmt.Println("This function has (a) switch statement(s):")
		for _, val := range switchCheck {
			fmt.Print(val, " ")
		}
	} else {
		fmt.Println()
		fmt.Println("There is no function with switch statements")
	}
	return caseWFunc
}

func main() {
	// filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run gen.go - test.go
	filename := os.Args[2]

	typeList := []string{
		"bool", "bType", "int", "iType", "int8", "i8Type", "int16", "i16Type", "int32", "i32Type", "int64", "i64Type", "uint", "uType",
		"uint8", "u8Type", "uint16", "u16Type", "uint32", "u32Type", "uint64", "u64Type", "uintptr", "uintptrType",
		"float32", "f32Type", "float64", "f64Type", "complex64", "c64Type", "complex128", "c128Type", "string", "strType",
		"unsafe", "Pointer", "unsafePointerType",
	}

	listFunctions2 := buildAstDataStr(filename)

	var funcList []string
	for s := range listFunctions2 {
		if listFunctions2[s].funcName != "" {
			funcList = append(funcList, listFunctions2[s].funcName)
		}
	}

	// check leaf of SelectorExpr and unsafe Pointer
	modListFunctions2 := checkSelectorExpr(listFunctions2)

	// Check Generic Replacement
	genCheck := checkGenerics(modListFunctions2, funcList, typeList)
	for s := range genCheck {
		if len(genCheck[s]) > 1 {
			fmt.Print("These functions have a same structure and the code are reused: ")
			for _, value := range genCheck[s] {
				fmt.Print(value, " ")
			}
			fmt.Println()
		}
	}

	// This variable is for checking switch statement
	caseWFunc := checkSwitchStatement(modListFunctions2)

	// Check reused cases in switch statement
	caseListCheck := checkReusedCases(caseWFunc, funcList, typeList)

	if len(caseListCheck) > 0 {
		fmt.Println()
		fmt.Println()
		for _, val := range caseListCheck {
			fmt.Print("This Function ", val.path, " has switch-statement and reused cases : \n")
			for _, cases := range val.value {
				fmt.Print(cases, " ")
			}
			fmt.Println()
			fmt.Println()
		}
	}

	// create a text file
	//createTextFile(filename, listFunctions2)
	createTextFile(filename, modListFunctions2)

}
