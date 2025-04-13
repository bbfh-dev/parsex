package parsex

import (
	"fmt"
)

type Executable func(args []string) error

type Program struct {
	// Pointer to a [var ... struct{}] containing all program options.
	//
	// All input flags will be stored in the struct,
	// --help menu will be automatically generated.
	//
	// Use the following field tags:
	//
	// `alt:"<single letter alternative use>" desc:"<description of the option>`
	Data any
	// The function to be called as program logic.
	// Checkout [parsex.Batch()] & [parsex.BatchSeq()] helper functions
	Executable Executable
	// The name of the executable / command
	Name string
	// Will be displayed in the --help menu
	Desc string
	// (Optional) Only needed for the primary executable program. SemVer is adviced. Use [Program.SetVersion(...)] to edit
	version string
	// (Optional) Positional arguments for the program
	posArgs []string
	// (Optional) Other subcommands
	branches []*Program
}

// NOTE: Check [parsex.Program] docs
func New(data any, exec Executable, name string, desc string) *Program {
	return &Program{
		Data:       data,
		Executable: exec,
		Name:       name,
		Desc:       desc,
		version:    "",
		posArgs:    []string{},
		branches:   []*Program{},
	}
}

// Sets the version of the program. Enables built-in `--version` flag.
//
// Designed to be used in main programs, not branches/subcommands.
func (program *Program) SetVersion(version string) *Program {
	program.version = version
	return program
}

// Specifies program positional arguments.
//
// - Include `?` in the argument if it's optional;
//
// - Suffix with `...` in the argument if its variadic;
//
// Prints the arguments in --help menu as provided.
func (program *Program) RequireArgs(args ...string) *Program {
	program.posArgs = args
	return program
}

// Registers a branch/subcommand, which is just another [Program] object.
//
// Prints the command in --help menu.
func (program *Program) RegisterCommand(cmd *Program) *Program {
	program.branches = append(program.branches, cmd)
	return program
}

// Processes options, validates them and then runs the [.Executable] callback
func (program *Program) Run(vArgs ...string) error {
	parser := newParser()

	if err := parser.LoadOptionsFrom(program.Data); err != nil {
		return fmt.Errorf("%s (internal): parsex.Program.Run(...): %w", program.Name, err)
	}

	if err := parser.Parse(vArgs); err != nil {
		return fmt.Errorf("%s (user input): %w", program.Name, err)
	}

	if err := program.Executable(parser.Args); err != nil {
		return fmt.Errorf("%s: %w", program.Name, err)
	}

	return nil
}
