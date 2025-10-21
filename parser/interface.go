package parser

import (
	"fmt"

	"SPL-compiler/lexer"
	"SPL-compiler/token"
)

type LexerAdapter struct {
	L   *lexer.Lexer
	AST *ASTNode
}

func (la *LexerAdapter) Lex(lval *yySymType) int {
	tok := la.L.NextToken()
	if tok.Type == token.ILLEGAL {
		panic("Lexical error at line " + fmt.Sprint(tok.Line) +
			", column " + fmt.Sprint(tok.Column) +
			": '" + tok.Literal + "'")
	}

	switch tok.Type {
	case token.GLOB:
		return GLOB
	case token.PROC:
		return PROC
	case token.FUNC:
		return FUNC
	case token.MAIN:
		return MAIN
	case token.LOCAL:
		return LOCAL
	case token.VAR:
		return VAR
	case token.RETURN:
		return RETURN
	case token.HALT:
		return HALT
	case token.PRINT:
		return PRINT
	case token.WHILE:
		return WHILE
	case token.DO:
		return DO
	case token.UNTIL:
		return UNTIL
	case token.IF:
		return IF
	case token.ELSE:
		return ELSE
	case token.SEMICOLON:
		return SEMICOLON
	case token.LPAREN:
		return LPAREN
	case token.RPAREN:
		return RPAREN
	case token.LBRACE:
		return LBRACE
	case token.RBRACE:
		return RBRACE
	case token.ASSIGN:
		return ASSIGN

	case token.NEG:
		lval.Str = tok.Literal
		return NEG
	case token.NOT:
		lval.Str = tok.Literal
		return NOT
	case token.EQ:
		lval.Str = tok.Literal
		return EQ
	case token.GT:
		lval.Str = tok.Literal
		return GT
	case token.OR:
		lval.Str = tok.Literal
		return OR
	case token.AND:
		lval.Str = tok.Literal
		return AND
	case token.PLUS:
		lval.Str = tok.Literal
		return PLUS
	case token.MINUS:
		lval.Str = tok.Literal
		return MINUS
	case token.MULT:
		lval.Str = tok.Literal
		return MULT
	case token.DIV:
		lval.Str = tok.Literal
		return DIV
	case token.IDENT:
		lval.Str = tok.Literal
		return IDENT
	case token.INT:
		lval.Str = tok.Literal
		return NUMBER
	case token.STRING:
		lval.Str = tok.Literal
		return STRING
	case token.EOF:
		return 0
	}

	return 0
}

func (la *LexerAdapter) Error(msg string) {
	panic("Parse error: " + msg)
}
