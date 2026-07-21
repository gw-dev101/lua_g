package lexer

import (
	"fmt"
	"testing"
)

func TestColorizedLexer(t *testing.T) {
	input := `local a = 10
if a > 5 then
  print("greater")
else
  print("lesser")
end

if a == 100 then
  print("no")
end`

	// ANSI Escape Codes for Terminal Colors
	const (
		Reset   = "\033[0m"
		Red     = "\033[31m"
		Green   = "\033[32m"
		Yellow  = "\033[33m"
		Blue    = "\033[34m"
		Magenta = "\033[35m"
		Cyan    = "\033[36m"
	)

	l := NewLexer(input)
	currentLine := 1

	fmt.Println("\n--- Colorized Token Output ---")

	for {
		tok := l.NextToken()
		if tok.Type == TokenTypeEOF {
			break
		}

		// Handle visual line breaks to match original source layout
		for currentLine < tok.Span.Start.Line {
			fmt.Println()
			currentLine++
		}

		// Choose color based on token type
		var color string
		switch tok.Type {
		case TokenTypeKeyword:
			color = Cyan
		case TokenTypeIdentifier:
			color = Blue
		case TokenTypeNumber:
			color = Yellow
		case TokenTypeString:
			color = Green // Keeps quotes inside the string representation if preferred
			// Note: If you want to show the literal quotes, wrap it: fmt.Sprintf(`"%s"`, tok.Literal)
		case TokenTypeOperator:
			color = Red
		case TokenTypePunctuation:
			color = Magenta
		default:
			color = Reset
		}

		// Print the token literal colorized (add space or formatting if preferred)
		// For strings, manually re-wrapping quotes helps visually identify them
		if tok.Type == TokenTypeString {
			fmt.Printf("%s\"%s\"%s ", color, tok.Literal, Reset)
		} else {
			fmt.Printf("%s%s%s ", color, tok.Literal, Reset)
		}
	}
	fmt.Println("\n------------------------------")
}
func TestLexer_NextToken(t *testing.T) {
	// 1. Define what we want to assert for each token
	type expectedToken struct {
		Type    TokenType
		Literal string
		// You can also add expected Line and Column here later
		// if you want to test position tracking!
	}

	// 2. Set up our test cases
	tests := []struct {
		name     string
		input    string
		expected []expectedToken
	}{
		{
			name:  "Basic assignment",
			input: `local a = 10`,
			expected: []expectedToken{
				{TokenTypeKeyword, "local"},
				{TokenTypeIdentifier, "a"},
				{TokenTypeOperator, "="},
				{TokenTypeNumber, "10"},
				{TokenTypeEOF, ""},
			},
		},
		{
			name: "If statement with strings",
			input: `if a > 5 then
  print("greater")
end`,
			expected: []expectedToken{
				{TokenTypeKeyword, "if"},
				{TokenTypeIdentifier, "a"},
				{TokenTypeOperator, ">"},
				{TokenTypeNumber, "5"},
				{TokenTypeKeyword, "then"},
				{TokenTypeIdentifier, "print"},
				{TokenTypePunctuation, "("},
				{TokenTypeString, "greater"},
				{TokenTypePunctuation, ")"},
				{TokenTypeKeyword, "end"},
				{TokenTypeEOF, ""},
			},
		},
		{
			name:  "Operators and multi-character operators",
			input: `a == 100 ~= 50`,
			expected: []expectedToken{
				{TokenTypeIdentifier, "a"},
				{TokenTypeOperator, "=="},
				{TokenTypeNumber, "100"},
				{TokenTypeOperator, "~="},
				{TokenTypeNumber, "50"},
				{TokenTypeEOF, ""},
			},
		},
	}

	// 3. Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)

			for i, expTok := range tt.expected {
				tok := l.NextToken()

				// Assert the Token Type
				if tok.Type != expTok.Type {
					t.Fatalf("[%s] token %d type wrong. expected=%v, got=%v (literal: %q)",
						tt.name, i, expTok.Type, tok.Type, tok.Literal)
				}

				// Assert the Token Literal
				if tok.Literal != expTok.Literal {
					t.Fatalf("[%s] token %d literal wrong. expected=%q, got=%q",
						tt.name, i, expTok.Literal, tok.Literal)
				}
			}

			// Optional: Ensure the lexer doesn't produce EXTRA tokens
			// after we expect it to hit EOF.
			extraTok := l.NextToken()
			if extraTok.Type != TokenTypeEOF {
				t.Errorf("[%s] expected EOF, but got extra token: %v (%q)",
					tt.name, extraTok.Type, extraTok.Literal)
			}
		})
	}

	// 3. Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)

			for i, expTok := range tt.expected {
				tok := l.NextToken()

				// Assert the Token Type
				if tok.Type != expTok.Type {
					t.Fatalf("[%s] token %d type wrong. expected=%v, got=%v (literal: %q)",
						tt.name, i, expTok.Type, tok.Type, tok.Literal)
				}

				// Assert the Token Literal
				if tok.Literal != expTok.Literal {
					t.Fatalf("[%s] token %d literal wrong. expected=%q, got=%q",
						tt.name, i, expTok.Literal, tok.Literal)
				}
			}

			// Optional: Ensure the lexer doesn't produce EXTRA tokens
			// after we expect it to hit EOF.
			extraTok := l.NextToken()
			if extraTok.Type != TokenTypeEOF {
				t.Errorf("[%s] expected EOF, but got extra token: %v (%q)",
					tt.name, extraTok.Type, extraTok.Literal)
			}
		})
	}
}
