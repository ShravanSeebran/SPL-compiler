package analyser

import (
	"SPL-compiler/parser"
	"fmt"
	"strings"
)

var (
	placeIndex int
	labelIndex int
)

func initialiseGenerator() {
	placeIndex = 0
	labelIndex = 0
}

func GenerateProgram(root *parser.ASTNode) {
	initialiseGenerator()
	checkNode(root)
}

func newPlace() string {
	newIndex := fmt.Sprintf("t%d", placeIndex)
	placeIndex++
	return newIndex
}

func newLabel() string {
	newIndex := fmt.Sprintf("l%d", labelIndex)
	labelIndex++
	return newIndex
}

func generateCode(node *parser.ASTNode) []string {
	if node == nil {
		return []string{}
	}

	switch node.Type {
	case SPL_PROG:
		return generateProgram(node)
	// case PROCDEFS:
	// 	return generateProcDefs(node)
	// case PDEF:
	// 	return generatePDef(node)
	// case FUNCDEFS:
	// 	return generateFuncDefs(node)
	// case FDEF:
	// 	return generateFDef(node)
	// case VAR:
	// 	return generateVar(node)
	// case NAME:
	// 	return generateName(node)
	// case BODY:
	// 	return generateBody(node)
	// case PARAM:
	// 	return generateParam(node)
	// case MAXTHREE:
	// 	return generateMaxThree(node)
	// case MAINPROG:
	// 	return generateMainProg(node)
	// case ATOM:
	// 	return generateAtom(node)
	// case ALGO:
	// 	return generateAlgo(node)
	// case INSTR:
	// 	return generateInstr(node)
	// case ASSIGN:
	// 	return generateAssign(node)
	// case LOOP:
	// 	return generateLoop(node)
	// case BRANCH:
	// 	return generateBranch(node)
	// case OUTPUT:
	// 	return generateOutput(node)
	// case INPUT:
	// 	return generateInput(node)
	// case TERM:
	// 	return generateTerm(node)
	// case UNOP:
	// 	return generateUnOp(node)
	// case BINOP:
	// 	return generateBinOp(node)
	default:
		panic("ungenerateed-node-type: " + node.Type)
	}
}

func generateProgram(node *parser.ASTNode) []string {
	return generateMainProg(node.Children[3])
}

func generateMainProg(node *parser.ASTNode) []string {
	return generateAlgo(node.Children[1])
}

func generateAlgo(node *parser.ASTNode) []string {
	output := make([]string, 0)
	for _, child := range node.Children {
		output = append(output, generateCode(child)...)
	}
	return output
}

func generateInstr(node *parser.ASTNode) []string {
	if node.Name == "halt" {
		return []string{"STOP"}
	} else if node.Name == "print" {
		return []string{fmt.Sprintf("PRINT %s", getOutput(node.Children[0]))}
	} else if node.Name == "call" {
		fname := symbolTable[int(node.Children[0].ID)].uniqueID
		code, places := generateInput(node.Children[1])
		return append(code, fmt.Sprintf("CALL %s(%s)", fname, strings.Join(places, " ")))
	} else {
		return generateCode(node.Children[0])
	}
}

func getOutput(node *parser.ASTNode) string {
	if node.Name == "atom" {
		return getAtom(node.Children[0])
	} else {
		return node.Name
	}
}

func generateInput(node *parser.ASTNode) ([]string, []string) {
	assignments := make([]string, 0)
	places := make([]string, 0)
	for _, child := range node.Children {
		place := newPlace()
		value := getAtom(child)
		assignments = append(assignments, fmt.Sprintf("%s = %s", place, value))
		places = append(places, place)
	}
	return assignments, places
}

func getAtom(node *parser.ASTNode) string {
	if len(node.Children) > 0 {
		return getVar(node.Children[0])
	} else {
		return node.Name
	}
}

func getVar(node *parser.ASTNode) string {
	return symbolTable[int(node.ID)].uniqueID
}

