package main

import (
	"strings"
)

func checkUnsafeUsages(str string) bool {
	return contains([]string{"unsafe", "Pointer"}, str)
}

func buildUnsafeList(modListFunctions []basicStr) []basicFunc {
	var unsafeList []basicFunc

	for ind := range modListFunctions {
		firstBool := false
		for _, val := range modListFunctions[ind].value {
			for i := range val.value {
				if strings.EqualFold(val.value[i], "unsafe") {
					firstBool = true
				} else if firstBool && strings.EqualFold(val.value[i], "Pointer") {
					if !checkDuplicateInFunc(unsafeList, modListFunctions[ind].funcName, modListFunctions[ind].funcToken) {
						if !strings.EqualFold(modListFunctions[ind].funcName, "") {
							unsafeList = append(unsafeList, basicFunc{modListFunctions[ind].funcName, modListFunctions[ind].funcToken})
						}
					}
					firstBool = false
				}
			}
		}
	}

	return unsafeList
}
