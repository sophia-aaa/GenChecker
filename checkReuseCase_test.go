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

func basicTestExecution_Switch(filename string) []basicStr {
	listFunctions := buildAstDataStr(filename)
	var funcList []string
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			funcList = append(funcList, listFunctions[s].funcName)
		}
	}
	return checkSelectorExpr(listFunctions)
}
func BenchmarkTestCheckSwitchStatement_getset(b *testing.B) {
	testCase := basicTestExecution_Switch(fileNameList[0])
	for i := 0; i < b.N; i++ {
		checkSwitchStatement(fileNameList[0], testCase)
	}
}

func BenchmarkTestCheckSwitchStatement_array(b *testing.B) {
	testCase := basicTestExecution_Switch(fileNameList[1])
	for i := 0; i < b.N; i++ {
		checkSwitchStatement(fileNameList[1], testCase)
	}
}

func BenchmarkTestCheckSwitchStatement_array_getset(b *testing.B) {
	testCase := basicTestExecution_Switch(fileNameList[2])
	for i := 0; i < b.N; i++ {
		checkSwitchStatement(fileNameList[2], testCase)
	}
}

func BenchmarkTestCheckSwitchStatement_eng_map(b *testing.B) {
	testCase := basicTestExecution_Switch(fileNameList[3])
	for i := 0; i < b.N; i++ {
		checkSwitchStatement(fileNameList[3], testCase)
	}
}
func BenchmarkTestCheckSwitchStatement_eng_reduce(b *testing.B) {
	testCase := basicTestExecution_Switch(fileNameList[4])
	for i := 0; i < b.N; i++ {
		checkSwitchStatement(fileNameList[4], testCase)
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
		testCase := basicTestExecution_Switch(d.filename)
		resBool, resultCase := checkSwitchStatement(d.filename, testCase)
		for i := range resultCase {
			if !strings.EqualFold(resultCase[i].funcName, d.want.caseResult[i].funcName) || !reflect.DeepEqual(resultCase[i].caseFiltered, d.want.caseResult[i].caseList) {
				t.Errorf("result: %v, %v, %v \n wanted: %v, %v, %v", resBool, resultCase[i].funcName, resultCase[i].caseFiltered, d.want.result, d.want.caseResult[i].funcName, d.want.caseResult[i].caseList)
			}
		}
		if resBool != d.want.result || reflect.DeepEqual(resultCase, d.want.caseResult) {
		}
	}
}
