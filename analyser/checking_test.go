package analyser

import (
	"SPL-compiler/lexer"
	"SPL-compiler/parser"
	"fmt"
	"os"
	"testing"
)

func TestChecker(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Testing ...", `
		  glob { x y }
		  proc {}
		  func {}
		  main { 
			var { x z } 
			halt 
		  }
		`},
		{"Testing ...", `
		  glob { x y }
		  proc {}
		  func {}
		  main { 
			var { x z } 
			x = ( x plus z );
			halt 
		  }
		`},
	}
	for _, tt := range tests {
		fmt.Println("\n-------------- ", tt.name, " --------------")
		fmt.Println(tt.input)
		ast := generateAST(tt.input)
		testChecker(ast)

	}
}

func generateAST(input string) *parser.ASTNode {
	l := lexer.New(input)
	lexerAdapter := &parser.LexerAdapter{L: l}
	result, err := parser.Parse(lexerAdapter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parsing error: %v\n", err)
		os.Exit(1)
	}

	if result != nil {
		fmt.Println("Parsing finished, AST generated.")
	} else {
		fmt.Fprintf(os.Stderr, "Parsing finished, but no AST was generated.\n")
	}

	return parser.ResultAST
}

func testChecker(ast *parser.ASTNode) {
	AnalyseProgram(ast)
	TypeCheckProgram(ast)
	fmt.Println("Type checking finished.")
}
