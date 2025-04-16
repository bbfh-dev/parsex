package parsex

import (
	"strings"

	"github.com/bbfh-dev/parsex/internal/cerr"
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
		return cerr.UnknownOption{
			Name:   runtime.name,
			Option: arg,
		}
	}
	if option.IsFlag() {
		option.SetFlag()
		return nil
	}

	*i++
	if *i >= len(inputArgs) {
		return cerr.OptionNeedsValue{
			Name:   runtime.name,
			Option: arg,
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
			return cerr.OptionNeedsValue{
				Name:   runtime.name,
				Option: arg,
			}
		}
		return runtime.setOption(name, inputArgs[*i])
	}

	// Each character in the cluster should map to a flag.
	for _, char := range optionStr {
		flagKey := string(char)
		mapped, exists := runtime.genOptionAlts[flagKey]
		if !exists {
			return cerr.UnknownOptionCluster{
				Name:   runtime.name,
				Option: arg,
			}
		}
		option, exists := runtime.genOptions.Get(mapped)
		if !exists || !option.IsFlag() {
			return cerr.ClusterOnlyFlags{
				Name:   runtime.name,
				Option: arg,
				Mapped: mapped,
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
		return cerr.UnknownOption{
			Name:   runtime.name,
			Option: "--" + name,
		}
	}
	if option.IsFlag() {
		option.SetFlag()
		return nil
	}
	if err := option.Set(value); err != nil {
		return cerr.SettingOption{
			Name:   runtime.name,
			Option: "--" + name,
			Err:    err,
		}
	}
	return nil
}
