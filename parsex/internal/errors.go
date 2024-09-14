package internal

import "fmt"

// Add prefix to an error if err != nil
func PrefixErr(prefix string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", prefix, err)
}
