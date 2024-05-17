package lexer

import (
	"testing"
)

type ExpectedToken struct {
	Type    TokenType
	Literal string
}

func TestLexerComments(t *testing.T) {
	input := `
	int i = 0;
	// This is a comment
	/* a!=i;
	+-*%/.
	void
	[]*/
	`

	expected := []ExpectedToken{
		{Type: "int", Literal: "int"},
		{Type: "IDENT", Literal: "i"},
		{Type: "=", Literal: "="},
		{Type: "INT", Literal: "0"},
		{Type: ";", Literal: ";"},
		{Type: "COMMENT", Literal: "// This is a comment"},
		{Type: "COMMENT", Literal: `/* a!=i;
	+-*%/.
	void
	[]*/`},
		{Type: "EOF", Literal: ""},
	}

	lexer := NewLexer(input)

	validateTokens(expected, lexer, t)
}

func TestLexerMiscCharacters(t *testing.T) {
	input := `
	int i = 0;
	a!=i;
	+-*/%.
	void
	[]
	`

	expected := []ExpectedToken{
		{Type: "int", Literal: "int"},
		{Type: "IDENT", Literal: "i"},
		{Type: "=", Literal: "="},
		{Type: "INT", Literal: "0"},
		{Type: ";", Literal: ";"},
		{Type: "IDENT", Literal: "a"},
		{Type: "!=", Literal: "!="},
		{Type: "IDENT", Literal: "i"},
		{Type: ";", Literal: ";"},
		{Type: "+", Literal: "+"},
		{Type: "-", Literal: "-"},
		{Type: "*", Literal: "*"},
		{Type: "/", Literal: "/"},
		{Type: "%", Literal: "%"},
		{Type: ".", Literal: "."},
		{Type: "void", Literal: "void"},
		{Type: "[", Literal: "["},
		{Type: "]", Literal: "]"},
		{Type: "EOF", Literal: ""},
	}

	lexer := NewLexer(input)

	validateTokens(expected, lexer, t)
}

func TestLexerStrings(t *testing.T) {
	input := `
		"Test one "
		"another \ttest \n"
		"a final test
		"test",i
		""
		`

	expected := []ExpectedToken{
		{Type: "STRING", Literal: "Test one "},
		{Type: "STRING", Literal: "another \\ttest \\n"},
		{Type: "STRING", Literal: "a final test"},
		{Type: "STRING", Literal: "test"},
		{Type: ",", Literal: ","},
		{Type: "IDENT", Literal: "i"},
		{Type: "STRING", Literal: ""},
		{Type: "EOF", Literal: ""},
	}
	lexer := NewLexer(input)

	validateTokens(expected, lexer, t)
	if !lexer.HasErrors() {
		t.Error("expected 1 error, found none")
	} else {
		if len(lexer.Errors()) != 1 {
			t.Errorf("expected 1 error, found %d", len(lexer.Errors()))
		}
		if lexer.Errors()[0].Error() != "[4:13] non-terminated string" {
			t.Errorf("error expected '[4:13] non-terminated string', got '%s'", lexer.Errors()[0].Error())
		}
	}
}

func TestLexerFactorialProgram(t *testing.T) {
	input := `
	int factorial(int n) {
		if (n == 0)
		  return 1;
		else
		  return n * factorial(n - 1);
	  }
	  
	  int main() {
		int i;
		for (i = 0; i <= 5; i++)
		  printf("Factorial of %d is %d\n", i, factorial(i));
		return 0;
	  }
	  `

	expected := []ExpectedToken{
		{Type: "int", Literal: "int"},
		{Type: "IDENT", Literal: "factorial"},
		{Type: "(", Literal: "("},
		{Type: "int", Literal: "int"},
		{Type: "IDENT", Literal: "n"},
		{Type: ")", Literal: ")"},
		{Type: "{", Literal: "{"},
		{Type: "if", Literal: "if"},
		{Type: "(", Literal: "("},
		{Type: "IDENT", Literal: "n"},
		{Type: "==", Literal: "=="},
		{Type: "INT", Literal: "0"},
		{Type: ")", Literal: ")"},
		{Type: "return", Literal: "return"},
		{Type: "INT", Literal: "1"},
		{Type: ";", Literal: ";"},
		{Type: "else", Literal: "else"},
		{Type: "return", Literal: "return"},
		{Type: "IDENT", Literal: "n"},
		{Type: "*", Literal: "*"},
		{Type: "IDENT", Literal: "factorial"},
		{Type: "(", Literal: "("},
		{Type: "IDENT", Literal: "n"},
		{Type: "-", Literal: "-"},
		{Type: "INT", Literal: "1"},
		{Type: ")", Literal: ")"},
		{Type: ";", Literal: ";"},
		{Type: "}", Literal: "}"},
		{Type: "int", Literal: "int"},
		{Type: "IDENT", Literal: "main"},
		{Type: "(", Literal: "("},
		{Type: ")", Literal: ")"},
		{Type: "{", Literal: "{"},
		{Type: "int", Literal: "int"},
		{Type: "IDENT", Literal: "i"},
		{Type: ";", Literal: ";"},
		{Type: "for", Literal: "for"},
		{Type: "(", Literal: "("},
		{Type: "IDENT", Literal: "i"},
		{Type: "=", Literal: "="},
		{Type: "INT", Literal: "0"},
		{Type: ";", Literal: ";"},
		{Type: "IDENT", Literal: "i"},
		{Type: "<=", Literal: "<="},
		{Type: "INT", Literal: "5"},
		{Type: ";", Literal: ";"},
		{Type: "IDENT", Literal: "i"},
		{Type: "++", Literal: "++"},
		{Type: ")", Literal: ")"},
		{Type: "printf", Literal: "printf"},
		{Type: "(", Literal: "("},
		{Type: "STRING", Literal: "Factorial of %d is %d\\n"},
		{Type: ",", Literal: ","},
		{Type: "IDENT", Literal: "i"},
		{Type: ",", Literal: ","},
		{Type: "IDENT", Literal: "factorial"},
		{Type: "(", Literal: "("},
		{Type: "IDENT", Literal: "i"},
		{Type: ")", Literal: ")"},
		{Type: ")", Literal: ")"},
		{Type: ";", Literal: ";"},
		{Type: "return", Literal: "return"},
		{Type: "INT", Literal: "0"},
		{Type: ";", Literal: ";"},
		{Type: "}", Literal: "}"},
		{Type: "EOF", Literal: ""},
	}
	lexer := NewLexer(input)

	validateTokens(expected, lexer, t)
}

func validateTokens(expected []ExpectedToken, lexer *Lexer, t *testing.T) {
	tokens := lexer.Tokens()
	if len(tokens) != len(expected) {
		t.Errorf("expected %d tokens, got %d tokens", len(expected), len(tokens))
	}
	for idx, tok := range tokens {
		if idx >= len(expected) {
			break
		}
		if expected[idx].Type != tok.Type {
			t.Errorf("Type error: expected %s, got %s", expected[idx].Type, tok.Type)
		} else if expected[idx].Literal != tok.Literal {
			t.Errorf("Literal error: expected %s, got %s", expected[idx].Literal, tok.Literal)
		}
	}
}
