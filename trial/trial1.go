

import (
	"container/list"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

package main

import (
"container/list"
"fmt"
"go/ast"
"go/parser"
"go/token"
"log"
"os"
)

type elem struct {
	name  string
	value []string
}

type function struct {
	funcName string
	value    []elem
}

func main() {
	//	filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run gen.go - test.go
	filename := os.Args[2]
	fmt.Println(filename)
	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	node_list := list.New()
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
	astNode := ""
	var astValue []string
	fmt.Println(len(astValue))
	ast.Inspect(astTree, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.FuncDecl:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			node_list.PushBack(elem{x.Name.String(), []string{}})

			/*
				if astNode == "" {
					astNode += x.Name.String()
				} else {
					astNode += " -> " + x.Name.String()
				}*/

			//	node_list.PushBack(elem{"FuncDecl", []string{x.Name.String()}})
			//if !strings.Contains(x.Name.String(), "et") {
			//if ident != 0 {
			//	fmt.Println("the number of ident is ", ident)
			//}
			//fmt.Println("function name is:", x.Name)
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
		// TODO: detailed!!!!! must to know sequence
		case *ast.SelectorExpr:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "SelectorExpr"
			} else {
				astNode += " -> " + "SelectorExpr"
			}
			//e := elem{"SelectorExpr", []string{}}
			//node_list.PushBack(e)

			//fmt.Print(fset.Position(x.Pos()))
			//fmt.Println(" ---> SelectorExpr", x, x.Sel, x.X)
			if x.Sel.String() == "SliceHeader" {
				fmt.Println("There is SliceHeader")
			}
		case *ast.AssignStmt:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "AssignStmt"
			} else {
				astNode += " -> " + "AssignStmt"
			}
			//e := elem{"AssignStmt", []string{}}
			//node_list.PushBack(e)
			//fmt.Print(fset.Position(x.Pos()))
			//fmt.Println(" ---> AssignStmt", x, x.Lhs, x.Rhs)
		case *ast.BlockStmt:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "BlockStmt"
			} else {
				astNode += " -> " + "BlockStmt"
			}
			//e := elem{"BlockStmt", []string{}}
			//node_list.PushBack(e)
			//	fmt.Print(fset.Position(x.Pos()))
			//	fmt.Println(" ---> BlockStmt", x, len(x.List))
			for i := 0; i < len(x.List); i++ {
				//fmt.Println("item value is ", item.Value)
				//		fmt.Print(x.List[i], "\t")
			}
		case *ast.CallExpr:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "CallExpr"
			} else {
				astNode += " -> " + "CallExpr"
			}
			//e := elem{"CallExpr", []string{}}
			//node_list.PushBack(e)
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println(" ---> CallExpr", x)
		case *ast.ReturnStmt:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "ReturnStmt"
			} else {
				astNode += " -> " + "ReturnStmt"
			}
			//e := elem{"ReturnStmt", []string{}}
			//node_list.PushBack(e)
		//	fmt.Print(fset.Position(x.Pos()))
		//	fmt.Println(" ---> ReturnStmt", x.Return, x.Results)
		case *ast.SliceExpr:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "SliceExpr"
			} else {
				astNode += " -> " + "SliceExpr"
			}
			//e := elem{"SliceExpr", []string{}}
			//node_list.PushBack(e)
		//	fmt.Print(fset.Position(x.Pos()))
		//	fmt.Println(" ---> SliceExpr", x, x.Slice3)
		case *ast.StarExpr:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "StarExpr"
			} else {
				astNode += " -> " + "StarExpr"
			}
			//e := elem{"StarExpr", []string{}}
			//node_list.PushBack(e)
		//	fmt.Print(fset.Position(x.Pos()))
		//		fmt.Println(" ---> StarExpr", x, x.Star, x.X)
		case *ast.IndexExpr:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "IndexExpr"
			} else {
				astNode += " -> " + "IndexExpr"
			}
			//e := elem{"IndexExpr", []string{}}
			//node_list.PushBack(e)
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println(" ---> IndexExpr", x, x.Index)
		case *ast.SwitchStmt:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "SwitchStmt"
			} else {
				astNode += " -> " + "SwitchStmt"
			}
			//e := elem{"SwitchStmt", []string{}}
			//node_list.PushBack(e)
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println(" ---> Switch", x)
		case *ast.InterfaceType:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "InterfaceType"
			} else {
				astNode += " -> " + "InterfaceType"
			}
			//e := elem{"InterfaceType", []string{}}
			//node_list.PushBack(e)
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println(" ---> Interface", x, x.Interface)
		case *ast.DeclStmt:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "DeclStmt"
			} else {
				astNode += " -> " + "DeclStmt"
			}
			//e := elem{"DeclStmt", []string{}}
			//node_list.PushBack(e)
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println("!!! Declstatement !!!", x.Decl)
		case *ast.FieldList:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "FieldList"
			} else {
				astNode += " -> " + "FieldList"
			}
			//e := elem{"FieldList", []string{}}
			//node_list.PushBack(e)
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println(" ---> FieldList", x, x.NumFields())
		case *ast.CaseClause: //input
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "CaseClause"
			} else {
				astNode += " -> " + "CaseClause"
			}
			//e := elem{"CaseClause", []string{}}
			//node_list.PushBack(e)
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println(" ------> CaseClause", x.Case)
		case *ast.BasicLit:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "BasicLit"
			} else {
				astNode += " -> " + "BasicLit"
			}
			//e := elem{"BasicLit", []string{}}
			//node_list.PushBack(e)
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println("BasicLit", x.Value)
		case *ast.StructType:
			if len(astValue) != 0 {
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "StructType"
			} else {
				astNode += " -> " + "StructType"
			}
			//e := elem{"StructType", []string{}}
			//node_list.PushBack(e)
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println("Struct", x.Struct)
		case *ast.Ident: /*
				if len(astValue) != 0 {
					node_list.PushBack(elem{astNode, astValue})
					astNode = ""
					astValue = []string{}
				}
				if astNode == "" {
					astNode += "Ident"
				} else {
					astNode += " -> " + "Ident"
				}*/
			astValue = append(astValue, x.Name)
			//e := elem{"Ident", []string{}}
			//node_list.PushBack(e)
			ident++
			//		fmt.Print(fset.Position(x.Pos()))
			//		fmt.Println(" ---------> Ident", x.Name)
			//fmt.Println("ast.Ident is: ", x)
		}

		return true
	})
	//		print_nodes := list.New()
	for item := node_list.Front(); item != nil; item = item.Next() {
		fmt.Println(item.Value)
		//node_list.PushBack(ast.Print(fset, item))
	}
	/*
		for item := print_nodes.Front(); item != nil; item = item.Next() {
			//fmt.Println("item value is ", item.Value)
			fmt.Println(item.Value)
		}*/
	//fmt.Println(ast.Print(fset, node_list.Front()))

}
