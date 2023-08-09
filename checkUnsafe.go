package main

import (
	"go/token"
	"strings"
)

type unsafeCheck struct {
	funcName  string
	funcToken token.Pos // to distinguish same name function
}

func checkUnsafeUsages(str string) bool {
	return contains([]string{"unsafe", "Pointer"}, str)
}

func checkDuplicateInUnsafe(list []unsafeCheck, funcName string, funcPos token.Pos) bool {
	for ind := range list {
		if strings.EqualFold(funcName, list[ind].funcName) && (rune(funcPos) == rune(list[ind].funcToken)) {
			return true
		}
	}
	return false
}

func buildUnsafeList(modListFunctions []basicStr) []unsafeCheck {
	var unsafeList []unsafeCheck

	for ind := range modListFunctions {
		firstBool := false
		for _, val := range modListFunctions[ind].value {
			for i := range val.value {
				if strings.EqualFold(val.value[i], "unsafe") {
					firstBool = true
				} else if firstBool && strings.EqualFold(val.value[i], "Pointer") {
					if !checkDuplicateInUnsafe(unsafeList, modListFunctions[ind].funcName, modListFunctions[ind].funcToken) {
						if !strings.EqualFold(modListFunctions[ind].funcName, "") {
							unsafeList = append(unsafeList, unsafeCheck{modListFunctions[ind].funcName, modListFunctions[ind].funcToken})
						}
					}
					firstBool = false
				}
			}
		}
	}

	return unsafeList
}
