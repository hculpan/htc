package lexer

import (
	"errors"
	"fmt"
	"unicode"
)

// TokenType represents the type of token.
type TokenType string

// Token represents a lexical token.
type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	Position int
}

// Token types
const (
	EOF          = "EOF"
	ILLEGAL      = "ILLEGAL"
	IDENT        = "IDENT"
	INT          = "INT"
	STRING       = "STRING"
	ASSIGN       = "="
	INCREMENT    = "++"
	DECREMENT    = "--"
	PLUS_EQUALS  = "+="
	MINUS_EQUALS = "-="
	PLUS         = "+"
	MINUS        = "-"
	ASTERISK     = "*"
	SLASH        = "/"
	PERCENT      = "%"
	BANG         = "!"
	EQ           = "=="
	NEQ          = "!="
	LT           = "<"
	GT           = ">"
	LE           = "<="
	GE           = ">="
	LPAREN       = "("
	RPAREN       = ")"
	LBRACKET     = "["
	RBRACKET     = "]"
	LBRACE       = "{"
	RBRACE       = "}"
	COMMA        = ","
	PERIOD       = "."
	SEMICOLON    = ";"
	IF           = "if"
	ELSE         = "else"
	WHILE        = "while"
	RETURN       = "return"
	INT_TYPE     = "int"
	VOID_TYPE    = "void"
	FOR          = "for"
	PRINTF       = "printf"
	COMMENT      = "COMMENT"
)

// Lexer represents a lexical scanner.
type Lexer struct {
	input         string
	position      int  // current position in input (points to current char)
	readPosition  int  // current reading position in input (after current char)
	ch            byte // current char under examination
	line          int
	tokenPosition int
	errors        []error
}

// NewLexer initializes a new instance of Lexer.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.line = 1
	l.readChar()
	l.errors = []error{}
	return l
}

func (l *Lexer) Tokens() []Token {
	result := []Token{}
	var tok Token
	for {
		tok = l.NextToken()
		result = append(result, tok)
		if tok.Type == EOF {
			break
		}
	}
	return result
}

func (l *Lexer) Errors() []error {
	return l.errors
}

func (l *Lexer) HasErrors() bool {
	return len(l.errors) != 0
}

func (l *Lexer) addError(msg string) {
	localMsg := fmt.Sprintf("[%d:%d] ", l.line, l.tokenPosition)
	localMsg = localMsg + msg
	l.errors = append(l.errors, errors.New(localMsg))
}

// readChar reads the next character and advances the positions in the input.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.tokenPosition++
}

