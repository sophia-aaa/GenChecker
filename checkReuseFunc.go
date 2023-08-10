package main

import (
	"fmt"
	"go/token"
	"strings"
)

type funcNamePos struct {
	funcName string
	funcPos  token.Pos
}

func checkGenerics(listFunctions []basicStr, funcList []string, typeList []string) [][]funcNamePos {

	var genFunc []funcNamePos
	var genCheck [][]funcNamePos
	flag := false

	bool1 := false
	bool2 := false
	for i := 0; i < len(listFunctions); i++ {
		if listFunctions[i].funcName != "" {
			/*
				toThink Pervasiveness
				var funcNameI []string
				funcNameI = funcNameDivider(funcNameI, listFunctions[i].funcName)
			*/

			// todo remove
			if strings.EqualFold(listFunctions[i].funcName, "RemoveColumn") {
				bool1 = true
				fmt.Println(bool1)

			} else {
				bool1 = false
			}

			if len(genCheck) != 0 {
				if !contains2D(genCheck, listFunctions[i].funcName, listFunctions[i].funcToken) {
					genFunc = []funcNamePos{{listFunctions[i].funcName, listFunctions[i].funcToken}}
				} else {
					continue
				}
			} else {
				genFunc = []funcNamePos{{listFunctions[i].funcName, listFunctions[i].funcToken}}
			}
			for j := i + 1; j < len(listFunctions); j++ {
				if listFunctions[j].funcName != "" {
					/*
						toThink Pervasiveness
						ambiguous how far acceptable for ident name
						var funcNameJ []string
						funcNameJ = funcNameDivider(funcNameJ, listFunctions[j].funcName)
						for _, val := range funcNameI {
							if !contains(funcNameJ, val) {
								funcNameJ = append(funcNameJ, val)
							}
						}
					*/
					// todo remove
					if bool1 && strings.EqualFold(listFunctions[j].funcName, "SetLines") {
						bool2 = true
						fmt.Println(bool2, " ", listFunctions[j].funcName)
					} else {
						bool2 = false
					}

					if len(listFunctions[i].value) == len(listFunctions[j].value) {
						// Compare details between functions
						for idx := range listFunctions[i].value {
							if bool1 && bool2 {
								fmt.Println(listFunctions[i].value[idx].path, listFunctions[j].value[idx].path)
								fmt.Println(listFunctions[i].value[idx].value, listFunctions[j].value[idx].value)
							}
							if strings.Compare(listFunctions[i].value[idx].path, listFunctions[j].value[idx].path) == 0 {
								if len(listFunctions[i].value[idx].value) == len(listFunctions[j].value[idx].value) {
									for idxValue := range listFunctions[i].value[idx].value {
										if bool1 && bool2 {
											fmt.Println(strings.Compare(listFunctions[i].value[idx].value[idxValue], listFunctions[j].value[idx].value[idxValue]))
										}
										if strings.Compare(listFunctions[i].value[idx].value[idxValue], listFunctions[j].value[idx].value[idxValue]) == 0 {
											flag = true
										} else {
											if bool1 && bool2 {
												fmt.Println((contains(funcList, listFunctions[i].value[idx].value[idxValue]) && contains(funcList, listFunctions[j].value[idx].value[idxValue])) ||
													(contains(typeList, listFunctions[i].value[idx].value[idxValue]) && contains(typeList, listFunctions[j].value[idx].value[idxValue])))
											}
											// compare ast.Ident Name
											if (contains(funcList, listFunctions[i].value[idx].value[idxValue]) && contains(funcList, listFunctions[j].value[idx].value[idxValue])) ||
												// to think Pervasiveness (contains(funcNameJ, listFunctions[i].value[idx].value[idxValue]) && contains(funcNameJ, listFunctions[j].value[idx].value[idxValue])) ||
												contains(typeList, listFunctions[i].value[idx].value[idxValue]) && contains(typeList, listFunctions[j].value[idx].value[idxValue]) {
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
								// ********************************* TODO how????????
								if !flag {
									break
								}
								if bool1 && bool2 {
									fmt.Println("1 flag is ", flag)
								}
								// TODO
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
								if bool1 && bool2 {
									fmt.Println("2 flag is ", flag)
								}
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
								if bool1 && bool2 {
									fmt.Println("3 flag is ", flag)
								}
							} else if strings.Contains(listFunctions[i].value[idx].path, "*ast.SelectorExpr") {
								// listFunctions[i].value[idx] looks like "... -> *ast.SelectorExpr [... unsafe Pointer]"
								// and listFunctions[j].value[idx] looks like " ... [TYPE]"
								if len(listFunctions[j].value[idx].value) > 0 && !contains(typeList, listFunctions[j].value[idx].value[0]) {
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
								if bool1 && bool2 {
									fmt.Println("4 flag is ", flag)
								}
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
								if bool1 && bool2 {
									fmt.Println("5 flag is ", flag)
								}
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
				if bool1 && bool2 {
					fmt.Println("6 flag is ", flag)
				}
				if flag == true {
					genFunc = append(genFunc, funcNamePos{listFunctions[j].funcName, listFunctions[j].funcToken})
					flag = false
				}
			}
			genCheck = append(genCheck, genFunc)
			//fmt.Println(genFunc)
		}
	}
	return genCheck
}
