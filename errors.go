package parsex

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

const errUnknownType = "unknown error"

type errKind int

const (
	ErrKindMustbePointer errKind = iota
	ErrKindPointToStruct

	ErrKindExecIsNil
	ErrKindExecution

	ErrKindNotEnoughArgs

	ErrKindUnknownOption
	ErrKindOptionNeedsValue
	ErrKindSettingOption
	ErrKindUnknownCluster
	ErrKindMistypedCluster
)

type ErrProgramData struct {
	ErrKind errKind
	Name    string
	Type    reflect.Type
}

func (err ErrProgramData) Error() string {
	switch err.ErrKind {
	case ErrKindMustbePointer:
		return fmt.Sprintf(
			"%s: Program.Data must be a pointer. Got %q instead",
			err.Name,
			err.Type,
		)
	case ErrKindPointToStruct:
		return fmt.Sprintf(
			"%s: Program.Data must point to a struct{}. Got %q instead",
			err.Name,
			err.Type,
		)
	}

	return errUnknownType
}

type ErrExecution struct {
	ErrKind errKind
	Name    string
	Err     error
}

func (err ErrExecution) Error() string {
	switch err.ErrKind {
	case ErrKindExecIsNil:
		return fmt.Sprintf("%s: runtime.Exec function is nil", err.Name)
	case ErrKindExecution:
		return fmt.Sprintf("%s: %s", err.Name, err.Err.Error())
	}

	return errUnknownType
}

type ErrInput struct {
	ErrKind     errKind
	Name        string
	RequiredLen int
	ProvidedLen int
	ExecArgs    []string
	ArgPrinter  func(io.Writer)
}

func (err ErrInput) Error() string {
	switch err.ErrKind {
	case ErrKindNotEnoughArgs:
		var builder strings.Builder
		fmt.Fprintf(&builder,
			"%s: %d positional argument(s) is/are required: `",
			err.Name,
			err.RequiredLen,
		)
		err.ArgPrinter(&builder)
		fmt.Fprintf(
			&builder,
			"`, but provided only %d `%s`",
			err.ProvidedLen,
			err.ExecArgs,
		)
		return builder.String()
	}

	return errUnknownType
}

type ErrOption struct {
	ErrKind errKind
	Name    string
	Option  string
	Err     error
}

func (err ErrOption) Error() string {
	switch err.ErrKind {
	case ErrKindUnknownOption:
		return fmt.Sprintf(
			"%s: unknown option %q. Refer to --help for usage information",
			err.Name,
			err.Option,
		)
	case ErrKindOptionNeedsValue:
		return fmt.Sprintf(
			"%s: option %q requires a value. Refer to --help for usage information",
			err.Name,
			err.Option,
		)
	case ErrKindSettingOption:
		return fmt.Sprintf(
			"%s: setting option %q: %s",
			err.Name,
			err.Option,
			err.Err.Error(),
		)
	case ErrKindUnknownCluster:
		return fmt.Sprintf(
			"%s: unknown option or cluster %q. Refer to --help for usage information",
			err.Name,
			err.Option,
		)
	case ErrKindMistypedCluster:
		return fmt.Sprintf(
			"%s: cluster %q can only contain flags but %q is found. Refer to --help for usage information",
			err.Name,
			err.Option,
			err.Err.Error(),
		)
	}

	return errUnknownType
}
