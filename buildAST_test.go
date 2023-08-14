package main

import (
	"reflect"
	"testing"
)

func BenchmarkTestBuildAst_getset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildAstDataStr(fileNameList[0])
	}
}

func BenchmarkTestBuildAst_array(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildAstDataStr(fileNameList[1])
	}
}

func BenchmarkTestBuildAst_array_getset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildAstDataStr(fileNameList[2])
	}
}

func BenchmarkTestBuildAst_eng_map(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildAstDataStr(fileNameList[3])
	}
}

func BenchmarkTestBuildAst_eng_reduce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildAstDataStr(fileNameList[4])
	}
}

func TestBuildAst(t *testing.T) {
	data := []struct {
		funcName string
		expected []basicStr
	}{
		{"hiddenDanger/main.go", []basicStr{{"", 0, []elem{{"*ast.File", []string{"main"}}}},
			{"", 15, []elem{{"*ast.GenDecl -> *ast.ImportSpec -> *ast.BasicLit", []string{"STRING", "\"fmt\""}}}},
			{"main", 29, []elem{
				{"", []string{"main"}},
				{"*ast.FuncType -> *ast.FieldList -> *ast.\n\t\t\tBlockStmt -> *ast.ExprStmt -> *ast.CallExpr -> *ast.SelectorExpr", []string{"fmt", "Println"}},
				{"*ast.BasicLit", []string{"STRING", "\"Hello world!\""}},
			}},
		}},
	}
	if result := buildAstDataStr(data[0].funcName); reflect.DeepEqual(result, data[0].expected) {
		t.Errorf("Result: %v, Expected: %v", result, data[0].expected)
	}

}
