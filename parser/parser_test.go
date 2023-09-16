package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

// parseInput takes a string input, tokenizes and parses it, then returns the resulting program.
func parseInput(input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	return p.ParseProgram()
}

// ----- Tests for "let" statements -----

// TestLetStatementsParsing verifies the correct parsing of 'let' statements in the Monkey language.
func TestLetStatementsParsing(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	program := parseInput(input)
	assertNumberOfStatements(t, program, 3)

	expectedIdentifiers := []string{"x", "y", "foobar"}
	for i, ident := range expectedIdentifiers {
		statement := program.Statements[i]
		assertLetStatement(t, statement, ident)
	}
}

// assertNumberOfStatements checks if a program contains the expected number of statements.
func assertNumberOfStatements(t *testing.T, program *ast.Program, num int) {
	if len(program.Statements) != num {
		t.Fatalf("Expected %d statements, but got %d", num, len(program.Statements))
	}
}

// assertLetStatement validates that a given statement is a correctly parsed 'let' statement.
func assertLetStatement(t *testing.T, statement ast.Statement, name string) {
	if statement.TokenLiteral() != "let" {
		t.Fatalf("Expected statement with 'let', but got %q", statement.TokenLiteral())
	}

	letStmt, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Fatalf("Expected *ast.LetStatement, but got %T", statement)
	}

	if letStmt.Name.Value != name {
		t.Errorf("Expected variable name to be %s, but got %s", name, letStmt.Name.Value)
	}
}

// ----- Tests for "return" statements -----

// TestReturnStatementsParsing verifies the correct parsing of 'return' statements in the Monkey language.
func TestReturnStatementsParsing(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	program := parseInput(input)
	assertNumberOfStatements(t, program, 3)

	for _, stmt := range program.Statements {
		assertReturnStatement(t, stmt)
	}
}

// assertReturnStatement validates that a given statement is a correctly parsed 'return' statement.
func assertReturnStatement(t *testing.T, stmt ast.Statement) {
	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Expected *ast.ReturnStatement, but got %T", stmt)
	}

	if returnStmt.TokenLiteral() != "return" {
		t.Fatalf("Expected 'return', but got %q", returnStmt.TokenLiteral())
	}
}

// ----- Tests for parser errors -----

// TestInvalidStatements checks the parser's ability to handle invalid input.
func TestParserErrors(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`let x 5;`},
		{`let = 10;`},
		{`let 838383;`},
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
