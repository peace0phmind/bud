package structure

import "reflect"

func ConvertTo[T any](from any) (T, error) {
	return ConvertToWithOption[T](from, defaultMapOption)
}

func MustConvertTo[T any](from any) T {
	if t, err := ConvertTo[T](from); err != nil {
		panic(err)
	} else {
		return t
	}
}

func ConvertToWithOption[T any](from any, option *MapOption) (t T, err error) {
	result, err := ConvertToTypeWithOption(from, reflect.TypeOf((*T)(nil)).Elem(), option)
	if err != nil {
		return t, err
	} else {
		return result.(T), nil
	}
}

func MustConvertToWithOption[T any](from any, option *MapOption) T {
	result, err := ConvertToWithOption[T](from, option)
	if err != nil {
		panic(err)
	} else {
		return result
	}
}

func ConvertToType(from any, toType reflect.Type) (any, error) {
	return ConvertToTypeWithOption(from, toType, defaultMapOption)
}

func ConvertToTypeWithOption(from any, toType reflect.Type, option *MapOption) (any, error) {
	result := reflect.New(toType).Elem()
	err := MapToValueWithOption(from, result, option)
	if err != nil {
		return nil, err
	}
	return result.Interface(), nil
}
