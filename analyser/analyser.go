package analyser

import "SPL-compiler/parser"

var (
	symbolTable     SymbolTable
	auxStack        *AuxillaryStack
	currentScope    int
	GLOBAL_SCOPE    int
	PROCEDURE_SCOPE int
	FUNCTION_SCOPE  int
)

func initialiseAnalyser() {
	symbolTable = make(SymbolTable)
	auxStack = Empty()
	currentScope = 0
}

func analyseProgram(root *parser.ASTNode) {
	initialiseAnalyser()
	visitNode(root)
}

func visitNode(node *parser.ASTNode) {
	if node == nil {
		return
	}

	switch node.Type {
	case SPL_PROG:
		handleProgram(node)
	case VARIABLES:
		handleVariables(node)
	case PROCDEFS:
		handleProcDefs(node)
	case PDEF:
		handlePDef(node)
	case FUNCDEFS:
		handleFuncDefs(node)
	case FDEF:
		handleFDef(node)
	case VAR:
		handleVar(node)
	case NAME:
		handleName(node)
	case BODY:
		handleBody(node)
	case PARAM:
		handleParam(node)
	case MAXTHREE:
		handleMaxThree(node)
	case MAINPROG:
		handleMainProg(node)
	case ATOM:
		handleAtom(node)
	case ALGO:
		handleAlgo(node)
	case INSTR:
		handleInstr(node)
	case ASSIGN:
		handleAssign(node)
	case LOOP:
		handleLoop(node)
	case BRANCH:
		handleBranch(node)
	case OUTPUT:
		handleOutput(node)
	case INPUT:
		handleInput(node)
	case TERM:
		handleTerm(node)
	case UNOP:
		handleUnOp(node)
	case BINOP:
		handleBinOp(node)
	default:
		panic("unhandled-node-type: " + node.Type)
	}
}

func handleProgram(node *parser.ASTNode) {
	// GLOBAL SCOPE
	currentScope = auxStack.enter(int(node.ID)) // everywhere scope

	GLOBAL_SCOPE = int(node.Children[0].ID)
	PROCEDURE_SCOPE = int(node.Children[1].ID)
	FUNCTION_SCOPE = int(node.Children[2].ID)

	for _, child := range node.Children {
		currentScope = auxStack.enter(int(child.ID))
		visitNode(child)
		currentScope = auxStack.exit()
	}
	// NOTE: Alternative fix for handleVar checking GLOBAL_SCOPE
	// uncomment the following line and comment the above one
	// currentScope = auxStack.exit()

	currentScope = auxStack.exit() // should be -1
}

func handleVariables(node *parser.ASTNode) {
	// NOTE: Fix bug where indexing empty children array
	if len(node.Children) > 0 {
		declareVar(node.Children[0])
		visitNode(node.Children[1])
	} else {
		return
	}
}

func handleProcDefs(node *parser.ASTNode) {
	for _, child := range node.Children {
		currentScope = auxStack.enter(int(child.ID))
		visitNode(child)
		currentScope = auxStack.exit()
	}
}

func handleFuncDefs(node *parser.ASTNode) {
	for _, child := range node.Children {
		currentScope = auxStack.enter(int(child.ID))
		visitNode(child)
		currentScope = auxStack.exit()
	}
}

func handleMainProg(node *parser.ASTNode) {
	for _, child := range node.Children {
		visitNode(child)
	}
}

func handlePDef(node *parser.ASTNode) {
	declareName(node.Children[0]) // name
	visitNode(node.Children[1])   // param
	visitNode(node.Children[2])   // body
}

func handleFDef(node *parser.ASTNode) {
	declareName(node.Children[0]) // name
	visitNode(node.Children[1])   // param
	visitNode(node.Children[2])   // body
	visitNode(node.Children[3])   // atom
}

func handleParam(node *parser.ASTNode) {
	visitNode(node.Children[0]) // maxthree
}

func handleMaxThree(node *parser.ASTNode) {
	for _, child := range node.Children {
		declareVar(child)
	}
}

func declareVar(node *parser.ASTNode) {
	varname := node.Name
	nodeID, ok := auxStack.lookup(varname)
	if ok {
		lookupScope := symbolTable[nodeID].scopeLevel
		if lookupScope == currentScope || lookupScope == PROCEDURE_SCOPE ||
			lookupScope == FUNCTION_SCOPE {
			// Redeclaration error
			panic("name-rule-violation: conflict of " + varname)
		}
	}

	// NOTE: Fixed bug of never binding the variables to the current AuxillaryStack!
	auxStack.bind(varname, int(node.ID))
	symbolTable[int(node.ID)] = SemanticInfo{
		nodeID:          int(node.ID),
		symbolName:      varname,
		scopeLevel:      currentScope,
		declarationNode: int(node.ID),
	}
}

