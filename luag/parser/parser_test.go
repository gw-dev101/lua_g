package parser

import (
	"luag/lexer"
	"testing"
)

func TestPrintTokens(t *testing.T) {
	input := `
print("greater")
else
`

	l := lexer.NewLexer(input)

	for {
		token := l.NextToken()
		t.Logf("type=%v literal=%q", token.Type, token.Literal)

		if token.Type == lexer.TokenTypeEOF {
			break
		}
	}
}
func TestParser(t *testing.T) {
	input := `
local a = 10
if a > 5 then
    print("greater")
else
    print("lesser")
end

if a == 100 then
    print("no")
end`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	chunk := p.ParseChunk()

	for _, err := range p.Errors() {
		t.Logf("parser error: %v", err)
	}

	if len(chunk.Statements) != 3 {
		t.Logf("Expected 3 statements, got %d", len(chunk.Statements))
		actualStatements := make([]string, len(chunk.Statements))
		for i, stmt := range chunk.Statements {
			actualStatements[i] = StringifyStatement(stmt)
		}
		t.Logf("Actual statements: %v", actualStatements)
	}

	// Check the first statement (local a = 10)
	localStmt, ok := chunk.Statements[0].(*LocalStatement)
	if !ok {
		t.Fatalf("Expected first statement to be LocalStatement, got %T", chunk.Statements[0])
	}
	if localStmt.Name != "a" {
		t.Errorf("Expected variable name 'a', got '%s'", localStmt.Name)
	}

	// Check the second statement (if a > 5 then ...)
	ifStmt, ok := chunk.Statements[1].(*IfStatement)
	if !ok {
		t.Fatalf("Expected second statement to be IfStatement, got %T", chunk.Statements[1])
	}

	// Check the condition of the if statement
	binaryExpr, ok := ifStmt.Condition.(*BinaryExpression)
	if !ok {
		t.Fatalf("Expected condition to be BinaryExpression, got %T", ifStmt.Condition)
	}
	if binaryExpr.Operator != ">" {
		t.Errorf("Expected operator '>', got '%s'", binaryExpr.Operator)
	}

	// Ensure body statements (e.g., print(...)) are parsed
	if len(ifStmt.ThenBody) != 1 {
		t.Fatalf("Expected then-body to have 1 statement, got %d", len(ifStmt.ThenBody))
	}
	if _, ok := ifStmt.ThenBody[0].(*FunctionCallStatement); !ok {
		t.Fatalf("Expected then-body statement to be FunctionCallStatement, got %T", ifStmt.ThenBody[0])
	}
	if len(ifStmt.ElseBody) != 1 {
		t.Fatalf("Expected else-body to have 1 statement, got %d", len(ifStmt.ElseBody))
	}
	if _, ok := ifStmt.ElseBody[0].(*FunctionCallStatement); !ok {
		t.Fatalf("Expected else-body statement to be FunctionCallStatement, got %T", ifStmt.ElseBody[0])
	}
}
func TestFunctionDefStatement(t *testing.T) {
	input := `
function add(a, b)
    return a + b
end`

	l := lexer.NewLexer(input)
	p := NewParser(l)

	chunk := p.ParseChunk()

	// 1. Check for parser errors first
	for _, err := range p.Errors() {
		t.Errorf("parser error: %v", err)
	}

	// 2. Ensure we parsed exactly 1 statement
	if len(chunk.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(chunk.Statements))
	}

	// 3. Verify it's a FunctionDefStatement
	funcStmt, ok := chunk.Statements[0].(*FunctionDefStatement)
	if !ok {
		t.Fatalf("Expected statement to be *FunctionDefStatement, got %T", chunk.Statements[0])
	}

	// 4. Check function name
	if funcStmt.Name != "add" {
		t.Errorf("Expected function name 'add', got %q", funcStmt.Name)
	}

	// 5. Check parameter list
	expectedParams := []string{"a", "b"}
	if len(funcStmt.Parameters) != len(expectedParams) {
		t.Fatalf("Expected %d parameters, got %d", len(expectedParams), len(funcStmt.Parameters))
	}
	for i, param := range expectedParams {
		if funcStmt.Parameters[i] != param {
			t.Errorf("Expected parameter %d to be %q, got %q", i, param, funcStmt.Parameters[i])
		}
	}

	// 6. Check the function body (should contain 1 ReturnStatement)
	if len(funcStmt.Body) != 1 {
		t.Fatalf("Expected body to have 1 statement, got %d", len(funcStmt.Body))
	}

	returnStmt, ok := funcStmt.Body[0].(*ReturnStatement)
	if !ok {
		t.Fatalf("Expected body statement to be *ReturnStatement, got %T", funcStmt.Body[0])
	}

	// 7. Check the return expression (binary expression a + b)
	binaryExpr, ok := returnStmt.ReturnValue.(*BinaryExpression)
	if !ok {
		t.Fatalf("Expected ReturnValue to be *BinaryExpression, got %T", returnStmt.ReturnValue)
	}
	if binaryExpr.Operator != "+" {
		t.Errorf("Expected operator '+', got %q", binaryExpr.Operator)
	}
}
