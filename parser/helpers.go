package parser

import (
	"errors"
	"fmt"
	"os"

	"SPL-compiler/lexer"
)

func Validate(input string) (root *ASTNode, err error) {
	defer func() {
		if r := recover(); r != nil {
			// TODO: Better error handling
			err = errors.New("syntax error")
		}
	}()

	root = GenerateAST(input)
	return root, nil
}

func GenerateAST(input string) *ASTNode {
	l := lexer.New(input)
	lexerAdapter := &LexerAdapter{L: l}
	result, err := Parse(lexerAdapter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parsing error: %v\n", err)
		os.Exit(1)
	}

	if result != nil {
		fmt.Println("Parsing finished, AST generated.")
	} else {
		fmt.Fprintf(os.Stderr, "Parsing finished, but no AST was generated.\n")
	}

	return ResultAST
}

func PrettyPrintASTNode(n *ASTNode, prefix string, isTail bool) {
	if n == nil {
		return
	}

	connector := "├── "
	if isTail {
		connector = "└── "
	}
	fmt.Printf("%s%s[%d] %s", prefix, connector, n.ID, n.Type)
	if n.Name != "" {
		fmt.Printf(": %s", n.Name)
	}
	fmt.Println()

	childPrefix := prefix
	if isTail {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	for i, child := range n.Children {
		isLast := i == len(n.Children)-1
		PrettyPrintASTNode(child, childPrefix, isLast)
	}
}

func GetNodeByID(node *ASTNode, id int) *ASTNode {
	if int(node.ID) == id {
		return node
	}
	for _, child := range node.Children {
		if result := GetNodeByID(child, id); result != nil {
			return result
		}
	}
	return nil
}

func GetDefNodeByNameID(node *ASTNode, id int) *ASTNode {
	for _, child := range node.Children {
		if int(child.ID) == id {
			return node
		}
	}
	for _, child := range node.Children {
		if result := GetDefNodeByNameID(child, id); result != nil {
			return result
		}
	}
	return nil
}
