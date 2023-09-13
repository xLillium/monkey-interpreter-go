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
	if l.nextPos >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.nextPos]
	}
	l.currentPos = l.nextPos
	l.nextPos++
}
