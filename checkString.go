package main

import (
	"go/token"
	"reflect"
	"strconv"
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

func istValueInArr(strArr []string, str string) bool {
	for _, val := range strArr {
		if val == str {
			return true
		}
		if strings.Contains(strings.ToLower(val), strings.ToLower(str)) {
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

func funcNameDivider(str string) []string {
	var funcWordList []string
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

func checkDuplicateInFunc(list []basicFunc, funcName string, funcPos token.Pos) bool {
	for ind := range list {
		if strings.EqualFold(funcName, list[ind].funcName) && (funcPos == list[ind].funcToken) {
			return true
		}
	}
	return false
}

func checkDuplicateInFuncGen(list []funcNamePos, funcName string, funcPos token.Pos) bool {
	for ind := range list {
		if strings.EqualFold(funcName, list[ind].funcName) && (funcPos == list[ind].funcPos) {
			return true
		}
	}
	return false
}

func checkDuplicateInFuncString(list []basicFunc, funcName string, funcPos string) bool {
	funcPosInt, _ := strconv.Atoi(funcPos)
	for ind := range list {
		if strings.EqualFold(funcName, list[ind].funcName) && funcPosInt == int(list[ind].funcToken) {
			return true
		}
	}
	return false
}

func checkReplaceFunc(lists []pattern3Result, nameList []funcNamePos) bool {
	for _, val := range lists {
		if reflect.DeepEqual(val, nameList) {
			return true
		}
	}
	return false
}
