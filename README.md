# 🐵 Monkey Interpreter in Go

A complete implementation of the Monkey programming language, written in Go. This project is inspired by the book [Writing an Interpreter in Go](https://interpreterbook.com/). by Thorsten Ball.

## 🚀 Features

- Lexical analysis (lexer)

- Parsing into an Abstract Syntax Tree (AST)

- Evaluation of expressions and statements

- Support for variables, functions, closures, and built-in functions

- Interactive Read-Eval-Print Loop (REPL)

- Modular and extensible architecture

📦 Project Structure

```
monkey-interpreter-go/
├── ast/        # Abstract Syntax Tree definitions
├── lexer/      # Lexical analyzer
├── parser/     # Parser for the Monkey language
├── token/      # Token definitions
├── repl/       # Read-Eval-Print Loop implementation
├── main.go     # Entry point for the interpreter
├── go.mod      # Go module file
```

🛠️ Getting Started

Prerequisites

Go 1.18 or higher

Installation

Clone the repository:
```
git clone https://github.com/xLillium/monkey-interpreter-go.git
cd monkey-interpreter-go
```
Build the project:

```go build -o monkey```

Run the REPL:

```./monkey```

🧪 Example Usage
```
>> let add = fn(a, b) { a + b; };
>> add(2, 3);
5
>> let factorial = fn(n) {
..   if (n == 0) {
..     1
..   } else {
..     n * factorial(n - 1);
..   }
.. };
>> factorial(5);
120
```    


🙏 Acknowledgments

Thanks to Thorsten Ball for his excellent book, which served as the foundation for this project.
