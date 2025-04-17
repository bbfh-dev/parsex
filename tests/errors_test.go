package parsex_test

import (
	"reflect"
	"testing"

	"github.com/bbfh-dev/parsex"
	"github.com/bbfh-dev/parsex/internal/cerr"
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
	case cerr.DuringPreprocessing:
		assert.Assert(test, reflect.TypeOf(err.Err) == reflect.TypeOf(cerr.DataMustBePointer{}))
	default:
		test.Fatal("error must be DuringPreprocessing{}")
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
	case cerr.DuringPreprocessing:
		assert.Assert(test, reflect.TypeOf(err.Err) == reflect.TypeOf(cerr.DataMustPointToStruct{}))
	default:
		test.Fatal("error must be DuringPreprocessing{}")
	}
}
