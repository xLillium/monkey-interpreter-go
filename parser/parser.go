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
	// TODO: Loop through tokens and parse statements and expressions
	return nil
}
