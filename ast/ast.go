// Package ast defines the abstract syntax tree for the language.
package ast

import "monkey/token"

// Node represents a single node in the AST. Every node is expected
// to provide its associated token's literal representation.
type Node interface {
	TokenLiteral() string
}

// Statement represents a single statement in the Monkey language.
// All statement nodes will implement this interface.
type Statement interface {
	Node
	statementNode()
}

// Expression represents a single expression in the Monkey language.
// All expression nodes will implement this interface.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root of every AST the parser will produce. It contains
// a slice of statements, representing the Monkey program.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the literal representation of the token
// associated with the program.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Identifier represents an identifier in Monkey,
// which holds a token of type token.IDENT and its actual value.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// LetStatement represents a let statement in Monkey.
// It holds a token of type token.LET, the name of the identifier,
// and the expression representing its value.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
