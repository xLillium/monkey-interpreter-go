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
