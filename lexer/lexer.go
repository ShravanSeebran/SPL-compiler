package lexer

import (
	"SPL-compiler/token"
	"fmt"
)

// Lexer represents the lexical analyzer
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips whitespace characters (space, tab, newline)
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier following the pattern [a...z]{a...z}*{0...9}*
func (l *Lexer) readIdentifier() string {
	position := l.position

	// First character must be lowercase letter
	if !isLowercaseLetter(l.ch) {
		return ""
	}

	l.readChar()

	// Subsequent characters can be lowercase letters or digits
	for isLowercaseLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// readNumber reads a number following the pattern (0 | [1...9][0...9]*)
func (l *Lexer) readNumber() string {
	position := l.position

	if l.ch == '0' {
		l.readChar()
		return l.input[position:l.position]
	}

	// First digit must be 1-9
	if l.ch >= '1' && l.ch <= '9' {
		l.readChar()
		// Subsequent digits can be 0-9
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

// readString reads a string literal between quotation marks with max length 15
func (l *Lexer) readString() (string, error) {
	position := l.position + 1 // skip opening quote
	l.readChar()               // move past opening quote

	length := 0
	for l.ch != '"' && l.ch != 0 {
		if length >= 15 {
			return "", fmt.Errorf("string literal exceeds maximum length of 15 characters")
		}
		length++
		l.readChar()
	}

	if l.ch == 0 {
		return "", fmt.Errorf("unterminated string literal")
	}

	result := l.input[position:l.position]
	l.readChar() // move past closing quote
	return result, nil
}

// Token struct with line and column info for lexer package
type Token struct {
	Type    token.TokenType
	Literal string
	Line    int
	Column  int
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch, l.line, l.column)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.line, l.column)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.line, l.column)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line, l.column)
	case '>':
		tok = newToken(token.GT, l.ch, l.line, l.column)
	case '"':
		str, err := l.readString()
		if err != nil {
			tok.Type = token.ILLEGAL
			tok.Literal = err.Error()
			tok.Line = l.line
			tok.Column = l.column
		} else {
			tok.Type = token.STRING
			tok.Literal = str
			tok.Line = l.line
			tok.Column = l.column
		}
		return tok // readString already advances position
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLowercaseLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Line = l.line
			tok.Column = l.column
			if tok.Literal == "" {
				tok.Type = token.ILLEGAL
			} else {
				tok.Type = token.LookupIdent(tok.Literal)
			}
			return tok // readIdentifier already advances position
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Line = l.line
			tok.Column = l.column
			if tok.Literal == "" {
				tok.Type = token.ILLEGAL
			} else {
				tok.Type = token.INT
			}
			return tok // readNumber already advances position
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

// newToken creates a new token with the given parameters
func newToken(tokenType token.TokenType, ch byte, line, column int) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    line,
		Column:  column,
	}
}

// Helper functions
func isLowercaseLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// TokenizeInput tokenizes the entire input and returns a slice of tokens
func TokenizeInput(input string) []Token {
	lexer := New(input)
	var tokens []Token

	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	return tokens
}

// PrintTokens prints all tokens in a formatted way (useful for debugging)
func PrintTokens(tokens []Token) {
	fmt.Println("Tokens:")
	fmt.Println("-------")
	for i, tok := range tokens {
		fmt.Printf("%d: Type: %-10s Literal: %-10s Line: %d Column: %d\n",
			i+1, tok.Type, tok.Literal, tok.Line, tok.Column)
	}
}

// Expose constants for main package to use
const (
	EOF    = token.EOF
	STRING = token.STRING
)
