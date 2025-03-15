package parsex_test

// These tests were written by AI but they seem good ™

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/bbfh-dev/parsex/parsex"
)

// TestValidString simply checks that the string is returned as-is.
func TestValidString(t *testing.T) {
	input := "hello world"
	output, err := parsex.ValidString(input)
	if err != nil {
		t.Errorf("ValidString(%q) returned unexpected error: %v", input, err)
	}
	if output != input {
		t.Errorf("ValidString(%q) = %v; want %v", input, output, input)
	}
}

// TestValidJSON tests both valid and invalid JSON inputs.
func TestValidJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"ValidObject", "{}", nil},
		{"ValidArray", "[1,2,3]", nil},
		{"ValidNull", "null", nil},
		{"InvalidJSON", "not json", parsex.ErrInvalidJSON},
		{"InvalidJSON2", "{", parsex.ErrInvalidJSON},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			output, err := parsex.ValidJSON(testCase.input)
			if testCase.wantErr != nil {
				if err == nil || !errors.Is(err, testCase.wantErr) {
					t.Errorf(
						"ValidJSON(%q) error = %v; want %v",
						testCase.input,
						err,
						testCase.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Errorf("ValidJSON(%q) returned unexpected error: %v", testCase.input, err)
				}
				if output != testCase.input {
					t.Errorf("ValidJSON(%q) = %v; want %v", testCase.input, output, testCase.input)
				}
			}
		})
	}
}

// TestValidPath verifies that an absolute path is returned.
// Note: It's hard to simulate an error from filepath.Abs in a normal test.
func TestValidPath(t *testing.T) {
	input := "testdata"
	output, err := parsex.ValidPath(input)
	if err != nil {
		t.Errorf("ValidPath(%q) returned error: %v", input, err)
	}
	abs, _ := filepath.Abs(input)
	if output != abs {
		t.Errorf("ValidPath(%q) = %v; want %v", input, output, abs)
	}
}

// TestValidExistingPath tests both an existing path and a non-existent one.
func TestValidExistingPath(t *testing.T) {
	// Existing path: use os.TempDir() which should exist.
	existing := os.TempDir()
	output, err := parsex.ValidExistingPath(existing)
	if err != nil {
		t.Errorf("ValidExistingPath(%q) unexpected error: %v", existing, err)
	}
	abs, _ := filepath.Abs(existing)
	if output != abs {
		t.Errorf("ValidExistingPath(%q) = %v; want %v", existing, output, abs)
	}

	// Non-existent path: create a path that is unlikely to exist.
	nonExist := filepath.Join(os.TempDir(), "non_existent_path_12345")
	_, err = parsex.ValidExistingPath(nonExist)
	if err == nil {
		t.Errorf("ValidExistingPath(%q) expected error for non-existent path", nonExist)
	}
}

// TestValidInt checks valid integer strings and cases that should fail.
func TestValidInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"123", 123, false},
		{"-456", -456, false},
		{"0", 0, false},
		{"abc", 0, true},
		{"12.34", 0, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.input, func(t *testing.T) {
			output, err := parsex.ValidInt(testCase.input)
			if testCase.wantErr {
				if err == nil {
					t.Errorf("ValidInt(%q) expected error, got nil", testCase.input)
				}
			} else {
				if err != nil {
					t.Errorf("ValidInt(%q) unexpected error: %v", testCase.input, err)
				}
				// Since ValidInt returns an int, we can do a type assertion.
				val, ok := output.(int)
				if !ok {
					t.Errorf("ValidInt(%q) returned type %T; want int", testCase.input, output)
				}
				if val != testCase.expected {
					t.Errorf("ValidInt(%q) = %v; want %v", testCase.input, val, testCase.expected)
				}
			}
		})
	}
}

// TestValidFloat checks both valid and invalid float conversions.
func TestValidFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		wantErr  bool
	}{
		{"123.45", 123.45, false},
		{"-0.001", -0.001, false},
		{"0", 0.0, false},
		{"abc", 0.0, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.input, func(t *testing.T) {
			output, err := parsex.ValidFloat(testCase.input)
			if testCase.wantErr {
				if err == nil {
					t.Errorf("ValidFloat(%q) expected error, got nil", testCase.input)
				}
			} else {
				if err != nil {
					t.Errorf("ValidFloat(%q) unexpected error: %v", testCase.input, err)
				}
				// Since floating point comparisons can be tricky, compare using strconv.FormatFloat.
				outStr := strconv.FormatFloat(output.(float64), 'f', 5, 64)
				expStr := strconv.FormatFloat(testCase.expected, 'f', 5, 64)
				if outStr != expStr {
					t.Errorf("ValidFloat(%q) = %v; want %v", testCase.input, output, testCase.expected)
				}
			}
		})
	}
}

// TestValidBool checks all valid boolean representations and a few invalid ones.
func TestValidBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		wantErr  bool
	}{
		{"true", true, false},
		{"1", true, false},
		{"yes", true, false},
		{"false", false, false},
		{"0", false, false},
		{"no", false, false},
		{"maybe", false, true},
		{"TRUE", false, true}, // case-sensitive test
		{"", false, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.input, func(t *testing.T) {
			output, err := parsex.ValidBool(testCase.input)
			if testCase.wantErr {
				if err == nil {
					t.Errorf("ValidBool(%q) expected error, got nil", testCase.input)
				}
			} else {
				if err != nil {
					t.Errorf("ValidBool(%q) unexpected error: %v", testCase.input, err)
				}
				if output != testCase.expected {
					t.Errorf("ValidBool(%q) = %v; want %v", testCase.input, output, testCase.expected)
				}
			}
		})
	}
}

// TestValidList uses the ValidList helper with the ValidInt validator.
// It tests both a list that entirely converts to ints and one with an invalid element.
func TestValidListWithIntValidator(t *testing.T) {
	intListValidator := parsex.ValidList(parsex.ValidInt)

	// Valid list input.
	validInput := "1,2,3"
	output, err := intListValidator(validInput)
	if err != nil {
		t.Errorf("ValidList with ValidInt failed for input %q: %v", validInput, err)
	}
	intSlice, ok := output.([]any)
	if !ok {
		t.Errorf("ValidList with ValidInt returned type %T; want []any", output)
	}
	expectedInts := []int{1, 2, 3}
	if len(intSlice) != len(expectedInts) {
		t.Errorf(
			"ValidList with ValidInt returned %v items; want %v",
			len(intSlice),
			len(expectedInts),
		)
	}
	for i, v := range intSlice {
		val, ok := v.(int)
		if !ok {
			t.Errorf("Element %d of ValidList output has type %T; want int", i, v)
		}
		if val != expectedInts[i] {
			t.Errorf("Element %d of ValidList output = %v; want %v", i, val, expectedInts[i])
		}
	}

	// Invalid list input: one element cannot be converted to int.
	invalidInput := "1,a,3"
	_, err = intListValidator(invalidInput)
	if err == nil {
		t.Errorf("ValidList with ValidInt expected error for input %q, got nil", invalidInput)
	}
}
