package parsex

// Input is a map of flags/options that were provided to the program.
//
// All keys are equalent to Arg.Name, flags have an empty value.
//
// Use Input.Has("<key>") to check if a flag is passed.
//
// Access it the same way as a map, e.g. Input["<key>"]
type Input map[string]any

// Check that a flag was passed
func (input Input) Has(key string) bool {
	_, ok := input[key]
	return ok
}

// Get the option value or default if unset
func (input Input) Default(key string, defaultValue any) any {
	val, ok := input[key]
	if ok {
		return val
	}
	return defaultValue
}

func Get[T comparable](input Input, key string) T {
	return input[key].(T)
}
