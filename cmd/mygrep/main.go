package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

const (
	digits  = "0123456789"
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special = "_"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}
	var ok bool

	switch pattern {
	case "\\d": // contains any digit
		ok = bytes.ContainsAny(line, digits)
	case "\\w": // contains any alphanumeric (a-z, A-Z, 0-9, _)
		ok = bytes.ContainsAny(line, digits+letters+special)
	default:
		ok = bytes.ContainsAny(line, pattern)
	}
	fmt.Printf("[DEBUG] pattern: '%s', match: '%v'\n", pattern, ok)

	return ok, nil
}
