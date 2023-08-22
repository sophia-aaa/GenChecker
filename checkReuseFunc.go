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

	for i := 0; i < len(listFunctions); i++ {
		if listFunctions[i].funcName != "" {
			/*
				toThink Pervasiveness
				var funcNameI []string
				funcNameI = funcNameDivider(funcNameI, listFunctions[i].funcName)
			*/
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
					if len(listFunctions[i].value) == len(listFunctions[j].value) {
						// Compare details between functions
						for idx := range listFunctions[i].value {
							if strings.EqualFold(listFunctions[i].value[idx].path, listFunctions[j].value[idx].path) {
								if len(listFunctions[i].value[idx].value) == len(listFunctions[j].value[idx].value) {
									for idxValue := range listFunctions[i].value[idx].value {
										if strings.EqualFold(listFunctions[i].value[idx].value[idxValue], listFunctions[j].value[idx].value[idxValue]) {
											flag = true
										} else {
											// values are not type or word, but string values
											if strings.Contains(listFunctions[i].value[idx].value[idxValue], " ") || strings.Contains(listFunctions[j].value[idx].value[idxValue], " ") ||
												strings.Contains(listFunctions[i].value[idx].value[idxValue], "_") || strings.Contains(listFunctions[j].value[idx].value[idxValue], "_") {
												flag = false
												break
											} else {
												// compare ast.Ident Name
												if strings.EqualFold(listFunctions[i].funcName, listFunctions[i].value[idx].value[idxValue]) && strings.EqualFold(listFunctions[j].funcName, listFunctions[j].value[idx].value[idxValue]) {
													flag = true
												} else if contains(typeList, listFunctions[i].value[idx].value[idxValue]) && contains(typeList, listFunctions[j].value[idx].value[idxValue]) {
													flag = true
												} else {
													flag = false
													break
												}
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
								if !flag {
									break
								}
							} else if (strings.Contains(listFunctions[i].value[idx].path, "*ast.SelectorExpr") && !strings.Contains(listFunctions[j].value[idx].path, "*ast.SelectorExpr")) ||
								(!strings.Contains(listFunctions[i].value[idx].path, "*ast.SelectorExpr") && strings.Contains(listFunctions[j].value[idx].path, "*ast.SelectorExpr")) {
								length := 0
								if (len(listFunctions[i].value[idx].value) == 0) || (len(listFunctions[j].value[idx].value) == 0) {
									flag = false
									break
								} else {
									if len(listFunctions[i].value[idx].value) <= len(listFunctions[j].value[idx].value) {
										length = len(listFunctions[i].value[idx].value) // will be index a
										a, b := 0, 0
										for ; a < length; a++ {

											if !strings.EqualFold(listFunctions[i].value[idx].value[a], listFunctions[j].value[idx].value[b]) {
												if strings.Contains(listFunctions[i].value[idx].value[a], " ") || strings.Contains(listFunctions[j].value[idx].value[b], " ") ||
													strings.Contains(listFunctions[i].value[idx].value[a], "_") || strings.Contains(listFunctions[j].value[idx].value[b], "_") {
													flag = false
													break
												} else if contains(typeList, listFunctions[i].value[idx].value[a]) {
													if !checkUnsafeUsages(listFunctions[j].value[idx].value[b]) {
														flag = false
														break
													} else if b+1 < length && !checkUnsafeUsages(listFunctions[j].value[idx].value[b+1]) {
														flag = false
														break
													} else {
														b++
													}
												} else if strings.EqualFold(listFunctions[i].funcName, listFunctions[i].value[idx].value[a]) && strings.EqualFold(listFunctions[j].funcName, listFunctions[j].value[idx].value[b]) {
													flag = true
												} else if contains(typeList, listFunctions[i].value[idx].value[a]) && contains(typeList, listFunctions[j].value[idx].value[b]) {
													flag = true
												} else {
													flag = false
													break
												}
											}
											b++
										}
										if b != len(listFunctions[j].value[idx].value) {
											flag = false
											break
										} else {
											flag = true
										}
									} else {
										length = len(listFunctions[j].value[idx].value)
										a, b := 0, 0
										for ; a < length; a++ {
											if strings.EqualFold(listFunctions[i].funcName, "ExtraData") {
												fmt.Println("landing 10")
											}
											if !strings.EqualFold(listFunctions[j].value[idx].value[a], listFunctions[i].value[idx].value[b]) {
												if contains(typeList, listFunctions[j].value[idx].value[a]) {
													if strings.Contains(listFunctions[j].value[idx].value[a], " ") || strings.Contains(listFunctions[i].value[idx].value[b], " ") ||
														strings.Contains(listFunctions[j].value[idx].value[a], "_") || strings.Contains(listFunctions[i].value[idx].value[b], "_") {
														flag = false
														break
													} else if !checkUnsafeUsages(listFunctions[i].value[idx].value[b]) {
														flag = false
														break
													} else if b+1 < length && !checkUnsafeUsages(listFunctions[i].value[idx].value[b+1]) {
														flag = false
														break
													} else {
														b++
													}
												} else if strings.EqualFold(listFunctions[i].funcName, listFunctions[i].value[idx].value[b]) && strings.EqualFold(listFunctions[j].funcName, listFunctions[j].value[idx].value[a]) {
													flag = true
												} else if contains(typeList, listFunctions[i].value[idx].value[b]) && contains(typeList, listFunctions[j].value[idx].value[a]) {
													flag = true
												} else {
													flag = false
													break
												}
											}
											b++
										}
										if b != len(listFunctions[i].value[idx].value) {
											flag = false
											break
										} else {
											flag = true
										}

									}
								}
								flag = true
							} else {
								flag = false
								break
							}
						}
						if !flag {
							continue
						}
					} else {
						flag = false
						continue // continue the progress comparing a next function to the compared one
					}
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
