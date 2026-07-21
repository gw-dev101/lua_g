package lexer

import (
	"fmt"
	"os"

	"unicode"
)

type TokenType int

const (
	// Token types
	TokenTypeEOF TokenType = iota
	TokenTypeIdentifier
	TokenTypeNumber
	TokenTypeString
	TokenTypeOperator
	TokenTypeKeyword
	TokenTypePunctuation
)
const (
	// Keywords
	KeywordIf       = "if"
	KeywordThen     = "then"
	KeywordElse     = "else"
	KeywordEnd      = "end"
	KeywordLocal    = "local"
	KeywordFunction = "function"
	KeywordReturn   = "return"
)

type Position struct {
	Offset int
	Line   int
	Column int
}

type Span struct {
	Start Position
	End   Position // exclusive
}

type Token struct {
	Type    TokenType
	Literal string
	Span    Span
}

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for NUL, signifies end of input
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: TokenTypeOperator, Literal: string(ch) + string(l.ch), Span: Span{Start: Position{Offset: l.position, Line: l.line, Column: l.column}, End: Position{Offset: l.position, Line: l.line, Column: l.column}}}
		} else {
			tok = Token{Type: TokenTypeOperator, Literal: string(l.ch), Span: Span{Start: Position{Offset: l.position, Line: l.line, Column: l.column}, End: Position{Offset: l.position, Line: l.line, Column: l.column}}}
		}
	case '>':
		tok = Token{Type: TokenTypeOperator, Literal: string(l.ch), Span: Span{Start: Position{Offset: l.position, Line: l.line, Column: l.column}, End: Position{Offset: l.position, Line: l.line, Column: l.column}}}
	case '<':
		tok = Token{Type: TokenTypeOperator, Literal: string(l.ch), Span: Span{Start: Position{Offset: l.position, Line: l.line, Column: l.column}, End: Position{Offset: l.position, Line: l.line, Column: l.column}}}
	case '+', '-', '*', '/':
		tok = Token{Type: TokenTypeOperator, Literal: string(l.ch), Span: Span{Start: Position{Offset: l.position, Line: l.line, Column: l.column}, End: Position{Offset: l.position, Line: l.line, Column: l.column}}}
	case '(', ')', '{', '}', ',', ';':
		tok = Token{Type: TokenTypePunctuation, Literal: string(l.ch), Span: Span{Start: Position{Offset: l.position, Line: l.line, Column: l.column}, End: Position{Offset: l.position, Line: l.line, Column: l.column}}}
	case '"':
		tok.Literal = l.readString()
		tok.Type = TokenTypeString
		tok.Span = Span{Start: Position{Offset: l.position, Line: l.line, Column: l.column}, End: Position{Offset: l.position, Line: l.line, Column: l.column}}
	case 0:
		tok.Literal = ""
		tok.Type = TokenTypeEOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = TokenTypeIdentifier
			tok.Span = Span{Start: Position{Offset: l.position, Line: l.line, Column: l.column}, End: Position{Offset: l.position, Line: l.line, Column: l.column}}
			if isKeyword(tok.Literal) {
				tok.Type = TokenTypeKeyword
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = TokenTypeNumber
			tok.Span = Span{Start: Position{Offset: l.position, Line: l.line, Column: l.column}, End: Position{Offset: l.position, Line: l.line, Column: l.column}}
			return tok
		} else {
			fmt.Fprintf(os.Stderr, "Illegal character: %q at line %d, column %d\n", l.ch, l.line, l.column)
			l.readChar()
			return l.NextToken()
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	position := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func isKeyword(literal string) bool {
	switch literal {
	case KeywordIf, KeywordThen, KeywordElse, KeywordEnd, KeywordLocal:
		return true
	default:
		return false
	}
}
