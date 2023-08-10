package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var Tree2str []string
var Tree2cases []basicCases

type elem struct {
	path  string
	value []string
}
type basicCases struct {
	depth int
	value string
}

// Basic structure for every function in the input file
type basicStr struct {
	funcName  string
	funcToken token.Pos // to distinguish same name function
	value     []elem
}

type basicCaseStr struct {
	caseName string
	value    []elem
}

func buildAstDataStr(filename string) []basicStr {
	// Commentgroups and comments are ignored
	// because they interfere with comparing the content of the code.
	// Ellipsis values are position lines rather than real values,
	// so this function does not express them.

	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var elem_list []elem
	nameFunction := ""
	var tokPos token.Pos
	astNode := ""
	var listFunctions []basicStr
	var astValue []string
	genDeclBool := false
	fileStart := false

	ast.Inspect(astTree, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.CommentGroup, *ast.Comment:
			astNode += ""
		case *ast.ImportSpec:
			if astNode == "" {
				astNode += reflect.TypeOf(x).String()
			} else {
				astNode += " -> " + reflect.TypeOf(x).String()
			}
			fileStart = true
		case *ast.GenDecl:
			if !genDeclBool {
				if astNode != "" || len(astValue) != 0 { // if a new root meets
					elem_list = append(elem_list, elem{astNode, astValue})
					listFunctions = append(listFunctions, basicStr{nameFunction, tokPos, elem_list})
					elem_list = []elem{}
					astNode = ""
					astValue = []string{}
				}
				nameFunction = ""
				tokPos = x.Pos()
			} else { // within function
				if len(astValue) != 0 {
					elem_list = append(elem_list, elem{astNode, astValue})
					astNode = ""
					astValue = []string{}
				}
			}
			if astNode == "" {
				astNode += reflect.TypeOf(x).String()
			} else {
				astNode += " -> " + reflect.TypeOf(x).String()
			}
		case *ast.DeclStmt:
			genDeclBool = true
		case *ast.FuncDecl:
			//	fmt.Println(x, "\t\t", reflect.TypeOf(x).String())
			//	fmt.Println(fset.Position(x.Pos()), fset.Position(x.End()))
			if fileStart {
				if astNode != "" || len(astValue) != 0 {
					elem_list = append(elem_list, elem{astNode, astValue})
					elem_list = []elem{}
					astNode = ""
					astValue = []string{}
				}
				listFunctions = append(listFunctions, basicStr{nameFunction, tokPos, elem_list})
				fileStart = false
				elem_list = []elem{}
				astNode = ""
				astValue = []string{}
			} else if astNode != "" || len(astValue) != 0 { // if a new root meets
				elem_list = append(elem_list, elem{astNode, astValue})
				listFunctions = append(listFunctions, basicStr{nameFunction, tokPos, elem_list})
				elem_list = []elem{}
				astNode = ""
				astValue = []string{}
			}
			nameFunction = x.Name.String()
			tokPos = x.Pos()
		//fmt.Println(x.Name.String(), "\t", rune(x.Pos()))
		case *ast.AssignStmt:
			if astNode != "" || len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += reflect.TypeOf(x).String()
			} else {
				astNode += " -> " + reflect.TypeOf(x).String()
			}
			astValue = append(astValue, x.Tok.String())
		case *ast.BinaryExpr:
			if astNode != "" || len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += reflect.TypeOf(x).String()
			} else {
				astNode += " -> " + reflect.TypeOf(x).String()
			}
			astValue = append(astValue, x.Op.String())
			genDeclBool = false
		case *ast.UnaryExpr:
			if astNode != "" || len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += reflect.TypeOf(x).String()
			} else {
				astNode += " -> " + reflect.TypeOf(x).String()
			}
			astValue = append(astValue, x.Op.String())
			genDeclBool = false
		case *ast.Ident:
			//fmt.Println(fset.Position(x.Pos()), reflect.TypeOf(x).String(), "\t", x.Name)
			astValue = append(astValue, x.Name)
			genDeclBool = false
		case *ast.BasicLit:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += reflect.TypeOf(x).String()
			} else {
				astNode += " -> " + reflect.TypeOf(x).String()
			}
			str := x.Kind.String() + " " + x.Value
			astValue = append(astValue, str)
			str = ""
			elem_list = append(elem_list, elem{astNode, astValue})
			astNode = ""
			astValue = []string{}
		case *ast.InterfaceType:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			astNode += reflect.TypeOf(x).String()
			if len(x.Methods.List) > 0 {
				for _, val := range x.Methods.List {
					if len(val.Names) > 0 {
						for _, valName := range val.Names {
							astValue = append(astValue, valName.String())
						}
					}
				}
			}
			elem_list = append(elem_list, elem{astNode, astValue})
			astNode = ""
			astValue = []string{}
			// I couldn't find a way to express "incomplete", so I customised it like this.
			if x.Incomplete == true {
				elem_list = append(elem_list, elem{"*ast.InterfaceType.Incomplete", []string{"true"}})
			} else {
				elem_list = append(elem_list, elem{"*ast.InterfaceType.Incomplete", []string{"false"}})
			}
		case *ast.CaseClause, *ast.SwitchStmt:
			// *ast.CaseClause -> *ast.SelectorExpr
			if !strings.EqualFold(astNode, "") || len(astValue) > 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += reflect.TypeOf(x).String()
			} else {
				astNode += " -> " + reflect.TypeOf(x).String()
			}
		default:
			if x != nil {
				//fmt.Println(fset.Position(x.Pos()), reflect.TypeOf(x).String())
				//			fmt.Println(fset.Position(x.Pos()))
				if len(astValue) != 0 {
					elem_list = append(elem_list, elem{astNode, astValue})
					astNode = ""
					astValue = []string{}
				}
				if astNode == "" {
					astNode += reflect.TypeOf(x).String()
				} else {
					astNode += " -> " + reflect.TypeOf(x).String()
				}
			}
		}
		return true
	})
	if nameFunction != "" { // add last function
		elem_list = append(elem_list, elem{astNode, astValue})
		listFunctions = append(listFunctions, basicStr{nameFunction, tokPos, elem_list})
		elem_list = []elem{}
		astNode = ""
		astValue = []string{}
	}

	/*for ind := range listFunctions {
		fmt.Println(listFunctions[ind].funcName, listFunctions[ind].value)
	}*/
	return listFunctions
}

