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
		for currentLine < tok.Line {
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
