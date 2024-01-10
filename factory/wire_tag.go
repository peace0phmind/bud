package factory

import (
	"errors"
	"fmt"
	_struct "github.com/peace0phmind/bud/struct"
	"reflect"
	"strings"
)

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

	return _struct.WalkWithTagName(v, string(WireTag), func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type, wireValue string) error {
		wireValueType, value, err := ParseTagValue[WireValue](wireValue)
		if err != nil {
			panic(err)
		}

		switch wireValueType {
		case WireValueSelf, WireValueAuto, WireValueType, WireValueName:
			if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Interface {
				if fieldValue.IsNil() {
					switch wireValueType {
					case WireValueSelf:
						return _struct.SetField(fieldValue, v)
					case WireValueAuto:
						return _struct.SetField(fieldValue, _context._getByNameOrType(structField.Name, structField.Type))
					case WireValueType:
						return _struct.SetField(fieldValue, _context._get(structField.Type))
					case WireValueName:
						name := value
						if len(name) == 0 {
							name = structField.Name
						}
						return _struct.SetField(fieldValue, _context._getByName(strings.TrimSpace(name)))
					}
				}
			} else {
				return wireError(structField, rootTypes, wireValue)
			}
		case WireValueValue:
			if len(value) > 0 {
				values := splitAndTrimValue(value, ".")

				if len(values) == 2 {
					name, fieldName := values[0], values[1]

					namedObj := _context._getByName(name)
					namedObjVal := reflect.ValueOf(namedObj)

					if namedObjVal.Kind() == reflect.Ptr {
						namedObjVal = namedObjVal.Elem()
					}
					namedObjField := namedObjVal.FieldByName(fieldName)

					if namedObjField.IsValid() && namedObjField.CanInterface() {
						return _struct.SetField(fieldValue, namedObjField.Interface())
					} else {
						return errors.New(fmt.Sprintf("Set value error: field name %s, value rule: %s", structField.Name, wireValue))
					}
				}
			}

			return errors.New(fmt.Sprintf("wire value error, want 'name:objName.objFieldName', but got '%s'", wireValue))
		}

		return nil
	})
}
