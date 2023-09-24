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
		expectedValue    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
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

		if !testLiteralExpression(t, expression.Right, test.expectedValue) {
			return
		}
	}
}

// TestParsingInfixExpressions tests the parsing of infix expressions
func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input              string
		expectedLeftValue  interface{}
		expectedOperator   string
		expectedRightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	for _, test := range infixTests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		// Check for errors first.
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		// Ensure that statement is an ExpressionStatement.
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, test.expectedLeftValue,
			test.expectedOperator, test.expectedRightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParsingBooleanExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}
	for _, test := range tests {
		program := parseInput(t, test.input)
		if len(program.Statements) != 1 {
			t.Fatalf("Expected a single statement, but got %d", len(program.Statements))
		}
		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		boolean, ok := statement.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("statement.Expression is not *ast.Boolean. got=%T", statement.Expression)
		}
		if boolean.Value != test.expected {
			t.Errorf("boolean.Value not %t. got=%t", test.expected, boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

// ----- Helper functions -----

// testInfixExpression checks if an expression is an InfixExpression
func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},

	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	// Check for errors first.
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	// Check if the left expression is correct.
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	// Check if the operator is correct.
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	// Check if the right expression is correct.
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

// testLiteralExpression checks if an expression is a literal expression.
func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

// testIntegerLiteral checks if an expression is an IntegerLiteral.
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	// Check for errors first.
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	// Check if the integer's value is correct.
	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}

	// Check if the integer's token literal is correct.
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", value, integer.TokenLiteral())
		return false
	}
	return true
}

// testIdentifier checks if an expression is an Identifier.
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	// Check for errors first.
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	// Check if the identifier's value is correct.
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	// Check if the identifier's token literal is correct.
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

// testBooleanLiteral checks if an expression is a Boolean.
func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	// Check for errors first.
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	// Check if the boolean's value is correct.
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	// Check if the boolean's token literal is correct.
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}
