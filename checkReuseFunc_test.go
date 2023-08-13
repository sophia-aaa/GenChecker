package main

import (
	"strings"
	"testing"
)

func basicTestExecution_Func(filename string) ([]basicStr, []string) {
	listFunctions := buildAstDataStr(filename)
	var funcList []string
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			funcList = append(funcList, listFunctions[s].funcName)
		}
	}
	return checkSelectorExpr(listFunctions), funcList
}

func BenchmarkTestCheckFunc_getset(b *testing.B) {
	modListFunctions, funcList := basicTestExecution_Func(fileNameList[0])
	for i := 0; i < b.N; i++ {
		checkGenerics(modListFunctions, funcList, typeList)
	}
}

func BenchmarkTestCheckFunc_array(b *testing.B) {
	modListFunctions, funcList := basicTestExecution_Func(fileNameList[1])
	for i := 0; i < b.N; i++ {
		checkGenerics(modListFunctions, funcList, typeList)
	}
}

func BenchmarkTestCheckFunc_array_getset(b *testing.B) {
	modListFunctions, funcList := basicTestExecution_Func(fileNameList[2])
	for i := 0; i < b.N; i++ {
		checkGenerics(modListFunctions, funcList, typeList)
	}
}

func BenchmarkTestCheckFunc_eng_map(b *testing.B) {
	modListFunctions, funcList := basicTestExecution_Func(fileNameList[3])
	for i := 0; i < b.N; i++ {
		checkGenerics(modListFunctions, funcList, typeList)
	}
}
func BenchmarkTestCheckFunc_eng_reduce(b *testing.B) {
	modListFunctions, funcList := basicTestExecution_Func(fileNameList[4])
	for i := 0; i < b.N; i++ {
		checkGenerics(modListFunctions, funcList, typeList)
	}
}

func TestCheckFunc(t *testing.T) {

	data := []struct {
		filename string
		want     [][]string
	}{
		{fileNameList[0], [][]string{{"Bools", "Ints", "Int8s", "Int16s", "Int32s", "Int64s", "Uints", "Uint8s",
			"Uint16s", "Uint32s", "Uint64s", "Uintptrs", "Float32s", "Float64s", "Complex64s", "Complex128s", "Strings", "UnsafePointers"},
			{"SetB", "SetI", "SetI8", "SetI16", "SetI32", "SetI64", "SetU", "SetU8", "SetU16", "SetU32", "SetU64", "SetUintptr",
				"SetF32", "SetF64", "SetC64", "SetC128", "SetStr", "SetUnsafePointer"},
			{"GetB", "GetI", "GetI8", "GetI16", "GetI32", "GetI64", "GetU", "GetU8", "GetU16", "GetU32", "GetU64",
				"GetUintptr", "GetF32", "GetF64", "GetC64", "GetC128", "GetStr", "GetUnsafePointer"}}},
		{fileNameList[1], [][]string{{"Len", "Cap"}}},
		{fileNameList[2], [][]string{}},
		{fileNameList[3], [][]string{}},
		{fileNameList[4], [][]string{}},
	}

	for _, d := range data {

		modListFunctions, funcList := basicTestExecution_Func(d.filename)
		genCheck := checkGenerics(modListFunctions, funcList, typeList)

		count := 0
		for i := range genCheck {
			if len(genCheck[i]) > 1 {
				if len(genCheck[i]) != len(d.want[count]) {
					t.Errorf("%v\n%v\n\n", genCheck[i], d.want[count])
					t.Errorf("These two results are not equally long.\nresult: %v \t wanted: %v", len(genCheck[i]), len(d.want[count]))

				} else if len(genCheck[i]) > 0 && len(d.want[count]) > 0 {
					//pattern1 = true
					for ind, val := range genCheck[i] {
						if !strings.EqualFold(val.funcName, d.want[count][ind]) {
							t.Errorf("result: %v\nwanted: %v\n", val.funcName, d.want[count][ind])
						}
					}
					count++
				}
			}
		}

	}
}
