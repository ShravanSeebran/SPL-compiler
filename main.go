package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"SPL-compiler/analyser"
	"SPL-compiler/lexer"
	"SPL-compiler/parser"
)

func main() {
	filename := getFilenameFromUser()
	program := readFromFile(filename)
	if lexer.Validate(program) {
		fmt.Println("Tokens accepted")
	} else {
		// TODO: Better error handling
		fmt.Println("Lexical error:")
		return
	}

	root, err := parser.Validate(program)
	if err != nil {
		fmt.Println("Syntax error:", err)
		return
	} else {
		fmt.Println("Syntax accepted")
	}

	if err := analyser.ValidateScoping(root); err != nil {
		fmt.Println("Naming error:", err)
		return
	} else {
		fmt.Println("Variable Naming and Function Naming accepted")
	}

	if err := analyser.ValidateTypeChecking(root); err != nil {
		fmt.Println("Type error:", err)
		return
	} else {
		fmt.Println("Types accepted")
	}

	if err := analyser.ValidateNoRecursion(root); err != nil {
		fmt.Println("Recursion detected error:", err)
		return
	} else {
		fmt.Println("No Recursion detected")
	}

	intermediateCode, err := analyser.ValidateCodeGeneration(root)
	if err != nil {
		fmt.Println("Intermediate Code Generation error:", err)
		return
	}
	generateHTML(intermediateCode, "output.html")

	basicCode, err := analyser.ValidateTranslateToBasic(intermediateCode)
	if err != nil {
		fmt.Println("BASIC Code Translation error:", err)
		return
	} else {
		writeToFile("output.txt", strings.Join(basicCode, "\n"))
		fmt.Println("Basic code generated successfully")
	}
}

func getFilenameFromUser() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	var filename string
	fmt.Print("Enter the filename of the SPL program (with the .txt at the end): ")
	fmt.Scanln(&filename)
	return filename
}

func readFromFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
	return string(content)
}

func writeToFile(filename string, content string) {
	err := os.WriteFile(filename, []byte(content), 0o644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	}
}

// WriteInstructionsHTML writes a standalone HTML file displaying the given instructions.
// filename should include .html (e.g. "instructions.html").
// Returns an error if writing/parsing fails.
func generateHTML(instructions []string, filename string) error {
	const tpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Intermediate Instructions</title>
<style>
  body {
    font-family: "JetBrains Mono", monospace;
    background-color: #f9fafb;
    color: #111827;
    margin: 40px;
  }
  h1 {
    text-align: center;
  }
  .meta {
    text-align: center;
    color: #6b7280;
    font-size: 0.9em;
    margin-bottom: 20px;
  }
  pre {
    background-color: #f3f4f6;
    padding: 16px;
    border-radius: 8px;
    overflow-x: auto;
    line-height: 1.6;
  }
  .line {
    counter-increment: line;
  }
  .line::before {
    content: counter(line) ": ";
    color: #9ca3af;
  }
</style>
</head>
<body>
  <h1>Intermediate Instructions</h1>
  <pre>
{{range .Instructions}}<div class="line">{{.}}</div>
{{end}}
  </pre>
</body>
</html>`

	t := template.Must(template.New("instructions").Parse(tpl))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		Instructions []string
		GeneratedAt  string
	}{
		Instructions: instructions,
		GeneratedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	return t.Execute(file, data)
}
