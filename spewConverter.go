package main

import (
	"fmt"
	"strings"
)

func tokenConversion(tok string) string {
	switch tok {
	case "ILLEGAL,", "ILLEGAL":
		return "token.ILLEGAL,"
	case "EOF,", "EOF":
		return "token.EOF,"
	case "COMMENT,", "COMMENT":
		return "token.COMMENT,"
	case "IDENT,", "IDENT":
		return "token.IDENT,"
	case "INT,", "INT":
		return "token.INT,"
	case "FLOAT,", "FLOAT":
		return "token.FLOAT,"
	case "IMAG,", "IMAG":
		return "token.IMAG,"
	case "CHAR,", "CHAR":
		return "token.CHAR,"
	case "STRING,", "STRING":
		return "token.STRING,"
	case "+,", "+":
		return "token.ADD,"
	case "-,", "-":
		return "token.SUB,"
	case "*,", "*":
		return "token.MUL,"
	case "/,", "/":
		return "token.QUO,"
	case "%,", "%":
		return "token.REM,"
	case "&,", "&":
		return "token.AND,"
	case "|,", "|":
		return "token.OR,"
	case "^,", "^":
		return "token.XOR,"
	case "<<,", "<<":
		return "token.SHL,"
	case ">>,", ">>":
		return "token.SHR,"
	case "&^,", "&^":
		return "token.AND_NOT,"
	case "+=,", "+=":
		return "token.ADD_ASSIGN,"
	case "-=,", "-=":
		return "token.SUB_ASSIGN,"
	case "*=,", "*=":
		return "token.MUL_ASSIGN,"
	case "/=,", "/=":
		return "token.QUO_ASSIGN,"
	case "%=,", "%=":
		return "token.REM_ASSIGN,"
	case "&=,", "&=":
		return "token.AND_ASSIGN,"
	case "|=,", "|=":
		return "token.OR_ASSIGN,"
	case "^=,", "^=":
		return "token.XOR_ASSIGN,"
	case "<<=,", "<<=":
		return "token.SHL_ASSIGN,"
	case ">>=,", ">>=":
		return "token.SHR_ASSIGN,"
	case "&^=,", "&^=":
		return "token.AND_NOT_ASSIGN,"
	case "&&,", "&&":
		return "token.LAND,"
	case "||,", "||":
		return "token.LOR,"
	case "<-,", "<-":
		return "token.ARROW,"
	case "++,", "++":
		return "token.INC,"
	case "--,", "--":
		return "token.DEC,"
	case "==,", "==":
		return "token.EQL,"
	case "<,", "<":
		return "token.LSS,"
	case ">,", ">":
		return "token.GTR,"
	case "=,", "=":
		return "token.ASSIGN,"
	case "!,", "!":
		return "token.NOT,"
	case "!=,", "!=":
		return "token.NEQ,"
	case "<=,", "<=":
		return "token.LEQ,"
	case ">=,", ">=":
		return "token.GEQ,"
	case ":=,", ":=":
		return "token.DEFINE,"
	case "...,", "...":
		return "token.ELLIPSIS,"
	case "(,", "(":
		return "token.LPAREN,"
	case "[,", "[":
		return "token.LBRACK,"
	case "{,", "{":
		return "token.LBRACE,"
	case ",,", ",":
		return "token.COMMA,"
	case ".,", ".":
		return "token.PERIOD,"
	case "),", ")":
		return "token.RPAREN,"
	case "],", "]":
		return "token.RBRACK,"
	case "},", "}":
		return "token.RBRACE,"
	case ";,", ";":
		return "token.SEMICOLON,"
	case ":,", ":":
		return "token.COLON,"
	case "break,", "break":
		return "token.BREAK,"
	case "case,", "case":
		return "token.CASE,"
	case "chan,", "chan":
		return "token.CHAN,"
	case "const,", "const":
		return "token.CONST,"
	case "continue,", "continue":
		return "token.CONTINUE,"
	case "default,", "default":
		return "token.DEFAULT,"
	case "defer,", "defer":
		return "token.DEFER,"
	case "else,", "else":
		return "token.ELSE,"
	case "fallthrough,", "fallthrough":
		return "token.FALLTHROUGH,"
	case "for,", "for":
		return "token.FOR,"
	case "func,", "func":
		return "token.FUNC,"
	case "go,", "go":
		return "token.GO,"
	case "goto,", "goto":
		return "token.GOTO,"
	case "if,", "if":
		return "token.IF,"
	case "import,", "import":
		return "token.IMPORT,"
	case "interface,", "interface":
		return "token.INTERFACE,"
	case "map,", "map":
		return "token.MAP,"
	case "package,", "package":
		return "token.PACKAGE,"
	case "range,", "range":
		return "token.RANGE,"
	case "return,", "return":
		return "token.RETURN,"
	case "select,", "select":
		return "token.SELECT,"
	case "struct,", "struct":
		return "token.STRUCT,"
	case "switch,", "switch":
		return "token.SWITCH,"
	case "type,", "type":
		return "token.TYPE,"
	case "var,", "var":
		return "token.VAR,"
	case "~,", "~":
		return "token.TILDE,"
	}
	return ""
}

