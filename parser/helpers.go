package parser

import (
	"fmt"
)

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
