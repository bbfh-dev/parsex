package parsex_test

import (
	"testing"

	"github.com/bbfh-dev/parsex"
	"gotest.tools/assert"
)

func TestErrDataIsPointer(test *testing.T) {
	program := parsex.Program{
		Data: 123,
		Name: "",
		Desc: "",
		Exec: nil,
	}.Runtime()
	err := program.Run([]string{})
	switch err := err.(type) {
	case parsex.ErrProgramData:
		assert.DeepEqual(test, err.ErrKind, parsex.ErrKindMustbePointer)
	default:
		test.Fatal("error must be ErrProgramData{}")
	}
}

func TestErrDataIsStruct(test *testing.T) {
	data := 123
	program := parsex.Program{
		Data: &data,
		Name: "",
		Desc: "",
		Exec: nil,
	}.Runtime()
	err := program.Run([]string{})
	switch err := err.(type) {
	case parsex.ErrProgramData:
		assert.DeepEqual(test, err.ErrKind, parsex.ErrKindPointToStruct)
	default:
		test.Fatal("error must be ErrProgramData{}")
	}
}

// TODO: Test all other errors
