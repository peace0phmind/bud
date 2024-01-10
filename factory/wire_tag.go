package factory

import (
	"errors"
	"fmt"
	"github.com/peace0phmind/bud/stream"
	_struct "github.com/peace0phmind/bud/struct"
	"reflect"
	"strings"
)

// WireTag is a constant that defines the annotation string used for wire injection in Go code.
const WireTag = "wire"

// WireValueSelf is a constant that defines the annotation string used for self injection in Go code.
const WireValueSelf = "self"

// WireValueAuto is a constant that defines the annotation string used for automatic wire injection in Go code. It is used in the AutoWire function.
const WireValueAuto = "auto"

// WireValueName is a constant that defines the annotation string used for name injection in Go code.
const WireValueName = "name"

// WireValueValue is a constant that defines the annotation string used for value injection in Go code.
const WireValueValue = "value"

func wireError(structField reflect.StructField, rootTypes []reflect.Type, wireRule string) error {
	fieldPath := _struct.GetFieldPath(structField, rootTypes)
	return errors.New(fmt.Sprintf("The field of 'wire' must be defined as a pointer to an object or an interface. %s, tag value: %s", fieldPath, wireRule))
}

func AutoWire(v any) error {
	if v == nil {
		return nil
	}

	vt := reflect.TypeOf(v)
	_context._addWiring(vt)
	defer _context._deleteWiring(vt)

	return _struct.WalkWithTagName(v, WireTag, func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type, wireRule string) error {
		lowRule := strings.ToLower(wireRule)

		if WireValueSelf == lowRule {
			if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Interface {
				if fieldValue.IsNil() {
					return _struct.SetField(fieldValue, v)
				}
				return nil
			} else {
				return wireError(structField, rootTypes, wireRule)
			}
		}

		if WireValueAuto == lowRule {
			if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Interface {
				if fieldValue.IsNil() {
					return _struct.SetField(fieldValue, _context._get(structField.Type))
				}
				return nil
			} else {
				return wireError(structField, rootTypes, wireRule)
			}
		}

		if strings.HasPrefix(lowRule, WireValueName) {
			if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Interface {
				if fieldValue.IsNil() {
					name := strings.TrimSpace(wireRule[len(WireValueName):])

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
				return nil
			} else {
				return wireError(structField, rootTypes, wireRule)
			}
		}

		if strings.HasPrefix(lowRule, WireValueValue) {
			valueRule := strings.TrimSpace(wireRule[len(WireValueValue):])

			if strings.HasPrefix(valueRule, ":") {
				valueRule = strings.TrimSpace(valueRule[1:])

				values := stream.Of(strings.Split(valueRule, ".")).
					Map(func(s string) string { return strings.TrimSpace(s) }).
					Filter(func(s string) (bool, error) { return len(s) > 0, nil }).MustToSlice()

				if len(values) == 2 {
					name := values[0]
					fieldName := values[1]

					namedObj := _context._getByName(strings.TrimSpace(name))

					namedObjVal := reflect.ValueOf(namedObj)

					if namedObjVal.Kind() == reflect.Ptr {
						namedObjVal = namedObjVal.Elem()
					}

					namedObjField := namedObjVal.FieldByName(fieldName)

					if namedObjField.IsValid() && namedObjField.CanInterface() {
						return _struct.SetField(fieldValue, namedObjField.Interface())
					} else {
						return errors.New(fmt.Sprintf("Set value error: field name %s, value rule: %s", structField.Name, valueRule))
					}
				}
			}

			return errors.New(fmt.Sprintf("wire value error, want 'name:objName.objFieldName', but got '%s'", wireRule))
		}

		return nil
	})
}
