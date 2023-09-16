// The parser package is responsible for parsing tokens produced
// by the lexer and constructing the AST (Abstract Syntax Tree).
package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

// New initializes a new Parser instance. It will set the current
// and peek tokens by reading from the lexer twice.
func New(l *lexer.Lexer) *Parser {
	// Initialize the parser with the lexer and set the current and peek tokens.
	p := &Parser{lexer: l}
	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances the tokens by one.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// ParseProgram is the entry point of the parser. It constructs
// the AST by parsing statements and expressions from the input.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	statement := &ast.LetStatement{Token: p.curToken}

	if !p.advanceIfPeekIs(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.advanceIfPeekIs(token.ASSIGN) {
		return nil
	}

	// TODO: Skip until we encounter a semicolon for simplicity now. We'll handle expressions later.
	for p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.EOF {
		p.nextToken()
	}

	return statement
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
	return false
}
