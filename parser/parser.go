// The parser package is responsible for parsing tokens produced
// by the lexer and constructing the AST (Abstract Syntax Tree).
package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

// New initializes a new Parser instance. It will set the current
// and peek tokens by reading from the lexer twice.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
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
	return statement
}

// Token navigation and validation functions.
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
func (p *Parser) advanceIfPeekIs(t token.TokenType) bool {
	if p.nextTokenIs(t) {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
	p.skipStatement()
	return false
}

func (p *Parser) skipStatement() {
	for p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.EOF {
		p.nextToken()
	}
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}
