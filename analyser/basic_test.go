package analyser

import (
	"fmt"
	"strings"
	"testing"

	"SPL-compiler/lexer"
	"SPL-compiler/parser"
)

func TestBasic(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Testing ...", `
		glob { }
	proc { }
	func { f(n) { local {a b c} b = n; return b} }
	main {
		var {x}
			x = 5;
		while (x > 0) {
			print x;
			x = (x minus 1)
		};
		halt
	}
				`},
	}
	for _, tt := range tests {
		fmt.Println("\n-------------- ", tt.name, " --------------")
		fmt.Println(tt.input)
		lexer.PrintTokensInline(lexer.TokenizeInput(tt.input))
		ast := parser.GenerateAST(tt.input)
		parser.PrettyPrintASTNode(ast, "", true)
		testBasic(ast)
	}
}

func testBasic(ast *parser.ASTNode) {
	AnalyseProgram(ast)
	TypeCheckProgram(ast)
	CheckRecursion(ast)
	lines := GenerateProgram(ast)
	PrettyPrintSymbolTable(symbolTable)
	fmt.Println("Generated Code:")
	fmt.Println(strings.Join(lines, "\n"))
	fmt.Println("\n--- Translated to Basic ---")
	TranslateToBasic(lines)
	fmt.Println(strings.Join(lines, "\n"))
}
