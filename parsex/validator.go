package parsex

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

var ErrInvalidJSON = errors.New("provided JSON is not valid")
var ErrInvalidBool = errors.New("provided boolean is not valid")

type Validator func(in string) (out any, err error)

func ValidString(in string) (any, error) {
	return in, nil
}

func ValidJSON(in string) (any, error) {
	if !gjson.Valid(in) {
		return in, ErrInvalidJSON
	}
	return in, nil
}

func ValidPath(in string) (any, error) {
	path, err := filepath.Abs(in)
	if err != nil {
		return path, err
	}

	return path, nil
}

func ValidExistingPath(in string) (any, error) {
	path, err := ValidPath(in)
	if err != nil {
		return in, err
	}

	if _, err := os.Stat(path.(string)); os.IsNotExist(err) {
		return path, err
	}

	return path, nil
}

func ValidInt(in string) (any, error) {
	return strconv.Atoi(in)
}

func ValidFloat(in string) (any, error) {
	return strconv.ParseFloat(in, 64)
}

func ValidBool(in string) (any, error) {
	switch in {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, ErrInvalidBool
	}
}

func ValidList(validator Validator) Validator {
	return func(in string) (any, error) {
		split := strings.Split(in, ",")
		out := make([]any, len(split))

		for i, item := range split {
			value, err := validator(item)
			if err != nil {
				return out, err
			}
			out[i] = value
		}

		return out, nil
	}
}