// This function helps to extract Nodes to add, delete, or replace codes via the ast.Inspect or astutil packages
func spewConversion(str string) string {
	// first modified
	// add second method
	// to improve on ,, }, token.Token
	trimString := strings.Split(str, "\n")
	trimmedStr := ""
	for _, val := range trimString {
		if strings.Contains(val, "token.Pos") || strings.Contains(val, "TokPos") || strings.Contains(val, "<nil>") {
			continue
		}
		trimmedStr += val + "\n"

	}
	//fmt.Println("************************")
	//fmt.Println(trimmedStr)
	//fmt.Println("************************")
	// split must be with this symbol ) *******
	substrings := strings.Split(trimmedStr, ")")

	var candidate []string
	identSelCheck := false
	//codeBuilder := ""
	newString := ""
	for _, val := range substrings {

		fmt.Println("------------------------------")
		fmt.Println(val)
		//out := []rune(val)
		for ind, c := range val {
			//ch, _ := utf8.DecodeRuneInString(c)
			if c == '(' || c == ')' || c == ' ' || c == '\t' {
				continue
			} else if ind != 0 && rune(val[ind-1]) != ']' && c == '*' && ind < len(val) {
				newString += string(rune('&'))
			} else if ind == 0 && c == '*' {
				newString += string(rune('&'))
			} else if c == '{' {
				newString += string(rune(' ')) + string(c)
			} else if c == ':' && ind != 0 && rune(val[ind-1]) != ' ' {
				newString += string(c) + string(rune(' '))
			} else {
				newString += string(c)
			}
		}
		if len(newString) > 2 && strings.EqualFold(newString[0:2], "0x") {
			newString = ""
			continue
		} else if (strings.Contains(newString, "&ast.Ident") || strings.Contains(newString, "Sel")) && !strings.Contains(newString, "SelectorExpr") {
			candidate = append(candidate, newString)
			identSelCheck = true
		} else if identSelCheck {
			identString := " {\n Name: " + string(rune('"')) + newString + string(rune('"')) + ", \n}"
			candidate = append(candidate, identString)
			identString = ""
			identSelCheck = false
		} else if strings.Contains(newString, "len=") {
			newString = ""
			continue
		} else if strings.Contains(newString, "token.Token") || strings.Contains(newString, "&ast.BasicLit") || strings.Contains(newString, "[]ast.Expr") {
			newString += " "
			candidate = append(candidate, newString)
		} else {
			newString += " "
			candidate = append(candidate, newString)
		}
		fmt.Println("---->")
		fmt.Println(newString)
		newString = ""
	}

	interim := ""
	for _, val := range candidate {
		interim += val
	}

	// last process for spew Converter
	trimString = strings.Split(interim, "\n")
	result := ""
	toAdd := ""
	flagBool := false
	flagValue := false
	flagToken := false
	for _, val := range trimString {
		if strings.Contains(val, "token.Token") {
			flagToken = true
			strCheck := strings.Split(val, " ")
			tokFlag := false
			for _, strVal := range strCheck {
				if strings.EqualFold(strVal, "token.Token") {
					tokFlag = true
				} else if tokFlag {
					toAdd += tokenConversion(strVal)
					//fmt.Println(val)
					tokFlag = false
				} else {
					toAdd += strVal + " "
				}
			}
		}
		if /*strings.Contains(val, "Value:") || */ strings.Contains(val, "Incomplete:") || strings.Contains(val, "Text:") {
			flagValue = true
			strCheck := strings.Split(val, " ")
			for ind, strVal := range strCheck {
				if ind == 0 {
					toAdd += strVal + string(rune(' '))
				} else if ind == len(strCheck)-1 {
					toAdd += strVal + string(rune(','))
				}
			}
		}
		if strings.Contains(val, "Value: string") {
			flagBool = true
			strCheck := strings.Split(val, " ")
			for ind, strVal := range strCheck {
				if strings.Contains(strVal, "string") {
					continue
				} else if ind == len(strCheck)-1 {
					toAdd += strVal + ","
				} else {
					toAdd += strVal + " "
				}
			}
		}
		if strings.Contains(val, "0x") {
			// remove pointer address
			flagBool = true
			strCheck := strings.Split(val, " ")
			for _, strVal := range strCheck {
				if strings.Contains(strVal, "0x") {
					continue
				} else {
					toAdd += strVal + " "
				}
			}
		}
		// make } with ,
		if strings.Contains(val, "}") && !strings.Contains(val, "},") {
			flagBool = true
			for _, c := range val {
				//ch, _ := utf8.DecodeRuneInString(c)
				if c == ' ' {
					continue
				} else {
					toAdd += string(rune(c))
				}
			}
			if strings.EqualFold(toAdd, "}") {
				toAdd += ","
			}
		}
		// remove pointer address
		if strings.Contains(val, "len=") && strings.Contains(val, "cap=") {
			flagBool = true
			strCheck := strings.Split(val, " ")
			for _, strVal := range strCheck {
				if strings.Contains(strVal, "len=") && strings.Contains(strVal, "cap=") {
					continue
				} else {
					toAdd += strVal + " "
				}
			}
		}

		if flagBool || flagToken || flagValue {
			result += toAdd + "\n"
			toAdd = ""
			flagBool = false
			flagValue = false
			flagToken = false
		} else {
			result += val + "\n"
		}

	}

	//fmt.Println(result)
	return result
}
