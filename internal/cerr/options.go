package cerr

import (
	"fmt"
	"reflect"
)

type OptionUnsupportedType struct {
	Type reflect.Type
}

func (err OptionUnsupportedType) Error() string {
	return fmt.Sprintf("(internal) can't set value of an option of type %q", err.Type)
}

type UnknownOption struct {
	Name   string
	Option string
}

func (err UnknownOption) Error() string {
	return fmt.Sprintf(
		"%s: unknown option %q. Refer to --help for usage information",
		err.Name,
		err.Option,
	)
}

type UnknownOptionCluster struct {
	Name   string
	Option string
}

func (err UnknownOptionCluster) Error() string {
	return fmt.Sprintf(
		"%s: unknown option or cluster %q. Refer to --help for usage information",
		err.Name,
		err.Option,
	)
}

type OptionNeedsValue struct {
	Name   string
	Option string
}

func (err OptionNeedsValue) Error() string {
	return fmt.Sprintf(
		"%s: option %q requires a value. Refer to --help for usage information",
		err.Name,
		err.Option,
	)
}

type ClusterOnlyFlags struct {
	Name   string
	Option string
	Mapped string
}

func (err ClusterOnlyFlags) Error() string {
	return fmt.Sprintf(
		"%s: cluster %q can only contain flags but %q is found. Refer to --help for usage information",
		err.Name,
		err.Option,
		err.Mapped,
	)
}

type SettingOption struct {
	Name   string
	Option string
	Err    error
}

func (err SettingOption) Error() string {
	return fmt.Sprintf(
		"%s: setting option %q: %s",
		err.Name,
		err.Option,
		err.Err.Error(),
	)
}
