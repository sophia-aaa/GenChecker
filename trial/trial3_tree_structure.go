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

var Tree2str []string

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
}

// Reference: https://golangdocs.com/golang-ast-package
type visitor int

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	var str string
	// TODO to test fmt.Printf -> fmt.Sprinft and save in []str
	// int(v) is a depth of a current node
	if reflect.TypeOf(n).String() == "*ast.Ident" {
		str = fmt.Sprintf("%s%d: ", strings.Repeat("\t", int(v)), int(v))
		Tree2str = append(Tree2str, str)
		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), n)
		Tree2str = append(Tree2str, str)
	} else if reflect.TypeOf(n).String() == "*ast.FuncDecl" {
		funcName := fmt.Sprintf("%v", n)
		count := 0
		var strIdx int
		var endIdx int
		for i := 0; i < len(str); i++ {
			if str[i] == ' ' {
				count++
			}
			if str[i] == ' ' && count == 2 {
				strIdx = i + 1
			}
			if str[i] == ' ' && count == 3 {
				endIdx = i
			}
		}
		str = fmt.Sprintf("%s%d: ", strings.Repeat("\t", int(v)), int(v))
		Tree2str = append(Tree2str, str)
		str = fmt.Sprintf("%s %v\n", reflect.TypeOf(n).String(), funcName[strIdx:endIdx])
		Tree2str = append(Tree2str, str)
	} else {
		str = fmt.Sprintf("%s%d: ", strings.Repeat("\t", int(v)), int(v))
		Tree2str = append(Tree2str, str)
		str = fmt.Sprintf("%T\n", n)
		Tree2str = append(Tree2str, str)
	}
	return v + 1
}
