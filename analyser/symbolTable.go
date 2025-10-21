package analyser

type SemanticInfo struct {
	nodeID          int
	symbolName      string
	dataType        string
	scopeLevel      int
	declarationNode int
}

type SymbolTable map[int]SemanticInfo
