package parsex

import (
	"fmt"
	"reflect"
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
	if opt.Type.Kind() == reflect.Bool {
		return opt.Flag()
	}
	return fmt.Sprintf("%s <%s>", opt.Flag(), opt.Type.String())
}
