package main

import "strings"

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

func contains2D(strArr [][]string, str string) bool {
	for i := 0; i < len(strArr); i++ {
		for j := 0; j < len(strArr[i]); j++ {
			if strArr[i][j] == str {
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
