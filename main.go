package main

import (
	"go/ast"
	_ "golang.org/x/exp/slices"
	"os"
)

type pattern3Result struct {
	nodes    *ast.FuncDecl
	funcList []funcNamePos
}

type patternReplace struct {
	nodes      []*ast.FuncDecl
	funcRemove []funcNamePos
}

func main() {
	// filename, err := os.ReadFile(os.Args[2])
	// command must be like this: go run . - "hiddenDanger/getset.go"
	filename := os.Args[2]

	refactoring(filename)

}
