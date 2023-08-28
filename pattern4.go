package main

import (
	"fmt"
	"strings"
)

/*
written by Jang, Hyang Gi on 28.August.2023

from function patternReduce
it will be printed out No parameterised methods
*/

func patternReduce(fnc basicStr) {
	flag0, flag1, flag2 := false, false, false
	receiverParam := []string{}
	paramList := []string{}
	var caseVar string
	beforeBlockStm := true
	noParametrisationMethod := false

	for ind, elems := range fnc.value {

		if strings.Contains(elems.path, "*ast.BlockStmt") {
			beforeBlockStm = false
			flag1 = true
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
			if beforeBlockStm && (strings.Contains(elems.path, "*ast.Field") || strings.Contains(elems.path, "*ast.StarExpr -> *ast.SelectorExpr") ||
				strings.Contains(elems.path, "*ast.InterfaceType")) {
				if strings.Contains(elems.path, "*ast.InterfaceType") {
					if len(elems.path) == 0 {
						paramList = append(paramList, "Interface")
					}
				} else if len(elems.path) != 0 {
					for _, value := range elems.value {
						paramList = append(paramList, value)
					}
				}
			}
		}

		// Check Case Clause
		if flag1 {
			if strings.Contains(elems.path, "*ast.CaseClause") && len(elems.value) == 1 {
				if !flag2 {
					caseVar = elems.value[0]
					flag2 = true
					continue
				}
			}
		}

		// If there is a case variable in the astValue?
		if flag2 {
			if istValueInArr(elems.value, caseVar) && len(elems.value) > 1 {
				if strings.Contains(elems.value[1], caseVar) && contains(paramList, elems.value[0]) {
					if !noParametrisationMethod {
						noParametrisationMethod = true
						break
					}
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