// Reference for type visitor int and func Visit : https://golangdocs.com/golang-ast-package
type visitor int

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	var str string       // for txt file
	var str4cases string // for basic structure
	var depth int
	// int(v) is a depth of a current node
	if reflect.TypeOf(n).String() == "*ast.Ident" {
		str = fmt.Sprintf("%s%d ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)

		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), n)
		str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), n)
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
	} else if reflect.TypeOf(n).String() == "*ast.AssignStmt" {
		str = fmt.Sprintf("%s%d ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)

		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), n.(*ast.AssignStmt).Tok.String())
		str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), n.(*ast.AssignStmt).Tok.String())
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
	} else if reflect.TypeOf(n).String() == "*ast.CaseClause" {
		str = fmt.Sprintf("%s%d ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)
		if n.(*ast.CaseClause).List == nil {
			str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), "default")
			str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), "default")
			Tree2str = append(Tree2str, str)
			Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
		} else {
			str = fmt.Sprintf("%T\n", n)
			str4cases = fmt.Sprintf("%T", n)
			Tree2str = append(Tree2str, str)
			Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
		}
	} else if reflect.TypeOf(n).String() == "*ast.BasicLit" {
		str = fmt.Sprintf("%s%d ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)

		str = fmt.Sprintf("%s %v %v\n", reflect.TypeOf(n).String(), n.(*ast.BasicLit).Kind.String(), n.(*ast.BasicLit).Value)
		str4cases = fmt.Sprintf("%s %v %v", reflect.TypeOf(n).String(), n.(*ast.BasicLit).Kind.String(), n.(*ast.BasicLit).Value)
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
	} else if reflect.TypeOf(n).String() == "*ast.BinaryExpr" {
		str = fmt.Sprintf("%s%d ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)

		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), n.(*ast.BinaryExpr).Op.String())
		str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), n.(*ast.BinaryExpr).Op.String())
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
	} else if reflect.TypeOf(n).String() == "*ast.UnaryExpr" {
		str = fmt.Sprintf("%s%d ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)

		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), n.(*ast.UnaryExpr).Op.String())
		str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), n.(*ast.UnaryExpr).Op.String())
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
	} else if reflect.TypeOf(n).String() == "*ast.FuncDecl" {
		funcName := fmt.Sprintf("%v", n)
		valExtracted := strings.Split(funcName, " ")
		str = fmt.Sprintf("%s%d ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)
		if len(valExtracted) > 2 {
			str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), valExtracted[2])
			str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), valExtracted[2])
		} else { // Function has no name!
			str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), "No Name")
			str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), "No Name")
		}
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
	} else {
		str = fmt.Sprintf("%s%d ", strings.Repeat("\t", int(v)), int(v))
		Tree2str = append(Tree2str, str)
		depth = int(v)
		str = fmt.Sprintf("%T\n", n)
		str4cases = fmt.Sprintf("%T", n)
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
	}
	return v + 1
}

type testFuncCase struct {
	funcName string
	caseCol  []testCase
}

type testCase struct {
	caseName  string
	caseValue []basicCases
}

