package parser

import (
	"fmt"
	"luag/lexer"
	"strconv"
)

// Statement and Expression deliberately use different marker methods.
// This prevents AST nodes from accidentally satisfying both interfaces.

type Statement interface {
	statementNode()
}

type Expression interface {
	expressionNode()
}

type Chunk struct {
	Statements []Statement
}

type LocalStatement struct {
	Name  string
	Value Expression
}

func (*LocalStatement) statementNode() {}

type IfStatement struct {
	Condition Expression
	ThenBody  []Statement
	ElseBody  []Statement
}

func (*IfStatement) statementNode() {}

type FunctionCallStatement struct {
	Name string
	Args []Expression
}
type FunctionDefStatement struct {
	Name       string
	Parameters []string
	Body       []Statement
}

func (*FunctionDefStatement) statementNode()  {}
func (*FunctionCallStatement) statementNode() {}

type ReturnStatement struct {
	ReturnValue Expression
}

func (*ReturnStatement) statementNode() {}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (*BinaryExpression) expressionNode() {}

type NumberLiteral struct {
	Value float64
}

func (*NumberLiteral) expressionNode() {}

type StringLiteral struct {
	Value string
}

func (*StringLiteral) expressionNode() {}

type Identifier struct {
	Value string
}

func (*Identifier) expressionNode() {}

type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
	errors       []error
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []error{},
	}

	// Load currentToken and peekToken.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) addError(format string, args ...any) {
	p.errors = append(p.errors, fmt.Errorf(format, args...))
}

// expectCurrent checks currentToken and consumes it when it matches.
func (p *Parser) expectCurrent(
	tokenType lexer.TokenType,
	literal string,
) bool {
	if p.currentToken.Type != tokenType {
		p.addError(
			"expected token type %v, got %v (%q)",
			tokenType,
			p.currentToken.Type,
			p.currentToken.Literal,
		)
		return false
	}

	if literal != "" && p.currentToken.Literal != literal {
		p.addError(
			"expected token %q, got %q",
			literal,
			p.currentToken.Literal,
		)
		return false
	}

	p.nextToken()
	return true
}

func (p *Parser) ParseChunk() *Chunk {
	chunk := &Chunk{
		Statements: []Statement{},
	}

	for p.currentToken.Type != lexer.TokenTypeEOF {
		startToken := p.currentToken

		stmt := p.parseStatement()
		if stmt != nil {
			chunk.Statements = append(chunk.Statements, stmt)
			continue
		}

		p.addError(
			"unexpected token %q of type %v at statement level",
			p.currentToken.Literal,
			p.currentToken.Type,
		)

		// Guarantee progress when parsing failed.
		//
		// Some parsers may already have consumed part of the malformed
		// statement, so only advance when we are still on the same token.
		if p.currentToken == startToken {
			p.nextToken()
		}
	}

	return chunk
}

// Retained in case external code calls parser.ParseChunk(parser).
func ParseChunk(p *Parser) *Chunk {
	return p.ParseChunk()
}

func (p *Parser) parseStatement() Statement {
	switch p.currentToken.Type {
	case lexer.TokenTypeKeyword:
		switch p.currentToken.Literal {
		case lexer.KeywordIf:
			return p.parseIfStatement()

		case lexer.KeywordLocal:
			return p.parseLocalStatement()

		case lexer.KeywordFunction:
			return p.parseFunctionDefStatement()

		default:
			return nil
		}

	case lexer.TokenTypeIdentifier:
		if p.peekToken.Type == lexer.TokenTypePunctuation &&
			p.peekToken.Literal == "(" {
			return p.parseFunctionCallStatement()
		}

		return nil

	default:
		return nil
	}
}
func (p *Parser) parseLocalStatement() Statement {
	// currentToken is "local".
	p.nextToken()

	if p.currentToken.Type != lexer.TokenTypeIdentifier {
		p.addError(
			"expected identifier after local, got %q",
			p.currentToken.Literal,
		)
		return nil
	}

	varName := p.currentToken.Literal
	p.nextToken()

	if !p.expectCurrent(lexer.TokenTypeOperator, "=") {
		return nil
	}

	value := p.parseExpression(1)
	if value == nil {
		p.addError(
			"expected expression after '=' in local declaration %q",
			varName,
		)
		return nil
	}

	// Do not call nextToken here.
	//
	// parseExpression already leaves currentToken positioned on the first
	// token after the expression.
	return &LocalStatement{
		Name:  varName,
		Value: value,
	}
}

