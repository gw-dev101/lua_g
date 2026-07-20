package parser

import (
	"luag/lexer"
	"strconv"
)

type Node interface{}
type Statement interface{ Node() }
type Expression interface{ Node() }

type Chunk struct {
	Statements []Statement
}

type LocalStatement struct {
	Name  string
	Value Expression
}

func (l *LocalStatement) Node() {}

type IfStatement struct {
	Condition Expression
	ThenBody  []Statement
	ElseBody  []Statement // Added to support 'else'
}

func (i *IfStatement) Node() {}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

// Node implements [Expression].
func (b *BinaryExpression) Node() {
	panic("unimplemented")
}

type NumberLiteral struct{ Value float64 }

func (n *NumberLiteral) Node() {}

type StringLiteral struct{ Value string }

func (s *StringLiteral) Node() {}

type Identifier struct{ Value string }

func (i *Identifier) Node() {}

// FunctionCallStatement represents a standalone statement like print("greater")
type FunctionCallStatement struct {
	Name string
	Args []Expression
}

func (f *FunctionCallStatement) Node() {}

type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
}

func (p *Parser) ParseChunk() *Chunk {
	return ParseChunk(p)
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) expectCurrent(tType lexer.TokenType, literal string) bool {
	if p.currentToken.Type != tType {
		return false
	}
	if literal != "" && p.currentToken.Literal != literal {
		return false
	}
	p.nextToken()
	return true
}

func ParseChunk(p *Parser) *Chunk {
	chunk := &Chunk{Statements: []Statement{}}
	for p.currentToken.Type != lexer.TokenTypeEOF {
		stmt := parse_statement(p)
		if stmt != nil {
			chunk.Statements = append(chunk.Statements, stmt)
		} else {
			// If we can't parse it as a known statement, advance to prevent infinite loops
			p.nextToken()
		}
	}
	return chunk
}

func parse_statement(p *Parser) Statement {
	switch p.currentToken.Type {
	case lexer.TokenTypeKeyword:
		switch p.currentToken.Literal {
		case lexer.KeywordIf:
			return parse_if_statement(p)
		case lexer.KeywordLocal:
			return parse_local_statement(p)
		default:
			return nil
		}
	case lexer.TokenTypeIdentifier:
		// Check for a simple function call like print(...)
		if p.peekToken.Type == lexer.TokenTypeOperator && p.peekToken.Literal == "(" {
			return parse_function_call(p)
		}
		return nil
	default:
		return nil
	}
}

func parse_local_statement(p *Parser) Statement {
	p.nextToken() // consume 'local'

	if p.currentToken.Type != lexer.TokenTypeIdentifier {
		return nil
	}
	varName := p.currentToken.Literal
	p.nextToken() // consume identifier

	if !p.expectCurrent(lexer.TokenTypeOperator, "=") {
		return nil
	}

	value := parse_expression(p, 0)
	return &LocalStatement{Name: varName, Value: value}
}

func parse_if_statement(p *Parser) Statement {
	p.nextToken() // consume 'if'

	condition := parse_expression(p, 0)

	if !p.expectCurrent(lexer.TokenTypeKeyword, lexer.KeywordThen) {
		return nil
	}

	thenBody := []Statement{}
	// Loop until we hit 'else' or 'end'
	for p.currentToken.Type != lexer.TokenTypeEOF {
		if p.currentToken.Type == lexer.TokenTypeKeyword &&
			(p.currentToken.Literal == lexer.KeywordEnd || p.currentToken.Literal == "else") {
			break
		}
		stmt := parse_statement(p)
		if stmt != nil {
			thenBody = append(thenBody, stmt)
		} else {
			p.nextToken()
		}
	}

	elseBody := []Statement{}
	// If we stopped because of an 'else', parse the else block
	if p.currentToken.Type == lexer.TokenTypeKeyword && p.currentToken.Literal == "else" {
		p.nextToken() // consume 'else'
		for p.currentToken.Type != lexer.TokenTypeEOF {
			if p.currentToken.Type == lexer.TokenTypeKeyword && p.currentToken.Literal == lexer.KeywordEnd {
				break
			}
			stmt := parse_statement(p)
			if stmt != nil {
				elseBody = append(elseBody, stmt)
			} else {
				p.nextToken()
			}
		}
	}

	if !p.expectCurrent(lexer.TokenTypeKeyword, lexer.KeywordEnd) {
		return nil // Missing 'end'
	}

	return &IfStatement{Condition: condition, ThenBody: thenBody, ElseBody: elseBody}
}

func parse_function_call(p *Parser) Statement {
	funcName := p.currentToken.Literal
	p.nextToken() // consume identifier
	p.nextToken() // consume '('

	args := []Expression{}
	if p.currentToken.Type != lexer.TokenTypeOperator || p.currentToken.Literal != ")" {
		args = append(args, parse_expression(p, 0))
	}

	if !p.expectCurrent(lexer.TokenTypeOperator, ")") {
		return nil
	}

	return &FunctionCallStatement{Name: funcName, Args: args}
}

// Simple operator precedence definition
func getPrecedence(op string) int {
	switch op {
	case "==", "~=", "<", ">", "<=", ">=":
		return 1
	case "+", "-":
		return 2
	case "*", "/":
		return 3
	default:
		return 0
	}
}

func parse_expression(p *Parser, minPrecedence int) Expression {
	var left Expression

	// Parse primary token
	switch p.currentToken.Type {
	case lexer.TokenTypeNumber:
		left = &NumberLiteral{Value: mustParseFloat(p.currentToken.Literal)}
		p.nextToken()
	case lexer.TokenTypeString:
		left = &StringLiteral{Value: p.currentToken.Literal}
		p.nextToken()
	case lexer.TokenTypeIdentifier:
		left = &Identifier{Value: p.currentToken.Literal}
		p.nextToken()
	default:
		return nil
	}

	// Pratt parsing / precedence climbing loop for binary operations like '>' or '=='
	for p.currentToken.Type == lexer.TokenTypeOperator {
		op := p.currentToken.Literal
		prec := getPrecedence(op)
		if prec < minPrecedence {
			break
		}

		p.nextToken() // consume operator
		right := parse_expression(p, prec+1)
		left = &BinaryExpression{Left: left, Operator: op, Right: right}
	}

	return left
}

func mustParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f

}
