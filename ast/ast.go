// Package ast defines the abstract syntax tree for the language.
package ast

import (
	"bytes"
	"monkey/token"
)

// Node represents a single node in the AST. Every node is expected
// to provide its associated token's literal representation.
type Node interface {
	TokenLiteral() string
	String() string
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

// String returns the string representation of the program.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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

// String returns the string representation of the identifier.
func (i *Identifier) String() string {
	return i.Value
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

// String returns the string representation of the let statement.
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token token.Token // the 'return' token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String returns the string representation of the return statement.
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
