package analyser

import (
	"fmt"
	"strings"
	"testing"

	"SPL-compiler/lexer"
	"SPL-compiler/parser"
)

func TestGenerator(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// 		{"Testing ...", `
		// glob { x y }
		// proc {}
		// func {}
		// main {
		// var { x z }
		// x = ( x plus z );
		// halt
		// }
		// 		`},
		// 		{"Testing ...", `
		// glob { x y }
		// proc {
		// 	f(a b) { local { c } x = a }
		// 	g(a b) { local { c } x = a }
		// }
		// func {
		// 	h(a b) { local { c } x = a ; return x }
		// }
		// main {
		// var { x z }
		// x = ( x plus ((neg z) mult 3) );
		// y = ( x div 3 );
		// print y;
		// print x;
		// print "yay it works!";
		// halt
		// }
		// 		`},
		{"Testing ...", `
glob { x y }
proc { }
func { 
	f(a) { local {} x = a; return x } 
}
main {
halt
}
		`},
	}
	// double (n) {
	//   local { res }
	//   res = ( n mult 2 );
	//   return res
	// }
	// x = double(y);
	for _, tt := range tests {
		fmt.Println("\n-------------- ", tt.name, " --------------")
		fmt.Println(tt.input)
		lexer.PrintTokensInline(lexer.TokenizeInput(tt.input))
		ast := generateAST(tt.input)
		testGenerator(ast)
	}
}

func testGenerator(ast *parser.ASTNode) {
	AnalyseProgram(ast)
	TypeCheckProgram(ast)
	lines := GenerateProgram(ast)
	PrettyPrintSymbolTable(symbolTable)
	fmt.Println("Generated Code:")
	fmt.Println(strings.Join(lines, "\n"))
}
