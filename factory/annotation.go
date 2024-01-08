package factory

import (
	_struct "github.com/peace0phmind/bud/struct"
	"reflect"
	"strings"
)

// AnnotationWire is a constant that defines the annotation string used for wire injection in Go code.
const AnnotationWire = "wire"

// AnnotationWireSelf is a constant that defines the annotation string used for self injection in Go code.
const AnnotationWireSelf = "self"

const AnnotationWireAuto = "auto"

// AnnotationDefault is a constant that defines the annotation string used as the default value in Go code.
const AnnotationDefault = "default"

func UpdateDefault(v any) error {

	return nil
}

func AutoWire(v any) error {
	if v == nil {
		return nil
	}

	return _struct.WalkWithTagName(v, AnnotationWire, func(fieldValue reflect.Value, structField reflect.StructField, rootFields []reflect.StructField, wireType string) error {
		switch strings.ToLower(wireType) {
		case AnnotationWireSelf:
			if fieldValue.IsNil() && fieldValue.CanAddr() {
				selfType := reflect.TypeOf(v)
				if selfType.ConvertibleTo(fieldValue.Type()) {
					_struct.SetField(fieldValue, v)
				}
			}
		case AnnotationWireAuto:

		}

		return nil
	})
}
