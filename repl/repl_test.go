// Package repl contains tests for the REPL of the Monkey language.
package repl

import (
	"bytes"
	"testing"
)

// TestREPL_SingleLineInput tests the REPL's handling of single line inputs.
func TestREPL_SingleLineInput(t *testing.T) {
	in := bytes.NewBufferString("let x = 5;\n")
	var out bytes.Buffer

	Start(in, &out)
	expectedOutput := `🐒💻>> {Type:LET Literal:let}
{Type:IDENT Literal:x}
{Type:= Literal:=}
{Type:INT Literal:5}
{Type:; Literal:;}
🐒💻>> `
	gotOutput := out.String()

	if expectedOutput != gotOutput {
		t.Errorf("Expected %q but got %q", expectedOutput, gotOutput)
	}
}

// TestREPL_IllegalToken tests the REPL's handling of illegal tokens.
func TestREPL_IllegalToken(t *testing.T) {
	in := bytes.NewBufferString("@#$%^&\n")
	var out bytes.Buffer

	Start(in, &out)

	expectedOutput := `🐒💻>> {Type:ILLEGAL Literal:@}
{Type:ILLEGAL Literal:#}
{Type:ILLEGAL Literal:$}
{Type:ILLEGAL Literal:%}
{Type:ILLEGAL Literal:^}
{Type:ILLEGAL Literal:&}
🐒💻>> `
	gotOutput := out.String()

	if expectedOutput != gotOutput {
		t.Errorf("Expected %q but got %q", expectedOutput, gotOutput)
	}
}

// TestREPL_Exit tests the REPL's behavior when receiving the "exit" command.
func TestREPL_Exit(t *testing.T) {
	in := bytes.NewBufferString("exit\n")
	var out bytes.Buffer

	Start(in, &out)
	expectedOutput := "🐒💻>> 🐒🍌 Monkey says goodbye! 🍌🐒\n"
	gotOutput := out.String()

	if expectedOutput != gotOutput {
		t.Errorf("Expected %q but got %q", expectedOutput, gotOutput)
	}
}
