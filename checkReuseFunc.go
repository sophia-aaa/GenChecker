package main

import "strings"

func checkGenerics(listFunctions []basicStr, funcList []string, typeList []string) [][]string {
	var genCheck [][]string
	var genFunc []string
	flag := false

	for i := 0; i < len(listFunctions); i++ {
		if listFunctions[i].funcName != "" {
			if len(genCheck) != 0 {
				if !contains2D(genCheck, listFunctions[i].funcName) {
					genFunc = []string{listFunctions[i].funcName}
				} else {
					continue
				}
			} else {
				genFunc = []string{listFunctions[i].funcName}
			}
			for j := i + 1; j < len(listFunctions); j++ {
				if listFunctions[j].funcName != "" {
					if len(listFunctions[i].value) == len(listFunctions[j].value) {
						// Compare details between functions
						for idx := range listFunctions[i].value {
							if strings.Compare(listFunctions[i].value[idx].path, listFunctions[j].value[idx].path) == 0 {
								if len(listFunctions[i].value[idx].value) == len(listFunctions[j].value[idx].value) {
									for idxValue := range listFunctions[i].value[idx].value {
										if strings.Compare(listFunctions[i].value[idx].value[idxValue], listFunctions[j].value[idx].value[idxValue]) == 0 {
											flag = true
										} else {
											if (contains(funcList, listFunctions[i].value[idx].value[idxValue]) && contains(funcList, listFunctions[j].value[idx].value[idxValue])) ||
												(contains(typeList, listFunctions[i].value[idx].value[idxValue]) && contains(typeList, listFunctions[j].value[idx].value[idxValue])) {
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
							} else if listFunctions[i].value[idx].path+" -> *ast.SelectorExpr" == listFunctions[j].value[idx].path {
								// This is a special case for modified elem list
								// listFunctions[j].value[idx].path is listFunctions[i].value[idx].path " -> *ast.SelectorExpr"
								// listFunctions[j].value[idx].value must look like [... unsafe Pointer] and
								// listFunctions[i].value[idx].value must have type variable in its elem value so like [... type]
								len1 := len(listFunctions[i].value[idx].value)
								len2 := len(listFunctions[j].value[idx].value)
								if len2 > 1 && !contains(typeList, listFunctions[i].value[idx].value[len1-1]) ||
									(!checkUnsafeUsages(listFunctions[j].value[idx].value[len2-1]) &&
										!checkUnsafeUsages(listFunctions[j].value[idx].value[len2-2])) {
									flag = false
									break
								}
								flag = true
							} else if listFunctions[j].value[idx].path+" -> *ast.SelectorExpr" == listFunctions[i].value[idx].path {
								// This is a special case for modified elem list
								// listFunctions[i].value[idx].path is listFunctions[j].value[idx].path " -> *ast.SelectorExpr"
								// listFunctions[i].value[idx].value must look like [... unsafe Pointer] and
								// listFunctions[j].value[idx].value must have type variable in its elem value so like [... type]
								len1 := len(listFunctions[j].value[idx].value)
								len2 := len(listFunctions[i].value[idx].value)
								if len2 > 2 && !contains(typeList, listFunctions[j].value[idx].value[len1-1]) ||
									(!checkUnsafeUsages(listFunctions[i].value[idx].value[len2-1]) &&
										!checkUnsafeUsages(listFunctions[i].value[idx].value[len2-2])) {
									flag = false
									break
								}
								flag = true
							} else if strings.Contains(listFunctions[i].value[idx].path, "*ast.SelectorExpr") {
								// listFunctions[i].value[idx] looks like "... -> *ast.SelectorExpr [... unsafe Pointer]"
								// and listFunctions[j].value[idx] looks like " ... [TYPE]"
								if !contains(typeList, listFunctions[j].value[idx].value[0]) {
									flag = false
									break
								}
								for _, val := range listFunctions[i].value[idx].value {
									if !checkUnsafeUsages(val) {
										flag = false
										break
									}
								}
								flag = true
							} else if strings.Contains(listFunctions[j].value[idx].path, "*ast.SelectorExpr") {
								// listFunctions[i].value[idx] looks like " ... [TYPE]"
								// and listFunctions[j].value[idx] looks like "... -> *ast.SelectorExpr [... unsafe Pointer]"
								if !contains(typeList, listFunctions[i].value[idx].value[0]) {
									flag = false
									break
								}
								for _, val := range listFunctions[j].value[idx].value {
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
					}
				}
				if flag == true {
					genFunc = append(genFunc, listFunctions[j].funcName)
					flag = false
				}
			}
			genCheck = append(genCheck, genFunc)
		}
	}
	return genCheck
}
