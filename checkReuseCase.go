package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type checkCases struct {
	funcName string
	cases    []basicCaseStr
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
										// listFunctions[j].value[idx].path is listFunctions[i].value[idx].path " -> *ast.SelectorExpr"
										// listFunctions[j].value[idx].value must look like [... unsafe Pointer] and
										// listFunctions[i].value[idx].value must have type variable in its elem value so like [... type]
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
										// listFunctions[i].value[idx].path is listFunctions[j].value[idx].path " -> *ast.SelectorExpr"
										// listFunctions[i].value[idx].value must look like [... unsafe Pointer] and
										// listFunctions[j].value[idx].value must have type variable in its elem value so like [... type]
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
										// listFunctions[i].value[idx] looks like "... -> *ast.SelectorExpr [unsafe Pointer]"
										// and listFunctions[j].value[idx] looks like " ... [TYPE]"
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
										// listFunctions[i].value[idx] looks like " ... [TYPE]"
										// and listFunctions[j].value[idx] looks like "... -> *ast.SelectorExpr [unsafe Pointer]"
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
								break // break and continue the progress comparing a next function to the compared one

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

func checkSwitchStatement(filename string, listFunctions []basicStr) (bool, []basicCaseStr) {
	var switchCheck []string
	var existsSwitch bool

	for i := range listFunctions {
		for j := range listFunctions[i].value {
			if strings.Contains(listFunctions[i].value[j].path, "*ast.SwitchStmt") {
				if !contains(switchCheck, listFunctions[i].funcName) {
					switchCheck = append(switchCheck, listFunctions[i].funcName)
				}
			}
		}
	}

	if len(switchCheck) >= 1 {
		existsSwitch = true
	}

	if existsSwitch {
		fset := token.NewFileSet()
		astTree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		var v visitor
		ast.Walk(v, astTree)
		if len(Tree2str) > 0 {
			createTextFileFromString("tree", Tree2str)
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
			modifiedFuncCheck = append(modifiedFuncCheck, checkCases{funcCheck[idx].funcName, checkSelectorExprCase(funcCheck[idx].cases)})
		}

		return true, checkSwitchCases(modifiedFuncCheck, funcList, typeList)

	}

	return false, []basicCaseStr{}
}

func checkSwitchCases(modifiedFuncCheck []checkCases, funcList []string, typeList []string) []basicCaseStr {

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
	// *********************
	var checkCaseClause []elem
	var funcCaseClause []basicCaseStr
	// filter case clause
	for idx := range modifiedFuncCheck {
		for _, val := range modifiedFuncCheck[idx].cases {
			if strings.Contains(val.funcName, "case") &&
				(val.value[0].path == "*ast.CaseClause" ||
					val.value[0].path == "*ast.CaseClause -> *ast.SelectorExpr") {
				//fmt.Println(val.value[0].value)
				checkCaseClause = append(checkCaseClause, elem{val.funcName, val.value[0].value})
			}
		}
		funcCaseClause = append(funcCaseClause, basicCaseStr{modifiedFuncCheck[idx].funcName, checkCaseClause})
		checkCaseClause = []elem{}
	}

	// check whether the case clauses are type variables
	flag4Case := make([]bool, len(funcCaseClause))
	flag4outer := false
	for idx := range funcCaseClause {
		for i, val := range funcCaseClause[idx].value {
			if len(val.value) > 0 {
				for _, value := range val.value {
					if !contains(typeList, value) {
						if i != len(funcCaseClause[idx].value)-1 {
							flag4Case[idx] = false
							flag4outer = true
							break
						}
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

	if len(caseListCheck) > 0 {
		for idx := range caseListCheck {
			lengthCase := len(caseListCheck[idx].value)
			/*				fmt.Println(flag4Case[idx])
							fmt.Println(lengthCase == lengthList[idx])
							fmt.Println(lengthCase == lengthList[idx]-1)
							fmt.Println(modifiedFuncCheck[idx].cases[1].value[0].path != modifiedFuncCheck[idx].cases[len(modifiedFuncCheck[idx].cases)-1].value[0].path)*/
			if flag4Case[idx] && ((lengthCase == lengthList[idx]) ||
				((lengthCase == lengthList[idx]-1) && // in this case, there exists a default case clause and it will not be considered.
					(modifiedFuncCheck[idx].cases[1].value[0].path != modifiedFuncCheck[idx].cases[len(modifiedFuncCheck[idx].cases)-1].value[0].path))) {
				fmt.Println()
				fmt.Println("The function ", modifiedFuncCheck[idx].funcName, " has same case statement for different types in switch statement and the cases are: ")
				for _, val := range funcCaseClause[idx].value {
					fmt.Print(val.value, " ")
				}
				fmt.Println()
			}
			fmt.Println()
		}
		return funcCaseClause
	}

	// there are no reused cases
	return []basicCaseStr{}
}
