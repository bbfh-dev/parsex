package parsex

import (
	"log"
	"reflect"

	"github.com/bbfh-dev/parsex/v2/internal"
	"github.com/iancoleman/strcase"
)

func (runtime *runtimeType) preprocess() error {
	if runtime.data == nil {
		return nil
	}

	typePtr := reflect.TypeOf(runtime.data)
	if typePtr.Kind() != reflect.Pointer {
		return ErrProgramData{
			ErrKind: ErrKindMustbePointer,
			Name:    runtime.name,
			Type:    typePtr,
		}
	}

	typeElem := typePtr.Elem()
	if typeElem.Kind() != reflect.Struct {
		return ErrProgramData{
			ErrKind: ErrKindPointToStruct,
			Name:    runtime.name,
			Type:    typePtr,
		}
	}

	valueElem := reflect.ValueOf(runtime.data).Elem()
	numOfFields := typeElem.NumField()

	runtime.genOptions.Clear()
	runtime.genOptions.Add("help", internal.HelpOption)
	if runtime.version != "" {
		runtime.genOptions.Add("version", internal.VersionOption)
	}

	for i := range numOfFields {
		fieldType := typeElem.Field(i)
		fieldValue := valueElem.Field(i)
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			log.Printf("WARNING: (Parsex reflection) not possible to modify field %+v", fieldType)
			continue
		}

		name := strcase.ToKebab(fieldType.Name)
		runtime.genOptions.Add(
			name,
			internal.Option{
				Name: name,
				Alt:  fieldType.Tag.Get("alt"),
				Desc: fieldType.Tag.Get("desc"),
				Type: fieldType.Type,
				Ref:  &fieldValue,
			},
		)
		if alt := fieldType.Tag.Get("alt"); alt != "" {
			runtime.genOptionAlts[alt] = name
		}
	}

	return nil
}
