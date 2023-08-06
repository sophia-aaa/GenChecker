package main

import (
	"fmt"
	_ "golang.org/x/exp/slices"
	"os"
)

func main() {
	// filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run gen.go - test.go
	filename := os.Args[2]

	var pattern1 bool
	var pattern2 bool
	var pattern3 bool // pattern4 overlapped but somewhat diffrent
	var pattern4 bool

	listFunctions := buildAstDataStr(filename)

	var funcList []string
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			funcList = append(funcList, listFunctions[s].funcName)
		}
	}

	// check leaf of SelectorExpr and unsafe Pointer
	modListFunctions := checkSelectorExpr(listFunctions)
	createTextFile(filename, modListFunctions)
	unsafeList := buildUnsafeList(modListFunctions)

	if len(unsafeList) > 0 {
		fmt.Println("This function contains unsafe.Pointer:")
		for s := range unsafeList {
			if unsafeList[s].funcName != "" {
				fmt.Print(unsafeList[s].funcName, " ")
			}
		}

		fmt.Println()
	}

	// Check Generic Replacement
	genCheck := checkGenerics(modListFunctions, funcList, typeList)
	for s := range genCheck {
		if len(genCheck[s]) > 1 {
			pattern1 = true
			fmt.Print("These functions have a same structure and the code are reused: ")
			for _, value := range genCheck[s] {
				fmt.Print(value, " ")
			}
			fmt.Println()
		}
	}

	checkDataFunc := checkData(modListFunctions)
	if len(checkDataFunc) > 0 {
		pattern2 = true
		fmt.Print("\nThere exists (a) function(s) with reflect.SliceHeader and Interface of return value. It recommends to use Generics Slice : ")
		for _, val := range checkDataFunc {
			fmt.Print(val, " ")
		}
		fmt.Println()
	}

	// This variable is for checking switch statement
	existsSwitch, caseList := checkSwitchStatement(filename, modListFunctions)
	if existsSwitch {
		for _, val := range caseList {
			fmt.Println(val.funcName, "\t", val.value)
		}

	}

	if pattern1 {
		fmt.Println("pattern1")
		//replacePattern(filename, 1)
	}
	if pattern2 {
		fmt.Println("pattern2")
		replacePattern(filename, 2)
	}
	if pattern3 {
		fmt.Println("pattern3")
	}
	if pattern4 {
		fmt.Println("pattern4")
	}
}
