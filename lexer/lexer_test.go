package lexer

import (
	"testing"

	"monkey/token"
)

func TestLexerInitialization(t *testing.T) {
	input := "=+(),;"
	lexer := New(input)

	// Check if the input string is correctly set
	if lexer.input != input {
		t.Fatalf("Expected lexer input to be %s, got %s", input, lexer.input)
	}

	// Check the initial value of currentChar
	if lexer.currentChar != input[0] {
		t.Fatalf("Expected current char to be '%c', got '%c'", input[0], lexer.currentChar)
	}

	// Check the initial position values
	if lexer.currentPos != 0 {
		t.Fatalf("Expected current position to be 0, got %d", lexer.currentPos)
	}
	if lexer.nextPos != 1 {
		t.Fatalf("Expected next position to be 1, got %d", lexer.nextPos)
	}

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

// Struct for asserting the expected token types and literals that the lexer should produce
type tokenTest struct {
	expectedType    token.TokenType
	expectedLiteral string
}

// Helper function to run through a series of tokenTests,
// checking the lexer's output against the expected tokens
func runNextTokenTests(tests []tokenTest, lexer *Lexer, t *testing.T) {
	for i, testToken := range tests {
		token := lexer.NextToken()

		if token.Type != testToken.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, testToken.expectedType, token.Type)
		}
		if token.Literal != testToken.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, testToken.expectedLiteral, token.Literal)
		}

	}
}

func TestNextToken_SimpleTokens(t *testing.T) {
	input := "=+(){},;"
	lexer := New(input)

	tests := []tokenTest{
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

	runNextTokenTests(tests, lexer, t)
}

func TestNextToken_MonkeySourceCode(t *testing.T) {
	input := `let five = 5;
let ten = 10;
let add = fn(x, y) {
x + y;
};
let result = add(five, ten);
`
	lexer := New(input)

	tests := []tokenTest{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	runNextTokenTests(tests, lexer, t)
}
