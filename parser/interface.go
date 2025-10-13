package parser

import (
	"SPL-compiler/lexer"
	"SPL-compiler/token"
	"fmt"
)

type LexerAdapter struct {
	L *lexer.Lexer
}

func (la *LexerAdapter) Lex(lval *YySymType) int {
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
		return PROC_KW
	case token.FUNC:
		return FUNC_KW
	case token.MAIN:
		return MAIN_KW
	case token.LOCAL:
		return LOCAL
	case token.VAR:
		return VAR_KW
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
	case token.NEG:
		return NEG
	case token.NOT:
		return NOT
	case token.EQ:
		return EQ
	case token.GT:
		return GREATER
	case token.OR:
		return OR
	case token.AND:
		return AND
	case token.PLUS:
		return PLUS
	case token.MINUS:
		return MINUS
	case token.MULT:
		return MULT
	case token.DIV:
		return DIV
	case token.INT:
		lval.Str = tok.Literal
		return NUMBER
	case token.STRING:
		lval.Str = tok.Literal
		return STRING
	case token.IDENT:
		lval.Str = tok.Literal
		return IDENT
	case token.ASSIGN:
		return '='
	case token.SEMICOLON:
		return ';'
	case token.LPAREN:
		return '('
	case token.RPAREN:
		return ')'
	case token.LBRACE:
		return '{'
	case token.RBRACE:
		return '}'
	case token.EOF:
		return 0
	}

	return 0
}

func (la *LexerAdapter) Error(msg string) {
	panic("Parse error: " + msg)
}
