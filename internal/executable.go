package internal

// [Executable] is a function that can be called by [parsex.Program].
// Usually it's a function pointing to your own program/command logic.
type Executable func(args []string) error

// An [Executable] with context such as program arguments.
type ContextExecutable struct {
	Function Executable
	Args     []string
}

func NewContextExecutable(exec Executable) *ContextExecutable {
	return &ContextExecutable{
		Function: exec,
		Args:     []string{},
	}
}

func (executable *ContextExecutable) Clear() {
	executable.Args = []string{}
}

// Saves the argument in the executable context
func (executable *ContextExecutable) AddArg(arg string) *ContextExecutable {
	executable.Args = append(executable.Args, arg)
	return executable
}

// Runs the [Executable] with the stored arguments
func (executable *ContextExecutable) Run() error {
	return executable.Function(executable.Args)
}
