package analyser

import (
	"fmt"
	"sort"
	"strings"
)

type SemanticInfo struct {
	nodeID          int
	symbolName      string
	dataType        string
	scopeLevel      int
	declarationNode int
}

type SymbolTable map[int]SemanticInfo

func PrettyPrintSymbolTable(st SymbolTable) {
	if len(st) == 0 {
		fmt.Println("Symbol table is empty")
		return
	}

	keys := make([]int, 0, len(st))
	for k := range st {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	widths := struct {
		nodeID, symbolName, dataType, scopeLevel, declarationNode int
	}{
		nodeID:          6,
		symbolName:      11,
		dataType:        9,
		scopeLevel:      11,
		declarationNode: 16,
	}

	for _, info := range st {
		if len(info.symbolName) > widths.symbolName {
			widths.symbolName = len(info.symbolName)
		}
		if len(info.dataType) > widths.dataType {
			widths.dataType = len(info.dataType)
		}
		if len(fmt.Sprint(info.nodeID)) > widths.nodeID {
			widths.nodeID = len(fmt.Sprint(info.nodeID))
		}
		if len(fmt.Sprint(info.scopeLevel)) > widths.scopeLevel {
			widths.scopeLevel = len(fmt.Sprint(info.scopeLevel))
		}
		if len(fmt.Sprint(info.declarationNode)) > widths.declarationNode {
			widths.declarationNode = len(fmt.Sprint(info.declarationNode))
		}
	}

	separator := fmt.Sprintf("+-%s-+-%s-+-%s-+-%s-+-%s-+",
		strings.Repeat("-", widths.nodeID),
		strings.Repeat("-", widths.symbolName),
		strings.Repeat("-", widths.dataType),
		strings.Repeat("-", widths.scopeLevel),
		strings.Repeat("-", widths.declarationNode))

	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*s | %-*s | %-*s | %-*s |\n",
		widths.nodeID, "NodeID",
		widths.symbolName, "Symbol Name",
		widths.dataType, "Data Type",
		widths.scopeLevel, "Scope Level",
		widths.declarationNode, "Declaration Node")
	fmt.Println(separator)

	for _, key := range keys {
		info := st[key]
		fmt.Printf("| %-*d | %-*s | %-*s | %-*d | %-*d |\n",
			widths.nodeID, info.nodeID,
			widths.symbolName, info.symbolName,
			widths.dataType, info.dataType,
			widths.scopeLevel, info.scopeLevel,
			widths.declarationNode, info.declarationNode)
	}
	fmt.Println(separator)
}
