package factory

import (
	_struct "github.com/peace0phmind/bud/struct"
	"reflect"
)

// DefaultTag is a constant that defines the annotation string used as the default value in Go code.
const DefaultTag = "default"

func SetDefault(v any) error {
	if v == nil {
		return nil
	}

	return _struct.WalkWithTagName(v, DefaultTag, func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type, defaultValue string) error {

		return nil
	})
}
