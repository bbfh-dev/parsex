package parsex

import (
	"fmt"
	"os"
	"strings"
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
	branchKeys []string
	branches   map[string]*Program
	// generated
	execArgs   []string
	optionKeys []string          // Used to keep the options sorted
	options    map[string]option // Map is used for more performant access
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
		branchKeys: []string{},
		branches:   map[string]*Program{},
		execArgs:   []string{},
		optionKeys: []string{},
		options: map[string]option{
			"help":    helpOption,
			"version": versionOption,
		},
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
	program.branchKeys = append(program.branchKeys, cmd.Name)
	program.branches[cmd.Name] = cmd
	return program
}

// Processes options, validates them and then runs the [.Executable] callback
func (program *Program) Run(vArgs ...string) error {
	if err := program.loadOptions(); err != nil {
		return fmt.Errorf("%s (internal): parsex.Program.Run(...): %w", program.Name, err)
	}

	for i := 0; i < len(vArgs); i++ {
		if !strings.HasPrefix(vArgs[i], "-") {
			// If it's a subcommand simply branch out
			branch, ok := program.branches[vArgs[i]]
			if ok {
				return branch.Run(vArgs[i+1:]...)
			}

			// Default to a argument
			program.execArgs = append(program.execArgs, vArgs[i])
			continue
		}
		arg := strings.TrimLeft(vArgs[i], "-")

		switch arg {

		case "help":
			program.printHelp(os.Stdout)
			return nil

		case "version":
			program.printVersion(os.Stdout)
			return nil
		}
	}

	if program.Executable == nil {
		return fmt.Errorf(
			"%s (internal): parsex.New(...): Executable function is nil",
			program.Name,
		)
	}

	if err := program.Executable(program.execArgs); err != nil {
		return fmt.Errorf("%s: %w", program.Name, err)
	}

	return nil
}
