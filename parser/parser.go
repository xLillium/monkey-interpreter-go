// The parser package is responsible for parsing tokens produced
// by the lexer and constructing the AST (Abstract Syntax Tree).
package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// Precedence levels are used to dictate the order in which operators are parsed.
// In many parsers, this helps ensure that mathematical operations like multiplication
// and division are executed before addition and subtraction, for instance.
// The iota keyword in Go auto-increments, providing an easy way to assign increasing
// values to each item in the constant list.
const (
	_ int = iota
	LOWEST
)

type (
	// prefixParseFn represents a function for parsing prefix expressions.
	prefixParseFn func() ast.Expression
)

// Parser represents the Monkey language parser structure.
type Parser struct {
	lexer          *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
}

// New initializes a new Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, prefixParseFns: make(map[token.TokenType]prefixParseFn)}

	// Set up initial tokens for curToken and peekToken.
	p.nextToken()
	p.nextToken()

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	return p
}

// Errors returns a slice of error messages encountered during parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// ParseProgram is the entry point of the parser. It constructs
// the AST by parsing statements and expressions from the input.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}
	return program
}

// parseStatement dispatches the correct parsing function based on the current token type.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	statement := &ast.LetStatement{Token: p.curToken}

	if !p.advanceIfPeekIs(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.advanceIfPeekIs(token.ASSIGN) {
		return nil
	}

	// TODO: Skip until we encounter a semicolon for simplicity now. We'll handle expressions later.
	p.skipStatement()
	return statement
}
func (p *Parser) parseReturnStatement() ast.Statement {
	statement := &ast.ReturnStatement{Token: p.curToken}
	p.skipStatement()
	return statement
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	statement := &ast.ExpressionStatement{Token: p.curToken}
	statement.Expression = p.parseExpression(LOWEST)
	if p.nextTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

// Token navigation and validation functions.

// nextToken advances the parser to the next token.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// currentTokenIs checks if the current token has a specific type.
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// nextTokenIs checks if the next token has a specific type.
func (p *Parser) nextTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// advanceIfPeekIs advances to the next token if the peek token matches the given type.
// If not, it logs an error and skips to the end of the statement.
func (p *Parser) advanceIfPeekIs(t token.TokenType) bool {
	if p.nextTokenIs(t) {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
	p.skipStatement()
	return false
}

// skipStatement skips tokens until a semicolon or EOF is encountered.
// This is useful for error recovery.
func (p *Parser) skipStatement() {
	for p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.EOF {
		p.nextToken()
	}
}

// addError logs a parsing error.
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}