func (p *Parser) parseIfStatement() Statement {
	// currentToken is "if".
	p.nextToken()

	condition := p.parseExpression(1)
	if condition == nil {
		p.addError("expected condition after 'if'")
		return nil
	}

	// parseExpression leaves currentToken on "then".
	if !p.expectCurrent(lexer.TokenTypeKeyword, lexer.KeywordThen) {
		return nil
	}

	thenBody := p.parseBlockUntil(
		lexer.KeywordElse,
		lexer.KeywordEnd,
	)

	elseBody := []Statement{}

	if p.currentToken.Type == lexer.TokenTypeKeyword &&
		p.currentToken.Literal == lexer.KeywordElse {
		p.nextToken()

		elseBody = p.parseBlockUntil(lexer.KeywordEnd)
	}

	if !p.expectCurrent(lexer.TokenTypeKeyword, lexer.KeywordEnd) {
		p.addError("missing 'end' for if statement")
		return nil
	}

	return &IfStatement{
		Condition: condition,
		ThenBody:  thenBody,
		ElseBody:  elseBody,
	}
}

func (p *Parser) parseBlockUntil(keywords ...string) []Statement {
	statements := []Statement{}

	for p.currentToken.Type != lexer.TokenTypeEOF {
		if p.currentToken.Type == lexer.TokenTypeKeyword &&
			containsString(keywords, p.currentToken.Literal) {
			break
		}

		startToken := p.currentToken

		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
			continue
		}

		p.addError(
			"unexpected token %q inside block",
			p.currentToken.Literal,
		)

		if p.currentToken == startToken {
			p.nextToken()
		}
	}

	return statements
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func (p *Parser) parseFunctionCallStatement() Statement {
	funcName := p.currentToken.Literal

	// print -> (
	p.nextToken()

	if !p.expectCurrent(lexer.TokenTypePunctuation, "(") {
		return nil
	}

	args := []Expression{}

	if !isPunctuation(p.currentToken, ")") {
		for {
			arg := p.parseExpression(1)
			if arg == nil {
				p.addError(
					"expected expression in argument list of %q",
					funcName,
				)
				return nil
			}

			args = append(args, arg)

			// parseExpression leaves currentToken on either "," or ")".
			if !isPunctuation(p.currentToken, ",") {
				break
			}

			// Consume the comma and move to the next argument.
			p.nextToken()
		}
	}

	// Do not call nextToken before this.
	if !p.expectCurrent(lexer.TokenTypePunctuation, ")") {
		return nil
	}

	return &FunctionCallStatement{
		Name: funcName,
		Args: args,
	}
}

func isPunctuation(token lexer.Token, literal string) bool {
	return token.Type == lexer.TokenTypePunctuation &&
		token.Literal == literal
}

// getPrecedence returns the precedence and whether the operator is supported.
func getPrecedence(op string) (int, bool) {
	switch op {
	case "==", "~=", "<", ">", "<=", ">=":
		return 1, true

	case "+", "-":
		return 2, true

	case "*", "/":
		return 3, true

	default:
		return 0, false
	}
}

