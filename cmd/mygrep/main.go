package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
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

	return matchPattern(line, pattern)
}

// matchPattern attempts to match the pattern against the line
func matchPattern(line []byte, pattern string) (bool, error) {
	// Split the pattern into sub-patterns
	patterns := splitPattern(pattern)
	fmt.Printf("[DEBUG] split patterns: %v\n", patterns)

	lineStr := string(line)

	// Iterate over the line to check for matches starting from each position
	for i := range lineStr {
		fmt.Printf("[DEBUG] checking from position %d\n", i)
		matched, _ := recursiveMatch(lineStr, patterns, i, 0)
		if matched {
			fmt.Printf("[DEBUG] match found starting at position %d\n", i)
			return true, nil
		}
	}
	fmt.Println("[DEBUG] no match found in line")
	return false, nil
}

// splitPattern splits the pattern into recognizable sub-patterns
func splitPattern(pattern string) []string {
	var patterns []string
	for i := 0; i < len(pattern); {
		switch {
		// Handle escaped characters
		case pattern[i] == '\\':
			if i+1 < len(pattern) {
				patterns = append(patterns, pattern[i:i+2])
				i += 2
			}
		// Handle range patterns
		case isRangePattern(pattern[i:]):
			end := strings.Index(pattern[i:], "]")
			if end != -1 {
				patterns = append(patterns, pattern[i:i+end+1])
				i += end + 1
			} else {
				patterns = append(patterns, string(pattern[i]))
				i++
			}
		// Handle literals
		default:
			patterns = append(patterns, string(pattern[i]))
			i++
		}
	}
	return patterns
}

// recursiveMatch checks the string against the pattern starting from the given positions
func recursiveMatch(line string, patterns []string, linePos, patPos int) (bool, error) {
	// If we've processed all sub-patterns, return true
	if patPos == len(patterns) {
		return true, nil
	}
	// If we've reached the end of the line, return false
	if linePos >= len(line) {
		return false, nil
	}

	// Get the current sub-pattern
	pat := patterns[patPos]
	fmt.Printf("[DEBUG] matching pattern '%s' at line position %d\n", pat, linePos)

	switch {
	// Handle '\d' pattern
	case pat == patternDigit:
		if unicode.IsDigit(rune(line[linePos])) {
			return recursiveMatch(line, patterns, linePos+1, patPos+1)
		}
	// Handle '\w' pattern
	case pat == patternWordChar:
		if unicode.IsLetter(rune(line[linePos])) || rune(line[linePos]) == '_' {
			return recursiveMatch(line, patterns, linePos+1, patPos+1)
		}
	// Handle range patterns
	case isRangePattern(pat):
		if matchRangePattern([]byte{line[linePos]}, pat) {
			return recursiveMatch(line, patterns, linePos+1, patPos+1)
		}
	// Handle literal patterns
	default:
		if strings.HasPrefix(line[linePos:], pat) {
			return recursiveMatch(line, patterns, linePos+len(pat), patPos+1)
		}
	}

	return false, nil
}

// isRangePattern checks if a pattern is a range pattern like [abc] or [^abc].
func isRangePattern(pattern string) bool {
	return strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]")
}

// matchRangePattern matches line against a range pattern, supporting both inclusive and exclusive ranges.
func matchRangePattern(line []byte, pattern string) bool {
	inside := pattern[1 : len(pattern)-1]
	negate := false

	// Check if the pattern is a negated range (starts with [^)
	if strings.HasPrefix(inside, "^") {
		// [^abc] Not Range (a or b or c)
		fmt.Println("[DEBUG] pattern is 'Not Range'")
		// Remove the '^' character
		inside = inside[1:]
		negate = true
	} else {
		// [abc] Range (a or b or c)
		fmt.Println("[DEBUG] pattern is 'Range'")
	}

	// Negated range pattern: return true if line does not contain any of the characters inside
	if negate {
		result := !bytes.ContainsAny(line, inside)
		fmt.Printf("[DEBUG] negate result: %v\n", result)
		return result
	}

	// Inclusive range pattern: return true if line contains any of the characters inside
	result := bytes.ContainsAny(line, inside)
	fmt.Printf("[DEBUG] range result: %v\n", result)
	return result
}
