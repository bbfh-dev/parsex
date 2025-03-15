package parsex_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bbfh-dev/parsex/parsex"
	"gotest.tools/assert"
)

var CLIErr = errors.New("This should have never been called!")

func testCLI(t *testing.T) *parsex.CLI {
	return parsex.New(
		"example",
		"An example description for the command",
		func(in parsex.Input) error {
			return CLIErr
		},
	).AddOptions(
		parsex.FlagOption{
			Name:     "all",
			Keywords: []string{"all", "all_this_works_too", "q"},
			Desc:     "Something allllllll",
		},
		parsex.FlagOption{
			Name:     "verbose",
			Keywords: []string{"verbose", "v"},
			Desc:     "Print verbose debug information",
		},
		parsex.FlagOption{
			Name:     "debug",
			Keywords: []string{"debug", "D"},
			Desc:     "Debug this",
		},
		parsex.ParamOption{
			Name:     "input",
			Keywords: []string{"input", "i"},
			Desc:     "Input file",
			Check:    parsex.ValidPath,
			Optional: true,
		},
		parsex.ParamOption{
			Name:     "length",
			Keywords: []string{"length"},
			Desc:     "Some length",
			Check:    parsex.ValidInt,
			Optional: true,
		},
		parsex.ParamOption{
			Name:     "json",
			Keywords: []string{"json"},
			Desc:     "Valid JSON",
			Check:    parsex.ValidJSON,
			Optional: true,
		},
		parsex.New("subcommand", "This is a subcommand", func(in parsex.Input) error {
			assert.Equal(t, in.Int("length"), 15)
			assert.Equal(t, in.Has("verbose"), true)
			assert.Equal(t, in.Has("debug"), true)
			assert.Equal(t, in.Has("abc"), true)
			assert.DeepEqual(t, in.Args(), []string{"-just an argument"})
			return nil
		}).AddArguments().AddOptions(
			parsex.FlagOption{
				Name:     "abc",
				Keywords: []string{"abc"},
				Desc:     "This is an example",
			},
		).Build(),
	).AddArguments("?output").SetVersion("v1.0.0-beta").Build()
}

func TestCLI(t *testing.T) {
	testFileCLI := testCLI(t)

	assert.NilError(t, testFileCLI.FromSlice([]string{
		"--length",
		"15",
		"-vD",
		"subcommand",
		"-abc",
		"--",
		"-just an argument",
	}).Run())

	assert.Assert(t, testFileCLI.FromSlice([]string{
		"--lenngt",
	}).Run() != nil)

	assert.Assert(t, testFileCLI.FromSlice([]string{
		"--length helloworld",
	}).Run() != nil)

	err := testFileCLI.FromSlice([]string{
		"-vD",
	}).Run()
	if err != CLIErr {
		t.Fatal(err)
	}

	err = testFileCLI.FromSlice([]string{
		"-vDa",
	}).Run()
	if err == nil {
		t.Fail()
	}

	err = testFileCLI.FromSlice([]string{
		"subcommand",
		"--a",
	}).Run()
	if err == nil {
		t.Fail()
	}
}

func TestCLIArr(t *testing.T) {
	var testFileCLI = parsex.New(
		"example",
		"An example description for the command",
		func(in parsex.Input) error {
			assert.DeepEqual(t, in.ListOfBools("bools"), []bool{true, false, false, true, true})
			assert.DeepEqual(t, in.ListOfStrings("strings"), []string{"hello"})
			assert.DeepEqual(t, in.ListOfStrings("json"), []string{"{}", `{"hello":123}`})
			assert.DeepEqual(t, in.ListOfFloats("floats"), []float64{12.5, 51})
			return nil
		},
	).AddOptions(
		parsex.ParamOption{
			Name:     "bools",
			Keywords: []string{"bools", "b"},
			Desc:     "...",
			Check:    parsex.ValidList(parsex.ValidBool),
			Optional: true,
		},
		parsex.ParamOption{
			Name:     "strings",
			Keywords: []string{"strings", "s"},
			Desc:     "...",
			Check:    parsex.ValidList(parsex.ValidString),
			Optional: true,
		},
		parsex.ParamOption{
			Name:     "json",
			Keywords: []string{"json"},
			Desc:     "...",
			Check:    parsex.ValidList(parsex.ValidJSON),
			Optional: true,
		},
		parsex.ParamOption{
			Name:     "floats",
			Keywords: []string{"f"},
			Desc:     "...",
			Check:    parsex.ValidList(parsex.ValidFloat),
			Optional: true,
		},
	).SetVersion("v1.0.0-beta").Build()

	assert.NilError(t, testFileCLI.FromSlice([]string{
		"-b",
		"true,false,0,true,true",
		"-s",
		"hello",
		`--json={},{"hello":123}`,
		"--f",
		"12.5,51",
	}).Run())
}

func TestArguments(t *testing.T) {
	var testFileCLI = parsex.New(
		"example",
		"An example description for the command",
		func(in parsex.Input) error {
			return nil
		},
	).AddArguments("filepath", "example").SetVersion("v1.0.0-beta").Build()

	assert.NilError(t, testFileCLI.FromSlice([]string{"/tmp/file", "123"}).Run())

	err := testFileCLI.FromSlice([]string{"/tmp/file"}).Run()
	if err == nil {
		t.Fail()
	}
}

func TestOptionalArguments(t *testing.T) {
	var testFileCLI = parsex.New(
		"example",
		"An example description for the command",
		func(in parsex.Input) error {
			return nil
		},
	).AddArguments("?ignore", "filepath", "?example").SetVersion("v1.0.0-beta").Build()

	assert.NilError(t, testFileCLI.FromSlice([]string{"ignored", "/tmp/file", "123"}).Run())
	assert.NilError(t, testFileCLI.FromSlice([]string{"ignored", "/tmp/file"}).Run())
	assert.NilError(t, testFileCLI.FromSlice([]string{"/tmp/file"}).Run())
}

func TestHelp(t *testing.T) {
	testFileCLI := testCLI(t)
	fmt.Println("--- Displaying help")
	assert.NilError(t, testFileCLI.FromSlice([]string{"--help"}).Run())
	fmt.Println("--- Displaying subcommand help")
	assert.NilError(t, testFileCLI.FromSlice([]string{"subcommand", "--help"}).Run())
}