// parseExpression implements precedence climbing.
//
// It consumes every token belonging to the expression and leaves currentToken
// positioned on the first token after the expression.
func (p *Parser) parseExpression(minPrecedence int) Expression {
	left := p.parsePrimaryExpression()
	if left == nil {
		return nil
	}

	for p.currentToken.Type == lexer.TokenTypeOperator {
		operator := p.currentToken.Literal

		precedence, supported := getPrecedence(operator)
		if !supported || precedence < minPrecedence {
			break
		}

		p.nextToken()

		right := p.parseExpression(precedence + 1)
		if right == nil {
			p.addError(
				"expected expression after operator %q",
				operator,
			)
			return nil
		}

		left = &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

// a function definition statement is of the form:
// function <name>(<params>) <body> end , returns isn't actually necessary
func (p *Parser) parseFunctionDefStatement() Statement {
	// currentToken is "function".
	p.nextToken()

	if p.currentToken.Type != lexer.TokenTypeIdentifier {
		p.addError(
			"expected identifier after function, got %q",
			p.currentToken.Literal,
		)
		return nil
	}

	funcName := p.currentToken.Literal
	p.nextToken()

	if !p.expectCurrent(lexer.TokenTypePunctuation, "(") {
		return nil
	}

	params := []string{}

	if !isPunctuation(p.currentToken, ")") {
		for {
			if p.currentToken.Type != lexer.TokenTypeIdentifier {
				p.addError(
					"expected identifier in parameter list of %q",
					funcName,
				)
				return nil
			}

			params = append(params, p.currentToken.Literal)
			p.nextToken()

			if !isPunctuation(p.currentToken, ",") {
				break
			}

			p.nextToken()
		}
	}

	if !p.expectCurrent(lexer.TokenTypePunctuation, ")") {
		return nil
	}

	body := p.parseBlockUntil(lexer.KeywordEnd)

	if !p.expectCurrent(lexer.TokenTypeKeyword, lexer.KeywordEnd) {
		p.addError("missing 'end' for function definition %q", funcName)
		return nil
	}

	return &FunctionDefStatement{
		Name:       funcName,
		Parameters: params,
		Body:       body,
	}
}
func (p *Parser) parsePrimaryExpression() Expression {
	switch p.currentToken.Type {
	case lexer.TokenTypeNumber:
		value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
		if err != nil {
			p.addError(
				"invalid number literal %q",
				p.currentToken.Literal,
			)
			return nil
		}

		p.nextToken()

		return &NumberLiteral{
			Value: value,
		}

	case lexer.TokenTypeString:
		value := p.currentToken.Literal
		p.nextToken()

		return &StringLiteral{
			Value: value,
		}

	case lexer.TokenTypeIdentifier:
		value := p.currentToken.Literal
		p.nextToken()

		return &Identifier{
			Value: value,
		}
	//if its a function definition statement, parse it
	case lexer.TokenTypeKeyword:
		return nil
	}

	return nil
}

// Debug stringification

func StringifyStatement(stmt Statement) string {
	switch s := stmt.(type) {
	case *LocalStatement:
		return fmt.Sprintf(
			"LocalStatement(Name: %s, Value: %s)",
			s.Name,
			StringifyExpression(s.Value),
		)

	case *IfStatement:
		thenBody := stringifyStatements(s.ThenBody)
		elseBody := stringifyStatements(s.ElseBody)

		return fmt.Sprintf(
			"IfStatement(Condition: %s, ThenBody: [%s], ElseBody: [%s])",
			StringifyExpression(s.Condition),
			thenBody,
			elseBody,
		)

	case *FunctionCallStatement:
		args := ""

		for index, arg := range s.Args {
			if index > 0 {
				args += ", "
			}

			args += StringifyExpression(arg)
		}

		return fmt.Sprintf(
			"FunctionCallStatement(Name: %s, Args: [%s])",
			s.Name,
			args,
		)

	default:
		return "UnknownStatement"
	}
}

func stringifyStatements(statements []Statement) string {
	result := ""

	for index, stmt := range statements {
		if index > 0 {
			result += "; "
		}

		result += StringifyStatement(stmt)
	}

	return result
}

func StringifyExpression(expr Expression) string {
	switch e := expr.(type) {
	case *NumberLiteral:
		return fmt.Sprintf(
			"NumberLiteral(Value: %s)",
			strconv.FormatFloat(e.Value, 'f', -1, 64),
		)

	case *StringLiteral:
		return fmt.Sprintf(
			"StringLiteral(Value: %q)",
			e.Value,
		)

	case *Identifier:
		return fmt.Sprintf(
			"Identifier(Value: %s)",
			e.Value,
		)

	case *BinaryExpression:
		return fmt.Sprintf(
			"BinaryExpression(Left: %s, Operator: %q, Right: %s)",
			StringifyExpression(e.Left),
			e.Operator,
			StringifyExpression(e.Right),
		)

	case nil:
		return "nil"

	default:
		return "UnknownExpression"
	}
}
