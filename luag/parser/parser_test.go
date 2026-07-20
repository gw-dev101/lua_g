package parser

import (
	"luag/lexer"
	"testing"
)

func TestParser(t *testing.T) {
	input := `local a = 10
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

	if len(chunk.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(chunk.Statements))
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
