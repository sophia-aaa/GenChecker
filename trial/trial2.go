

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

type function1 struct {
	funcName string
	value    *list.List
}
type function2 struct {
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

	ast.Print(fset, astTree)

	node_list := list.New()
	var elem_list []elem
	nameFunction := ""
	//nameCheck := 0
	astNode := ""
	var listFunctions1 []function1
	var listFunctions2 []function2
	var astValue []string
	fmt.Println(len(astValue))
	ast.Inspect(astTree, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.FuncDecl:
			fmt.Println(x.Name)
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
			node_list.PushBack(elem{nameFunction, []string{}})

		// TODO: detailed!!!!! must to know sequence
		case *ast.SelectorExpr:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "SelectorExpr"
			} else {
				astNode += " -> " + "SelectorExpr"
			}
			if x.Sel.String() == "SliceHeader" {
				fmt.Println("There is SliceHeader")
			}
		case *ast.AssignStmt:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "AssignStmt"
			} else {
				astNode += " -> " + "AssignStmt"
			}
		case *ast.BlockStmt:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "BlockStmt"
			} else {
				astNode += " -> " + "BlockStmt"
			}
		case *ast.CallExpr:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "CallExpr"
			} else {
				astNode += " -> " + "CallExpr"
			}
		case *ast.ReturnStmt:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "ReturnStmt"
			} else {
				astNode += " -> " + "ReturnStmt"
			}
		case *ast.SliceExpr:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "SliceExpr"
			} else {
				astNode += " -> " + "SliceExpr"
			}
		case *ast.StarExpr:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "StarExpr"
			} else {
				astNode += " -> " + "StarExpr"
			}
		case *ast.IndexExpr:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "IndexExpr"
			} else {
				astNode += " -> " + "IndexExpr"
			}
		case *ast.SwitchStmt:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "SwitchStmt"
			} else {
				astNode += " -> " + "SwitchStmt"
			}
		case *ast.InterfaceType:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "InterfaceType"
			} else {
				astNode += " -> " + "InterfaceType"
			}
		case *ast.DeclStmt:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "DeclStmt"
			} else {
				astNode += " -> " + "DeclStmt"
			}
		case *ast.FieldList:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "FieldList"
			} else {
				astNode += " -> " + "FieldList"
			}
		case *ast.CaseClause: //input
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "CaseClause"
			} else {
				astNode += " -> " + "CaseClause"
			}
		case *ast.BasicLit:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "BasicLit"
			} else {
				astNode += " -> " + "BasicLit"
			}
		case *ast.StructType:
			if len(astValue) != 0 {
				elem_list = append(elem_list, elem{astNode, astValue})
				node_list.PushBack(elem{astNode, astValue})
				astNode = ""
				astValue = []string{}
			}
			if astNode == "" {
				astNode += "StructType"
			} else {
				astNode += " -> " + "StructType"
			}
		case *ast.Ident:
			astValue = append(astValue, x.Name)
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
	for item := node_list.Front(); item != nil; item = item.Next() {
		fmt.Println(item.Value)
	}
	fmt.Println(len(listFunctions1))
	fmt.Println(len(listFunctions2))

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
		str = "function name is " + listFunctions2[s].funcName + " "
		_, err2 := f.WriteString(str)
		if err2 != nil {
			log.Fatal(err2)
		}
		for _, value := range listFunctions2[s].value {
			str = fmt.Sprintln(value)
			_, err2 := f.WriteString(str)
			if err2 != nil {
				log.Fatal(err2)
			}
		}
		_, err2 = f.WriteString(nextLine)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	for s := range listFunctions2 {
		fmt.Print("function name is ", listFunctions2[s].funcName, " ")
		fmt.Printf("%d\n", len(listFunctions2[s].value))
		for _, value := range listFunctions2[s].value {
			fmt.Printf("%s\n", value)
		}
		fmt.Println()
		fmt.Println()
	}

}
