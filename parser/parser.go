// The parser package is responsible for parsing tokens produced
// by the lexer and constructing the AST (Abstract Syntax Tree).
package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// Precedence levels are used to dictate the order in which operators are parsed.
// In many parsers, this helps ensure that mathematical operations like multiplication
// and division are executed before addition and subtraction, for instance.
// The iota keyword in Go auto-increments, providing an easy way to assign increasing
// values to each item in the constant list.
const (
	_ int = iota
	LOWEST
	SUM    // +
	PREFIX // -X or !X
)

var precedences = map[token.TokenType]int{
	token.PLUS: SUM,
}

type (
	// prefixParseFn represents a function for parsing prefix expressions.
	prefixParseFn func() ast.Expression
	// infixParseFn represents a function for parsing infix expressions.
	infixParseFn func(ast.Expression) ast.Expression
)

// Parser represents the Monkey language parser structure.
type Parser struct {
	lexer          *lexer.Lexer
	current        token.Token
	peek           token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New initializes a new Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:          l,
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	// Set up initial tokens for curToken and peekToken.
	p.advanceToken()
	p.advanceToken()

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	return p
}

// Errors returns a slice of error messages encountered during parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

// registerPrefix registers a prefix parsing function for a given token type.
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers an infix parsing function for a given token type.
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// ParseProgram is the entry point of the parser. It constructs
// the AST by parsing statements and expressions from the input.
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !parser.tokenIs(parser.current, token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.advanceToken()
	}
	return program
}

// parseStatement dispatches the correct parsing function based on the current token type.
func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	statement := &ast.LetStatement{Token: p.current}

	if !p.advanceIfPeekIs(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

	if !p.advanceIfPeekIs(token.ASSIGN) {
		return nil
	}

	// TODO: Skip until we encounter a semicolon for simplicity now. We'll handle expressions later.
	p.skipToStatementEnd()
	return statement
}
func (p *Parser) parseReturnStatement() ast.Statement {
	statement := &ast.ReturnStatement{Token: p.current}
	p.skipToStatementEnd()
	return statement
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	statement := &ast.ExpressionStatement{Token: p.current}
	statement.Expression = p.parseExpression(LOWEST)
	if p.tokenIs(p.peek, token.SEMICOLON) {
		p.advanceToken()
	}
	return statement
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	integerLiteral := &ast.IntegerLiteral{Token: p.current}
	value, err := strconv.ParseInt(p.current.Literal, 0, 64)

	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.current.Literal))
		return nil
	}

	integerLiteral.Value = value
	return integerLiteral
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.current.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.current.Type)
		return nil
	}
	leftExp := prefix()

	// TODO: Skip until we encounter a semicolon for simplicity now. We'll handle expressions later.
	for !p.tokenIs(p.peek, token.SEMICOLON) && precedence <= p.peekPrecedence() {
		infix := p.infixParseFns[p.peek.Type]
		if infix == nil {
			return leftExp
		}
		p.advanceToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.current,
		Operator: p.current.Literal,
	}
	p.advanceToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.current,
		Operator: p.current.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.advanceToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.addError(fmt.Sprintf("no prefix parse function for %s found", t))
}

// Token navigation and validation functions.

// advanceToken advances to the next token.
func (p *Parser) advanceToken() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

// advanceIfPeekIs advances to the next token if the peek token matches the given type.
// If not, it logs an error and skips to the end of the statement.
func (parser *Parser) advanceIfPeekIs(t token.TokenType) bool {
	if parser.tokenIs(parser.peek, t) {
		parser.advanceToken()
		return true
	}
	parser.addError(fmt.Sprintf("expected next token to be %s, got %s instead", t, parser.peek.Type))
	parser.skipToStatementEnd()
	return false
}

// tokenIs checks if the given token has a specific type.
func (p *Parser) tokenIs(token token.Token, tokenType token.TokenType) bool {
	return token.Type == tokenType
}

// currentPrecedence returns the precedence of the current token.
func (p *Parser) currentPrecedence() int {
	if precedence, ok := precedences[p.current.Type]; ok {
		return precedence
	}
	return LOWEST
}

// peekPrecedence returns the precedence of the next token.
func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peek.Type]; ok {
		return prec
	}
	return LOWEST
}

// skipToStatementEnd skips tokens until a semicolon or EOF is encountered.
// This is useful for error recovery.
func (p *Parser) skipToStatementEnd() {
	for p.current.Type != token.SEMICOLON && p.current.Type != token.EOF {
		p.advanceToken()
	}
}

// addError logs a parsing error.
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}
