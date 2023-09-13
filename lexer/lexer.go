package lexer

type Lexer struct {
	input       string
	currentChar byte
	currentPos  int
	nextPos     int
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) readChar() {
}
