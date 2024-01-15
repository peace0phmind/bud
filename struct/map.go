package _struct

import (
	"fmt"
	"github.com/peace0phmind/bud/util"
	"reflect"
)

var typeAliasMap = util.Cache[reflect.Type, reflect.Type]{}

func init() {
	AddTypeAliasMap[int8, int]()
	AddTypeAliasMap[int16, int]()
	AddTypeAliasMap[int32, int]()
	AddTypeAliasMap[int64, int]()

	AddTypeAliasMap[uint8, uint]()
	AddTypeAliasMap[uint16, uint]()
	AddTypeAliasMap[uint32, uint]()
	AddTypeAliasMap[uint64, uint]()

	AddTypeAliasMap[float32, float64]()
}

type MapOption struct {
	ZeroFields bool
}

type Mapper func(from reflect.Value, to reflect.Value) error
type Converter func(input any, to reflect.Type) (output any, err error)

func AddTypeAliasMap[Name any, Alias any]() {
	nameType := reflect.TypeOf((*Name)(nil)).Elem()
	aliasType := reflect.TypeOf((*Alias)(nil)).Elem()

	if v, got := typeAliasMap.GetOrStore(nameType, aliasType); got {
		panic(fmt.Sprintf("type '%v' already set alias to '%v'", nameType, v))
	}
}

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
			to.Set(reflect.Zero(fromVal.Type()))
		}
		return nil
	}

	if !fromVal.IsValid() {
		// If the input value is invalid, then we just set the value
		// to be the zero value.
		to.Set(reflect.Zero(fromVal.Type()))
		return nil
	}

	return nil
}

func ConvertToType(from any, toType reflect.Type) (any, error) {
	return ConvertToTypeWithOption(from, toType, defaultMapOption)
}

func ConvertToTypeWithOption(from any, toType reflect.Type, option *MapOption) (any, error) {
	return nil, nil
}

func ConvertTo[T any](from any) (T, error) {
	return ConvertToWithOption[T](from, defaultMapOption)
}

func ConvertToWithOption[T any](from any, option *MapOption) (t T, err error) {
	result, err := ConvertToTypeWithOption(from, reflect.TypeOf((*T)(nil)).Elem(), option)
	if err != nil {
		return t, err
	} else {
		return result.(T), nil
	}
}
