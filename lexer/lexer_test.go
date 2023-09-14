package lexer

import (
	"testing"

	"monkey/token"
)

func TestLexerInitialization(t *testing.T) {
	input := "=+(),;"
	lexer := New(input)

	if lexer.input != input {
		t.Fatalf("Expected lexer input to be %s, got %s", input, lexer.input)
	}
}

func TestReadCharProgression(t *testing.T) {
	input := "abc"
	lexer := New(input)

	if lexer.currentChar != 'a' {
		t.Fatalf("Expected current char to be 'a', got '%c'", lexer.currentChar)
	}
	if lexer.currentPos != 0 {
		t.Fatalf("Expected current position to be 0, got %d", lexer.currentPos)
	}
	if lexer.nextPos != 1 {
		t.Fatalf("Expected next position to be 1, got %d", lexer.nextPos)
	}

	lexer.readChar()
	if lexer.currentChar != 'b' {
		t.Fatalf("Expected current char to be 'b', got '%c'", lexer.currentChar)
	}

	lexer.readChar()
	if lexer.currentChar != 'c' {
		t.Fatalf("Expected current char to be 'c', got '%c'", lexer.currentChar)
	}

	lexer.readChar()
	if lexer.currentChar != 0 {
		t.Fatalf("Expected current char to be 0 after reaching end, got '%c'", lexer.currentChar)
	}
}

func TestNextToken_SimpleTokens(t *testing.T) {
	input := "=+(){},;"
	lexer := New(input)

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		token := lexer.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, token.Type)
		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, token.Literal)
		}

	}
}
