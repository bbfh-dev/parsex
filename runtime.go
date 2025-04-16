package parsex

import (
	"fmt"
	"os"
	"strings"

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
	genOptions    *internal.OrderedMap[internal.Option]
	genOptionAlts map[string]string
}

func newRuntime(program *Program) *runtimeType {
	return &runtimeType{
		data:          program.Data,
		exec:          internal.NewContextExecutable(program.Exec),
		name:          program.Name,
		desc:          program.Desc,
		version:       "",
		posArgs:       []string{},
		branches:      internal.NewOrderedMap[*runtimeType](),
		genOptions:    internal.NewOrderedMap[internal.Option](),
		genOptionAlts: map[string]string{},
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

// Run processes options, validates them, and then executes the command.
func (runtime *runtimeType) Run(inputArgs []string) error {
	if err := runtime.preprocess(); err != nil {
		return fmt.Errorf("%s (internal/preprocessor): %w", runtime.name, err)
	}

	for i := 0; i < len(inputArgs); i++ {
		arg := inputArgs[i]
		if !strings.HasPrefix(arg, "-") {
			// Branch or positional argument.
			if branch, exists := runtime.branches.Get(arg); exists {
				return branch.Run(inputArgs[i+1:])
			}
			runtime.exec.AddArg(arg)
			continue
		}

		// Handle help and version shortcuts.
		switch arg {
		case "--help":
			runtime.printHelp(os.Stdout)
			return nil
		case "--version":
			runtime.PrintVersion(os.Stdout)
			return nil
		}

		var err error
		if strings.HasPrefix(arg, "--") {
			err = runtime.processLongOption(arg, &i, inputArgs)
		} else {
			err = runtime.processShortOption(arg, &i, inputArgs)
		}
		if err != nil {
			return err
		}
	}

	if runtime.exec.Function == nil {
		return fmt.Errorf("%s (internal): runtime.Exec function is nil", runtime.name)
	}
	if err := runtime.exec.Run(); err != nil {
		return fmt.Errorf("%s: %w", runtime.name, err)
	}
	return nil
}

// processLongOption handles options starting with "--".
func (runtime *runtimeType) processLongOption(arg string, i *int, inputArgs []string) error {
	// Remove "--" prefix.
	optionStr := arg[2:]
	var name, value string

	if parts := strings.SplitN(optionStr, "=", 2); len(parts) == 2 {
		name = parts[0]
		value = parts[1]
		return runtime.setOption(name, value)
	}

	name = optionStr
	option, exists := runtime.genOptions.Get(name)
	if !exists {
		return fmt.Errorf(
			"%s: unknown option %q. Refer to --help for usage information",
			runtime.name,
			arg,
		)
	}
	if option.IsFlag() {
		option.SetFlag()
		return nil
	}

	*i++
	if *i >= len(inputArgs) {
		return fmt.Errorf(
			"%s: option %q requires a value. Refer to --help for usage information",
			runtime.name,
			arg,
		)
	}
	value = inputArgs[*i]
	return runtime.setOption(name, value)
}

// processShortOption handles options starting with a single "-".
// It supports both single options and clusters.
func (runtime *runtimeType) processShortOption(arg string, i *int, inputArgs []string) error {
	// Remove "-" prefix.
	optionStr := arg[1:]

	if strings.Contains(optionStr, "=") {
		parts := strings.SplitN(optionStr, "=", 2)
		name, value := parts[0], parts[1]
		return runtime.setOption(name, value)
	}

	name := optionStr
	option, exists := runtime.genOptions.Get(name)
	if exists {
		if option.IsFlag() {
			option.SetFlag()
			return nil
		}
		*i++
		if *i >= len(inputArgs) {
			return fmt.Errorf(
				"%s: option %q requires a value. Refer to --help for usage information",
				runtime.name,
				arg,
			)
		}
		return runtime.setOption(name, inputArgs[*i])
	}

	// Each character in the cluster should map to a flag.
	for _, char := range optionStr {
		flagKey := string(char)
		mapped, exists := runtime.genOptionAlts[flagKey]
		if !exists {
			return fmt.Errorf(
				"%s: unknown option or cluster %q. Refer to --help for usage information",
				runtime.name,
				arg,
			)
		}
		option, exists := runtime.genOptions.Get(mapped)
		if !exists || !option.IsFlag() {
			return fmt.Errorf(
				"%s: cluster %q can only contain flags but %q is found. Refer to --help for usage information",
				runtime.name,
				arg,
				mapped,
			)
		}
		option.SetFlag()
	}

	return nil
}

// setOption retrieves the option by name and applies the value.
func (runtime *runtimeType) setOption(name, value string) error {
	option, exists := runtime.genOptions.Get(name)
	if !exists {
		return fmt.Errorf(
			"%s: unknown option %q. Refer to --help for usage information",
			runtime.name,
			"--"+name,
		)
	}
	if option.IsFlag() {
		option.SetFlag()
		return nil
	}
	if err := option.Set(value); err != nil {
		return fmt.Errorf("%s: setting option %q: %w", runtime.name, "--"+name, err)
	}
	return nil
}
