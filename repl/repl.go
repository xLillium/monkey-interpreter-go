// Package repl provides a Read-Eval-Print Loop (REPL) for the Monkey language.
// The REPL allows users to type Monkey code and immediately see the lexical tokens.
package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = "🐒💻>> "

// Start initializes the REPL for the Monkey language.
// It reads input line by line, lexically analyzes it, and prints out the recognized tokens.
// The loop continues until an end-of-file marker is encountered or "exit" is typed
//
// Parameters:
// in : An io.Reader from which input lines are read.
// out : An io.Writer to which the lexical tokens are written.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			fmt.Fprintln(out, "🐒🍌 Monkey says goodbye! 🍌🐒")
			return
		}

		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
