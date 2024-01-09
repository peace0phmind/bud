package factory

import (
	"errors"
	"fmt"
	_struct "github.com/peace0phmind/bud/struct"
	"reflect"
	"strings"
)

// AnnotationWire is a constant that defines the annotation string used for wire injection in Go code.
const AnnotationWire = "wire"

// AnnotationWireSelf is a constant that defines the annotation string used for self injection in Go code.
const AnnotationWireSelf = "self"

// AnnotationWireAuto is a constant that defines the annotation string used for automatic wire injection in Go code. It is used in the AutoWire function.
const AnnotationWireAuto = "auto"

// AnnotationWireName is a constant that defines the annotation string used for name injection in Go code.
const AnnotationWireName = "name"

// AnnotationDefault is a constant that defines the annotation string used as the default value in Go code.
const AnnotationDefault = "default"

func SetDefault(v any) error {

	return nil
}

func AutoWire(v any) error {
	if v == nil {
		return nil
	}

	vt := reflect.TypeOf(v)
	_context._addWire(vt)
	defer _context._deleteWire(vt)

	return _struct.WalkWithTagName(v, AnnotationWire, func(fieldValue reflect.Value, structField reflect.StructField, rootFields []reflect.StructField, wireRule string) error {
		lowRule := strings.ToLower(wireRule)

		if AnnotationWireSelf == lowRule {
			if fieldValue.IsNil() {
				return _struct.SetField(fieldValue, v)
			}
			return nil
		}

		if AnnotationWireAuto == lowRule {
			if fieldValue.IsNil() {
				return _struct.SetField(fieldValue, _context._get(structField.Type))
			}
			return nil
		}

		if strings.HasPrefix(lowRule, AnnotationWireName) {
			if fieldValue.IsNil() {
				name := strings.TrimSpace(wireRule[len(AnnotationWireName):])

				if len(name) == 0 || name == ":" {
					name = structField.Name
				} else if strings.HasPrefix(name, ":") {
					name = strings.TrimSpace(name[1:])
				} else {
					name = ""
				}

				if len(name) > 0 {
					return _struct.SetField(fieldValue, _context._getByName(strings.TrimSpace(name)))
				}

				return errors.New(fmt.Sprintf("wire format error, want 'name:NamedSingleton' or 'name:' or 'name', but got '%s'", wireRule))
			}
		}

		return nil
	})
}
