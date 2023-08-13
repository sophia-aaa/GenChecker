package main

import (
	"fmt"
	_ "golang.org/x/exp/slices"
	"os"
	"strings"
)

func main() {
	// filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run . - "dataset/getset.go"
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

	fmt.Println(len(funcList), " Function list: ")
	for ind, val := range funcList {
		if len(funcList) == 1 {
			fmt.Print("{ ", val, " }")
		} else if ind == len(funcList)-1 {
			fmt.Println(val, "}")
		} else if ind == 0 {
			fmt.Print("{ ", val, ", ")
		} else {
			fmt.Print(val, ", ")
		}
	}
	fmt.Println()
	fmt.Println()

	var methodList []string
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			if len(listFunctions[s].value) > 0 &&
				strings.Contains(listFunctions[s].value[0].path, "*ast.FieldList -> *ast.Field") {
				methodList = append(methodList, listFunctions[s].funcName)
			}
		}
	}

	fmt.Println(len(methodList), " Method list: ")
	for ind, val := range methodList {
		if len(methodList) == 1 {
			fmt.Print("{ ", val, " }")
		} else if ind == len(methodList)-1 {
			fmt.Println(val, "}")
		} else if ind == 0 {
			fmt.Print("{ ", val, ", ")
		} else {
			fmt.Print(val, ", ")
		}
	}
	fmt.Println()
	fmt.Println()

	// check leaf of SelectorExpr and unsafe Pointer
	modListFunctions := checkSelectorExpr(listFunctions)
	createTextFile(filename, modListFunctions)
	unsafeList := buildUnsafeList(modListFunctions)

	if len(unsafeList) > 0 {
		fmt.Println(len(unsafeList), " This function contains unsafe.Pointer:")
		for s := range unsafeList {
			if unsafeList[s].funcName != "" {
				if len(unsafeList) == 1 {
					fmt.Print("{ ", unsafeList[s].funcName, " }")
				} else if s == len(unsafeList)-1 {
					fmt.Println(unsafeList[s].funcName, "}")
				} else if s == 0 {
					fmt.Print("{ ", unsafeList[s].funcName, ", ")
				} else {
					fmt.Print(unsafeList[s].funcName, ", ")
				}
			}
		}
		fmt.Println()
		fmt.Println()
	}

	// Check Generic Replacement
	genCheck := checkGenerics(modListFunctions, funcList, typeList)
	for s := range genCheck {
		if len(genCheck[s]) > 1 {
			pattern1 = true
			fmt.Print("These functions have a same structure and the code are reused:\n")
			for ind, val := range genCheck[s] {
				if len(genCheck[s]) == 1 {
					fmt.Print("{ ", val.funcName, " }")
				} else if ind == len(genCheck[s])-1 {
					fmt.Println(val.funcName, "}")
				} else if ind == 0 {
					fmt.Print("{ ", val.funcName, ", ")
				} else {
					fmt.Print(val.funcName, ", ")
				}
			}
			fmt.Println()
		}
	}

	checkDataFunc := checkDataPattern(modListFunctions)
	if len(checkDataFunc) > 0 {
		pattern2 = true
		fmt.Print("\nThere exists (a) function(s) with reflect.SliceHeader and Interface of return value. It recommends to use Generics Slice : ")
		for ind, val := range checkDataFunc {
			if len(checkDataFunc) == 1 {
				fmt.Print("{ ", val, " }")
			} else if ind == len(checkDataFunc)-1 {
				fmt.Println(val, "}")
			} else if ind == 0 {
				fmt.Print("{ ", val, ", ")
			} else {
				fmt.Print(val, ", ")
			}

		}
		fmt.Println()
	}

	// This variable is for checking switch statement
	existsSwitch, caseList := checkSwitchStatement(filename, modListFunctions)
	if existsSwitch {
		if len(caseList) > 0 {
			fmt.Println("This function has switch statement: ")
			for ind, val := range caseList {
				if len(caseList) == 1 {
					fmt.Print("{ ", val.funcName, " }")
				} else if ind == len(caseList)-1 {
					fmt.Println(val.funcName, "}")
				} else if ind == 0 {
					fmt.Print("{ ", val.funcName, ", ")
				} else {
					fmt.Print(val.funcName, ", ")
				}
			}
		}
	}

	if pattern1 {
		fmt.Println("pattern1")
		//replacePattern(filename, 1)
	}
	if pattern2 {
		fmt.Println("pattern2")
		//replacePattern(filename, 2)
	}
	if pattern3 {
		fmt.Println("pattern3")
	}
	if pattern4 {
		fmt.Println("pattern4")
	}
}
