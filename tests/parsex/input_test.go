package parsex_test

import (
	"testing"

	"github.com/bbfh-dev/parsex/parsex"
)

func TestInput(test *testing.T) {
	in := parsex.Input{"test": "hello"}
	if !in.Has("test") {
		test.Fatal("test must be present")
	}

	if in["test"] != "hello" {
		test.Fatal("test must be equal to 'hello'")
	}
}
