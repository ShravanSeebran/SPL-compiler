package analyser

import (
	"SPL-compiler/lexer"
	"SPL-compiler/parser"
	"fmt"
	"strings"
	"testing"
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
		{
			"Basic Program Structure",
			`glob { 
				counter 
				max 
			} 
			main { 
				var { temp } 
				counter = 0; 
				max = 10; 
				halt 
			}`,
		},
		{
			"Function Definition",
			`func { 
				add(x y) { 
					local { result } 
					result = (x plus y); 
					return result 
				} 
			}`,
		},
		{
			"Control Flow",
			`if (x > 0) { 
				print "positive" 
			} else { 
				print "zero or negative" 
			}`,
		},
		{
			"Loop Example",
			`while (counter > 0) { 
				print counter; 
				counter = (counter minus 1); 
			}`,
		},
		{
			"Complete Small Program",
			`glob { x y }
			proc {
				printsum(a b) {
					local { sum }
					sum = (a plus b);
					print sum;
				}
			}
			main {
				var { result }
				x = 5;
				y = 3;
				printsum(x y);
				halt
			}`,
		},
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
