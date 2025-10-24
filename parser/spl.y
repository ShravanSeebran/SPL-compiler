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

%token GLOB PROC FUNC MAIN
%token LOCAL VAR
%token RETURN
%token HALT PRINT
%token WHILE DO UNTIL
%token IF ELSE
%token SEMICOLON LPAREN RPAREN LBRACE RBRACE
%token ASSIGN

%token <Str> NEG NOT
%token <Str> EQ GT OR AND PLUS MINUS MULT DIV
%token <Str> IDENT NUMBER STRING
%type <node> spl_prog variables var name procdefs pdef funcdefs fdef body param maxthree mainprog atom algo instr assign loop branch output input term unop binop
%left OR AND
%left PLUS MINUS
%left MULT DIV
%right NEG NOT

%start spl_prog
%type <node> spl_prog

%%

spl_prog
    : GLOB LBRACE variables RBRACE
      PROC LBRACE procdefs  RBRACE
      FUNC LBRACE funcdefs  RBRACE
      MAIN LBRACE mainprog  RBRACE
      {
        $$ = NewNode("SPL_PROG", "", $3, $7, $11, $15)
        ResultAST = $$
        yylex.(*LexerAdapter).AST = ResultAST
      }
    ;

variables
    : /* empty */        { $$ = NewNode("VARIABLES", "empty") }
    | var variables      { $$ = NewNode("VARIABLES", "", $1, $2) }
    ;

var  : IDENT { $$ = NewNode("VAR", $1) };

name : IDENT { $$ = NewNode("NAME", $1) };

procdefs
    : /* empty */        { $$ = NewNode("PROCDEFS", "empty") }
    | pdef procdefs      { $$ = NewNode("PROCDEFS", "", $1, $2) }
    ;

pdef
    : name LPAREN param RPAREN LBRACE body RBRACE
      { $$ = NewNode("PDEF", "", $1, $3, $6) }
    ;

funcdefs
    : /* empty */        { $$ = NewNode("FUNCDEFS", "empty") }
    | fdef funcdefs      { $$ = NewNode("FUNCDEFS", "", $1, $2) }
    ;

fdef
    : name LPAREN param RPAREN LBRACE body RETURN atom RBRACE
        { $$ = NewNode("FDEF", "", $1, $3, $6, $8) }
    ;

body
    : LOCAL LBRACE maxthree RBRACE algo
        { $$ = NewNode("BODY", "", $3, $5) }
    ;

param : maxthree { $$ = NewNode("PARAM", "", $1) };

maxthree
    : /* empty */         { $$ = NewNode("MAXTHREE", "empty") }
    | var                 { $$ = NewNode("MAXTHREE", "", $1) }
    | var var             { $$ = NewNode("MAXTHREE", "", $1, $2) }
    | var var var         { $$ = NewNode("MAXTHREE", "", $1, $2, $3) }
    ;

mainprog
    : VAR LBRACE variables RBRACE algo 
        { $$ = NewNode("MAINPROG", "", $3, $5) }
    ;

atom
    : var { $$ = $1 }
    | NUMBER { $$ = NewNode("Number", $1) }
    ;

algo 
    : instr SEMICOLON { $$ = NewNode("ALGO", "", $1) }
    | instr SEMICOLON algo { $$ = NewNode("ALGO", "", $1, $3) }
    ;

instr
    : HALT                      { $$ = NewNode("INSTR", "halt") }
    | PRINT output              { $$ = NewNode("INSTR", "print", $2) }
    | name LPAREN input RPAREN  { $$ = NewNode("INSTR", "call", $1, $3) }
    | assign                    { $$ = NewNode("INSTR", "assign", $1) }
    | loop                      { $$ = NewNode("INSTR", "loop", $1) }
    | branch                    { $$ = NewNode("INSTR", "branch", $1) }
    ;

assign
    : var ASSIGN name LPAREN input RPAREN { $$ = NewNode("ASSIGN", "call", $1, $3, $5) }
    | var ASSIGN term                     { $$ = NewNode("ASSIGN", "", $1, $3) }
    ;

loop
    : WHILE term LBRACE algo RBRACE          { $$ = NewNode("LOOP", "while", $2, $4) }
    | DO LBRACE algo RBRACE UNTIL term       { $$ = NewNode("LOOP", "do", $3, $6) }
    ;

branch
    : IF term LBRACE algo RBRACE                           { $$ = NewNode("BRANCH", "if", $2, $4) }
    | IF term LBRACE algo RBRACE ELSE LBRACE algo RBRACE   { $$ = NewNode("BRANCH", "ifelse", $2, $4, $8) }
    ;

output
    : atom   { $$ = NewNode("OUTPUT", "atom", $1) }
    | STRING { $$ = NewNode("OUTPUT", $1) }
    ;

input
    : /* empty */     { $$ = NewNode("INPUT", "empty") }
    | atom            { $$ = NewNode("INPUT", "", $1) }
    | atom atom       { $$ = NewNode("INPUT", "", $1, $2) }
    | atom atom atom  { $$ = NewNode("INPUT", "", $1, $2, $3) }
    ;

term
    : atom                           { $$ = NewNode("TERM", "atom", $1) }
    | LPAREN unop term RPAREN        { $$ = NewNode("TERM", "unop", $2, $3) }
    | LPAREN term binop term RPAREN  { $$ = NewNode("TERM", "binop", $2, $3, $4) }
    ;

unop
    : NEG { $$ = NewNode("UNOP", $1) }
    | NOT { $$ = NewNode("UNOP", $1) }
    ;

binop
    : EQ      { $$ = NewNode("BINOP", $1) }
    | GT      { $$ = NewNode("BINOP", $1) }
    | OR      { $$ = NewNode("BINOP", $1) }
    | AND     { $$ = NewNode("BINOP", $1) }
    | PLUS    { $$ = NewNode("BINOP", $1) }
    | MINUS   { $$ = NewNode("BINOP", $1) }
    | MULT    { $$ = NewNode("BINOP", $1) }
    | DIV     { $$ = NewNode("BINOP", $1) }
    ;

%%

func Parse(lex yyLexer) (*ASTNode, error) {
    if yyParse(lex) != 0 {
        return nil, fmt.Errorf("syntax error")
    }

    // If you stored the result on your lexer (recommended)
    if l, ok := lex.(*LexerAdapter); ok {
        return l.AST, nil
    }

    return nil, fmt.Errorf("no result produced")
}
