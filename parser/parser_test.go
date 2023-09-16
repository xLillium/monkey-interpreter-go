package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

// TestLetStatements ensures that "let" statements are parsed correctly.
func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()

	// Check if parsing resulted in a valid program.
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	// Check if the correct number of "let" statements are parsed.
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	// Expected identifiers in the "let" statements.
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	// Iterate through the parsed statements and check if they match expected values.
	for i, test := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, test.expectedIdentifier) {
			return
		}
	}
}

// testLetStatement checks if the given statement is a valid "let" statement with the provided name.
func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	// Check if token literal is "let".
	if statement.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}

	// Type assert the statement into a "let" statement.
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", statement)
		return false
	}

	// Type assert the statement into a "let" statement.
	if letStatement.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStatement.Name.Value)
		return false
	}

	// Check if the token literal of the statement's name matches the expected name.
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStatement.Name)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
    return 5;
    return 10;
    return 993322;
`
	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()

	// Check if the correct number of "let" statements are parsed.
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`let x 5;`},
		{`let = 10;`},
		{`let 838383;`},
		{`let x 5;
let = 10;
let 838383;
`},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		parser.ParseProgram()
		if len(parser.Errors()) == 0 {
			t.Errorf("parser.ParseProgram() should have returned errors")
		} else {
			for _, err := range parser.Errors() {
				t.Logf("error: %s", err)
			}
		}
	}
}
