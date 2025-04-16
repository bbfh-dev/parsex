package cerr

import (
	"fmt"
	"reflect"
)

type DataMustBePointer struct {
	Type reflect.Type
}

func (err DataMustBePointer) Error() string {
	return fmt.Sprintf("Program.Data must be a pointer. Got %q instead", err.Type)
}

type DataMustPointToStruct struct {
	Type reflect.Type
}

func (err DataMustPointToStruct) Error() string {
	return fmt.Sprintf("Program.Data must point to a struct{}. Got %q instead", err.Type)
}
