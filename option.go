package parsex

import "reflect"

type option struct {
	Name string
	Alt  string
	Desc string

	Type reflect.Type
	Ref  *reflect.Value
}
