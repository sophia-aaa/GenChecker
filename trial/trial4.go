package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"os"
)

func main() {
	/*fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "test.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		// handle function declarations without documentation
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			if fn.Name.IsExported() && fn.Doc.Text() == "" {
				// print warning
				fmt.Printf("exported function declaration without documentation found on line %d: \n\t%s\n", fset.Position(fn.Pos()).Line, fn.Name.Name)
				// create todo-comment
				comment := &ast.Comment{
					Text:  "// TODO: document exported function",
					Slash: fn.Pos() - 1,
				}
				// create CommentGroup and set it to the function's documentation comment
				cg := &ast.CommentGroup{
					List: []*ast.Comment{comment},
				}
				fn.Doc = cg
				fmt.Println()
			}
		}
		return true
	})

	// write newFile ast to file
	f, err := os.Create("new.go")
	defer f.Close()
	if err := printer.Fprint(f, fset, node); err != nil {
		log.Fatal(err)
	}*/

	fset := token.NewFileSet()
	filename := "test1.go"
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	//import "github.com/davecgh/go-spew/spew"
	//spew.Dump(f)
	//ast.Print(fset, f)
	if err != nil {
		log.Fatal(err)
	}
	once := true
	n := astutil.Apply(f, func(c *astutil.Cursor) bool {
		//ast.Print(fset, c.Node().(*ast.File).Decls[0].(*ast.GenDecl).Specs[0].End())

		if _, ok := c.Node().(*ast.GenDecl); ok {
			if once {
				// Generics Getter
				c.InsertAfter(&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "g",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.IndexExpr{
										X: &ast.Ident{
											Name: "GenHeader",
										},
										Index: &ast.Ident{
											Name: "T",
										},
									},
								},
							},
						},
					},
					Name: &ast.Ident{
						Name: "GetGI",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								&ast.Field{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "i",
										},
									},
									Type: &ast.Ident{
										Name: "int",
									},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								&ast.Field{
									Type: &ast.Ident{
										Name: "T",
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.IndexExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "g",
											},
											Sel: &ast.Ident{
												Name: "lists",
											},
										},
										Index: &ast.Ident{
											Name: "i",
										},
									},
								},
							},
						},
					},
				})

				// Generics setter
				c.InsertAfter(&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "g",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.IndexExpr{
										X: &ast.Ident{
											Name: "GenHeader",
										},
										Index: &ast.Ident{
											Name: "T",
										},
									},
								},
							},
						},
					},
					Name: &ast.Ident{
						Name: "SetGI",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								&ast.Field{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "i",
										},
									},
									Type: &ast.Ident{
										Name: "int",
									},
								},
								&ast.Field{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "x",
										},
									},
									Type: &ast.Ident{
										Name: "T",
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.IndexExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "g",
											},
											Sel: &ast.Ident{
												Name: "lists",
											},
										},
										Index: &ast.Ident{
											Name: "i",
										},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.Ident{
										Name: "x",
									},
								},
							},
						},
					},
				})

				// Generics struct
				c.InsertAfter(&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "GenHeader",
							},
							TypeParams: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "T",
											},
										},
										Type: &ast.Ident{
											Name: "any",
										},
									},
								},
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										&ast.Field{
											Names: []*ast.Ident{
												&ast.Ident{
													Name: "lists",
												},
											},
											Type: &ast.ArrayType{
												Elt: &ast.Ident{
													Name: "T",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				)
			}
			once = false
		}
		return true
	}, nil)

	newFile, err := os.Create("newFile.go")
	if err != nil {
		log.Fatal(err)
	}
	defer func(new *os.File) {
		err := new.Close()
		if err != nil {
		}
	}(newFile)

	if err := printer.Fprint(newFile, fset, n); err != nil {
		log.Fatal(err)
	}
}
