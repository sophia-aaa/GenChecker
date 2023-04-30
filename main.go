package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	//	filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run main.go - test.go
	filename := os.Args[2]
	fmt.Println(filename)
	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// print all of nodes from ast
	//ast.Print(fset, astTree)
	count := 0

	for _, tree := range astTree.Decls {
		nodes, ok := tree.(*ast.FuncDecl)
		if ok {
			fmt.Println("Functions:", nodes.Name.Name)
			ast.Print(fset, nodes)
			fmt.Println()
			count++
			//				continue
		}
		//fmt.Println(fn.Name.Name)
	}
	fmt.Println("count is: ", count)
}
