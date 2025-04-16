package cerr

import (
	"fmt"
	"io"
	"strings"
)

type DuringExecution struct {
	Name string
	Err  error
}

func (err DuringExecution) Error() string {
	return fmt.Sprintf("%s: %s", err.Name, err.Err.Error())
}

type DuringPreprocessing struct {
	Name string
	Err  error
}

func (err DuringPreprocessing) Error() string {
	return fmt.Sprintf("%s (internal/preprocessor): %s", err.Name, err.Err.Error())
}

type ExecIsNil struct {
	Name string
}

func (err ExecIsNil) Error() string {
	return fmt.Sprintf("%s (internal): runtime.Exec function is nil", err.Name)
}

type NotEnoughArgs struct {
	Name          string
	LenOfRequired int
	LenOfProvided int
	ArgPrinter    func(io.Writer)
	Args          []string
}

func (err NotEnoughArgs) Error() string {
	var builder strings.Builder
	fmt.Fprintf(&builder,
		"%s: %d positional argument(s) is/are required: `",
		err.Name,
		err.LenOfRequired,
	)
	err.ArgPrinter(&builder)
	fmt.Fprintf(
		&builder,
		"`, but provided only %d `%s`",
		err.LenOfProvided,
		err.Args,
	)
	return builder.String()
}