func declareName(node *parser.ASTNode) {
	name := node.Name
	nodeID, ok := auxStack.lookup(name)
	if ok {
		lookupScope := symbolTable[nodeID].scopeLevel
		if lookupScope == PROCEDURE_SCOPE || lookupScope == FUNCTION_SCOPE ||
			lookupScope == GLOBAL_SCOPE {
			// Redeclaration error
			panic("name-rule-violation: name redeclaration of " + name)
		}
	}

	// NOTE: Fixed bug of never binding the variables to the current AuxillaryStack!
	auxStack.bind(name, int(node.ID))
	symbolTable[int(node.ID)] = SemanticInfo{
		nodeID:          int(node.ID),
		symbolName:      name,
		scopeLevel:      currentScope,
		declarationNode: int(node.ID),
	}
}

func handleVar(node *parser.ASTNode) {
	varname := node.Name
	nodeID, ok := auxStack.lookup(varname)
	if !ok {
		// NOTE: Fixed bug of not checking GLOBAL_SCOPE
		for _, value := range symbolTable {
			if value.symbolName == varname {
				if value.scopeLevel == GLOBAL_SCOPE {
					symbolTable[int(node.ID)] = SemanticInfo{
						nodeID:          int(node.ID),
						symbolName:      varname,
						scopeLevel:      GLOBAL_SCOPE,
						declarationNode: GLOBAL_SCOPE,
					}
					return
				}
			}
		}

		// Undeclared variable error
		panic("undeclared-variable: " + varname)
	}

	lookupScope := symbolTable[nodeID].scopeLevel
	if lookupScope == PROCEDURE_SCOPE || lookupScope == FUNCTION_SCOPE {
		// Variable-function/procedure conflict error
		panic("undeclared-variable: " + varname)
	}

	symbolTable[int(node.ID)] = SemanticInfo{
		nodeID:          int(node.ID),
		symbolName:      varname,
		scopeLevel:      symbolTable[nodeID].scopeLevel,
		declarationNode: symbolTable[nodeID].declarationNode,
	}
}

func handleName(node *parser.ASTNode) {
	name := node.Name
	nodeID, ok := auxStack.lookup(name)
	if !ok {
		// Undeclared name error
		panic("undeclared-name: " + name)
	}

	lookupScope := symbolTable[nodeID].scopeLevel
	if !(lookupScope == PROCEDURE_SCOPE || lookupScope == FUNCTION_SCOPE) {
		// Name not a function/procedure error
		panic("undeclared-name: " + name)
	}

	symbolTable[int(node.ID)] = SemanticInfo{
		nodeID:          int(node.ID),
		symbolName:      name,
		scopeLevel:      symbolTable[nodeID].scopeLevel,
		declarationNode: symbolTable[nodeID].declarationNode,
	}
}

func handleBody(node *parser.ASTNode) {
	visitNode(node.Children[0]) // maxthree
	visitNode(node.Children[1]) // algo
}

func handleAlgo(node *parser.ASTNode) {
	for _, child := range node.Children {
		visitNode(child)
	}
}

func handleInstr(node *parser.ASTNode) {
	for _, child := range node.Children {
		visitNode(child)
	}
}

func handleOutput(node *parser.ASTNode) {
	visitNode(node.Children[0])
}

func handleInput(node *parser.ASTNode) {
	for _, child := range node.Children {
		visitNode(child)
	}
}

func handleAssign(node *parser.ASTNode) {
	for _, child := range node.Children {
		visitNode(child)
	}
}

func handleLoop(node *parser.ASTNode) {
	for _, child := range node.Children {
		visitNode(child)
	}
}

func handleBranch(node *parser.ASTNode) {
	for _, child := range node.Children {
		visitNode(child)
	}
}

func handleTerm(node *parser.ASTNode) {
	for _, child := range node.Children {
		visitNode(child)
	}
}

func handleUnOp(node *parser.ASTNode) {
	return
}

func handleBinOp(node *parser.ASTNode) {
	return
}

func handleAtom(node *parser.ASTNode) {
	for _, child := range node.Children {
		visitNode(child)
	}
}
