package main

import (
	"go/token"
	"strings"
)

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

func contains2D(strArr [][]funcNamePos, str string, pos token.Pos) bool {
	for i := 0; i < len(strArr); i++ {
		for j := 0; j < len(strArr[i]); j++ {
			if strings.EqualFold(strArr[i][j].funcName, str) && strArr[i][j].funcPos == pos {
				return true
			}
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

func funcNameDivider(funcWordList []string, str string) []string {
	funcWord := ""
	for _, val := range str {
		if 'A' <= val && val <= 'Z' {
			if !strings.EqualFold(funcWord, "") {
				funcWordList = append(funcWordList, funcWord)
			}
			funcWord = string(val)
		} else {
			funcWord += string(val)
		}
	}
	return funcWordList
}
