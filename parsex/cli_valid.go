package parsex

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Ensures that the argument is valid.
// The returning string is an optionally transformed input string
type Validator func(string) (any, bool)

// Require a string argument. No modifications
func ValidString(in string) (any, bool) {
	return in, true
}

// Require a valid path. Transforms into an absolute path
func ValidPath(in string) (any, bool) {
	path, _ := os.UserHomeDir()

	path, err := filepath.Abs(strings.Replace(in, "~", path, 1))
	if err != nil {
		return in, false
	}

	return path, true
}

// Require an integer argument. Converts into an int
func ValidInt(in string) (any, bool) {
	value, err := strconv.Atoi(in)
	return value, err == nil
}

// Require an integer argument. Converts into the int of specified base and size.
//
// It interprets the flag in the given base (0, 2 to 36) and bit size (0 to 64).
func ValidUint(base int, bitSize int) Validator {
	return func(in string) (any, bool) {
		value, err := strconv.ParseUint(in, base, bitSize)
		return value, err == nil
	}
}
