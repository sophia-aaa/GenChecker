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
	funcName string
	value    []elem
}

func buildAstDataStr(filename string) []basicStr {
	// Commentgroups and comments are ignored
	// because they interfere with comparing the content of the code.

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

	ast.Inspect(astTree, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.CommentGroup, *ast.Comment:
			astNode += ""

		case *ast.GenDecl:
			if !genDeclBool {
				fmt.Println("GenDecl ", nameFunction)
				if astNode != "" || len(astValue) != 0 { // if a new root meets
					fmt.Println("GenDecl ", nameFunction)
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
			fmt.Println("FuncDecl ", nameFunction, astNode != "" || len(astValue) != 0)
			//	fmt.Println(x, "\t\t", reflect.TypeOf(x).String())
			//	fmt.Println(fset.Position(x.Pos()), fset.Position(x.End()))
			if astNode != "" || len(astValue) != 0 { // if a new root meets
				elem_list = append(elem_list, elem{astNode, astValue})
				listFunctions = append(listFunctions, basicStr{nameFunction, tokPos, elem_list})
				elem_list = []elem{}
				astNode = ""
				astValue = []string{}
			}
			nameFunction = x.Name.String()
			tokPos = x.Pos()
			//fmt.Println(x.Name.String(), "\t", rune(x.Pos()))
		case *ast.BinaryExpr:
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
			astValue = append(astValue, x.Op.String())
			genDeclBool = false
		case *ast.UnaryExpr:
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

	var str string
	var str4cases string
	var depth int
	// int(v) is a depth of a current node
	if reflect.TypeOf(n).String() == "*ast.Ident" {
		str = fmt.Sprintf("%s%d: ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)

		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), n)
		str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), n)
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
	} else if reflect.TypeOf(n).String() == "*ast.FuncDecl" {
		funcName := fmt.Sprintf("%v", n)
		count := 0
		var strIdx int
		var endIdx int
		for i := 0; i < len(funcName); i++ {
			if funcName[i] == ' ' {
				count++
			}
			if funcName[i] == ' ' && count == 2 {
				strIdx = i + 1
			}
			if funcName[i] == ' ' && count == 3 {
				endIdx = i
			}
		}
		str = fmt.Sprintf("%s%d: ", strings.Repeat("\t", int(v)), int(v))
		depth = int(v)
		Tree2str = append(Tree2str, str)
		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), funcName[strIdx:endIdx])
		str4cases = fmt.Sprintf("%s %v", reflect.TypeOf(n).String(), funcName[strIdx:endIdx])
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})
		count = 0
	} else {
		str = fmt.Sprintf("%s%d: ", strings.Repeat("\t", int(v)), int(v))
		Tree2str = append(Tree2str, str)
		depth = int(v)
		str = fmt.Sprintf("%T\n", n)
		str4cases = fmt.Sprintf("%T", n)
		Tree2str = append(Tree2str, str)
		Tree2cases = append(Tree2cases, basicCases{depth, str4cases})

	}
	return v + 1
}

func buildAstCaseStr(Tree2cases []basicCases) []checkCases {
	var funcName string
	var path string
	var astValue string
	var astValueList []string
	var depthFirstCase int
	var caseName string
	var elemList []elem
	var caseCheck []basicCaseStr
	var funcCheck []checkCases
	flag := true
	count := 0

	for _, val := range Tree2cases {
		if strings.Contains(val.value, "*ast.FuncDecl") { // function begins
			if funcName != "" {
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
			}
			funcName = val.value[14:len(val.value)]
			//fmt.Println()
			//fmt.Println()
			//fmt.Println(">> ", funcName, " <<")
			flag = true
			continue
		}
		if val.depth == 2 && strings.Contains(val.value, "*ast.Ident") { // file name
			if len(astValueList) > 0 || path != "" {
				elemList = append(elemList, elem{path, astValueList})
				//fmt.Println(path, " ", astValueList)
				path = ""
				astValueList = []string{}
			}
			astValue = val.value[11:len(val.value)]
			astValueList = append(astValueList, astValue)
			//fmt.Println(astValueList)
			elemList = append(elemList, elem{"file name", astValueList}) // filename is not a convention from go. I added.
			path = ""
			astValueList = []string{}
			continue
		} else if strings.Contains(val.value, "*ast.Ident") {
			astValue = val.value[11:len(val.value)]
			astValueList = append(astValueList, astValue)
			continue
		}
		if flag && (strings.Contains(val.value, "*ast.CaseClause") || strings.Contains(val.value, "*ast.CaseClause -> *ast.SelectorExpr")) { // to set a depth
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
			//fmt.Println("case begins!!!\t", " ", depthFirstCase, " ", len(val.value))
			continue
		} else if depthFirstCase == val.depth && (strings.Contains(val.value, "*ast.CaseClause") || strings.Contains(val.value, "*ast.CaseClause -> *ast.SelectorExpr")) {
			if len(astValueList) > 0 || path != "" {
				//fmt.Println(path, " ", astValueList)
				elemList = append(elemList, elem{path, astValueList})
				caseCheck = append(caseCheck, basicCaseStr{caseName, elemList})
				path = ""
				elemList = []elem{}
				astValueList = []string{}
				count++
			}
			caseName = strconv.Itoa(count) + ".case"
			path = val.value
			//fmt.Println("case begins!!!\t", " ", depthFirstCase, " ", len(val.value))
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
	}
	if funcName != "" {
		elemList = append(elemList, elem{path, astValueList})
		caseCheck = append(caseCheck, basicCaseStr{caseName, elemList})
		funcCheck = append(funcCheck, checkCases{funcName, caseCheck})
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
						modAstValue = append(modAstValue, "unsafe")
						modAstValue = append(modAstValue, "Pointer")
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
			modFuncName = listFunctions[l].funcName
			for i := 0; i < len(listFunctions[l].value); i++ {
				modPath = listFunctions[l].value[i].path
				for _, val := range listFunctions[l].value[i].value {
					modAstValue = append(modAstValue, val)
				}
				if i < len(listFunctions[l].value)-1 && listFunctions[l].value[i+1].path == "*ast.SelectorExpr" {
					if listFunctions[l].value[i+1].value[0] == "unsafe" && listFunctions[l].value[i+1].value[1] == "Pointer" {
						modPath += " -> *ast.SelectorExpr"
						modAstValue = append(modAstValue, "unsafe")
						modAstValue = append(modAstValue, "Pointer")
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
