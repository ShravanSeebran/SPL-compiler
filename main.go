package main

import (
	"SPL-compiler/lexer"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("SPL (Students' Programming Language) Lexer")
	fmt.Println("==========================================")

	// Run examples
	if len(os.Args) == 1 {
		runExamples()
		return
	}

	//  Read from file
	if os.Args[1] == "-f" || os.Args[1] == "--file" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go -f <filename>")
			return
		}
		runFromFile(os.Args[2])
		return
	}

}

func runExamples() {
	examples := []struct {
		name  string
		input string
	}{
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

	for i, example := range examples {
		fmt.Printf("\n=== Example %d: %s ===\n", i+1, example.name)
		fmt.Printf("Input:\n%s\n\n", example.input)
		tokenizeAndPrint(example.input)
		fmt.Println("\n" + strings.Repeat("-", 60))
	}
}

func runFromFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("Tokenizing file: %s\n", filename)
	fmt.Printf("Content:\n%s\n\n", string(content))
	tokenizeAndPrint(string(content))
}

func tokenizeAndPrint(input string) {
	tokens := lexer.TokenizeInput(input)

	fmt.Println("Tokens:")
	fmt.Println("-------")

	for i, token := range tokens {
		if token.Type == lexer.EOF {
			fmt.Printf("%3d: %-12s %-15s (Line: %d, Col: %d)\n",
				i+1, "EOF", "EOF", token.Line, token.Column)
			break
		}

		// Color coding for different token types (if terminal supports it)
		typeStr := colorizeTokenType(string(token.Type))
		literal := token.Literal
		if token.Type == lexer.STRING {
			literal = fmt.Sprintf("\"%s\"", literal)
		}

		fmt.Printf("%3d: %-12s %-15s (Line: %d, Col: %d)\n",
			i+1, typeStr, literal, token.Line, token.Column)
	}

	// Print summary
	validTokens := len(tokens) - 1 // Exclude EOF
	fmt.Printf("\nSummary: %d tokens processed\n", validTokens)
}

// Simple color coding for terminal output (optional)
func colorizeTokenType(tokenType string) string {
	// ANSI color codes (may not work in all terminals)
	switch tokenType {
	case "IDENT":
		return "\033[36m" + tokenType + "\033[0m" // Cyan
	case "INT":
		return "\033[33m" + tokenType + "\033[0m" // Yellow
	case "STRING":
		return "\033[32m" + tokenType + "\033[0m" // Green
	default:
		if isKeyword(tokenType) {
			return "\033[35m" + tokenType + "\033[0m" // Magenta
		}
		return tokenType
	}
}

func isKeyword(tokenType string) bool {
	keywords := []string{"glob", "proc", "func", "main", "var", "local", "return",
		"halt", "print", "while", "do", "until", "if", "else", "eq", "or", "and",
		"plus", "minus", "mult", "div", "neg", "not"}

	for _, keyword := range keywords {
		if tokenType == keyword {
			return true
		}
	}
	return false
}
