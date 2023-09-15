// Package lexer implements lexical tokenization for the Monkey programming language.
package lexer

import "monkey/token"

// Lexer represents a lexical scanner for tokenizing the Monkey programming language.
type Lexer struct {
	input       string
	currentChar byte
	currentPos  int
	nextPos     int
}

// New returns a new instance of the Lexer, initialized with the provided input string.
func New(input string) *Lexer {
	l := &Lexer{input: input}

	if len(input) > 0 {
		l.currentChar = input[0]
		l.nextPos = 1
		// no need to set l.currentPos to 0, Go already instanciates it with its struct type's zero value
	}

	return l
}

// NextToken scans and returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.currentChar {
	case '=':
		tok = l.handleTwoCharToken(token.ASSIGN, '=', token.EQ)
	case '+':
		tok = l.handleSingleCharToken(token.PLUS)
	case '(':
		tok = l.handleSingleCharToken(token.LPAREN)
	case ')':
		tok = l.handleSingleCharToken(token.RPAREN)
	case '{':
		tok = l.handleSingleCharToken(token.LBRACE)
	case '}':
		tok = l.handleSingleCharToken(token.RBRACE)
	case ',':
		tok = l.handleSingleCharToken(token.COMMA)
	case ';':
		tok = l.handleSingleCharToken(token.SEMICOLON)
	case '-':
		tok = l.handleSingleCharToken(token.MINUS)
	case '/':
		tok = l.handleSingleCharToken(token.SLASH)
	case '*':
		tok = l.handleSingleCharToken(token.ASTERISK)
	case '<':
		tok = l.handleSingleCharToken(token.LT)
	case '>':
		tok = l.handleSingleCharToken(token.GT)
	case '!':
		tok = l.handleTwoCharToken(token.BANG, '=', token.NOT_EQ)
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.currentChar) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.currentChar) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.currentChar)}
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace advances the scanner until a non-whitespace character is encountered.
func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}

// handleTwoCharToken checks if the next character matches the expected character for a two-character token.
func (l *Lexer) handleTwoCharToken(defaultType token.TokenType, expectedChar byte, twoCharType token.TokenType) token.Token {
	if l.peekChar() == expectedChar {
		ch := l.currentChar
		l.readChar()
		return token.Token{Type: twoCharType, Literal: string(ch) + string(l.currentChar)}
	}
	return l.handleSingleCharToken(defaultType)
}

// handleSingleCharToken returns a token of the given type with the current character as its literal.
func (l *Lexer) handleSingleCharToken(t token.TokenType) token.Token {
	return token.Token{Type: t, Literal: string(l.currentChar)}
}

// readChar reads the next character from the input and updates the current and next positions.
func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.nextPos]
	}
	l.currentPos = l.nextPos
	l.nextPos++
}

// peekChar returns the next character from the input without advancing the current and next positions.
func (l *Lexer) peekChar() byte {
	if l.nextPos >= len(l.input) {
		return 0
	}
	return l.input[l.nextPos]
}

// readIdentifier scans an identifier from the input, capturing characters until a non-letter is encountered.
func (l *Lexer) readIdentifier() string {
	startPos := l.currentPos
	for isLetter(l.currentChar) {
		l.readChar()
	}
	return l.input[startPos:l.currentPos]
}

// readNumber scans a number from the input, capturing characters until a non-digit is encountered.
func (l *Lexer) readNumber() string {
	startPos := l.currentPos
	for isDigit(l.currentChar) {
		l.readChar()
	}
	return l.input[startPos:l.currentPos]
}

// Utility functions

// isDigit checks if the given byte is a valid digit.
func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

// isLetter checks if the given byte corresponds to a valid letter for identifiers in Monkey.
func isLetter(b byte) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_'
}
