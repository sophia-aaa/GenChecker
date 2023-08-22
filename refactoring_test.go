package main

import "testing"

func BenchmarkRefactoring_getset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		refactoring(fileNameList[0])
	}
}

func BenchmarkRefactoring_array(b *testing.B) {
	for i := 0; i < b.N; i++ {
		refactoring(fileNameList[1])
	}
}

func BenchmarkRefactoring_array_getset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		refactoring(fileNameList[2])
	}
}

func BenchmarkRefactoring_eng_map(b *testing.B) {
	for i := 0; i < b.N; i++ {
		refactoring(fileNameList[3])
	}
}
func BenchmarkRefactoring_eng_reduce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		refactoring(fileNameList[4])
	}
}
