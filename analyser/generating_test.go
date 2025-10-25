package analyser

import (
	"SPL-compiler/lexer"
	"SPL-compiler/parser"
	"fmt"
	"strings"
	"testing"
)

func TestGenerator(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Testing ...", `
		glob { x y }
		proc {
			f(a b) { local { c } x = a }
			}	
		func {}
		main {
		var { x z }
		x = ( x plus z );
		halt
		}
				`},
		{"Testing ...", `
		glob { x y }
		proc {
			f(a b) { local { c } x = a; b = c }
			g(a b) { local { c } x = a }
		}
		func {
			h(a b) { local { c } x = a ; return x }
		}
		main {
		var { x z }
		x = ( x plus ((neg z) mult 3) );
		y = ( x div 3 );
		print y;
		print x;
		print "yay it works!";
		halt
		}
				`},
		// 		{"Testing ...", `
		// glob { x y }
		// proc {
		// 	hello(a) { local { b } b = a; }
		// 	world(b) { local { c } c = b; }
		// }
		// func {
		// 	f(a) {
		//     local {} x = a; return x
		// 	}
		// 	double (n) {
		// 	  local { res }
		// 	  res = ( n mult 2 );
		//     res = f(n);
		// 	  return res
		// 	}
		// }
		// main {
		//   var { res z }
		//   res = hello(z);
		//   halt
		// }
		// 		`},
		// {"Testing expected error of having treating pdef as fdef...", `
		// glob { x y }
		// proc {
		// 	f(a b) { local { c } x = a; b = c }
		// }
		// func {}
		// main {
		// var { x z }
		// x = f(x z);
		// halt
		// }`},
		{"Testing ...", `
		glob { x y }
		proc {}
		func {
			f(a b) { local { c } x = a; b = 10; return b }
		}
		main {
		var { x z }
		x = f(x z);
			halt;
	    print x
		}`},
	}
	for _, tt := range tests {
		fmt.Println("\n-------------- ", tt.name, " --------------")
		fmt.Println(tt.input)
		lexer.PrintTokensInline(lexer.TokenizeInput(tt.input))
		ast := parser.GenerateAST(tt.input)
		parser.PrettyPrintASTNode(ast, "", true)
		testGenerator(ast)
	}
}

func testGenerator(ast *parser.ASTNode) {
	AnalyseProgram(ast)
	TypeCheckProgram(ast)
	CheckRecursion(ast)
	lines := GenerateProgram(ast)
	PrettyPrintSymbolTable(symbolTable)
	fmt.Println("Generated Code:")
	fmt.Println(strings.Join(lines, "\n"))
}
