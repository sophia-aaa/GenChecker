package main

import (
	"go/ast"
	"go/token"
)

var (
	// Name your slice here
	sliceName = "Raw"
	// Name your type here
	typeName = "T"

	typeList = []string{
		"bool", "bType", "int", "iType", "int8", "i8Type", "int16", "i16Type", "int32", "i32Type", "int64", "i64Type", "uint", "uType",
		"uint8", "u8Type", "uint16", "u16Type", "uint32", "u32Type", "uint64", "u64Type", "uintptr", "uintptrType",
		"float32", "f32Type", "float64", "f64Type", "complex64", "c64Type", "complex128", "c128Type", "string", "strType",
		"unsafe", "Pointer", "unsafePointer", "unsafePointerType", "reflect", "struct", "slice", "byte", "Get", "Set",
	}

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
	fileNameList = []string{
		"hiddenDanger/getset.go",
		"hiddenDanger/array.go",
		"hiddenDanger/array_getset.go",
		"hiddenDanger/eng_map.go",
		"hiddenDanger/eng_reduce.go",
	}
)
