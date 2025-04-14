package parsex

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

var errCanceled = errors.New("parsex:err/cancelled")

func (program *Program) loadOptions() error {
	typePtr := reflect.TypeOf(program.Data)
	if typePtr.Kind() != reflect.Pointer {
		return fmt.Errorf("Program.Data must be a pointer. Got %q instead", typePtr)
	}

	typeElem := typePtr.Elem()
	if typeElem.Kind() != reflect.Struct {
		return fmt.Errorf("Program.Data must point to a struct{}. Got %q instead", typeElem)
	}

	valueElem := reflect.ValueOf(program.Data).Elem()
	numOfFields := typeElem.NumField()
	program.optionKeys = make([]string, 0, numOfFields+2)

	program.optionKeys = append(program.optionKeys, "help")
	if program.version != "" {
		program.optionKeys = append(program.optionKeys, "version")
	}

	for i := range numOfFields {
		fieldType := typeElem.Field(i)
		fieldValue := valueElem.Field(i)
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			// If parsex can't edit it, then ignore
			continue
		}

		name := strcase.ToKebab(fieldType.Name)
		program.optionKeys = append(program.optionKeys, name)
		program.options[name] = option{
			Name: name,
			Alt:  fieldType.Tag.Get("alt"),
			Desc: fieldType.Tag.Get("desc"),
			Type: fieldType.Type,
			Ref:  &fieldValue,
		}
	}

	return nil
}

func (program *Program) parse(vArgs []string) error {
	for i := 0; i < len(vArgs); i++ {
		if !strings.HasPrefix(vArgs[i], "-") {
			continue // TODO: .
		}
		arg := strings.TrimLeft(vArgs[i], "-")

		switch arg {

		case "help":
			program.printHelp(os.Stdout)
			return errCanceled

		case "version":
			program.printVersion(os.Stdout)
			return errCanceled
		}
	}

	return nil
}
