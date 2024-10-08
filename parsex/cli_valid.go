package parsex

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Ensures that the argument is valid.
// The returning string is an optionally transformed input string
type Validator func(string) (string, bool)

// Require a string argument. No modifications
func ValidString(in string) (string, bool) {
	return in, true
}

// Require a valid path. Transforms into an absolute path
func ValidPath(in string) (string, bool) {
	path, _ := os.UserHomeDir()

	path, err := filepath.Abs(strings.Replace(in, "~", path, 1))
	if err != nil {
		return in, false
	}

	return path, true
}

// Require an integer argument. No modifications
func ValidInt(in string) (string, bool) {
	_, err := strconv.Atoi(in)
	return in, err == nil
}
