package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
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
	matched, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	fmt.Printf("[DEBUG] result: %v\n", matched)

	// Exit with code 1 if no match, otherwise default to 0 (success)
	if !matched {
		os.Exit(1)
	}

	// Exit with code 0 - success
}

// matchLine checks if the given line matches the pattern.
func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}
	fmt.Printf("[DEBUG] line '%s', pattern: '%s'\n", line, pattern)

	var matched bool

	switch {
	case isRangePattern(pattern):
		matched = matchRangePattern(line, pattern)
	case pattern == patternDigit:
		fmt.Println("[DEBUG] pattern is 'Digit'")
		matched = bytes.ContainsAny(line, digits)
	case pattern == patternWordChar:
		fmt.Println("[DEBUG] pattern is 'WordChar'")
		matched = bytes.ContainsAny(line, alphanumeric)
	default:
		matched = bytes.Contains(line, []byte(pattern))
	}

	return matched, nil
}

// isRangePattern checks if a pattern is a range pattern like [abc] or [^abc].
func isRangePattern(pattern string) bool {
	return strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]")
}

// matchRangePattern matches line against a range pattern, supporting both inclusive and exclusive ranges.
func matchRangePattern(line []byte, pattern string) bool {
	var ok bool
	inside := pattern[1 : len(pattern)-1]

	if strings.HasPrefix(inside, "^") {
		// [^abc] Not Range (a or b or c)
		fmt.Println("[DEBUG] pattern is 'Not Range'")
		inside = inside[1:]
		ok = !bytes.ContainsAny(line, inside)
	} else {
		// [abc] Range (a or b or c)
		fmt.Println("[DEBUG] pattern is 'Range'")
		ok = bytes.ContainsAny(line, inside)
	}

	return ok
}
