package main

import (
	"SPL-compiler/lexer"
	"SPL-compiler/parser"
	"fmt"
	"os"
)

func main() {
	program := readFromFile("spl.txt")
	if lexer.Validate(program) {
		fmt.Println("Tokens accepted")
	} else {
		// TODO: Better error handling
		fmt.Println("Lexical error:")
	}

	if err := parser.Validate(program); err != nil {
		fmt.Println("Syntax error:", err)
	} else {
		fmt.Println("Syntax accepted")
	}
}

func readFromFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
	return string(content)
}

func writeToFile(filename string, content string) {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	}
}
