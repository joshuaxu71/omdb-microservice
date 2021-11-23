package main

import (
	"fmt"
	"strings"
)

// Question 3
func main() {
	str := findFirstStringInBracket("te(st)(a)sd")
	fmt.Println(str)
}

func findFirstStringInBracket(str string) string {
	if len(str) > 0 {
		openBracketIdx := strings.Index(str, "(")
		closeBracketIdx := strings.Index(str, ")")

		// checking if there's a string in bracket
		// if openBracketIdx is >= 0, there's an open bracket
		// if closeBracketIdx is > openBracketIdx + 1, there's at least
		// a character in the bracket.
		// if there's both brackets but no character in between, it
		// will just return an empty string as well
		if openBracketIdx >= 0 && closeBracketIdx > openBracketIdx+1 {
			return str[openBracketIdx+1 : closeBracketIdx]
		}
	}
	return ""
}
