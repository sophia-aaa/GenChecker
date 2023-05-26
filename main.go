package main

import (
	"container/list"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
)

type elem struct {
	name  string
	value []string
}

type function1 struct {
	funcName string
	value    *list.List
}
type function2 struct {
	funcName string
	value    []elem
}

type values struct {
	lists []string
}

func resetValues() values {
	return values{lists: []string{}}

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

	/*test := values{lists: []string{"test1", "test2"}}
	fmt.Println(test)
	test = resetValues()
	fmt.Println("is reset?", test)*/

	//ast.Print(fset, astTree)

	node_list := list.New()
	var elem_list []elem
	nameFunction := ""
	//nameCheck := 0
	astNode := ""
	var listFunctions2 []function2
	var astValue []string
	fmt.Println(len(astValue))
	ast.Inspect(astTree, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.FuncDecl:
			//	fmt.Println(x.Name)
			//	fmt.Println(x, "\t\t", reflect.TypeOf(x).String())
			//	fmt.Println(fset.Position(x.Pos()), fset.Position(x.End()))
			if len(astValue) != 0 { // if a new root meets
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				func2 := function2{nameFunction, elem_list}
				listFunctions2 = append(listFunctions2, func2)
				elem_list = []elem{}
				astNode = ""
				astValue = []string{}
				//nameCheck = 0
			}
			nameFunction = x.Name.String()
			node_list.PushBack(elem{reflect.TypeOf(x).String(), []string{nameFunction}})

		case *ast.Ident:
			astValue = append(astValue, x.Name)

		case *ast.GenDecl, *ast.ImportSpec, *ast.BasicLit, *ast.CommentGroup, *ast.Comment:
			fmt.Print(reflect.TypeOf(x))
		default:
			if x != nil {
				//		fmt.Println(x, "\t\t", reflect.TypeOf(x).String())
				//			fmt.Println(fset.Position(x.Pos()))
				if len(astValue) != 0 {
					elem_list = append(elem_list, elem{astNode, astValue})
					node_list.PushBack(elem{astNode, astValue})
					astNode = ""
					astValue = []string{}
				}
				if astNode == "" {
					astNode += reflect.TypeOf(x).String()
				} else {
					astNode += " -> " + reflect.TypeOf(x).String()
				}
			}

			/*	if x.Name == nameFunction {
						nameCheck++
					}
					if x.Name != nameFunction && nameCheck != 0 {
						astValue = append(astValue, x.Name)
					}
				}

					if node_list != nil {
						fmt.Println(nameFunction, " added! node_list")
						listFunctions1 = append(listFunctions1, function1{nameFunction, node_list})
						node_list.Init()
					}
					if elem_list != nil {
						fmt.Println(nameFunction, " added! elem_list")
						listFunctions2 = append(listFunctions2, function2{nameFunction, elem_list})
						elem_list = []elem{}
					}*/
		}
		return true
	})

	if len(astValue) != 0 { // if a new root meets
		elem_list = append(elem_list, elem{astNode, astValue})
		node_list.PushBack(elem{astNode, astValue})
		func2 := function2{nameFunction, elem_list}
		listFunctions2 = append(listFunctions2, func2)
		elem_list = []elem{}
		astNode = ""
		astValue = []string{}
	}
	/*for item := node_list.Front(); item != nil; item = item.Next() {
		fmt.Println(item.Value)
	}*/

	for s := range listFunctions2 {
		if listFunctions2[s].funcName != "" {
			fmt.Println(listFunctions2[s].funcName)
			for _, value := range listFunctions2[s].value {
				fmt.Println("\t", value)
			}
		} else {
			for _, value := range listFunctions2[s].value {
				fmt.Println(value)
			}
		}
		fmt.Println()
	}

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
	var str string
	nextLine := "\n\n"
	for s := range listFunctions2 {
		if listFunctions2[s].funcName != "" {

			str = "function name is " + listFunctions2[s].funcName + " \n"
			_, err2 := f.WriteString(str)
			if err2 != nil {
				log.Fatal(err2)
			}
			for _, value := range listFunctions2[s].value {
				str = fmt.Sprintln("\t", value)
				_, err2 := f.WriteString(str)
				if err2 != nil {
					log.Fatal(err2)
				}
			}
			_, err2 = f.WriteString(nextLine)
			if err2 != nil {
				log.Fatal(err2)
			}

		} else {

			for _, value := range listFunctions2[s].value {
				str = fmt.Sprintln(value)
				_, err2 := f.WriteString(str)
				if err2 != nil {
					log.Fatal(err2)
				}
			}
			_, err2 := f.WriteString(nextLine)
			if err2 != nil {
				log.Fatal(err2)
			}
		}

	}
	/*for s := range listFunctions2 {
		fmt.Print("function name is ", listFunctions2[s].funcName, " ")
		fmt.Printf("%d\n", len(listFunctions2[s].value))
		for _, value := range listFunctions2[s].value {
			fmt.Printf("%s\n", value)
		}
		fmt.Println()
		fmt.Println()
	}*/

}
