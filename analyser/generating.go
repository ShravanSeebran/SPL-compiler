package analyser

import (
	"fmt"
	"strings"

	"SPL-compiler/parser"
)

var (
	placeIndex int
	labelIndex int
)

func initialiseGenerator() {
	placeIndex = 0
	labelIndex = 0
}

func GenerateProgram(root *parser.ASTNode) []string {
	initialiseGenerator()
	return generateCode(root)
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
	case MAINPROG:
		return generateMainProg(node)
	case ALGO:
		return generateAlgo(node)
	case INSTR:
		return generateInstr(node)
	case ASSIGN:
		return generateAssign(node)
	case LOOP:
		return generateLoop(node)
	case BRANCH:
		return generateBranch(node)
	// case TERM:
	// 	return generateTerm(node)
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
	switch node.Name {
	case "halt":
		return []string{"STOP"}
	case "print":
		return []string{fmt.Sprintf("PRINT %s", getOutput(node.Children[0]))}
	case "call":
		fname := symbolTable[int(node.Children[0].ID)].uniqueID
		code, places := generateInput(node.Children[1])
		return append(code, fmt.Sprintf("CALL %s(%s)", fname, strings.Join(places, " ")))
		// procNodeID := symbolTable[int(node.Children[0].ID)].declarationNode
		// procNode := parser.GetNodeByID(procNodeID)
		// inlineCode, places, := inlineProc(procNode)
		// outptut := make([]string, 0)
		// for i, param := range node.Children[1].Children {
		// 	outptut = append(outptut, fmt.Sprintf("%s = %s", places[i], getAtom(param)))
		// }
		// outptut = append(outptut, inlineCode...)
		// return output
	default:
		return generateCode(node.Children[0])
	}
}

func getOutput(node *parser.ASTNode) string {
	if node.Name == "atom" {
		return getAtom(node.Children[0])
	} else {
		return fmt.Sprintf(`"%s"`, node.Name)
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
			fmt.Sprintf("GOTO %s", labelStart),
			fmt.Sprintf("REM %s", labelExit),
		)

	} else {
		panic("expected 'while' or 'do' Loop node name")
	}
}

func generateBranch(node *parser.ASTNode) []string {
	switch node.Name {
	case "if":
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
	case "ifelse":
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
	default:
		panic("expected 'if' or 'ifelse' Branch node name")
	}
}

func generateCond(node *parser.ASTNode, labelT, labelF string) []string {
	switch node.Name {
	case "unop":
		if node.Children[0].Name != "not" {
			panic("expected 'not' UnOp in Cond node")
		}
		return generateCond(node.Children[1], labelF, labelT)
	case "binop":
		op := node.Children[1].Name
		switch op {
		case "and":
			arg2 := newLabel()
			codeL := generateCond(node.Children[0], arg2, labelF)
			codeR := generateCond(node.Children[2], labelT, labelF)
			part0 := append(codeL, fmt.Sprintf("REM %s", arg2))
			return append(part0, codeR...)
		case "or":
			arg2 := newLabel()
			codeL := generateCond(node.Children[0], labelT, arg2)
			codeR := generateCond(node.Children[2], labelT, labelF)
			part0 := append(codeL, fmt.Sprintf("REM %s", arg2))
			return append(part0, codeR...)
		}
		t1 := newPlace()
		t2 := newPlace()
		codeL := generateTerm(node.Children[0], t1)
		binop := getBinOp(node.Children[1])
		codeR := generateTerm(node.Children[2], t2)
		part0 := append(codeL, codeR...)
		part1 := append(
			part0,
			fmt.Sprintf("IF %s %s %s THEN %s", t1, binop, t2, labelT),
			fmt.Sprintf("GOTO %s", labelF),
		)
		return part1
	default:
		panic("expected 'unop' or 'binop' Cond node name")
	}
}

func generateTerm(node *parser.ASTNode, place string) []string {
	switch node.Name {
	case "atom":
		atom := getAtom(node.Children[0])
		return []string{fmt.Sprintf("%s = %s", place, atom)}
	case "unop":
		if node.Children[0].Name != "neg" {
			panic("expected 'neg' UnOp in Term node")
		}
		t0 := newPlace()
		unop := getUnOp(node.Children[0])
		code := generateTerm(node.Children[1], t0)
		return append(code, fmt.Sprintf("%s = %s%s", place, unop, t0))
	case "binop":
		if node.Children[1].Name == "and" || node.Children[1].Name == "or" {
			panic("expected non-boolean BinOp in Term node")
		}
		t0 := newPlace()
		t1 := newPlace()
		codeL := generateTerm(node.Children[0], t0)
		binop := getBinOp(node.Children[1])
		codeR := generateTerm(node.Children[2], t1)
		part0 := append(codeL, codeR...)
		return append(part0, fmt.Sprintf("%s = %s %s %s", place, t0, binop, t1))
	default:
		panic("expected 'atom', 'unop', or 'binop' Term node name")
	}
}

func getUnOp(node *parser.ASTNode) string {
	if node.Name == "neg" {
		return "-"
	}
	panic("expected 'neg' UnOp node name")
}

func getBinOp(node *parser.ASTNode) string {
	switch node.Name {
	case "eq":
		return "="
	case ">":
		return ">"
	case "plus":
		return "+"
	case "minus":
		return "-"
	case "mult":
		return "*"
	case "div":
		return "/"
	default:
		panic("expected 'eq', '>', 'plus', 'minus', 'mult', or 'div' BinOp node name")
	}
}
