package parsex_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/bbfh-dev/parsex/v2"
	"gotest.tools/assert"
)

var testOptions struct {
	Verbose    bool   `alt:"v" desc:"Print verbose debug information"`
	Debug      bool   `alt:"d" desc:"Run the program in DEBUG mode"`
	Input      string `desc:"Some input file"`
	SomeNumber int    `alt:"N" desc:"A valid integer" default:"69"`
}

var buildProgram = parsex.Program{
	Data: nil,
	Name: "build",
	Desc: "This is an example subcommand that builds something",
	Exec: func(args []string) error {
		log.Println("--- * Running test.BuildProgram with args", args)
		return nil
	},
}.Runtime().SetPosArgs("filename?")

var testProgram = parsex.Program{
	Data: &testOptions,
	Name: "example",
	Desc: "This is an example program",
	Exec: func(args []string) error {
		log.Println("--- * Running test.MainProgram with args", args)
		return nil
	},
}.Runtime().
	SetVersion("1.0.0-dev").
	SetPosArgs("arg1", "arg2", "argN...").
	RegisterCommand(buildProgram)

var ExpectedVersion = "example v1.0.0-dev\n"

var ExpectedHelp = ExpectedVersion + `
This is an example program

Usage:
    example [options] <arg1> <arg2> <argN...> 

Commands:
    build [options] <filename?> 

Options:
    --help
        # Print this help message
    --version
        # Print the program version
    --verbose, -v
        # Print verbose debug information
    --debug, -d
        # Run the program in DEBUG mode
    --input <string>
        # Some input file
    --some-number, -N <int> (default: 69)
        # A valid integer
`

func setup() {
	testOptions.Verbose = false
	testOptions.Debug = false
	testOptions.Input = ""
	testOptions.SomeNumber = 15
}

func TestProgramVersion(test *testing.T) {
	var buffer bytes.Buffer
	testProgram.PrintVersion(&buffer)
	assert.DeepEqual(test, buffer.String(), ExpectedVersion)
}

func TestProgramHelp(test *testing.T) {
	var buffer bytes.Buffer
	testProgram.SafePrintHelp(&buffer)
	assert.DeepEqual(test, buffer.String(), ExpectedHelp)
}

func TestProgramLongOptions(test *testing.T) {
	setup()
	assert.NilError(test, testProgram.Run([]string{
		"--verbose",
		"--debug",
		"--input=/tmp/filename",
		"arg1", "arg2", "arg3",
	}))
	assert.DeepEqual(test, testOptions.Verbose, true)
	assert.DeepEqual(test, testOptions.Debug, true)
	assert.DeepEqual(test, testOptions.Input, "/tmp/filename")
	assert.DeepEqual(test, testOptions.SomeNumber, 69)
}

func TestProgramClusterOptions(test *testing.T) {
	setup()
	assert.NilError(test, testProgram.Run([]string{
		"-vd",
		"-input=/tmp/filename",
		"-some-number", "15",
		"--", "arg1", "arg2", "--value",
	}))
	assert.DeepEqual(test, testOptions.Verbose, true)
	assert.DeepEqual(test, testOptions.Debug, true)
	assert.DeepEqual(test, testOptions.Input, "/tmp/filename")
	assert.DeepEqual(test, testOptions.SomeNumber, 15)
}

func TestProgramSubcommand(test *testing.T) {
	setup()
	assert.NilError(test, testProgram.Run([]string{
		"-vd",
		"build",
	}))
	assert.DeepEqual(test, testOptions.Verbose, true)
	assert.DeepEqual(test, testOptions.Debug, true)
}
