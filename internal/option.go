package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var flagType = reflect.TypeOf(true)

var HelpOption = Option{
	Name: "help",
	Alt:  "",
	Desc: "Print this help message",
	Type: flagType,
	Ref:  nil,
}

var VersionOption = Option{
	Name: "version",
	Alt:  "",
	Desc: "Print the program version",
	Type: flagType,
	Ref:  nil,
}

type Option struct {
	Name string
	Alt  string
	Desc string

	Type reflect.Type
	Ref  *reflect.Value
}

func (option Option) Flag() string {
	if option.Alt == "" {
		return "--" + option.Name
	}
	return fmt.Sprintf("--%s, -%s", option.Name, option.Alt)
}

func (option Option) String() string {
	if option.IsFlag() {
		return option.Flag()
	}
	return fmt.Sprintf("%s <%s>", option.Flag(), option.Type.String())
}

func (option Option) IsFlag() bool {
	return option.Type.Kind() == reflect.Bool
}

func (option Option) Set(value string) error {
	if option.Ref == nil {
		return errors.New("(internal) option.Ref is nil!")
	}

	switch option.Type.Kind() {

	case reflect.String:
		option.Ref.SetString(value)
		return nil

	case reflect.Int:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		option.Ref.SetInt(num)
		return nil

	case reflect.Float64, reflect.Float32:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		option.Ref.SetFloat(num)
		return nil
	}

	return fmt.Errorf("(internal) can't set value of an option of type %q", option.Type)
}
