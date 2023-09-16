package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"testing"
)

// parseInput takes a string input, tokenizes and parses it, then returns the resulting program.
func parseInput(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	return program
}

// checkParserErrors checks if the parser encountered any errors.
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

// ----- Tests for "let" statements -----

// TestLetStatementsParsing verifies the correct parsing of 'let' statements in the Monkey language.
func TestLetStatementsParsing(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	program := parseInput(t, input)
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

	program := parseInput(t, input)
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

// ----- Tests for string representation of AST nodes -----

// TestString verifies the correct string representation of AST nodes.
func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.IDENT, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

// ----- Tests for parsing expressions -----

// TestParseIdentifierExpression verifies the correct parsing of identifier expressions.
func TestParseIdentifierExpression(t *testing.T) {
	input := "myIdentifier;"

	program := parseInput(t, input)

	// Check for errors first.
	if len(program.Statements) != 1 {
		t.Fatalf("Expected a single statement, but got %d", len(program.Statements))
	}

	// Ensure that statement is an ExpressionStatement.
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// Ensure that the expression is an Identifier.
	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Identifier. got=%T", statement.Expression)
	}

	// Check the identifier's value.
	if ident.Value != "myIdentifier" {
		t.Errorf("ident.Value not %s. got=%s", "myIdentifier", ident.Value)
	}

	// Check the identifier's token literal.
	if ident.TokenLiteral() != "myIdentifier" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "myIdentifier", ident.TokenLiteral())
	}
}

func TestParseIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	program := parseInput(t, input)

	// Check for errors first.
	if len(program.Statements) != 1 {
		t.Fatalf("Expected a single statement, but got %d", len(program.Statements))
	}
	// Ensure that statement is an ExpressionStatement.
	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	// Ensure that the expression is an IntegerLiteral.
	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IntegerLiteral. got=%T", statement.Expression)
	}
	// Check the integer's value.
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	// Check the integer's token literal.
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

// TestParsingPrefixExpressions tests the parsing of prefix expressions
// such as ! and -.
func TestParsePrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input            string
		expectedOperator string
		expectedValue    int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, test := range prefixTests {
		program := parseInput(t, test.input)

		// Check for errors first.
		if len(program.Statements) != 1 {
			t.Fatalf("Expected a single statement, but got %d", len(program.Statements))
		}

		// Ensure that statement is an ExpressionStatement.
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("statement is not *ast.PrefixExpression. got=%T", statement.Expression)
		}

		if expression.Operator != test.expectedOperator {
			t.Fatalf("expression.Operator is not '%s'. got=%s", test.expectedOperator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Right, test.expectedValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", value, integer.TokenLiteral())
		return false
	}
	return true
}
