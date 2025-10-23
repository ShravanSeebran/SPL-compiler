package checker

import (
	"SPL-compiler/analyser"
	"SPL-compiler/parser"
)

func TypeCheckProgram(root *parser.ASTNode) {
	checkNode(root)
}

func checkNode(node *parser.ASTNode) {
	if node == nil {
		return
	}

	switch node.Type {
	case analyser.SPL_PROG:
		checkProgram(node)
	case analyser.VARIABLES:
		checkVariables(node)
	case analyser.PROCDEFS:
		checkProcDefs(node)
	case analyser.PDEF:
		checkPDef(node)
	case analyser.FUNCDEFS:
		checkFuncDefs(node)
	case analyser.FDEF:
		checkFDef(node)
	case analyser.VAR:
		checkVar(node)
	case analyser.NAME:
		checkName(node)
	case analyser.BODY:
		checkBody(node)
	case analyser.PARAM:
		checkParam(node)
	case analyser.MAXTHREE:
		checkMaxThree(node)
	case analyser.MAINPROG:
		checkMainProg(node)
	case analyser.ATOM:
		checkAtom(node)
	case analyser.ALGO:
		checkAlgo(node)
	case analyser.INSTR:
		checkInstr(node)
	case analyser.ASSIGN:
		checkAssign(node)
	case analyser.LOOP:
		checkLoop(node)
	case analyser.BRANCH:
		checkBranch(node)
	case analyser.OUTPUT:
		checkOutput(node)
	case analyser.INPUT:
		checkInput(node)
	case analyser.TERM:
		checkTerm(node)
	case analyser.UNOP:
		checkUnOp(node)
	case analyser.BINOP:
		checkBinOp(node)
	default:
		panic("unchecked-node-type: " + node.Type)
	}
}

func checkProgram(node *parser.ASTNode) {
	for _, child := range node.Children {
		checkNode(child)
	}
}

func checkVariables(node *parser.ASTNode) {
	if len(node.Children) > 0 {
		n := checkVar(node.Children[0])
		if n != "numeric" {
			panic("expected numeric type for variable")
		}
		checkNode(node.Children[1]) // VARIABLES
	} else {
		return
	}
}

func checkProcDefs(node *parser.ASTNode) {
	if len(node.Children) > 0 {
		checkNode(node.Children[0]) // PDEF
		checkNode(node.Children[1]) // PROCDEFS
	} else {
		return
	}
}

func checkFuncDefs(node *parser.ASTNode) {
	if len(node.Children) > 0 {
		checkNode(node.Children[0]) // FDEF
		checkNode(node.Children[1]) // FUNCDEFS
	} else {
		return
	}
}

func checkMainProg(node *parser.ASTNode) {
	for _, child := range node.Children {
		checkNode(child)
	}
}

func checkPDef(node *parser.ASTNode) {
	checkName(node.Children[0])  // name
	checkParam(node.Children[1]) // param
	checkBody(node.Children[2])  // body
}

func checkFDef(node *parser.ASTNode) {
	checkName(node.Children[0])  // name
	checkParam(node.Children[1]) // param
	checkBody(node.Children[2])  // body

	n := checkAtom(node.Children[3]) // atom
	if n != "numeric" {
		panic("expected numeric type for atom")
	}
}

func checkParam(node *parser.ASTNode) {
	checkNode(node.Children[0]) // maxthree
}

func checkMaxThree(node *parser.ASTNode) {
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			n := checkVar(child)
			if n != "numeric" {
				panic("expected numeric type for variable")
			}
		}
	} else {
		return
	}
}

func checkVar(node *parser.ASTNode) string {
	return "numeric"
}

func checkName(node *parser.ASTNode) {
	// NOTE: Find out about what type-less means
	// TODO: implement checkName
}

func checkBody(node *parser.ASTNode) {
	checkNode(node.Children[0]) // maxthree
	checkNode(node.Children[1]) // algo
}

func checkAlgo(node *parser.ASTNode) {
	for _, child := range node.Children {
		checkNode(child)
	}
}

func checkInstr(node *parser.ASTNode) {
	for _, child := range node.Children {
		checkNode(child)
	}
}

func checkOutput(node *parser.ASTNode) {
	if len(node.Children) > 0 {
		if n := checkVar(node.Children[0]); n != "numeric" {
			panic("expected numeric type for variable")
		}
	} else {
		return
	}
}

func checkInput(node *parser.ASTNode) {
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			if n := checkVar(child); n != "numeric" {
				panic("expected numeric type for variable")
			}
		}
	} else {
		return
	}
}

