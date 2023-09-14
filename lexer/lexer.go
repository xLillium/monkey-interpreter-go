package lexer

import "monkey/token"

type Lexer struct {
	input       string
	currentChar byte
	currentPos  int
	nextPos     int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}

	if len(input) > 0 {
		l.currentChar = input[0]
		l.nextPos = 1
		// no need to set l.currentPos to 0, Go already instanciates it with its struct type's zero value
	}

	return l
}

func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.nextPos]
	}
	l.currentPos = l.nextPos
	l.nextPos++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.currentChar {
	case '=':
		tok = token.Token{Type: token.ASSIGN, Literal: string(l.currentChar)}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: string(l.currentChar)}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.currentChar)}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.currentChar)}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.currentChar)}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.currentChar)}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.currentChar)}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.currentChar)}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		tok = token.Token{Type: token.ILLEGAL, Literal: string(l.currentChar)}
	}

	l.readChar()
	return tok
}
