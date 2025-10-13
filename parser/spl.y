%{
package parser

import (
	"fmt"
	"sync/atomic"
    "strings"
)

var ResultAST *ASTNode
var nodeCounter int64

func nextID() int64 {
	return atomic.AddInt64(&nodeCounter, 1)
}

// Base AST ASTNode
type ASTNode struct {
	ID       int64
	Type     string
	Name     string
	Children []*ASTNode
}

func NewNode(nodeType, name string, children ...*ASTNode) *ASTNode {
	return &ASTNode{
		ID:       nextID(),
		Type:     nodeType,
		Name:     name,
		Children: children,
	}
}

func PrintAST(node *ASTNode, indent int) {
	if node == nil {
		return
	}
	fmt.Printf("%s[%d] %s: %s\n", 
		strings.Repeat("  ", indent), node.ID, node.Type, node.Name)
	for _, c := range node.Children {
		PrintAST(c, indent+1)
	}
}
%}

%union {
    Str   string
    node  *ASTNode
}

%token GLOB
%token LOCAL
%token VAR_KW
%token PROC_KW
%token FUNC_KW
%token MAIN_KW
%token RETURN
%token HALT
%token PRINT
%token WHILE
%token DO
%token UNTIL
%token IF
%token ELSE

%token <Str> NEG
%token <Str> NOT
%token <Str> EQ GREATER OR AND PLUS MINUS MULT DIV
%token <Str> IDENT NUMBER STRING
%type <node> spl_prog variables var procdefs pdef funcdefs fdef body param maxthree mainprog atom algo instr_list instr assign loop branch output input term 
%type <Str> unop binop
%left OR AND
%left PLUS MINUS
%left MULT DIV
%right NEG NOT

%start spl_prog

%%

spl_prog
    : GLOB '{' variables '}' PROC_KW '{' procdefs '}' FUNC_KW '{' funcdefs '}' MAIN_KW '{' mainprog '}'
      {
        $$ = NewNode("Program", "SPL Program", $3, $7, $11, $15)
        ResultAST = $$
      }
    ;

variables
    : /* empty */        { $$ = NewNode("VARIABLES", "empty") }
    | var variables      { $$ = NewNode("VARIABLES", "", $1, $2) }
    ;

var
    : IDENT              { $$ = NewNode("VAR", $1) }
    ;

procdefs
    : /* empty */        { $$ = NewNode("PROCDEFS", "empty") }
    | pdef procdefs      { $$ = NewNode("PROCDEFS", "", $1, $2) }
    ;

pdef
    : IDENT '(' param ')' '{' body '}'
      { $$ = NewNode("ProcDef", $1, $3, $6) }
    ;

funcdefs
    : /* empty */        { $$ = NewNode("FUNCDEFS", "empty") }
    | fdef funcdefs      { $$ = NewNode("FUNCDEFS", "", $1, $2) }
    ;

fdef
    : IDENT '(' param ')' '{' body ';' RETURN atom '}'
        { $$ = NewNode("FDEF", $1, $3, $6, $9) }
    ;

body
    : LOCAL '{' maxthree '}' algo
        { $$ = NewNode("BODY", "", $3, $5) }
    ;

param
    : maxthree           { $$ = NewNode("PARAM", "", $1) }
    ;

maxthree
    : /* empty */        { $$ = NewNode("MAXTHREE", "empty") }
    | var                 { $$ = NewNode("MAXTHREE", "", $1) }
    | var var             { $$ = NewNode("MAXTHREE", "", $1, $2) }
    | var var var         { $$ = NewNode("MAXTHREE", "", $1, $2, $3) }
    ;

mainprog
    : VAR_KW '{' variables '}' algo { $$ = NewNode("MAINPROG", "", $3, $5) }
    ;

atom
    : var { $$ = $1 }
    | NUMBER { $$ = NewNode("Number", $1) }
    ;

algo : instr_list { $$ = NewNode("ALGO", "", $1) }

instr_list 
    : instr          { $$ = NewNode("INSTRLIST", "", $1) }
    | instr_list ';' instr { $$ = NewNode("INSTRLIST", "", $1, $3) }
    ;

instr
    : HALT              { $$ = NewNode("HALT", "") }
    | PRINT output      { $$ = NewNode("Print", "", $2) }
    | IDENT '(' input ')' { $$ = NewNode("ProcCall", $1, $3) }
    | assign            { $$ = NewNode("ASSIGN", "", $1) }
    | loop              { $$ = NewNode("LOOP", "", $1) }
    | branch            { $$ = NewNode("BRANCH", "", $1) }
    ;

assign
    : var '=' IDENT '(' input ')' { 
        $$ = NewNode("FuncAssign", $1.Name, NewNode("FuncCall", $3, $5)) 
      }
    | var '=' term { $$ = NewNode("TermAssign", $1.Name, $3) }
    ;

loop
    : WHILE term '{' algo '}'          { $$ = NewNode("WhileLoop", "", $2, $4) }
    | DO '{' algo '}' UNTIL term      { $$ = NewNode("DoUntilLoop", "", $3, $6) }
    ;

branch
    : IF term '{' algo '}'                      { $$ = NewNode("IfBranch", "", $2, $4) }
    | IF term '{' algo '}' ELSE '{' algo '}'   { $$ = NewNode("IfElseBranch", "", $2, $4, $8) }
    ;

output
    : atom   { $$ = NewNode("OUTPUT", "", $1) }
    | STRING { $$ = NewNode("OUTPUT", $1) }
    ;

input
    : /* empty */                { $$ = NewNode("INPUT", "empty") }
    | atom                        { $$ = NewNode("INPUT", "", $1) }
    | atom atom                   { $$ = NewNode("INPUT", "", $1, $2) }
    | atom atom atom               { $$ = NewNode("INPUT", "", $1, $2, $3) }
    ;

term
    : atom                      { $$ = $1 }
    | '(' unop term ')'         { $$ = NewNode("UnaryOp", $2, $3) }
    | '(' term binop term ')'   { $$ = NewNode("BinaryOp", $3, $2, $4) }
    ;

unop
    : NEG { $$ = $1 }
    | NOT { $$ = $1 }
    ;

binop
    : EQ      { $$ = $1 }
    | GREATER { $$ = $1 }
    | OR      { $$ = $1 }
    | AND     { $$ = $1 }
    | PLUS    { $$ = $1 }
    | MINUS   { $$ = $1 }
    | MULT    { $$ = $1 }
    | DIV     { $$ = $1 }
    ;

%%
