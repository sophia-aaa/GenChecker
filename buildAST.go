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

type basicFunc struct {
	funcName  string
	funcToken token.Pos // to distinguish same name function
}

type basicCaseStr struct {
	caseName string
	value    []elem
}

type caseResult struct {
	funcName  string
	funcToken string
	cases     []basicCaseStr
}

type caseFilteredResult struct {
	funcName     string
	funcToken    string
	caseFiltered []string
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
		case *ast.GenDecl:
			if !genDeclBool {
				if astNode != "" || len(astValue) != 0 { // if a new root meets
					elem_list = append(elem_list, elem{astNode, astValue})
					listFunctions = append(listFunctions, basicStr{nameFunction, tokPos, elem_list})
					elem_list = []elem{}
					astNode = ""
					astValue = []string{}

					nameFunction = ""
					tokPos = x.Pos()

				}
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
			if len(elem_list) > 0 || astNode != "" || len(astValue) != 0 { // if a new root meets
				if astNode != "" || len(astValue) != 0 {
					elem_list = append(elem_list, elem{astNode, astValue})
				}
				listFunctions = append(listFunctions, basicStr{nameFunction, tokPos, elem_list})
				elem_list = []elem{}
				astNode = ""
				astValue = []string{}
			}
			nameFunction = x.Name.String()

			tokPos = x.Pos()
		//fmt.Println(x.Name.String(), "\t", rune(x.Pos()))
		case *ast.AssignStmt:
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
			astValue = append(astValue, x.Tok.String())
		case *ast.BinaryExpr:
			if !strings.EqualFold(astNode, "") || len(astValue) != 0 {
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
			if !strings.EqualFold(astNode, "") || len(astValue) != 0 {
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
			genDeclBool = false
		case *ast.InterfaceType:
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
			/* I couldn't find a way to express "incomplete", so I customised it like this.
			if x.Incomplete == true {
				elem_list = append(elem_list, elem{"*ast.InterfaceType.Incomplete", []string{"true"}})
			} else {
				elem_list = append(elem_list, elem{"*ast.InterfaceType.Incomplete", []string{"false"}})
			}*/
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
		if !strings.EqualFold(astNode, "") || len(astValue) > 0 {
			elem_list = append(elem_list, elem{astNode, astValue})
		}
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
			str = fmt.Sprintf("%s %v %v\n", reflect.TypeOf(n).String(), valExtracted[2], n.Pos())
			str4cases = fmt.Sprintf("%s %v %v", reflect.TypeOf(n).String(), valExtracted[2], n.Pos())
		} else { // Function has no name!
			str = fmt.Sprintf("%s %v %v\n", reflect.TypeOf(n).String(), "No Name", n.Pos())
			str4cases = fmt.Sprintf("%s %v %v", reflect.TypeOf(n).String(), "No Name", n.Pos())
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

type funcCollection struct {
	funcName  string
	funcToken string
	caseCol   []caseCollection
}

type caseCollection struct {
	caseName  string
	caseValue []basicCases
}

func buildAstCaseStr(Tree2cases []basicCases, funcNameList []basicFunc) []caseResult {
	var funcName string
	var depthFirstCase int
	var caseName string

	var cases []caseCollection
	var caseValues []basicCases
	var resultFunc []funcCollection
	var funcToken string
	flag := true
	count := 0
	caseBegin := false

	for _, val := range Tree2cases {
		if val.depth == 1 && strings.Contains(val.value, "*ast.FuncDecl") { // function begins
			if checkDuplicateInFuncString(funcNameList, funcName, funcToken) {
				// new
				if len(caseValues) > 0 {
					cases = append(cases, caseCollection{caseName, caseValues})
				}
				resultFunc = append(resultFunc, funcCollection{funcName, funcToken, cases})
				caseValues = []basicCases{}
				cases = []caseCollection{}
				count = 0
				depthFirstCase = 0
			} else {
				caseValues = []basicCases{}
				cases = []caseCollection{}
				count = 0
				depthFirstCase = 0
			}
			valExtracted := strings.Split(val.value, " ")
			funcName = valExtracted[1]
			funcToken = valExtracted[2]

			flag = true
			continue
		}
		if flag && (strings.Contains(val.value, "*ast.CaseClause")) { // to set a depth
			// case in switch begins
			depthFirstCase = val.depth
			flag = false
			caseName = strconv.Itoa(count) + ".case"
			caseBegin = true
			//fmt.Println("case begins!!!\t", " ", depthFirstCase, " ", len(val.value))
			continue
		}
		if depthFirstCase == val.depth && strings.Contains(val.value, "*ast.CaseClause") {
			if len(caseValues) > 0 {
				// new
				cases = append(cases, caseCollection{caseName, caseValues})
				//resultFunc = append(resultFunc, testFuncCase{funcName, cases})
				caseValues = []basicCases{}
				count++
			}
			caseName = strconv.Itoa(count) + ".case"
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
		} else if caseBegin && val.depth < depthFirstCase {
			caseBegin = false
			if len(caseValues) > 0 {
				cases = append(cases, caseCollection{caseName, caseValues})
				caseValues = []basicCases{}
			}
		}

	}
	if funcName != "" {
		if len(cases) > 0 {
			//cases = append(cases, testCase{caseName, caseValues})
			resultFunc = append(resultFunc, funcCollection{funcName, funcToken, cases})
			caseValues = []basicCases{}
			cases = []caseCollection{}

		}
	}

	var listFunction []caseResult
	caseList := []basicCaseStr{}
	elemList := []elem{}
	depthCompare := 0
	path := ""
	values := []string{}
	for ind := range resultFunc {
		for _, val := range resultFunc[ind].caseCol {
			for i, value := range val.caseValue {
				if depthCompare <= value.depth {
					if strings.Contains(value.value, " ") {

						substring := strings.Split(value.value, " ")
						//fmt.Println("substring ", substring)
						if strings.EqualFold(substring[0], "*ast.Ident") {
							for i := 1; i < len(substring); i++ {
								values = append(values, substring[i])
							}

						} else {
							if strings.EqualFold(path, "") {
								path = substring[0]
							} else {
								path += " -> " + substring[0]
							}
							if strings.EqualFold(substring[0], "*ast.BasicLit") {
								basicLitValue := ""
								if len(substring) > 2 {
									for i := 2; i < len(substring); i++ {
										if i == len(substring)-1 {
											basicLitValue += substring[i]
										} else {
											basicLitValue += substring[i] + " "
										}
									}
									values = append(values, substring[1])
									values = append(values, basicLitValue)
								}

							} else {
								for i := 1; i < len(substring); i++ {
									values = append(values, substring[i])
								}
							}

						}
					} else {
						if strings.EqualFold(path, "") {
							path = value.value
						} else {
							path += " -> " + value.value
						}
					}
					depthCompare = value.depth
				} else {
					if !strings.EqualFold(path, "") || len(values) != 0 { // append elemList
						elemList = append(elemList, elem{path, values})
						path = ""
						values = []string{}
					}
					if strings.Contains(value.value, " ") {
						substring := strings.Split(value.value, " ")
						if !strings.EqualFold(substring[0], "*ast.Ident") {
							if strings.EqualFold(path, "") {
								path = substring[0]
							} else {
								path += " -> " + substring[0]
							}
						}
						if strings.EqualFold(substring[0], "*ast.BasicLit") {
							basicLitValue := ""
							if len(substring) > 2 {
								for i := 2; i < len(substring); i++ {
									if i == len(substring)-1 {
										basicLitValue += substring[i]
									} else {
										basicLitValue += substring[i] + " "
									}
								}
								values = append(values, substring[1])
								values = append(values, basicLitValue)
							}

						} else {
							for i := 1; i < len(substring); i++ {
								values = append(values, substring[i])
							}
						}

					} else {
						if strings.EqualFold(path, "") {
							path = value.value
						} else {
							path += " -> " + value.value
						}
					}
					depthCompare = value.depth
				}
				if i == len(val.caseValue)-1 {
					if !strings.EqualFold(path, "") || len(values) != 0 {
						elemList = append(elemList, elem{path, values})
						path = ""
						values = []string{}
					}
				}
			}
			caseList = append(caseList, basicCaseStr{val.caseName, elemList})
			elemList = []elem{}
		}
		listFunction = append(listFunction, caseResult{resultFunc[ind].funcName, resultFunc[ind].funcToken, caseList})
		caseList = []basicCaseStr{}
		elemList = []elem{}
		path = ""
		values = []string{}
	}
	return listFunction
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