func generateAssign(node *parser.ASTNode) []string {
	if node.Name == "call" {
		place := newPlace()
		vname := symbolTable[int(node.Children[0].ID)].uniqueID
		fname := symbolTable[int(node.Children[1].ID)].uniqueID
		code, places := generateInput(node.Children[2])

		return append(
			code,
			fmt.Sprintf("%s = CALL %s(%s)", place, fname, strings.Join(places, " ")),
			fmt.Sprintf("%s = %s", vname, place),
		)

	} else {
		place := newPlace()
		vname := symbolTable[int(node.Children[0].ID)].uniqueID
		code := generateTerm(node.Children[1], place)
		return append(code, fmt.Sprintf("%s = %s", vname, place))
	}
}

func generateLoop(node *parser.ASTNode) []string {
	// NOTE: Check for implicit GOTOS
	if node.Name == "while" {
		labelCond := newLabel()
		labelStart := newLabel()
		labelExit := newLabel()
		t0 := newPlace()
		cond := generateTerm(node.Children[0], t0)
		algo := generateAlgo(node.Children[1])
		part0 := append([]string{
			fmt.Sprintf("REM %s", labelCond),
		}, cond...)
		part1 := append(
			part0,
			fmt.Sprintf("IF %s THEN %s", t0, labelStart),
			fmt.Sprintf("GOTO %s", labelExit),
			fmt.Sprintf("REM %s", labelStart),
		)
		part2 := append(
			part1,
			algo...,
		)
		return append(
			part2,
			fmt.Sprintf("GOTO %s", labelCond),
			fmt.Sprintf("REM %s", labelExit))

		// NOTE: Stopped here
	} else if node.Name == "do" {
		labelStart := newLabel()
		labelExit := newLabel()
		algo := generateAlgo(node.Children[0])
		t0 := newPlace()
		cond := append(
			generateTerm(node.Children[1], t0),
			fmt.Sprintf("IF %s THEN %s", t0, labelExit),
		)
		part0 := append([]string{
			fmt.Sprintf("REM %s", labelStart),
		}, algo...)
		part1 := append(
			part0,
			cond...,
		)
		return append(
			part1,
			fmt.Sprintf("REM %s", labelExit),
		)

	} else {
		panic("expected 'while' or 'do' Loop node name")
	}
}

func generateBranch(node *parser.ASTNode) []string {
	if node.Name == "if" {
		l0 := newLabel()
		l1 := newLabel()
		t0 := newPlace()
		cond := generateTerm(node.Children[0], t0)
		part1 := append(
			cond,
			fmt.Sprintf("IF %s THEN %s", t0, l0),
			fmt.Sprintf("GOTO %s", l1),
			fmt.Sprintf("REM %s", l0),
		)
		part2 := append(
			part1,
			generateAlgo(node.Children[1])...,
		)
		return append(
			part2,
			fmt.Sprintf("REM %s", l1),
		)
	} else if node.Name == "ifelse" {
		l0 := newLabel()
		l1 := newLabel()
		t0 := newPlace()
		cond := generateTerm(node.Children[0], t0)
		part0 := append(
			cond,
			fmt.Sprintf("IF %s THEN %s", t0, l0),
		)
		part1 := append(
			part0,
			generateAlgo(node.Children[2])...,
		)
		part2 := append(
			part1,
			fmt.Sprintf("GOTO %s", l1),
			fmt.Sprintf("REM %s", l0),
		)
		part3 := append(
			part2,
			generateAlgo(node.Children[1])...,
		)
		return append(
			part3,
			fmt.Sprintf("REM %s", l1),
		)
	} else {
		panic("expected 'if' or 'ifelse' Branch node name")
	}
}

func generateTerm(node *parser.ASTNode, place string) []string {
	// if node.Name == "atom" {
	// } else if node.Name == "unop" {
	// } else if node.Name == "binop" {
	// } else {
	// 	panic("expected 'atom', 'unop', or 'binop' Term node name")
	// }
	return []string{}
}

func generateUnOp(node *parser.ASTNode) string {
	if node.Name == "neg" {
		return "numeric"
	} else if node.Name == "not" {
		return "boolean"
	} else {
		panic("expected 'neg' or 'not' UnOp node name")
	}
}

func generateBinOp(node *parser.ASTNode) string {
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
