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

	ck := getContextKeyFromType(reflect.TypeOf(v))
	_, ok := _context().wiringCache.GetOrStore(ck, true)
	if ok {
		panic(fmt.Sprintf("%s:%s is wiring, possible circular reference exists.", ck.Package, ck.Name))
	}
	defer _context().wiringCache.Delete(ck)

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
				return _struct.SetField(fieldValue, _context()._get(structField.Type))
			}
			return nil
		}

		if strings.HasPrefix(lowRule, AnnotationWireName) {
			rules := strings.Split(wireRule, ":")
			if len(rules) != 2 {
				return errors.New(fmt.Sprintf("wire format error, want 'name:NamedSingleton' got '%s'", wireRule))
			}
			return _struct.SetField(fieldValue, _context()._getByName(strings.TrimSpace(rules[1])))
		}

		return nil
	})
}
