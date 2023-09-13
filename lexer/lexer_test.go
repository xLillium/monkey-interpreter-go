package lexer

import (
	"testing"
)

func TestLexerInitialization(t *testing.T) {
	input := "=+(),;"
	lexer := New(input)

	if lexer.input != input {
		t.Fatalf("Expected lexer input to be %s, got %s", input, lexer.input)
	}
}

func TestReadChar(t *testing.T) {
	input := "abc"
	lexer := New(input)

	lexer.readChar()
	if lexer.currentChar != 'a' {
		t.Fatalf("Expected first char to be 'a', got '%c'", lexer.currentChar)
	}
	if lexer.currentPos != 0 {
		t.Fatalf("Expected position to be 0, got %d", lexer.currentPos)
	}
	if lexer.nextPos != 1 {
		t.Fatalf("Expected readPosition to be 1, got %d", lexer.nextPos)
	}

	lexer.readChar()
	if lexer.currentChar != 'b' {
		t.Fatalf("Expected second char to be 'b', got '%c'", lexer.currentChar)
	}

	lexer.readChar()
	if lexer.currentChar != 'c' {
		t.Fatalf("Expected third char to be 'c', got '%c'", lexer.currentChar)
	}

	lexer.readChar()
	if lexer.currentChar != 0 {
		t.Fatalf("Expected ch to be 0 after reaching end, got '%c'", lexer.currentChar)
	}
}
