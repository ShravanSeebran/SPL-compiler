package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // user-defined-name: [a...z]{a...z}*{0...9}*
	INT    = "INT"    // number: 0 | [1...9][0...9]*
	STRING = "STRING" // string: "..." max length 15

	// assignment Operators
	ASSIGN   = "="
	

	// SPL Binary Operators
	EQ   = "eq"    // equality
	GT   = ">"     // greater than
	OR   = "or"    // logical or
	AND  = "and"   // logical and
	PLUS = "plus"     // addition
	MINUS = "minus"    // subtraction
	MULT = "mult"  // multiplication
	DIV  = "div"   // division

	// SPL Unary Operators
	NEG = "neg" // negation
	NOT = "not" // logical not

	// Delimiters
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// SPL Keywords - Program Structure
	GLOB   = "glob"   // global variables section
	PROC   = "proc"   // procedures section
	FUNC   = "func"   // functions section
	MAIN   = "main"   // main program section
	VAR    = "var"    // variable declaration
	LOCAL  = "local"  // local variables
	RETURN = "return" // return statement

	// SPL Keywords - Control Flow
	HALT   = "halt"  // halt instruction
	PRINT  = "print" // print statement
	WHILE  = "while" // while loop
	DO     = "do"    // do-until loop
	UNTIL  = "until" // until condition
	IF     = "if"    // if statement
	ELSE   = "else"  // else clause
)

// Keywords map for identifying SPL keywords
var keywords = map[string]TokenType{
	"glob":   GLOB,
	"proc":   PROC,
	"func":   FUNC,
	"main":   MAIN,
	"var":    VAR,
	"local":  LOCAL,
	"return": RETURN,
	"halt":   HALT,
	"print":  PRINT,
	"while":  WHILE,
	"do":     DO,
	"until":  UNTIL,
	"if":     IF,
	"else":   ELSE,
	"eq":     EQ,
	"or":     OR,
	"and":    AND,
	"plus":   PLUS,
	"minus":  MINUS,
	"mult":   MULT,
	"div":    DIV,
	"neg":    NEG,
	"not":    NOT,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}