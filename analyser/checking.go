package analyser

import (
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
	case SPL_PROG:
		checkProgram(node)
	case VARIABLES:
		checkVariables(node)
	case PROCDEFS:
		checkProcDefs(node)
	case PDEF:
		checkPDef(node)
	case FUNCDEFS:
		checkFuncDefs(node)
	case FDEF:
		checkFDef(node)
	case VAR:
		checkVar(node)
	case NAME:
		checkName(node)
	case BODY:
		checkBody(node)
	case PARAM:
		checkParam(node)
	case MAXTHREE:
		checkMaxThree(node)
	case MAINPROG:
		checkMainProg(node)
	case ATOM:
		checkAtom(node)
	case ALGO:
		checkAlgo(node)
	case INSTR:
		checkInstr(node)
	case ASSIGN:
		checkAssign(node)
	case LOOP:
		checkLoop(node)
	case BRANCH:
		checkBranch(node)
	case OUTPUT:
		checkOutput(node)
	case INPUT:
		checkInput(node)
	case TERM:
		checkTerm(node)
	case UNOP:
		checkUnOp(node)
	case BINOP:
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
	for _, child := range node.Children {
		n := checkVar(child)
		if n != "numeric" {
			panic("expected numeric type for variable")
		}
	}
}

func checkVar(_ *parser.ASTNode) string {
	return "numeric"
}

func checkName(node *parser.ASTNode) {
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
	for _, child := range node.Children {
		if n := checkVar(child); n != "numeric" {
			panic("expected numeric type for variable")
		}
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
		if n := checkTerm(node.Children[1]); n != "numeric" {
			panic("expected numeric type for variable")
		}
	}
}

func checkLoop(node *parser.ASTNode) {
	if node.Name == "while" {
		b := checkTerm(node.Children[0])
		if b != "boolean" {
			panic("expected boolean type for WHILE condition")
		}
		checkNode(node.Children[1]) // ALGO
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
	if node.Name == "if" {
		b := checkTerm(node.Children[0])
		if b != "boolean" {
			panic("expected boolean type for IF condition")
		}
		checkNode(node.Children[1]) // ALGO
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
	if node.Name == "atom" {
		n := checkAtom(node.Children[0])
		if n != "numeric" {
			panic("expected numeric type for atom")
		}
		return n
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

func checkUnOp(node *parser.ASTNode) string {
	if node.Name == "neg" {
		return "numeric"
	} else if node.Name == "not" {
		return "boolean"
	} else {
		panic("expected 'neg' or 'not' UnOp node name")
	}
}

func checkBinOp(node *parser.ASTNode) string {
	if node.Name == "eq" {
		return "comparison"
	} else if node.Name == ">" {
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
		panic("expected 'eq', '>', 'or', 'and', 'plus', 'minus', 'mult', or 'div' BinOp node name")
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
