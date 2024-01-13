package factory

import (
	"errors"
	"fmt"
	_struct "github.com/peace0phmind/bud/struct"
	"reflect"
)

func wireError(structField reflect.StructField, rootTypes []reflect.Type, wireRule string) error {
	fieldPath := _struct.GetFieldPath(structField, rootTypes)
	return errors.New(fmt.Sprintf("The field of 'wire' must be defined as a pointer to an object or an interface. %s, tag value: %s", fieldPath, wireRule))
}

func getByWireTag(tagValue *TagValue[WireValue], t reflect.Type) (any, error) {
	switch tagValue.Tag {
	case WireValueAuto:
		if len(tagValue.Value) > 0 {
			return _context._getByNameOrType(tagValue.Value, t), nil
		} else {
			return _context._get(t), nil
		}
	case WireValueType:
		return _context._get(t), nil
	case WireValueName:
		if len(tagValue.Value) > 0 {
			return _context._getByName(tagValue.Value), nil
		}
	case WireValueValue:
		if len(tagValue.Value) > 0 {
			values := splitAndTrimValue(tagValue.Value, ".")

			if len(values) == 2 {
				name, fieldName := values[0], values[1]

				namedObj := _context._getByName(name)
				namedObjVal := reflect.ValueOf(namedObj)

				if namedObjVal.Kind() == reflect.Ptr {
					namedObjVal = namedObjVal.Elem()
				}
				namedObjField := namedObjVal.FieldByName(fieldName)

				if namedObjField.IsValid() && namedObjField.CanInterface() {
					return namedObjField.Interface(), nil
				}
			}
		}
	}

	return nil, errors.New(fmt.Sprintf("Tag value not supported: %+v", tagValue))
}

func AutoWire(self any) error {
	if self == nil {
		return nil
	}

	vt := reflect.TypeOf(self)
	_context._addWiring(vt)
	defer _context._deleteWiring(vt)

	return _struct.WalkWithTagName(self, WireTag, func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type, wireValue string) error {
		tv, err := ParseTagValue[WireValue](wireValue, func(tv *TagValue[WireValue]) {
			if (tv.Tag == WireValueName && len(tv.Value) == 0) ||
				(tv.Tag == WireValueAuto) {
				tv.Value = structField.Name
			}
		})

		if err != nil {
			panic(err)
		}

		switch tv.Tag {
		case WireValueSelf, WireValueAuto, WireValueType, WireValueName:
			if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Interface {
				if fieldValue.IsNil() {
					switch tv.Tag {
					case WireValueSelf:
						return _struct.SetField(fieldValue, self)
					default:
						if wiredValue, err1 := getByWireTag(tv, structField.Type); err1 == nil {
							return _struct.SetField(fieldValue, wiredValue)
						} else {
							return errors.New(fmt.Sprintf("%v on field: %s", err1, structField.Name))
						}
					}
				}
			} else {
				return wireError(structField, rootTypes, wireValue)
			}
		case WireValueValue:
			if wiredValue, err1 := getByWireTag(tv, structField.Type); err1 == nil {
				return _struct.SetField(fieldValue, wiredValue)
			} else {
				return errors.New(fmt.Sprintf("%v on field: %s. 'name:objName.objFieldName'", err1, structField.Name))
			}
		}

		return nil
	})
}
