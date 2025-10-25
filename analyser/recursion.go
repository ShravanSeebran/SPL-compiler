package analyser

import (
	"slices"

	"SPL-compiler/parser"
)

var rootNode *parser.ASTNode

func CheckRecursion(root *parser.ASTNode) {
	rootNode = root
	procdefs := root.Children[1]
	funcdefs := root.Children[2]

	for len(procdefs.Children) > 0 {
		if checkDefForRecursion(
			procdefs.Children[0],
			[]string{symbolTable[int(procdefs.Children[0].Children[0].ID)].uniqueID},
		) {
			panic("Recursion detected in procedure definitions")
		}
		procdefs = procdefs.Children[1]
	}

	for len(funcdefs.Children) > 0 {
		if checkDefForRecursion(
			funcdefs.Children[0],
			[]string{symbolTable[int(funcdefs.Children[0].Children[0].ID)].uniqueID},
		) {
			panic("Recursion detected in function definitions")
		}
		funcdefs = funcdefs.Children[1]
	}
}

func checkDefForRecursion(node *parser.ASTNode, names []string) bool {
	if node.Type == "FDEF" || node.Type == "PFDEF" {
		body := node.Children[2]
		algo := body.Children[1]
		if checkAlgoForRecursion(algo, names) {
			return true
		}
	}
	return false
}

func checkAlgoForRecursion(node *parser.ASTNode, names []string) bool {
	for len(node.Children) > 1 {
		instr := node.Children[0]
		switch instr.Name {
		case "call":
			calledName := instr.Children[0]
			if slices.Contains(names, calledName.Name) {
				return true
			}
			nameDefID := symbolTable[int(calledName.ID)].declarationNode
			nameDefNode := parser.GetNodeByID(rootNode, nameDefID)
			if checkDefForRecursion(nameDefNode, append(names, calledName.Name)) {
				return true
			}
		case "assign":
			if checkAssignForRecursion(instr.Children[0], names) {
				return true
			}
		}
		node = node.Children[1]
	}
	instr := node.Children[0]
	switch instr.Name {
	case "call":
		calledName := instr.Children[0]
		if slices.Contains(names, calledName.Name) {
			return true
		}
		nameDefID := symbolTable[int(calledName.ID)].declarationNode
		nameDefNode := parser.GetNodeByID(rootNode, nameDefID)
		if checkDefForRecursion(nameDefNode, append(names, calledName.Name)) {
			return true
		}
	case "assign":
		if checkAssignForRecursion(instr.Children[0], names) {
			return true
		}
	}
	return false
}

func checkAssignForRecursion(node *parser.ASTNode, names []string) bool {
	if node.Name == "call" {
		calledName := node.Children[1]
		if slices.Contains(names, calledName.Name) {
			return true
		}
		nameDefID := symbolTable[int(calledName.ID)].declarationNode
		nameDefNode := parser.GetNodeByID(rootNode, nameDefID)
		if checkDefForRecursion(nameDefNode, append(names, calledName.Name)) {
			return true
		}
	}
	return false
}
