package main

import (
	"fmt"
	"log"
	"os"
)

// show how the AST is constructed

func createTextFile(filename string, listFunctions []basicStr) {
	f, err := os.Create(filename + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	var str string
	nextLine := "\n\n"
	for s := range listFunctions {
		if listFunctions[s].funcName != "" {
			str = "function name is " + listFunctions[s].funcName + " \n"
			_, err2 := f.WriteString(str)
			if err2 != nil {
				log.Fatal(err2)
			}
			for _, value := range listFunctions[s].value {
				str = fmt.Sprintln("\t", value)
				_, err2 := f.WriteString(str)
				if err2 != nil {
					log.Fatal(err2)
				}
			}
			_, err2 = f.WriteString(nextLine)
			if err2 != nil {
				log.Fatal(err2)
			}

		} else {
			for _, value := range listFunctions[s].value {
				str = fmt.Sprintln(value)
				_, err2 := f.WriteString(str)
				if err2 != nil {
					log.Fatal(err2)
				}
			}
			_, err2 := f.WriteString(nextLine)
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	}
}

func createTextFileFromString(filename string, strList []string) {
	f, err := os.Create(filename + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)
	if len(strList) > 0 {
		for _, val := range strList {
			_, err1 := f.WriteString(val)
			if err1 != nil {
				log.Fatal(err1)
			}
		}
	}
}

func createCasesTextFile(filename string, funcList []caseResult) {
	f, err := os.Create(filename + "_cases.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	var str string

	for _, val := range funcList {
		str = "function name: " + val.funcName + "\n"
		_, err2 := f.WriteString(str)
		if err2 != nil {
			log.Fatal(err2)
		}
		for _, value := range val.cases {
			str = fmt.Sprintln("  ", value.caseName)
			_, err2 := f.WriteString(str)
			if err2 != nil {
				log.Fatal(err2)
			}
			for _, cases := range value.value {
				str = fmt.Sprintln("    ", cases.path, "\t", cases.value)
				_, err2 := f.WriteString(str)
				if err2 != nil {
					log.Fatal(err2)
				}
			}
		}
		str = fmt.Sprintln("")
		_, err2 = f.WriteString(str)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}
