package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	_ "golang.org/x/exp/slices"
	"log"
	"os"
	"reflect"
	"strings"
)

type elem struct {
	path  string
	value []string
}

// Basic structure for every function in the input file
type basicStr struct {
	funcName string
	value    []elem
}

type funcNameList struct {
	lists []string
}

type checkCases struct {
	funcName string
	cases    []basicStr
}

type basicCases struct {
	depth int
	value string
}

var Tree2str []string
var Tree2cases []basicCases

func createTextFile(filename string, strList []string) {
	f, err := os.Create(filename + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)
	if len(strList) > 0 {
		for _, val := range strList {
			_, err1 := f.WriteString(val)
			if err1 != nil {
				log.Fatal(err1)
			}
		}
	}
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

func main() {
	// filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run gen.go - test.go
	filename := os.Args[2]
	fmt.Println(filename)
	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var v visitor
	ast.Walk(v, astTree)
	if len(Tree2str) > 0 {
		createTextFile("tree", Tree2str)
	}
	//fmt.Println(len("*ast.Ident"))
	// TODO build basicStructure and element for SwitchStatement
	var funcName string
	var path string
	var astValue string
	var astValueList []string
	var depthFirstCase int
	flag := true
	if len(Tree2cases) > 0 {
		for _, val := range Tree2cases {
			if strings.Contains(val.value, "*ast.FuncDecl") { // function begins
				if len(astValueList) > 0 {
					fmt.Println(path, " ", astValueList)
					astValueList = []string{}
				}
				funcName = val.value[14:len(val.value)]
				fmt.Println()
				fmt.Println()
				fmt.Println(">> ", funcName, " <<")
				flag = true
				continue
			}
			if val.depth == 2 && strings.Contains(val.value, "*ast.Ident") { // file name
				if len(astValueList) > 0 {
					fmt.Println(path, " ", astValueList)
					astValueList = []string{}
				}
				astValue = val.value[11:len(val.value)]
				astValueList = append(astValueList, astValue)
				fmt.Println(astValueList)
				astValueList = []string{}
				continue
			} else if strings.Contains(val.value, "*ast.Ident") {
				astValue = val.value[11:len(val.value)]
				astValueList = append(astValueList, astValue)
				continue
			}
			if flag && strings.Contains(val.value, "*ast.CaseClause") { // to set a depth
				if len(astValueList) > 0 {
					fmt.Println(path, " ", astValueList)
					astValueList = []string{}
				}
				depthFirstCase = val.depth
				flag = false
				path = val.value
				fmt.Println("case begins!!!\t", " ", depthFirstCase, " ", len(val.value))
				continue
			} else if depthFirstCase == val.depth && strings.Contains(val.value, "*ast.CaseClause") {
				if len(astValueList) > 0 {
					fmt.Println(path, " ", astValueList)
					astValueList = []string{}
				}
				path = val.value
				fmt.Println("case begins!!!\t", " ", depthFirstCase, " ", len(val.value))
				continue
			}
			if len(astValueList) > 0 {
				fmt.Println(path, " ", astValueList)
				path = ""
				astValueList = []string{}
			}
			if path == "" {
				path += val.value
			} else {
				path = path + " -> " + val.value
			}
		}
	}

}
