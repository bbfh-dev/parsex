package parsex

import (
	"fmt"
	"strings"
)

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
