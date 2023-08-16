package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	_ "golang.org/x/exp/slices"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"os"
	"strings"
)

type pattern3Result struct {
	nodes    *ast.FuncDecl
	funcList []funcNamePos
}

type patternReplace struct {
	nodes      []*ast.FuncDecl
	funcRemove []funcNamePos
}

func main() {
	// filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run . - "hiddenDanger/getset.go"
	filename := os.Args[2]

	var patternObjectNode *ast.FuncDecl
	var patternSetNode *ast.FuncDecl
	var patternGetNode *ast.FuncDecl

	var pattern1 bool

	var pattern2 bool
	var pattern3 bool // pattern4 overlapped but somewhat diffrent
	var pattern4 bool

	var toReplace []pattern3Result
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
	unsafeList := buildUnsafePointerList(modListFunctions)

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

	patternB2S(filename, modListFunctions, unsafeList)
	// Check Generic Replacement
	genCheck := checkGenerics(modListFunctions, funcList, typeList)
	for s := range genCheck {
		if len(genCheck[s]) > 1 {
			//pattern1 = true
			fmt.Println()
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

	for _, fnc := range genCheck {
		if len(fnc) > 1 {
			for _, val := range modListFunctions {
				if strings.EqualFold(val.funcName, fnc[0].funcName) && val.funcToken == fnc[0].funcPos {
					patternObjectNode = patternObjectSlice(val)
					if patternObjectNode != nil {
						if !checkReplaceFunc(toReplace, fnc) {
							toReplace = append(toReplace, pattern3Result{patternObjectNode, fnc})
						}
						patternObjectNode = nil
					}
					patternSetNode = patternSet(val)
					if patternSetNode != nil {
						if !checkReplaceFunc(toReplace, fnc) {
							toReplace = append(toReplace, pattern3Result{patternSetNode, fnc})
						}
						patternSetNode = nil
					}
					patternGetNode = patternGet(val)
					if patternGetNode != nil {
						if !checkReplaceFunc(toReplace, fnc) {
							toReplace = append(toReplace, pattern3Result{patternGetNode, fnc})
						}
						patternGetNode = nil
					}
				}
			}
		}
	}
	var nodeList []*ast.FuncDecl
	var removeNameList []funcNamePos
	var toResult patternReplace
	for _, val := range toReplace {
		nodeList = append(nodeList, val.nodes)
		for _, value := range val.funcList {
			removeNameList = append(removeNameList, value)
		}

	}
	toResult = patternReplace{nodeList, removeNameList}

	if len(toResult.nodes) > 0 {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		//count := 0
		//var tmp *ast.Decl
		count := 0
		astutil.Apply(node, func(c *astutil.Cursor) bool {
			n := c.Node()

			if d, ok := n.(*ast.FuncDecl); ok {
				if checkDuplicateInFuncGen(toResult.funcRemove, d.Name.String(), d.Pos()) {
					if count == 0 {
						for _, val := range toResult.nodes {
							c.InsertBefore(val)
						}
					}
					c.Delete()
					count++

				}
			}
			return true
		}, nil)

		/*ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.File:
				for _, val := range x.Decls {
					if fnc, ok := val.(*ast.FuncDecl); ok {
						for _, value := range toResult {
							if strings.EqualFold(value.funcName0.funcName, fnc.Name.String()) {
								fnc = value.nodes
							}

						}

					}

				}

			case *ast.FuncDecl:

			}
			return true
		})*/

		// TODO
		// 파일을 굳이 저장안하고 노드만 저장한 이후에 파일 변경해도 될듯?
		newName := filename[0:len(filename)-3] + "_replaced.go"

		newFile, err := os.Create(newName)
		if err != nil {
			log.Fatal(err)
		}
		defer func(new *os.File) {
			err := new.Close()
			if err != nil {
			}
		}(newFile)

		if err := printer.Fprint(newFile, fset, node); err != nil {
			log.Fatal(err)
		}
	}

	for _, val := range modListFunctions {
		if patternGenSlice(val) != nil {
			fmt.Println("This function ", val.funcName, " has (a) function(s) with reflect.SliceHeader and Interface of return value. \nIt recommends to use Generics Slice.\nNo replacement because of Generics Replacement Suggestion.\n\n")
		}
	}
	fmt.Println()

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
