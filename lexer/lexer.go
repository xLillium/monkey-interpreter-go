package lexer

type Lexer struct {
	input string
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

