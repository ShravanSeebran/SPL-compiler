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
%type <node> spl_prog variables var procdefs pdef funcdefs fdef body param maxthree mainprog atom algo instr_list instr assign loop branch output input term unop binop name
%left OR AND
%left PLUS MINUS
%left MULT DIV
%right NEG NOT

%start spl_prog

%%

spl_prog
    : GLOB '{' variables '}' PROC_KW '{' procdefs '}' FUNC_KW '{' funcdefs '}' MAIN_KW '{' mainprog '}'
      {
        $$ = NewNode("spl_prog", "SPL Program", $3, $7, $11, $15)
        ResultAST = $$
      }
    ;

variables
    : /* empty */        { $$ = NewNode("variables", "empty") }
    | var variables      { $$ = NewNode("variables", "", $1, $2) }
    ;

var
    : IDENT              { $$ = NewNode("var", $1) }
    ;

name
    : IDENT              { $$ = NewNode("name", $1) }
    ;

procdefs
    : /* empty */        { $$ = NewNode("procdefs", "empty") }
    | pdef procdefs      { $$ = NewNode("procdefs", "", $1, $2) }
    ;

pdef
    : name '(' param ')' '{' body '}'
      { $$ = NewNode("pdef", "", $1, $3, $6) }
    ;

funcdefs
    : /* empty */        { $$ = NewNode("funcdefs", "empty") }
    | fdef funcdefs      { $$ = NewNode("funcdefs", "", $1, $2) }
    ;

fdef
    : name '(' param ')' '{' body ';' RETURN atom '}'
        { $$ = NewNode("fdef","", $1, $3, $6, $9) }
    ;

body
    : LOCAL '{' maxthree '}' algo
        { $$ = NewNode("body", "", $3, $5) }
    ;

param
    : maxthree           { $$ = NewNode("param", "", $1) }
    ;

maxthree
    : /* empty */        { $$ = NewNode("maxthree", "empty") }
    | var                { $$ = NewNode("maxthree", "", $1) }
    | var var            { $$ = NewNode("maxthree", "", $1, $2) }
    | var var var        { $$ = NewNode("maxthree", "", $1, $2, $3) }
    ;

mainprog
    : VAR_KW '{' variables '}' algo 
        { $$ = NewNode("mainprog", "", $3, $5) }
    ;

atom
    : var         { $$ = NewNode("atom", "", $1) }
    | NUMBER      { $$ = NewNode("atom", $1) }
    ;

algo 
    : instr_list  { $$ = NewNode("algo", "", $1) }
    ;

instr_list 
    : instr                   { $$ = NewNode("instr_list", "", $1) }
    | instr_list ';' instr    { $$ = NewNode("instr_list", "", $1, $3) }
    ;

instr
    : HALT                    { $$ = NewNode("instr", "HALT") }
    | PRINT output            { $$ = NewNode("instr", "PRINT", $2) }
    | name '(' input ')'     { $$ = NewNode("instr", "", NewNode("call", "", $1, $3)) } 
    | assign                  { $$ = NewNode("instr", "", $1) }
    | loop                    { $$ = NewNode("instr", "", $1) }
    | branch                  { $$ = NewNode("instr", "", $1) }
    ;

assign
    : var '=' name '(' input ')' { 
        $$ = NewNode("assign", "", $1, NewNode("call", "", $3, $5))
      }
    | var '=' term { $$ = NewNode("assign", "", $1, $3) }
    ;

loop
    : WHILE term '{' algo '}'         { $$ = NewNode("loop", "while", $2, $4) }
    | DO '{' algo '}' UNTIL term      { $$ = NewNode("loop", "do_until", $3, $6) }
    ;

branch
    : IF term '{' algo '}'                     { $$ = NewNode("branch", "if", $2, $4) }
    | IF term '{' algo '}' ELSE '{' algo '}'   { $$ = NewNode("branch", "if_else", $2, $4, $8) }
    ;

output
    : atom      { $$ = NewNode("output", "", $1) }
    | STRING    { $$ = NewNode("output", $1) }
    ;

input
    : /* empty */          { $$ = NewNode("input", "empty") }
    | atom                 { $$ = NewNode("input", "", $1) }
    | atom atom            { $$ = NewNode("input", "", $1, $2) }
    | atom atom atom       { $$ = NewNode("input", "", $1, $2, $3) }
    ;

term
    : atom                     { $$ = NewNode("term", "", $1) }
    | '(' unop term ')'        { $$ = NewNode("term", "", $2, $3) }
    | '(' term binop term ')'  { $$ = NewNode("term", "", $2, $3, $4) }
    ;

unop
    : NEG { $$ = NewNode("unop" , "neg") }
    | NOT { $$ = NewNode("unop", "not") }
    ;

binop
    : EQ      { $$ = NewNode("binop" , "eq") }
    | GREATER { $$ = NewNode("binop", "greater") }
    | OR      { $$ = NewNode("binop", "or") }
    | AND     { $$ = NewNode("binop", "and") }
    | PLUS    { $$ = NewNode("binop", "plus") }
    | MINUS   { $$ = NewNode("binop", "minus") }
    | MULT    { $$ = NewNode("binop", "mult") }
    | DIV     { $$ = NewNode("binop", "div") }
    ;

%%