package main

import (
	"reflect"
	"strings"
	"testing"
)

type funcNameAndList struct {
	funcName string
	caseList []string
}
type output struct {
	result     bool
	caseResult []funcNameAndList
}

func BenchmarkTestCheckSwitchStatement_array_getset(b *testing.B) {
	listFunctions := buildAstDataStr(fileNameList[2])
	var funcList []string
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			funcList = append(funcList, listFunctions[s].funcName)
		}
	}
	modListFunctions := checkSelectorExpr(listFunctions)
	for i := 0; i < b.N; i++ {
		checkSwitchStatement(fileNameList[2], modListFunctions)
	}
}

func BenchmarkTestCheckSwitchStatement_eng_reduce(b *testing.B) {
	listFunctions := buildAstDataStr(fileNameList[4])
	var funcList []string
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			funcList = append(funcList, listFunctions[s].funcName)
		}
	}
	modListFunctions := checkSelectorExpr(listFunctions)
	for i := 0; i < b.N; i++ {
		checkSwitchStatement(fileNameList[4], modListFunctions)
	}
}

func TestCheckSwitchStatement(t *testing.T) {

	data := []struct {
		filename string
		want     output
	}{
		{fileNameList[0], output{false, []funcNameAndList{}}},
		{fileNameList[1], output{false, []funcNameAndList{}}},
		{fileNameList[2], output{true, []funcNameAndList{
			{"Set",
				[]string{"reflect Bool", "reflect Int", "reflect Int8", "reflect Int16", "reflect Int32", "reflect Int64", "reflect Uint", "reflect Uint8", "reflect Uint16", "reflect Uint32", "reflect Uint64", "reflect Uintptr",
					"reflect Float32", "reflect Float64", "reflect Complex64", "reflect Complex128", "reflect String", "reflect UnsafePointer"},
			},
			{"Get",
				[]string{"reflect Bool", "reflect Int", "reflect Int8", "reflect Int16", "reflect Int32", "reflect Int64", "reflect Uint", "reflect Uint8", "reflect Uint16", "reflect Uint32", "reflect Uint64", "reflect Uintptr",
					"reflect Float32", "reflect Float64", "reflect Complex64", "reflect Complex128", "reflect String", "reflect UnsafePointer"},
			},
			{"Memset",
				[]string{"Bool", "Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Uintptr", "Float32", "Float64", "Complex64", "Complex128", "String", "UnsafePointer"},
			},
			{"memsetiter",
				[]string{"Bool", "Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Uintptr", "Float32", "Float64", "Complex64", "Complex128", "String", "UnsafePointer"},
			},
			{"Eq",
				[]string{"reflect Bool", "reflect Int", "reflect Int8", "reflect Int16", "reflect Int32", "reflect Int64", "reflect Uint", "reflect Uint8", "reflect Uint16", "reflect Uint32", "reflect Uint64", "reflect Uintptr",
					"reflect Float32", "reflect Float64", "reflect Complex64", "reflect Complex128", "reflect String", "reflect UnsafePointer"},
			},
		}}},
		{fileNameList[3], output{false, []funcNameAndList{}}},
		{fileNameList[4], output{true,
			[]funcNameAndList{
				{"ReduceFirst",
					[]string{"Bool", "Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Uintptr", "Float32", "Float64", "Complex64", "Complex128", "String", "UnsafePointer"},
				},
				{"ReduceLast",
					[]string{"Bool", "Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Uintptr", "Float32", "Float64", "Complex64", "Complex128", "String", "UnsafePointer"},
				},
				{"ReduceDefault",
					[]string{"Bool", "Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Uintptr", "Float32", "Float64", "Complex64", "Complex128", "String", "UnsafePointer"},
				},
				{"Reduce",
					[]string{"Bool", "Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "Uintptr", "Float32", "Float64", "Complex64", "Complex128", "String", "UnsafePointer"},
				},
			}}},
	}

	for _, d := range data {
		listFunctions := buildAstDataStr(d.filename)
		var funcList []string
		for s := range listFunctions {
			if listFunctions[s].funcName != "" {
				funcList = append(funcList, listFunctions[s].funcName)
			}
		}
		// check leaf of SelectorExpr and unsafe Pointer
		modListFunctions := checkSelectorExpr(listFunctions)
		resBool, resultCase := checkSwitchStatement(d.filename, modListFunctions)
		for i := range resultCase {
			if !strings.EqualFold(resultCase[i].funcName, d.want.caseResult[i].funcName) || !reflect.DeepEqual(resultCase[i].caseFiltered, d.want.caseResult[i].caseList) {
				t.Errorf("result: %v, %v, %v \n wanted: %v, %v, %v", resBool, resultCase[i].funcName, resultCase[i].caseFiltered, d.want.result, d.want.caseResult[i].funcName, d.want.caseResult[i].caseList)
			}
		}
		if resBool != d.want.result || reflect.DeepEqual(resultCase, d.want.caseResult) {
		}
	}
}
