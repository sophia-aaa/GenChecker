package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func checkSwitchStatement(filename string, listFunctions []basicStr) (bool, []caseFilteredResult) {
	var switchCheck []basicFunc
	var existsSwitch bool

	for i := range listFunctions {
		for j := range listFunctions[i].value {
			if strings.Contains(listFunctions[i].value[j].path, "*ast.SwitchStmt") ||
				strings.Contains(listFunctions[i].value[j].path, "*ast.TypeSwitchStmt") {

				if !checkDuplicateInFunc(switchCheck, listFunctions[i].funcName, listFunctions[i].funcToken) {
					switchCheck = append(switchCheck, basicFunc{listFunctions[i].funcName, listFunctions[i].funcToken})
				}
				continue
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

		// functions which have switch statements collected in funcCheck
		funcCheck := buildAstCaseStr(Tree2cases, switchCheck)

		if len(funcCheck) > 0 {
			fmt.Println("Before checking")

			if len(funcCheck) > 0 {
				fmt.Println("There is(are) (a) switch statement(s): ")
			}
			for _, fnc := range funcCheck {
				fmt.Print(fnc.funcName, " ")
			}
			fmt.Println()
		}
		return true, checkSwitchCases(filename, funcCheck, typeList)

	}

	return false, []caseFilteredResult{}
}

func checkSwitchCases(filename string, funcList []caseResult, typeList []string) []caseFilteredResult {
	//var funcResult []caseFilteredResult

	createCasesTextFile(filename, funcList)
	/*for _, val := range funcList {
		fmt.Println("function name: ", val.funcName)
		for i, value := range val.cases {
			if i < 2 {
				fmt.Println(value.caseName)
				for _, cases := range value.value {
					fmt.Print(cases.path, "\t", cases.value)
					fmt.Println()
				}
			}
		}
		fmt.Println()
	}*/

	var result []caseFilteredResult
	var caseList []string
	flag := false
	for _, fnc := range funcList {
		checkFuncName := funcNameDivider(fnc.funcName)
		for i := range fnc.cases {
			for j := i + 1; j < len(fnc.cases); j++ {
				if len(fnc.cases[i].value) == len(fnc.cases[j].value) {
					for ind := range fnc.cases[i].value {
						if strings.EqualFold(fnc.cases[i].value[ind].path, fnc.cases[j].value[ind].path) {
							if len(fnc.cases[i].value[ind].value) == len(fnc.cases[j].value[ind].value) {
								for indVal := range fnc.cases[i].value[ind].value {
									if strings.EqualFold(fnc.cases[i].value[ind].value[indVal], fnc.cases[j].value[ind].value[indVal]) {
										flag = true
									} else {
										if strings.Contains(fnc.cases[j].value[ind].value[indVal], " ") || strings.Contains(fnc.cases[j].value[ind].value[indVal], " ") ||
											strings.Contains(fnc.cases[j].value[ind].value[indVal], "_") || strings.Contains(fnc.cases[j].value[ind].value[indVal], "_") {
											flag = false
											break
										} else if strings.Contains(fnc.cases[i].value[ind].value[indVal], fnc.funcName) && strings.Contains(fnc.cases[j].value[ind].value[indVal], fnc.funcName) {
											flag = true
										} else if contains(typeList, fnc.cases[i].value[ind].value[indVal]) && contains(typeList, fnc.cases[j].value[ind].value[indVal]) {
											flag = true
										} else if contains(checkFuncName, fnc.cases[i].value[ind].value[indVal]) && contains(checkFuncName, fnc.cases[j].value[ind].value[indVal]) {
											flag = true
										} else {
											flag = false
											break
										}
									}
								}
								if !flag {
									break
								}
							} else {
								flag = false
								break
							}
						} else if (strings.Contains(fnc.cases[i].value[ind].path, "*ast.SelectorExpr") && !strings.Contains(fnc.cases[j].value[ind].path, "*ast.SelectorExpr")) ||
							(!strings.Contains(fnc.cases[i].value[ind].path, "*ast.SelectorExpr") && strings.Contains(fnc.cases[j].value[ind].path, "*ast.SelectorExpr")) {
							length := 0
							if (len(fnc.cases[i].value[ind].value) == 0) || (len(fnc.cases[j].value[ind].value) == 0) {
								flag = false
								break
							} else {
								if len(fnc.cases[i].value[ind].value) <= len(fnc.cases[j].value[ind].value) {
									length = len(fnc.cases[i].value[ind].value) // will be index a
									a, b := 0, 0
									for ; a < length; a++ {
										if !strings.EqualFold(fnc.cases[i].value[ind].value[a], fnc.cases[j].value[ind].value[b]) {
											if strings.Contains(fnc.cases[i].value[ind].value[a], " ") || strings.Contains(fnc.cases[j].value[ind].value[b], " ") ||
												strings.Contains(fnc.cases[i].value[ind].value[a], "_") || strings.Contains(fnc.cases[j].value[ind].value[b], "_") {
												flag = false
												break
											} else if contains(typeList, fnc.cases[i].value[ind].value[a]) {
												if strings.EqualFold(fnc.cases[j].value[ind].value[b], "unsafe") {
													if !checkUnsafeUsages(fnc.cases[j].value[ind].value[b]) {
														flag = false
														break
													} else if b+1 < length && !checkUnsafeUsages(fnc.cases[j].value[ind].value[b+1]) {
														flag = false
														break
													} else {
														b++
														flag = true
													}
												} else if contains(typeList, fnc.cases[j].value[ind].value[b]) {
													flag = true
												} else {
													flag = false
													break
												}
											} else if strings.Contains(fnc.cases[i].value[ind].value[a], fnc.funcName) && strings.Contains(fnc.cases[j].value[ind].value[b], fnc.funcName) {
												flag = true
											} else if contains(checkFuncName, fnc.cases[i].value[ind].value[a]) && contains(checkFuncName, fnc.cases[j].value[ind].value[b]) {
												flag = true
											} else {
												flag = false
												break
											}
										}
										b++
									}
									if !flag || a != len(fnc.cases[i].value[ind].value) || b != len(fnc.cases[j].value[ind].value) {
										flag = false
										break
									}
								} else {
									length = len(fnc.cases[j].value[ind].value)
									a, b := 0, 0
									for ; a < length; a++ {
										if !strings.EqualFold(fnc.cases[j].value[ind].value[a], fnc.cases[i].value[ind].value[b]) {
											if strings.Contains(fnc.cases[j].value[ind].value[a], " ") || strings.Contains(fnc.cases[i].value[ind].value[b], " ") ||
												strings.Contains(fnc.cases[j].value[ind].value[a], "_") || strings.Contains(fnc.cases[i].value[ind].value[b], "_") {
												flag = false
												break
											} else if contains(typeList, fnc.cases[j].value[ind].value[a]) {
												if strings.EqualFold(fnc.cases[i].value[ind].value[b], "unsafe") {
													if !checkUnsafeUsages(fnc.cases[i].value[ind].value[b]) {
														flag = false
														break
													} else if b+1 < length && !checkUnsafeUsages(fnc.cases[i].value[ind].value[b+1]) {
														flag = false
														break
													} else {
														b++
														flag = true
													}
												} else if contains(checkFuncName, fnc.cases[i].value[ind].value[b]) {
													flag = true
												} else {
													flag = false
													break
												}
											} else if strings.Contains(fnc.cases[j].value[ind].value[a], fnc.funcName) && strings.Contains(fnc.cases[i].value[ind].value[b], fnc.funcName) {
												flag = true
											} else if contains(checkFuncName, fnc.cases[i].value[ind].value[a]) && contains(checkFuncName, fnc.cases[j].value[ind].value[b]) {
												flag = true
											} else {
												flag = false
												break
											}
										}
										b++
									}
									if !flag || a != len(fnc.cases[j].value[ind].value) || b != len(fnc.cases[i].value[ind].value) {
										flag = false
										break
									}
								}
							}
						} else {
							flag = false
							break
						}
					}
					if !flag {
						break
					}
				} else {
					flag = false
					break
				}
				if flag {
					if strings.EqualFold(fnc.cases[i].value[0].value[0], "reflect") {
						substring := fnc.cases[i].value[0].value[0] + " " + fnc.cases[i].value[0].value[1]
						if !isSameString(caseList, substring) {
							caseList = append(caseList, substring)
						}
					} else {
						if !isSameString(caseList, fnc.cases[i].value[0].value[0]) {
							caseList = append(caseList, fnc.cases[i].value[0].value[0])
						}
					}
					if strings.EqualFold(fnc.cases[j].value[0].value[0], "reflect") {
						substring := fnc.cases[j].value[0].value[0] + " " + fnc.cases[j].value[0].value[1]
						if !isSameString(caseList, substring) {
							caseList = append(caseList, substring)
						}
					} else {
						if !isSameString(caseList, fnc.cases[j].value[0].value[0]) {
							caseList = append(caseList, fnc.cases[j].value[0].value[0])
						}
					}
				} else {
					break
				}
			}
			if !flag {
				break
			}
		}
		if flag && len(caseList) > 0 {
			result = append(result, caseFilteredResult{fnc.funcName, fnc.funcToken, caseList})
			caseList = []string{}
			flag = false
		} else {
			continue
		}

	}

	count := 0
	cntResult := 0
	if len(result) > 0 {
		for i := range funcList {
			if strings.EqualFold(funcList[i].funcName, result[cntResult].funcName) &&
				strings.EqualFold(funcList[i].funcToken, result[cntResult].funcToken) {
				if len(funcList[i].cases) > 0 && len(result[cntResult].caseFiltered) > 0 && len(funcList[i].cases) == len(result[cntResult].caseFiltered) {
					fmt.Println()
					fmt.Println("This function ", funcList[i].funcName, " has reused case clauses.")
					count++
					fmt.Println("Number of cases: ", len(result[cntResult].caseFiltered))
					fmt.Println(result[cntResult].caseFiltered)
					cntResult++
				}
			}
			if cntResult >= len(result) {
				break
			}
		}
		fmt.Println()
	}

	return result

}