func checkAssign(node *parser.ASTNode) {
	if node.Name == "call" {
		if n := checkVar(node.Children[0]); n != "numeric" {
			panic("expected numeric type for variable")
		}
		checkNode(node.Children[1]) // NAME
		checkNode(node.Children[2]) // INPUT
	} else {
		if n := checkVar(node.Children[0]); n != "numeric" {
			panic("expected numeric type for variable")
		}
		if n := checkVar(node.Children[1]); n != "numeric" {
			panic("expected numeric type for variable")
		}
	}
}

func checkLoop(node *parser.ASTNode) {
	// TODO: Need to update parser with this name
	if node.Name == "while" {
		b := checkTerm(node.Children[0])
		if b != "boolean" {
			panic("expected boolean type for WHILE condition")
		}
		checkNode(node.Children[1]) // ALGO
		// TODO: Need to update parser with this name
	} else if node.Name == "do" {
		checkNode(node.Children[0]) // ALGO
		b := checkTerm(node.Children[1])
		if b != "boolean" {
			panic("expected boolean type for DO condition")
		}
	} else {
		panic("expected 'while' or 'do' Loop node name")
	}
}

func checkBranch(node *parser.ASTNode) {
	// TODO: Need to update parser with this name
	if node.Name == "if" {
		b := checkTerm(node.Children[0])
		if b != "boolean" {
			panic("expected boolean type for IF condition")
		}
		checkNode(node.Children[1]) // ALGO
		// TODO: Need to update parser with this name
	} else if node.Name == "ifelse" {
		b := checkTerm(node.Children[0])
		if b != "boolean" {
			panic("expected boolean type for IF condition")
		}
		checkNode(node.Children[1]) // ALGO
		checkNode(node.Children[2]) // ALGO
	} else {
		panic("expected 'if' or 'ifelse' Branch node name")
	}
}

func checkTerm(node *parser.ASTNode) string {
	// TODO: Need to update parser with this name
	if node.Name == "atom" {
		n := checkAtom(node.Children[0])
		if n != "numeric" {
			panic("expected numeric type for atom")
		}
		return n
		// TODO: Need to update parser with this name
	} else if node.Name == "unop" {
		t := checkUnOp(node.Children[0])
		s := checkTerm(node.Children[1])
		if t == "numeric" && s == "numeric" {
			return "numeric"
		} else if t == "boolean" && s == "boolean" {
			return "boolean"
		} else {
			panic("expected numeric or boolean type for unop")
		}
		// TODO: Need to update parser with this name
	} else if node.Name == "binop" {
		t := checkTerm(node.Children[0])
		s := checkBinOp(node.Children[1])
		r := checkTerm(node.Children[2])
		if t == "numeric" && s == "numeric" && r == "numeric" {
			return "numeric"
		} else if t == "boolean" && s == "boolean" && r == "boolean" {
			return "boolean"
		} else if t == "numeric" && s == "comparison" && r == "numeric" {
			return "boolean"
		} else {
			panic("expected numeric or boolean type for binop")
		}
	} else {
		panic("expected 'atom', 'unop', or 'binop' Term node name")
	}
}

// NOTE: Check if node.Name is correct in parser
// like "neg" or "-" or "NEG"
func checkUnOp(node *parser.ASTNode) string {
	if node.Name == "neg" {
		return "numeric"
	} else if node.Name == "not" {
		return "boolean"
	} else {
		panic("expected 'neg' or 'not' UnOp node name")
	}
}

// NOTE: Check if node.Name is correct in parser
// like "eq" or "==" or "EQ"
func checkBinOp(node *parser.ASTNode) string {
	if node.Name == "eq" {
		return "comparison"
	} else if node.Name == "gt" {
		return "comparison"
	} else if node.Name == "or" {
		return "boolean"
	} else if node.Name == "and" {
		return "boolean"
	} else if node.Name == "plus" {
		return "numeric"
	} else if node.Name == "minus" {
		return "numeric"
	} else if node.Name == "mult" {
		return "numeric"
	} else if node.Name == "div" {
		return "numeric"
	} else {
		panic("expected 'eq', 'gt', 'or', 'and', 'plus', 'minus', 'mult', or 'div' BinOp node name")
	}
}

func checkAtom(node *parser.ASTNode) string {
	if len(node.Children) > 0 {
		n := checkVar(node.Children[0])
		if n != "numeric" {
			panic("expected numeric type for atom")
		}
		return n
	}
	return "numeric"
}
