package analyser

import (
	"SPL-compiler/parser"
	"fmt"
)

// func TestChecker(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		input string
// 	}{
// 		{"Testing ...", `
// 		  glob { x y }
// 		  proc {}
// 		  func {}
// 		  main {
// 			var { x z }
// 			halt
// 		  }
// 		`},
// 		{"Testing ...", `
// 		  glob { x y }
// 		  proc {}
// 		  func {}
// 		  main {
// 			var { x z }
// 			x = ( x plus z );
// 			halt
// 		  }
// 		`},
// 	}
// 	for _, tt := range tests {
// 		fmt.Println("\n-------------- ", tt.name, " --------------")
// 		fmt.Println(tt.input)
// 		ast := generateAST(tt.input)
// 		testChecker(ast)
//
// 	}
// }

func testChecker(ast *parser.ASTNode) {
	AnalyseProgram(ast)
	TypeCheckProgram(ast)
	fmt.Println("Type checking finished.")
}
