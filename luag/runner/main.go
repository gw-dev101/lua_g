package main

import (
	"fmt"
	"os"

	"luag/lexer"
	"luag/parser"
	"luag/runtime"
)

// Execute runs the full pipeline for a given input string and returns the output.
func Execute(input string) (string, error) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)

	// Optional: Handle or check for parser errors here if your parser supports it
	chunk := p.ParseChunk()

	r := runtime.NewRuntime()
	output := r.ExecuteChunkWithOutput(chunk)

	return output, nil
}

// main function that take a file path as an argument and executes the Lua code in that file.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: luag <file_path>")
		return
	}

	filePath := os.Args[1]
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	output, err := Execute(string(content))
	if err != nil {
		fmt.Printf("Error executing Lua code: %v\n", err)
		return
	}

	fmt.Println(output)
}
