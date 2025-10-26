package analyser

import (
	"fmt"
	"os"
	"testing"

	"SPL-compiler/lexer"
	"SPL-compiler/parser"
)

func TestSymbolTable(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Testing redeclaration of global variables in main scope", `
		  glob { x y }
		  proc {}
		  func {}
		  main {
			var { x z }
			halt
		  }
		`},
		{"Testing using global variable in local function scope", `
		  glob { x y }
		  proc { f(a b) { local { c } x = a } g(a b) { local { c } x = a } }
		  func {}
		  main {
			var { x z }
			x = ( x plus z );
			z = f(x z);
			halt
		  }
		`},
	}
	for _, tt := range tests {
		fmt.Println("\n-------------- ", tt.name, " --------------")
		fmt.Println(tt.input)
		ast := testParse(tt.input)
		testProgram(ast)

	}
}

func testParse(input string) *parser.ASTNode {
	l := lexer.New(input)
	lexerAdapter := &parser.LexerAdapter{L: l}
	result, err := parser.Parse(lexerAdapter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parsing error: %v\n", err)
		os.Exit(1)
	}

	if result != nil {
		fmt.Println("\n--- Abstract Syntax Tree ---")
		parser.PrintAST(result, 0)
	} else {
		fmt.Fprintf(os.Stderr, "Parsing finished, but no AST was generated.\n")
	}

	return parser.ResultAST
}

func testProgram(ast *parser.ASTNode) {
	AnalyseProgram(ast)
	fmt.Println("\n---Symbol Table ---")
	PrettyPrintSymbolTable(symbolTable)
}
