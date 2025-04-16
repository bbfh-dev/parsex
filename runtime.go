package parsex

import (
	"fmt"

	"github.com/bbfh-dev/parsex/internal"
)

// [runtimeType] is created from [Program] and it's what actually handles everything
type runtimeType struct {
	// These are set from [Program]
	data any
	exec *internal.ContextExecutable
	name string
	desc string
	// (Optional) Only needed for the primary executable program.
	// SemVer is adviced. Use [Program.SetVersion(...)] to edit
	version string
	// (Optional) Positional arguments for the program
	posArgs []string
	// (Optional) Other subcommands
	branches *internal.OrderedMap[*runtimeType]

	// --- Internal
	genOptions *internal.OrderedMap[internal.Option]
}

func newRuntime(program *Program) *runtimeType {
	return &runtimeType{
		data:       program.Data,
		exec:       internal.NewContextExecutable(program.Exec),
		name:       program.Name,
		desc:       program.Desc,
		version:    "",
		posArgs:    []string{},
		branches:   internal.NewOrderedMap[*runtimeType](),
		genOptions: internal.NewOrderedMap[internal.Option](),
	}
}

// Sets the version of the program (SemVer is adviced). Enables built-in `--version` flag.
//
// Designed to be used in main programs, not branches/subcommands.
func (runtime *runtimeType) SetVersion(version string) *runtimeType {
	runtime.version = version
	return runtime
}

// Specifies program positional arguments.
//
// - Include `?` in the argument if it's optional;
//
// - Suffix with `...` in the argument if its variadic;
//
// Prints the arguments in --help menu as provided.
func (runtime *runtimeType) SetPosArgs(args ...string) *runtimeType {
	runtime.posArgs = args
	return runtime
}

// Registers a branch/subcommand, which is just another [Program] object.
//
// Prints the command in --help menu.
func (runtime *runtimeType) RegisterCommand(command *runtimeType) *runtimeType {
	runtime.branches.Add(command.name, command)
	return runtime
}

// Processes options, validates them and then runs the [.Exec] function
func (runtime *runtimeType) Run(inArgs []string) error {
	if err := runtime.preprocess(); err != nil {
		return fmt.Errorf("%s (internal/preprocessor): %w", runtime.name, err)
	}

	if runtime.exec.Function == nil {
		return fmt.Errorf("%s (internal): runtime.Exec function is nil", runtime.name)
	}
	if err := runtime.exec.Run(); err != nil {
		return fmt.Errorf("%s: %w", runtime.name, err)
	}
	return nil
}
