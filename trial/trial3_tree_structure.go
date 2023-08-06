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
	"strconv"
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

type basicCases struct {
	depth int
	value string
}

var Tree2str []string
var Tree2cases []basicCases

func contains(strArr []string, str string) bool {
	for _, val := range strArr {
		if val == str {
			return true
		}
		if strings.Contains(strings.ToLower(str), strings.ToLower(val)) {
			return true
		}
	}
	return false
}
func isSameString(strArr []string, str string) bool {
	for _, val := range strArr {
		if strings.EqualFold(val, str) {
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
func buildAstCaseStr(Tree2cases []basicCases) []checkCases {
	var funcName string
	var path string
	var astValue string
	var astValueList []string
	var depthFirstCase int
	var caseName string
	var elemList []elem
	var caseCheck []basicStr
	var funcCheck []checkCases
	flag := true
	count := 0

	for _, val := range Tree2cases {
		if strings.Contains(val.value, "*ast.FuncDecl") { // function begins
			if funcName != "" {
				elemList = append(elemList, elem{path, astValueList})
				caseCheck = append(caseCheck, basicStr{caseName, elemList})
				funcCheck = append(funcCheck, checkCases{funcName, caseCheck})
				//fmt.Println(path, " ", astValueList)
				path = ""
				caseName = ""
				astValueList = []string{}
				elemList = []elem{}
				caseCheck = []basicStr{}
				count = 0
				depthFirstCase = 0
			}
			funcName = val.value[14:len(val.value)]
			//fmt.Println()
			//fmt.Println()
			//fmt.Println(">> ", funcName, " <<")
			flag = true
			continue
		}
		if val.depth == 2 && strings.Contains(val.value, "*ast.Ident") { // file name
			if len(astValueList) > 0 || path != "" {
				elemList = append(elemList, elem{path, astValueList})
				//fmt.Println(path, " ", astValueList)
				path = ""
				astValueList = []string{}
			}
			astValue = val.value[11:len(val.value)]
			astValueList = append(astValueList, astValue)
			//fmt.Println(astValueList)
			elemList = append(elemList, elem{"file name", astValueList}) // filename is not a convention from go. I added.
			path = ""
			astValueList = []string{}
			continue
		} else if strings.Contains(val.value, "*ast.Ident") {
			astValue = val.value[11:len(val.value)]
			astValueList = append(astValueList, astValue)
			continue
		}
		if flag && strings.Contains(val.value, "*ast.CaseClause") { // to set a depth
			// case in switch begins
			if len(astValueList) > 0 || path != "" {
				//fmt.Println(path, " ", astValueList)
				elemList = append(elemList, elem{path, astValueList})
				caseCheck = append(caseCheck, basicStr{"", elemList})
				path = ""
				elemList = []elem{}
				astValueList = []string{}
			}
			depthFirstCase = val.depth
			path = val.value
			flag = false
			caseName = strconv.Itoa(count) + ".case"
			//fmt.Println("case begins!!!\t", " ", depthFirstCase, " ", len(val.value))
			continue
		} else if depthFirstCase == val.depth && strings.Contains(val.value, "*ast.CaseClause") {
			if len(astValueList) > 0 || path != "" {
				//fmt.Println(path, " ", astValueList)
				elemList = append(elemList, elem{path, astValueList})
				caseCheck = append(caseCheck, basicStr{caseName, elemList})
				path = ""
				elemList = []elem{}
				astValueList = []string{}
				count++
			}
			caseName = strconv.Itoa(count) + ".case"
			path = val.value
			//fmt.Println("case begins!!!\t", " ", depthFirstCase, " ", len(val.value))
			continue
		}
		if len(astValueList) > 0 {
			//fmt.Println(path, " ", astValueList)
			elemList = append(elemList, elem{path, astValueList})
			path = ""
			astValueList = []string{}
		}
		if path == "" {
			path = val.value
		} else {
			path = path + " -> " + val.value
		}
	}
	if funcName != "" {
		elemList = append(elemList, elem{path, astValueList})
		caseCheck = append(caseCheck, basicStr{caseName, elemList})
		funcCheck = append(funcCheck, checkCases{funcName, caseCheck})
	}

	return funcCheck
}

func checkUnsafeUsages(str string) bool {
	return contains([]string{"unsafe", "Pointer"}, str)
}

func checkReusedCases(caseWFunc []checkCases, funcList []string, typeList []string) []elem {
	var caseListCheck []elem
	var caseReplacement []string
	caseFlag := false

	if len(caseWFunc) > 0 {
		for k := range caseWFunc {
			for i := 0; i < len(caseWFunc[k].cases); i++ {
				if strings.Contains(caseWFunc[k].cases[i].funcName, "case") {
					for j := i + 1; j < len(caseWFunc[k].cases); j++ {
						if strings.Contains(caseWFunc[k].cases[j].funcName, "case") {
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
							if !isSameString(caseReplacement, caseWFunc[k].cases[i].funcName) {
								caseReplacement = append(caseReplacement, caseWFunc[k].cases[i].funcName)
							}
							if !isSameString(caseReplacement, caseWFunc[k].cases[j].funcName) {
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

func createTextFile(filename string, strList []string) {
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
	if len(strList) > 0 {
		for _, val := range strList {
			_, err1 := f.WriteString(val)
			if err1 != nil {
				log.Fatal(err1)
			}
		}
	}
}

// Reference for type visitor int and func Visit : https://golangdocs.com/golang-ast-package
type visitor int

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	var str string
	var str4cases string
	var depth int
	// int(v) is a depth of a current node
	if reflect.TypeOf(n).String() == "*ast.Ident" {
		str = fmt.Sprintf("%s%d: ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)

		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), n)
		str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), n)
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
	} else if reflect.TypeOf(n).String() == "*ast.FuncDecl" {
		funcName := fmt.Sprintf("%v", n)
		count := 0
		var strIdx int
		var endIdx int
		for i := 0; i < len(funcName); i++ {
			if funcName[i] == ' ' {
				count++
			}
			if funcName[i] == ' ' && count == 2 {
				strIdx = i + 1
			}
			if funcName[i] == ' ' && count == 3 {
				endIdx = i
			}
		}
		str = fmt.Sprintf("%s%d: ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)
		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), funcName[strIdx:endIdx])
		str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), funcName[strIdx:endIdx])
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
		count = 0
	} else {
		str = fmt.Sprintf("%s%d: ", strings.Repeat("\t", int(v)), int(v))
		Tree2str = append(Tree2str, str)
		depth = int(v)
		str = fmt.Sprintf("%T\n", n)
		str4cases = fmt.Sprintf("%T", n)
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})

	}
	return v + 1
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

func main() {
	// filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run gen.go - example.go
	filename := os.Args[2]
	fmt.Println(filename)
	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	typeList := []string{
		"bool", "bType", "int", "iType", "int8", "i8Type", "int16", "i16Type", "int32", "i32Type", "int64", "i64Type", "uint", "uType",
		"uint8", "u8Type", "uint16", "u16Type", "uint32", "u32Type", "uint64", "u64Type", "uintptr", "uintptrType",
		"float32", "f32Type", "float64", "f64Type", "complex64", "c64Type", "complex128", "c128Type", "string", "strType",
		"unsafe", "Pointer", "UnsafePointer", "unsafePointerType",
	}
	ast.Print(fset, astTree)

	var v visitor
	ast.Walk(v, astTree)
	if len(Tree2str) > 0 {
		createTextFile("tree", Tree2str)
	}

	funcCheck := buildAstCaseStr(Tree2cases)

	var funcList []string
	for s := range funcCheck {
		if funcCheck[s].funcName != "" {
			funcList = append(funcList, funcCheck[s].funcName)
		}
	}
	var modifiedFuncCheck []checkCases
	for idx := range funcCheck {
		modifiedFuncCheck = append(modifiedFuncCheck, checkCases{funcCheck[idx].funcName, checkSelectorExpr(funcCheck[idx].cases)})
	}

	caseListCheck := checkReusedCases(modifiedFuncCheck, funcList, typeList)

	// count number of case clauses
	lengthList := make([]int, len(modifiedFuncCheck))
	numberOfCase := 0
	for idx := range modifiedFuncCheck {
		for _, val := range modifiedFuncCheck[idx].cases {
			if strings.Contains(val.funcName, "case") {
				numberOfCase++
			}
		}
		lengthList[idx] = numberOfCase
		numberOfCase = 0
	}

	var caseClauseList []string
	var checkCaseClause []elem
	var funcCaseClause []basicStr

	for idx := range modifiedFuncCheck {
		for _, val := range modifiedFuncCheck[idx].cases {
			if strings.Contains(val.funcName, "case") && val.value[0].path == "*ast.CaseClause" {
				caseClauseList = val.value[0].value
				checkCaseClause = append(checkCaseClause, elem{val.funcName, caseClauseList})
			}
		}
		funcCaseClause = append(funcCaseClause, basicStr{modifiedFuncCheck[idx].funcName, checkCaseClause})
		checkCaseClause = []elem{}
	}

	// check whether the case clauses are type variables
	flag4Case := make([]bool, len(funcCaseClause))
	flag4outer := false
	for idx := range funcCaseClause {
		fmt.Println(funcCaseClause[idx].funcName)
		for _, val := range funcCaseClause[idx].value {
			if len(val.value) > 0 {
				for _, value := range val.value {
					if !contains(typeList, value) {
						flag4Case[idx] = false
						flag4outer = true
						break
					}
					flag4Case[idx] = true
				}
			}
			if flag4outer {
				flag4outer = false
				break
			}
		}
	}
	fmt.Println(flag4Case)

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

	for idx := range caseListCheck {
		lengthCase := len(caseListCheck[idx].value)
		if flag4Case[idx] && ((lengthCase == lengthList[idx]) ||
			((lengthCase == lengthList[idx]-1) && // in this case, there exists a default case clause and it will not be considered.
				(modifiedFuncCheck[idx].cases[1].value[0].path != modifiedFuncCheck[idx].cases[len(modifiedFuncCheck[0].cases)-1].value[0].path))) {
			fmt.Println("The function ", modifiedFuncCheck[idx].funcName, " can be replaced by Generics.")
		}
	}

	/*	for idx := range modifiedFuncCheck {
		for _, w := range modifiedFuncCheck[idx].cases {
			fmt.Println(w.funcName)
			for _, w1 := range w.value {
				fmt.Println(w1.path)
			}

		}
	}*/

	//	fmt.Println(modifiedFuncCheck[0].cases[1].value[0].path, "\t", modifiedFuncCheck[0].cases[len(modifiedFuncCheck[0].cases)-1].value[0].path)

}
