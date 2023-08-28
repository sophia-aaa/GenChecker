package main

import (
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
