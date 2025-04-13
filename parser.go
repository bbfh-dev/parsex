package parsex

import (
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
)

type parser struct {
	// Used to keep the options sorted
	optionKeys []string
	// The actual option map. Map is used for more performant access
	options map[string]option
}

func newParser() *parser {
	return &parser{}
}

func (parser *parser) LoadOptionsFrom(data any) error {
	typePtr := reflect.TypeOf(data)
	if typePtr.Kind() != reflect.Pointer {
		return fmt.Errorf("Program.Data must be a pointer. Got %q instead", typePtr)
	}

	typeElem := typePtr.Elem()
	if typeElem.Kind() != reflect.Struct {
		return fmt.Errorf("Program.Data must point to a struct{}. Got %q instead", typeElem)
	}

	valueElem := reflect.ValueOf(data).Elem()
	numOfFields := typeElem.NumField()
	parser.optionKeys = make([]string, numOfFields)

	for i := range numOfFields {
		fieldType := typeElem.Field(i)
		fieldValue := valueElem.Field(i)
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			// If parsex can't edit it, then ignore
			continue
		}

		name := strcase.ToKebab(fieldType.Name)
		parser.optionKeys[i] = name
		parser.options[name] = option{
			Name: name,
			Alt:  fieldType.Tag.Get("alt"),
			Desc: fieldType.Tag.Get("desc"),
			Type: fieldType.Type,
			Ref:  &fieldValue,
		}
	}

	return nil
}

func (parser *parser) Parse(vArgs []string) error {
	return nil
}
