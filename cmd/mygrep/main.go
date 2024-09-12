package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

const (
	digits          = "0123456789"
	letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanumeric    = digits + letters + "_"
	patternDigit    = "\\d"
	patternWordChar = "\\w"
)

// main is the entry point of the program.
func main() {
	// Validate arguments
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	// Read input line from stdin
	line, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	// Match line against pattern
	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	// Exit with code 1 if no match, otherwise default to 0 (success)
	if !ok {
		os.Exit(1)
	} else {
		os.Exit(0)
	}

}

// matchLine checks if the given line matches the pattern.
func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var matched bool

	if bytes.ContainsRune([]byte(pattern), '[') && bytes.ContainsRune([]byte(pattern), ']') { // positive character group
		matched = bytes.ContainsAny(line, pattern)
	} else {
		switch pattern {
		case patternDigit: // contains any digit
			matched = bytes.ContainsAny(line, digits)
		case patternWordChar: // contains any alphanumeric (a-z, A-Z, 0-9, _)
			matched = bytes.ContainsAny(line, alphanumeric)
		default:
			matched = bytes.Contains(line, []byte(pattern))
		}
	}

	return matched, nil
}
