package runtime

import (
	"fmt"
	"luag/lexer"
	"luag/parser"
	"testing"
)

func TestRuntimeHelloWorld(t *testing.T) {
	input := `print("Hello, World!")`

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	chunk := p.ParseChunk()

	r := NewRuntime()
	output := r.ExecuteChunkWithOutput(chunk)
	fmt.Printf("Output:\n%s", output)
	if output != "Hello, World!\n" {
		t.Errorf("Expected output 'Hello, World!', got %q", output)
	}
}
func TestRuntimeIF(t *testing.T) {
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
	output := r.ExecuteChunkWithOutput(chunk)
	fmt.Printf("Output:\n%s", output)
	// Check if variable 'a' is set correctly
	if val, exists := r.Variables["a"]; !exists || val != float64(10) {
		t.Errorf("Expected variable 'a' to be 10, got %v (%T)", val, val)
	}

}
func TestRuntimeTableLiteralsmall(t *testing.T) {
	input := `local t = {1, 2, 3}
print(t[1])
print(t[2])
print(t[3])
`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	chunk := p.ParseChunk()

	r := NewRuntime()
	output := r.ExecuteChunkWithOutput(chunk)
	fmt.Printf("Output:\n%s", output)
	if output != "1\n2\n3\n" {
		t.Errorf("Expected output '1\\n2\\n3\\n', got %q", output)
	}
}
func TestRuntimeTableLiteralWithKeys(t *testing.T) {
	input := `local t = {a = 1, b = 2, c = 3}
print(t["a"])
print(t["b"])
print(t["c"])
`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	chunk := p.ParseChunk()

	r := NewRuntime()
	output := r.ExecuteChunkWithOutput(chunk)
	fmt.Printf("Output:\n%s", output)
	if output != "1\n2\n3\n" {
		t.Errorf("Expected output '1\\n2\\n3\\n', got %q", output)
	}
}
func TestRuntimeTableLiteralMixed(t *testing.T) {
	input := `local t = {1, 2, a = 3, b = 4}
print(t[1])
print(t[2])
print(t["a"])
print(t["b"])
`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	chunk := p.ParseChunk()

	r := NewRuntime()
	output := r.ExecuteChunkWithOutput(chunk)
	fmt.Printf("Output:\n%s", output)
	if output != "1\n2\n3\n4\n" {
		t.Errorf("Expected output '1\\n2\\n3\\n4\\n', got %q", output)
	}
}

func TestRuntimeNestedTableLiteral(t *testing.T) {
	input := `local t = {1, 2, {a = 3, b = 4}}
print(t[1])
print(t[2])
print(t[3]["a"])
print(t[3]["b"])
`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	chunk := p.ParseChunk()

	r := NewRuntime()
	output := r.ExecuteChunkWithOutput(chunk)
	fmt.Printf("Output:\n%s", output)
	if output != "1\n2\n3\n4\n" {
		t.Errorf("Expected output '1\\n2\\n3\\n4\\n', got %q", output)
	}
}
func TestRuntimeTableLiteralWithExpressions(t *testing.T) {
	input := `local t = {1 + 1, 2 * 2, a = 3 - 1, b = 4 / 2}
print(t[1])
print(t[2])
print(t["a"])
print(t["b"])
`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	chunk := p.ParseChunk()

	r := NewRuntime()
	output := r.ExecuteChunkWithOutput(chunk)
	fmt.Printf("Output:\n%s", output)
	if output != "2\n4\n2\n2\n" {
		t.Errorf("Expected output '2\\n4\\n2\\n2\\n', got %q", output)
	}
}
