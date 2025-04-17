package parsex

import (
	"errors"
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
		return ErrOption{
			ErrKind: ErrKindUnknownOption,
			Name:    runtime.name,
			Option:  arg,
			Err:     nil,
		}
	}
	if option.IsFlag() {
		option.SetFlag()
		return nil
	}

	*i++
	if *i >= len(inputArgs) {
		return ErrOption{
			ErrKind: ErrKindOptionNeedsValue,
			Name:    runtime.name,
			Option:  arg,
			Err:     nil,
		}
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
			return ErrOption{
				ErrKind: ErrKindOptionNeedsValue,
				Name:    runtime.name,
				Option:  arg,
				Err:     nil,
			}
		}
		return runtime.setOption(name, inputArgs[*i])
	}

	// Each character in the cluster should map to a flag.
	for _, char := range optionStr {
		flagKey := string(char)
		mapped, exists := runtime.genOptionAlts[flagKey]
		if !exists {
			return ErrOption{
				ErrKind: ErrKindUnknownCluster,
				Name:    runtime.name,
				Option:  arg,
				Err:     nil,
			}
		}
		option, exists := runtime.genOptions.Get(mapped)
		if !exists || !option.IsFlag() {
			return ErrOption{
				ErrKind: ErrKindMistypedCluster,
				Name:    runtime.name,
				Option:  arg,
				Err:     errors.New(mapped),
			}
		}
		option.SetFlag()
	}

	return nil
}

// setOption retrieves the option by name and applies the value.
func (runtime *runtimeType) setOption(name, value string) error {
	option, exists := runtime.genOptions.Get(name)
	if !exists {
		return ErrOption{
			ErrKind: ErrKindUnknownOption,
			Name:    runtime.name,
			Option:  "--" + name,
			Err:     nil,
		}
	}
	if option.IsFlag() {
		option.SetFlag()
		return nil
	}
	if err := option.Set(value); err != nil {
		return ErrOption{
			ErrKind: ErrKindSettingOption,
			Name:    runtime.name,
			Option:  "--" + name,
			Err:     err,
		}
	}
	return nil
}
