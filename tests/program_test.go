package parsex_test

import (
	"fmt"
	"testing"

	"github.com/bbfh-dev/parsex"
	"gotest.tools/assert"
)

var SubOptions struct{}

var Options struct {
	Verbose    bool   `alt:"v" desc:"Prints verbose debug information"`
	Input      string `alt:"i" desc:"Provide a filename"`
	SomeLength int    `desc:"I have no idea, some integer whatever"`
}

var Subcommand = parsex.New(&SubOptions, nil, "subcommand", "A subcommand with its own logic!")

var MainProgram = parsex.New(&Options, func(args []string) error {
	return nil
}, "example", "This is an example program").
	SetVersion("1.0.0-dev").
	RequireArgs("arg1", "arg2", "argN...").
	RegisterCommand(Subcommand)

// Human-checked test. Simply to check what the library generates.
func TestProgramVersion(test *testing.T) {
	fmt.Println("```")
	assert.NilError(test, MainProgram.Run("--version"))
	fmt.Println("```")
}

// Human-checked test. Simply to check what the library generates.
func TestProgramHelp(test *testing.T) {
	fmt.Println("```")
	assert.NilError(test, MainProgram.Run("--help"))
	fmt.Println("```")
}

func TestProgramSubcommand(test *testing.T) {
	testing := parsex.New(&SubOptions, func(args []string) error {
		assert.DeepEqual(test, args, []string{"arg1"})
		test.SkipNow()
		return nil
	}, "testing", "")
	assert.NilError(test, MainProgram.RegisterCommand(testing).Run(
		"--verbose",
		"testing",
		"arg1",
	))

	test.FailNow()
}

func TestProgramFlag(test *testing.T) {
	Options.Verbose = false
	assert.NilError(test, MainProgram.Run("--verbose"))
	assert.DeepEqual(test, Options.Verbose, true)
}

func TestProgramOption(test *testing.T) {
	Options.SomeLength = 0
	assert.NilError(test, MainProgram.Run("--some-length=15"))
	assert.DeepEqual(test, Options.SomeLength, 15)

	Options.SomeLength = 0
	assert.NilError(test, MainProgram.Run("--some-length", "69"))
	assert.DeepEqual(test, Options.SomeLength, 69)

	Options.Input = ""
	Options.SomeLength = 0
	assert.NilError(test, MainProgram.Run("-some-length", "420", "-i", "hello world"))
	assert.DeepEqual(test, Options.SomeLength, 420)
	assert.DeepEqual(test, Options.Input, "hello world")
}