func buildAstCaseStr(Tree2cases []basicCases) []checkCases {
	var funcName string
	var path string
	var astValueList []string
	var depthFirstCase int
	var caseName string
	var elemList []elem
	var caseCheck []basicCaseStr
	var funcCheck []checkCases

	var cases []testCase
	var caseValues []basicCases
	var resultFunc []testFuncCase
	flag := true
	count := 0
	caseBegin := false

	for _, val := range Tree2cases {
		if val.depth == 1 && strings.Contains(val.value, "*ast.FuncDecl") { // function begins
			if funcName != "" {

				fmt.Println(funcName, " and path is: ", path, "\t", caseName)
				elemList = append(elemList, elem{path, astValueList})
				caseCheck = append(caseCheck, basicCaseStr{caseName, elemList})
				funcCheck = append(funcCheck, checkCases{funcName, caseCheck})
				//fmt.Println(path, " ", astValueList)
				path = ""
				caseName = ""
				astValueList = []string{}
				elemList = []elem{}
				caseCheck = []basicCaseStr{}
				count = 0
				depthFirstCase = 0

				// new
				cases = append(cases, testCase{caseName, caseValues})
				resultFunc = append(resultFunc, testFuncCase{funcName, cases})
				caseValues = []basicCases{}
				cases = []testCase{}

			}
			valExtracted := strings.Split(val.value, " ")
			funcName = valExtracted[1]
			//fmt.Println()
			//fmt.Println()
			//fmt.Println(">> ", funcName, " <<")
			flag = true
			continue
		}
		if flag && (strings.Contains(val.value, "*ast.CaseClause")) { // to set a depth
			// case in switch begins
			if len(astValueList) > 0 || path != "" {
				//fmt.Println(path, " ", astValueList)
				elemList = append(elemList, elem{path, astValueList})
				caseCheck = append(caseCheck, basicCaseStr{"", elemList})
				path = ""
				elemList = []elem{}
				astValueList = []string{}

			}

			depthFirstCase = val.depth
			path = val.value
			flag = false
			caseName = strconv.Itoa(count) + ".case"
			caseBegin = true
			//fmt.Println("case begins!!!\t", " ", depthFirstCase, " ", len(val.value))
			continue
		}
		if depthFirstCase == val.depth && strings.Contains(val.value, "*ast.CaseClause") {
			if len(astValueList) > 0 || path != "" {
				//fmt.Println(path, " ", astValueList)
				elemList = append(elemList, elem{path, astValueList})
				caseCheck = append(caseCheck, basicCaseStr{caseName, elemList})
				path = ""
				elemList = []elem{}
				astValueList = []string{}
				count++
				// new
				cases = append(cases, testCase{caseName, caseValues})
				//resultFunc = append(resultFunc, testFuncCase{funcName, cases})
				caseValues = []basicCases{}

			}
			caseName = strconv.Itoa(count) + ".case"
			path = val.value

			if strings.Contains(val.value, "default") {
				caseBegin = false
				continue
			}
			//fmt.Println("case begins!!!\t", " ", depthFirstCase, " ", len(val.value))
			continue
		}
		if caseBegin && val.depth >= depthFirstCase {

			// new
			caseValues = append(caseValues, val)

			if strings.Contains(val.value, "*ast.Ident") {
				valExtracted := strings.Split(val.value, " ")
				astValueList = append(astValueList, valExtracted[1])
				continue
			}

			if strings.Contains(val.value, "*ast.BasicLit") {
				if !strings.EqualFold(path, "") && len(astValueList) > 0 {
					//fmt.Println(path, " ", astValueList)
					elemList = append(elemList, elem{path, astValueList})
					path = ""
					astValueList = []string{}
				}
				valExtracted := strings.Split(val.value, " ")
				fmt.Println("*ast.BasicLit", valExtracted)
				if path == "" {
					path = valExtracted[0]
				} else {
					path = path + " -> " + valExtracted[0]
				}
				if len(valExtracted) > 1 {
					astBasicLitValue := valExtracted[1] + " " + valExtracted[2]
					astValueList = append(astValueList, astBasicLitValue)
					elemList = append(elemList, elem{path, astValueList})
					path = ""
					astValueList = []string{}
				}
				continue
			}

			if strings.Contains(val.value, "*ast.AssignStmt") || strings.Contains(val.value, "*ast.BinaryExpr") || strings.Contains(val.value, "*ast.UnaryExpr") {
				if !strings.EqualFold(path, "") && len(astValueList) > 0 {
					//fmt.Println(path, " ", astValueList)
					elemList = append(elemList, elem{path, astValueList})
					path = ""
					astValueList = []string{}
				}
				valExtracted := strings.Split(val.value, " ")
				if path == "" {
					path = valExtracted[0]
				} else {
					path = path + " -> " + valExtracted[0]
				}
				astValueList = append(astValueList, valExtracted[1])
				elemList = append(elemList, elem{path, astValueList})
				path = ""
				astValueList = []string{}
				continue
			}
			if len(astValueList) > 0 {
				//fmt.Println(path, " ", astValueList)
				elemList = append(elemList, elem{path, astValueList})
				path = ""
				astValueList = []string{}
			}
			if path == "" {
				path = val.value
			} else {
				path = path + " -> " + val.value
			}
		} else if caseBegin && val.depth < depthFirstCase {
			caseBegin = false
			if len(caseValues) > 0 {
				cases = append(cases, testCase{caseName, caseValues})
				caseValues = []basicCases{}
			}
		}

	}
	if funcName != "" {
		elemList = append(elemList, elem{path, astValueList})
		caseCheck = append(caseCheck, basicCaseStr{caseName, elemList})
		funcCheck = append(funcCheck, checkCases{funcName, caseCheck})

		if len(cases) > 0 {
			//cases = append(cases, testCase{caseName, caseValues})
			resultFunc = append(resultFunc, testFuncCase{funcName, cases})
			caseValues = []basicCases{}
			cases = []testCase{}

		}
	}
	/*for ind := range funcCheck {
		fmt.Println("funcCheck[ind].funcName: ", funcCheck[ind].funcName)
		if strings.EqualFold(funcCheck[ind].funcName, "Set") {
			for _, val := range funcCheck[ind].cases {
				fmt.Println("val.funcName: ", val.caseName)
				for _, value := range val.value {
					fmt.Println(value)
				}
			}
			fmt.Println()
			fmt.Println()

		}
	}*/

	for ind := range resultFunc {
		fmt.Println(resultFunc[ind].funcName)
		for _, val := range resultFunc[ind].caseCol {
			fmt.Println("--- ", val.caseName, " ---")
			for _, value := range val.caseValue {
				fmt.Println(value)
			}
		}
	}
	return funcCheck
}

