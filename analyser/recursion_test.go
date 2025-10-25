package analyser

import (
	"SPL-compiler/parser"
	"fmt"
	"testing"
)

func TestRecursion(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Testing ...", `
glob { }
proc { }
func { 
	f() { local { x } x = double(x); return 0}
	double (n) {
	  local { res }
		res = f();
	  return res
	}
}
main {
  var { }
  halt
}
		`},
	}
	for _, tt := range tests {
		fmt.Println("\n-------------- ", tt.name, " --------------")
		fmt.Println(tt.input)
		ast := parser.GenerateAST(tt.input)
		testRecursion(ast)
	}
}

func testRecursion(ast *parser.ASTNode) {
	AnalyseProgram(ast)
	TypeCheckProgram(ast)
	CheckRecursion(ast)
	PrettyPrintSymbolTable(symbolTable)
}
