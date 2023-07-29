package main

import (
	"go/ast"
	"go/token"
)

var (
	genStruct = &ast.GenDecl{
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
	}

	genGet = &ast.FuncDecl{
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
	}

	genSet = &ast.FuncDecl{
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
	}
)
