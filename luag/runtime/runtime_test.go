package runtime

import (
	"luag/lexer"
	"luag/parser"
	"testing"
)

func TestRuntime(t *testing.T) {
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
	p := parser.NewParser(l)
	chunk := p.ParseChunk()

	r := NewRuntime()
	r.ExecuteChunk(chunk)

	// Check if variable 'a' is set correctly
	if val, exists := r.Variables["a"]; !exists || val != float64(10) {
		t.Errorf("Expected variable 'a' to be 10, got %v (%T)", val, val)
	}

}
