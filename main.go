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
	//node_list := list.New()
	// print all of nodes from ast
	//ast.Print(fset, astTree)
	/*count := 0

	for _, tree := range astTree.Decls {
		nodes, ok := tree.(*ast.FuncDecl)

		if ok {
			fmt.Println("Functions:", nodes.Name.Name)
			node_list.PushBack(nodes)
			//ast.Print(fset, nodes)
			fmt.Println()
			count++
			//				continue
		}
		//fmt.Println(fn.Name.Name)
	}
	fmt.Println("count is: ", count)*/
	count := 0
	ident := 0
	ast.Inspect(astTree, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.FuncDecl:
			//if !strings.Contains(x.Name.String(), "et") {
			if ident != 0 {
				fmt.Println("the number of ident is ", ident)
			}
			fmt.Println(count, "function name is:", x.Name)
			count++
			ident = 0
			//}
			/*
				case *ast.CaseClause:
					if len(x.List) > 0 { // len(x.List) == 0 -> default
						fmt.Println(len(x.List))
						fmt.Println("There is switch statement: ", x.List[0])
						fmt.Println("Body is: ", x.Body)
						node_list.PushBack(x.Body)
						//node_list.PushBack(x)
						//				if listElemAsIdent, ok := x.List[0].(*ast.Ident); ok {
						//				}
					}*/
		case *ast.Ident:
			ident++
			fmt.Println("ast.Ident is: ", x)
		//node_list.PushBack(x)
		// TODO: detailed!!!!! must to know sequence
		case *ast.AssignStmt:
			fmt.Println("!!! Assignstatement !!!")
		case *ast.BlockStmt:
			fmt.Println("!!! Blockstatement !!!")
		case *ast.DeclStmt:
			fmt.Println("!!! Declstatement !!!")

		}

		return true
	}) /*
		print_nodes := list.New()
		for item := node_list.Front(); item != nil; item = item.Next() {
			//fmt.Println("item value is ", item.Value)
			node_list.PushBack(ast.Print(fset, item))
		}

		for item := print_nodes.Front(); item != nil; item = item.Next() {
			//fmt.Println("item value is ", item.Value)
			fmt.Println(item.Value)
		}*/
	//fmt.Println(ast.Print(fset, node_list.Front()))

}