func checkSelectorExpr(listFunctions []basicStr) []basicStr {
	var modElemList []elem
	var modListFunctions []basicStr
	var modFuncName string
	var modTokPos token.Pos
	var modPath string
	var modAstValue []string
	if len(listFunctions) > 0 {
		for l := range listFunctions {
			modFuncName = listFunctions[l].funcName
			modTokPos = listFunctions[l].funcToken
			for i := 0; i < len(listFunctions[l].value); i++ {
				modPath = listFunctions[l].value[i].path
				for _, val := range listFunctions[l].value[i].value {
					modAstValue = append(modAstValue, val)
				}
				if i < len(listFunctions[l].value)-1 && listFunctions[l].value[i+1].path == "*ast.SelectorExpr" {
					if listFunctions[l].value[i+1].value[0] == "unsafe" && listFunctions[l].value[i+1].value[1] == "Pointer" {
						modPath += " -> *ast.SelectorExpr"
						for _, selVal := range listFunctions[l].value[i+1].value {
							modAstValue = append(modAstValue, selVal)
						}
						//						modAstValue = append(modAstValue, "unsafe")
						//						modAstValue = append(modAstValue, "Pointer")
						i++
					}
				}
				modElemList = append(modElemList, elem{modPath, modAstValue})
				modPath = ""
				modAstValue = []string{}
			}
			modListFunctions = append(modListFunctions, basicStr{modFuncName, modTokPos, modElemList})
			modFuncName = ""
			modElemList = []elem{}
		}
	}

	return modListFunctions
}

func checkSelectorExprCase(listFunctions []basicCaseStr) []basicCaseStr {
	var modElemList []elem
	var modListFunctions []basicCaseStr
	var modFuncName string
	var modPath string
	var modAstValue []string
	if len(listFunctions) > 0 {
		for l := range listFunctions {
			modFuncName = listFunctions[l].caseName
			for i := 0; i < len(listFunctions[l].value); i++ {
				modPath = listFunctions[l].value[i].path
				for _, val := range listFunctions[l].value[i].value {
					modAstValue = append(modAstValue, val)
				}

				if i < len(listFunctions[l].value)-1 && listFunctions[l].value[i+1].path == "*ast.SelectorExpr" {
					if listFunctions[l].value[i+1].value[0] == "unsafe" && listFunctions[l].value[i+1].value[1] == "Pointer" {
						modPath += " -> *ast.SelectorExpr"
						for _, selVal := range listFunctions[l].value[i+1].value {
							modAstValue = append(modAstValue, selVal)
						}
						i++
					}
				}
				modElemList = append(modElemList, elem{modPath, modAstValue})
				modPath = ""
				modAstValue = []string{}
			}
			modListFunctions = append(modListFunctions, basicCaseStr{modFuncName, modElemList})
			modFuncName = ""
			modElemList = []elem{}
		}
	}

	return modListFunctions
}