// NextToken lexes the next token from the input.
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	repeat := false
	for {
		repeat = false
		switch l.ch {
		case '\n':
			l.line++
			l.readChar()
			l.skipWhitespace()
			l.tokenPosition = 0
			repeat = true
		case '=':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: l.line, Position: l.tokenPosition}
			} else {
				tok = newToken(ASSIGN, l.ch, l.line, l.position)
			}
		case '+':
			if l.peekChar() == '+' {
				ch := l.ch
				l.readChar()
				tok = Token{Type: INCREMENT, Literal: string(ch) + string(l.ch), Line: l.line, Position: l.tokenPosition}
			} else if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = Token{Type: PLUS_EQUALS, Literal: string(ch) + string(l.ch), Line: l.line, Position: l.tokenPosition}
			} else {
				tok = newToken(PLUS, l.ch, l.line, l.position)
			}
		case '-':
			if l.peekChar() == '-' {
				ch := l.ch
				l.readChar()
				tok = Token{Type: DECREMENT, Literal: string(ch) + string(l.ch), Line: l.line, Position: l.tokenPosition}
			} else if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = Token{Type: MINUS_EQUALS, Literal: string(ch) + string(l.ch), Line: l.line, Position: l.tokenPosition}
			} else {
				tok = newToken(MINUS, l.ch, l.line, l.position)
			}
		case '*':
			tok = newToken(ASTERISK, l.ch, l.line, l.position)
		case '/':
			if l.peekChar() == '/' {
				literal := l.readLineComment()
				tok.Type = COMMENT
				tok.Literal = literal
				tok.Line = l.line
				tok.Position = l.tokenPosition
			} else if l.peekChar() == '*' {
				literal := l.readBlockComment()
				tok.Type = COMMENT
				tok.Literal = literal
				tok.Line = l.line
				tok.Position = l.tokenPosition
			} else {
				tok = newToken(SLASH, l.ch, l.line, l.position)
			}
		case '%':
			tok = newToken(PERCENT, l.ch, l.line, l.position)
		case '!':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = Token{Type: NEQ, Literal: string(ch) + string(l.ch), Line: l.line, Position: l.tokenPosition}
			} else {
				tok = newToken(BANG, l.ch, l.line, l.position)
			}
		case '<':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = Token{Type: LE, Literal: string(ch) + string(l.ch), Line: l.line, Position: l.tokenPosition}
			} else {
				tok = newToken(LT, l.ch, l.line, l.position)
			}
		case '>':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = Token{Type: GE, Literal: string(ch) + string(l.ch), Line: l.line, Position: l.tokenPosition}
			} else {
				tok = newToken(GT, l.ch, l.line, l.position)
			}
		case '(':
			tok = newToken(LPAREN, l.ch, l.line, l.position)
		case ')':
			tok = newToken(RPAREN, l.ch, l.line, l.position)
		case '{':
			tok = newToken(LBRACE, l.ch, l.line, l.position)
		case '}':
			tok = newToken(RBRACE, l.ch, l.line, l.position)
		case '[':
			tok = newToken(LBRACKET, l.ch, l.line, l.position)
		case ']':
			tok = newToken(RBRACKET, l.ch, l.line, l.position)
		case ',':
			tok = newToken(COMMA, l.ch, l.line, l.position)
		case '.':
			tok = newToken(PERIOD, l.ch, l.line, l.position)
		case ';':
			tok = newToken(SEMICOLON, l.ch, l.line, l.position)
		case 0:
			tok.Literal = ""
			tok.Type = EOF
		case '"':
			literal, err := l.readString()
			if err != nil {
				l.addError(err.Error())
			}
			tok.Type = STRING
			tok.Literal = literal
			tok.Line = l.line
			tok.Position = l.tokenPosition
		default:
			if isLetter(l.ch) {
				literal := l.readIdentifier()
				tok.Type = lookupIdent(literal)
				tok.Literal = literal
				tok.Line = l.line
				tok.Position = l.tokenPosition
				return tok
			} else if isDigit(l.ch) {
				tok.Type = INT
				tok.Literal = l.readNumber()
				tok.Line = l.line
				tok.Position = l.tokenPosition
				return tok
			} else {
				tok = newToken(ILLEGAL, l.ch, l.line, l.position)
			}
		}

		if !repeat {
			break
		}
	}
	l.readChar()
	return tok
}

// newToken creates a new token with the given type and character.
func newToken(tokenType TokenType, ch byte, line int, position int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line, Position: position}
}

func (l *Lexer) readLineComment() string {
	position := l.position
	for l.ch != '\r' && l.ch != '\n' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readBlockComment() string {
	position := l.position
	for {
		if l.ch == '*' && l.peekChar() == '/' {
			break
		} else if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
	l.readChar()
	l.readChar()
	return l.input[position:l.position]
}

func (l *Lexer) readString() (string, error) {
	position := l.position
	l.readChar()
	for l.ch != '"' {
		if l.ch == '\n' {
			return l.input[position+1 : l.position], errors.New("non-terminated string")
		}
		l.readChar()
	}

	return l.input[position+1 : l.position], nil
}

// readIdentifier reads an identifier starting with a letter.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number starting with a digit.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// peekChar returns the next character without advancing the position.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// skipWhitespace skips any whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// isLetter checks if the character is a letter.
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

// isDigit checks if the character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// lookupIdent returns the correct token type for a given identifier.
func lookupIdent(ident string) TokenType {
	switch ident {
	case "if":
		return IF
	case "else":
		return ELSE
	case "while":
		return WHILE
	case "return":
		return RETURN
	case "int":
		return INT_TYPE
	case "void":
		return VOID_TYPE
	case "for":
		return FOR
	case "printf":
		return PRINTF
	default:
		return IDENT
	}
}
