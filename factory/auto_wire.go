package factory

import (
	"encoding"
	"errors"
	"fmt"
	"github.com/peace0phmind/bud/stream"
	"github.com/peace0phmind/bud/structure"
	"reflect"
	"strings"
)

//go:generate go-enum --marshal --values --nocomments --nocase

// WireTag is a constant that defines the annotation string used for wire injection in Go code.
const WireTag = "wire"

// ValueTag is a shortcut of wire:"value:"
const ValueTag = "value"

// WireValue is a enum
// ENUM(self, auto, type, name, value, option)
type WireValue string

func splitAndTrimValue(value, sep string) []string {
	return stream.Of(strings.Split(strings.TrimSpace(value), sep)).
		Map(func(s string) string { return strings.TrimSpace(s) }).
		Filter(func(s string) (bool, error) { return len(s) > 0, nil }).MustToSlice()
}

type TagValue[T any] struct {
	Tag   T
	Value string
}

func (tv *TagValue[T]) String() string {
	return fmt.Sprintf("%v:%s", tv.Tag, tv.Value)
}

func ParseTagValue[T any](tagValue string, checkAndSet func(tv *TagValue[T])) (tv *TagValue[T], err error) {
	result := &TagValue[T]{}
	values := splitAndTrimValue(tagValue, ":")
	if len(values) == 0 {
		return nil, errors.New("tag value is empty")
	}

	if len(values) > 2 {
		return nil, errors.New("tag value contains multiple `:`")
	}

	if unmarshaler, ok := any(&result.Tag).(encoding.TextUnmarshaler); ok {
		if err = unmarshaler.UnmarshalText([]byte(values[0])); err != nil {
			return nil, err
		} else {
			if len(values) == 2 {
				result.Value = values[1]
			}

			if checkAndSet != nil {
				checkAndSet(result)
			}

			return result, nil
		}
	} else {
		panic("parse type muse implements encoding.TextUnmarshaler")
	}
}

func wireError(structField reflect.StructField, rootTypes []reflect.Type, wireRule string) error {
	fieldPath := structure.GetFieldPath(structField, rootTypes)
	return errors.New(fmt.Sprintf("The field of 'wire' must be defined as a pointer to an object or an interface. %s, tag value: %s", fieldPath, wireRule))
}

func getExpr(value string) (exprCode string, isExpr bool) {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		return strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}"), true
	}
	return value, false
}

func getByWireTag(tagValue *TagValue[WireValue], t reflect.Type) (any, error) {
	switch tagValue.Tag {
	case WireValueAuto:
		if len(tagValue.Value) > 0 {
			return _context.getByNameOrType(tagValue.Value, t), nil
		} else {
			return _context.getByType(t), nil
		}
	case WireValueType:
		return _context.getByType(t), nil
	case WireValueName:
		if len(tagValue.Value) > 0 {
			return _context.getByName(tagValue.Value), nil
		}
	case WireValueValue:
		if len(tagValue.Value) > 0 {

			exprCode, isExpr := getExpr(tagValue.Value)
			if isExpr {
				value, err := _context.evalExpr(exprCode)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Tag value %s expr eval err: %v", tagValue, err))
				}

				return structure.ConvertToType(value, t)
			} else {
				return structure.ConvertToType(exprCode, t)
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
	_context.wiring(vt)
	defer _context.wired(vt)

	return structure.WalkWithTagNames(self, []string{WireTag, ValueTag}, func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type, tags map[string]string) (err error) {
		if len(tags) > 1 {
			panic("Only one can exist at a time, either 'wire' or 'value'.")
		}

		var tv *TagValue[WireValue]
		if wireValue, ok := tags[WireTag]; ok {
			tv, err = ParseTagValue[WireValue](wireValue, func(tv *TagValue[WireValue]) {
				if (tv.Tag == WireValueName && len(tv.Value) == 0) ||
					(tv.Tag == WireValueAuto) {
					tv.Value = structField.Name
				}
			})
		}
		if wireValue, ok := tags[ValueTag]; ok {
			tv = &TagValue[WireValue]{Tag: WireValueValue, Value: wireValue}
		}

		if err != nil {
			panic(err)
		}

		switch tv.Tag {
		case WireValueSelf, WireValueAuto, WireValueType, WireValueName:
			if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Interface {
				if fieldValue.IsNil() {
					switch tv.Tag {
					case WireValueSelf:
						return structure.SetField(fieldValue, self)
					default:
						if wiredValue, err1 := getByWireTag(tv, structField.Type); err1 == nil {
							return structure.SetField(fieldValue, wiredValue)
						} else {
							return errors.New(fmt.Sprintf("%v on field: %s", err1, structField.Name))
						}
					}
				}
			} else {
				return wireError(structField, rootTypes, tv.String())
			}
		case WireValueValue:
			if wiredValue, err1 := getByWireTag(tv, structField.Type); err1 == nil {
				return structure.SetField(fieldValue, wiredValue)
			} else {
				return errors.New(fmt.Sprintf("%v on field: %s. 'name:objName.objFieldName'", err1, structField.Name))
			}
		}

		return nil
	})
}
