package structure

import (
	"errors"
	"fmt"
	"github.com/peace0phmind/bud/util"
	"reflect"
)

type MapOption struct {
	ZeroFields bool
}

type Mapper func(from reflect.Value, to reflect.Value) error

type mapperKey struct {
	from reflect.Type
	to   reflect.Type
}

var mapperCache = util.Cache[mapperKey, Mapper]{}

func RegisterMapper[From any, To any](mapper Mapper) {
	fromType := reflect.TypeOf((*From)(nil)).Elem()
	toType := reflect.TypeOf((*To)(nil)).Elem()

	key := mapperKey{from: fromType, to: toType}
	if _, got := mapperCache.GetOrStore(key, mapper); got {
		panic(fmt.Sprintf("type %+v already registed", key))
	}
}

func ReplaceMapper[From any, To any](mapper Mapper) {
	fromType := reflect.TypeOf((*From)(nil)).Elem()
	toType := reflect.TypeOf((*To)(nil)).Elem()

	mapperCache.Set(mapperKey{from: fromType, to: toType}, mapper)
}

func NewMapOption() *MapOption {
	return &MapOption{
		ZeroFields: true,
	}
}

var defaultMapOption = NewMapOption()

func Map(from, to any) error {
	return MapWithOption(from, to, defaultMapOption)
}

func MapWithOption(from, to any, option *MapOption) error {
	return MapToValueWithOption(from, reflect.ValueOf(to).Elem(), option)
}

func MapToValue(from any, to reflect.Value) error {
	return MapToValueWithOption(from, to, defaultMapOption)
}

func MapToValueWithOption(from any, to reflect.Value, option *MapOption) error {
	if option == nil {
		option = defaultMapOption
	}

	var fromVal reflect.Value
	if from != nil {
		fromVal = reflect.ValueOf(from)

		// We need to check here if input is a typed nil. Typed nils won't
		// match the "input == nil" below so we check that here.
		if fromVal.Kind() == reflect.Ptr && fromVal.IsNil() {
			from = nil
		}
	}

	if from == nil {
		// If the data is nil, then we don't set anything, unless ZeroFields is set
		// to true.
		if option.ZeroFields {
			to.Set(reflect.Zero(to.Type()))
		}
		return nil
	}

	return Value2ValueWithOption(fromVal, to, option)
}

func value2PtrValueWithOption(from reflect.Value, to reflect.Value, option *MapOption) error {
	toElemType := to.Type().Elem()

	if to.CanSet() {
		realTo := to
		if realTo.IsNil() || option.ZeroFields {
			realTo = reflect.New(toElemType)
		}

		if err := Value2ValueWithOption(reflect.Indirect(from), reflect.Indirect(realTo), option); err != nil {
			return err
		}

		to.Set(realTo)
	} else {
		if err := Value2ValueWithOption(reflect.Indirect(from), reflect.Indirect(to), option); err != nil {
			return err
		}
	}
	return nil
}

func Value2ValueWithOption(from reflect.Value, to reflect.Value, option *MapOption) error {
	if !from.IsValid() {
		// If the input value is invalid, then we just set the value
		// to be the zero value.
		to.Set(reflect.Zero(to.Type()))
		return nil
	}

	if to.Kind() == reflect.Ptr {
		return value2PtrValueWithOption(from, to, option)
	}

	fromType := from.Type()
	toType := to.Type()

	mapper, ok := mapperCache.Get(mapperKey{from: fromType, to: toType})
	if !ok {
		fromType, _ = typeAliasMap.GetOrDefault(fromType, fromType)
		toType, _ = typeAliasMap.GetOrDefault(toType, toType)

		mapper, ok = mapperCache.Get(mapperKey{from: fromType, to: toType})
		if !ok {
			return errors.New(fmt.Sprintf("no mapper found for type %+v to %+v", fromType, toType))
		}
	}

	if err := mapper(from, to); err != nil {
		return err
	}

	return nil
}