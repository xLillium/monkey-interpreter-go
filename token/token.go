// Package token defines the set of lexical tokens for the Monkey programming language.
package token

// TokenType represents the type of a lexical token.
type TokenType string

// Token represents a lexical token with a type and literal string value.
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// ILLEGAL represents a token/character that we don't know how to handle.
	ILLEGAL = "ILLEGAL"

	// EOF signals the end of parsing, representing the end of our input.
	EOF = "EOF"

	// IDENT and INT are used for user-defined identifiers (e.g. variable names) and integer literals.
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456789

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"

	// Delimiters such as comma, semicolon, and various brackets.
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	// EQ and NOT_EQ are used for equality checking.
	EQ     = "=="
	NOT_EQ = "!="
)

// keywords maps Monkey's keyword strings to their TokenType values.
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent checks the keywords table to see if the given identifier is a reserved keyword.
// If it's not found, the identifier is assumed to be a user-defined name and IDENT is returned.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
