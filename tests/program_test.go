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

var MainProgram = parsex.New(&Options, nil, "example", "This is an example program").
	SetVersion("1.0.0-dev").RequireArgs("arg1", "arg2", "argN...").RegisterCommand(Subcommand)

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
