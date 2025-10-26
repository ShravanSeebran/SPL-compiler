package analyser

import (
	"SPL-compiler/parser"
	"fmt"
	"slices"
)

var (
	placeIndex  int
	labelIndex  int
	letterIndex int
)

func initialiseGenerator() {
	placeIndex = 0
	labelIndex = 0
	letterIndex = 0
}

func getUniquePlace() string {
	if placeIndex == 25 {
		placeIndex = 0
		letterIndex++
	}
	if letterIndex == 25 {
		panic("Run out of unique place names")
	}
	currentChar := charSet[letterIndex]
	recentChar := charSet[placeIndex]
	candidate := fmt.Sprintf("%c%c", currentChar, recentChar)
	if keywordSet(candidate) {
		placeIndex++
		return getUniquePlace()
	}
	placeIndex++
	return candidate
}

func keywordSet(candidate string) bool {
	// TODO: Add more keywords
	keywords := []string{"at", "or", "if", "to", "on", "go", "as", "is", "do", "in"}
	return slices.Contains(keywords, candidate)
}

func ValidateCodeGeneration(root *parser.ASTNode) (intrs []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	instrs := GenerateProgram(root)
	return instrs, nil
}

func GenerateProgram(root *parser.ASTNode) []string {
	initialiseGenerator()
	return generateCode(root)
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
		code, argPlaces := generateInput(node.Children[1])
		procNodeID := symbolTable[int(node.Children[0].ID)].declarationNode
		procNode := parser.GetDefNodeByNameID(rootNode, procNodeID)
		inlineCode := inlineProc(procNode)
		output := make([]string, 0)
		output = append(output, code...)
		for i, param := range procNode.Children[1].Children[0].Children {
			output = append(output, fmt.Sprintf("%s = %s", getVar(param), argPlaces[i]))
		}
		output = append(output, inlineCode...)
		return output
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
		place := getUniquePlace()
		value := getAtom(child)
		assignments = append(assignments, fmt.Sprintf("%s = %s", place, value))
		places = append(places, place)
	}
	return assignments, places
}

func inlineProc(node *parser.ASTNode) []string {
	if node.Type != "PDEF" {
		panic(fmt.Sprintf("expected 'pdef' Proc node name but got %s", node.Type))
	}
	return generateAlgo(node.Children[2].Children[1])
}

func inlineFunc(node *parser.ASTNode, place string) []string {
	if node.Type != "FDEF" {
		panic(fmt.Sprintf("expected 'fdef' Func node name but got %s", node.Type))
	}
	algo := generateAlgo(node.Children[2].Children[1])
	return append(algo, fmt.Sprintf("%s = %s", place, getAtom(node.Children[3])))
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
		place := getUniquePlace()
		vname := symbolTable[int(node.Children[0].ID)].uniqueID
		code, argPlaces := generateInput(node.Children[2])
		funcNodeID := symbolTable[int(node.Children[1].ID)].declarationNode
		funcNode := parser.GetDefNodeByNameID(rootNode, funcNodeID)
		inlineCode := inlineFunc(funcNode, place)
		output := make([]string, 0)
		output = append(output, code...)
		for i, param := range funcNode.Children[1].Children[0].Children {
			output = append(output, fmt.Sprintf("%s = %s", getVar(param), argPlaces[i]))
		}
		output = append(output, inlineCode...)

		return append(
			output,
			fmt.Sprintf("%s = %s", vname, place),
		)

	} else {
		place := getUniquePlace()
		vname := symbolTable[int(node.Children[0].ID)].uniqueID
		code := generateTerm(node.Children[1], place)
		return append(code, fmt.Sprintf("%s = %s", vname, place))
	}
}

func generateLoop(node *parser.ASTNode) []string {
	switch node.Name {
	case "while":
		labelCond := newLabel()
		labelStart := newLabel()
		labelExit := newLabel()
		cond := generateCond(node.Children[0], labelStart, labelExit)
		algo := generateAlgo(node.Children[1])
		part0 := append([]string{
			fmt.Sprintf("REM %s", labelCond),
		}, cond...)
		part1 := append(
			part0,
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

	case "do":
		labelStart := newLabel()
		labelExit := newLabel()
		algo := generateAlgo(node.Children[0])
		cond := generateCond(node.Children[1], labelExit, labelStart)
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

	default:
		panic("expected 'while' or 'do' Loop node name")
	}
}

func generateBranch(node *parser.ASTNode) []string {
	switch node.Name {
	case "if":
		labelStart := newLabel()
		labelExit := newLabel()
		cond := generateCond(node.Children[0], labelStart, labelExit)
		part1 := append(
			cond,
			fmt.Sprintf("REM %s", labelStart),
		)
		part2 := append(
			part1,
			generateAlgo(node.Children[1])...,
		)
		return append(
			part2,
			fmt.Sprintf("REM %s", labelExit),
		)
	case "ifelse":
		labelStart := newLabel()
		labelExit := newLabel()
		ifAlgo := generateAlgo(node.Children[1])
		elseAlgo := generateAlgo(node.Children[2])
		cond := generateCondElse(node.Children[0], labelStart, labelExit, elseAlgo)
		part1 := append(
			cond,
			fmt.Sprintf("REM %s", labelStart),
		)
		part2 := append(
			part1,
			ifAlgo...,
		)
		part3 := append(
			part2,
			fmt.Sprintf("REM %s", labelExit),
		)
		return part3
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
		t1 := getUniquePlace()
		t2 := getUniquePlace()
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

func generateCondElse(node *parser.ASTNode, labelT, labelF string, elseInstrs []string) []string {
	switch node.Name {
	case "unop":
		if node.Children[0].Name != "not" {
			panic("expected 'not' UnOp in Cond node")
		}
		return generateCondElse(node.Children[1], labelF, labelT, elseInstrs)
	case "binop":
		op := node.Children[1].Name
		switch op {
		case "and":
			arg2 := newLabel()
			codeL := generateCondElse(node.Children[0], arg2, labelF, elseInstrs)
			codeR := generateCond(node.Children[2], labelT, labelF)
			part0 := append(codeL, fmt.Sprintf("REM %s", arg2))
			return append(part0, codeR...)
		case "or":
			arg2 := newLabel()
			codeL := generateCond(node.Children[0], labelT, arg2)
			codeR := generateCondElse(node.Children[2], labelT, labelF, elseInstrs)
			part0 := append(codeL, fmt.Sprintf("REM %s", arg2))
			return append(part0, codeR...)
		}
		t1 := getUniquePlace()
		t2 := getUniquePlace()
		codeL := generateTerm(node.Children[0], t1)
		binop := getBinOp(node.Children[1])
		codeR := generateTerm(node.Children[2], t2)
		part0 := append(codeL, codeR...)
		part1 := append(
			part0,
			fmt.Sprintf("IF %s %s %s THEN %s", t1, binop, t2, labelT),
		)
		part2 := append(
			part1,
			elseInstrs...,
		)
		part3 := append(
			part2,
			fmt.Sprintf("GOTO %s", labelF),
		)
		return part3
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
		t0 := getUniquePlace()
		unop := getUnOp(node.Children[0])
		code := generateTerm(node.Children[1], t0)
		return append(code, fmt.Sprintf("%s = %s%s", place, unop, t0))
	case "binop":
		if node.Children[1].Name == "and" || node.Children[1].Name == "or" {
			panic("expected non-boolean BinOp in Term node")
		}
		t0 := getUniquePlace()
		t1 := getUniquePlace()
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
