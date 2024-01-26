package enum

import (
	"errors"
	"fmt"
	"reflect"
)

var enumTypes = []reflect.Kind{
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Float32,
	reflect.Float64,
	reflect.String,
}

func getKindByName(name string) (reflect.Kind, error) {
	for _, k := range enumTypes {
		if k.String() == name {
			return k, nil
		}
	}

	return reflect.Invalid, errors.New(fmt.Sprintf("unknown reflect.kind name %s", name))
}