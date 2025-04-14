package parsex

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type option struct {
	Name string
	Alt  string
	Desc string

	Type reflect.Type
	Ref  *reflect.Value
}

func (opt option) Flag() string {
	if opt.Alt == "" {
		return "--" + opt.Name
	}
	return fmt.Sprintf("--%s, -%s", opt.Name, opt.Alt)
}

func (opt option) String() string {
	if opt.IsFlag() {
		return opt.Flag()
	}
	return fmt.Sprintf("%s <%s>", opt.Flag(), opt.Type.String())
}

func (opt option) IsFlag() bool {
	return opt.Type.Kind() == reflect.Bool
}

func (opt option) Set(value string) error {
	if opt.Ref == nil {
		return errors.New("(internal) option.Ref is nil!")
	}

	switch opt.Type.Kind() {

	case reflect.String:
		opt.Ref.SetString(value)
		return nil

	case reflect.Int:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		opt.Ref.SetInt(num)
		return nil

	case reflect.Float64, reflect.Float32:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		opt.Ref.SetFloat(num)
		return nil
	}

	return fmt.Errorf("(internal) can't set value of an option of type %q", opt.Type)
}
