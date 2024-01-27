package structure

import (
	"errors"
	"fmt"
	"github.com/peace0phmind/bud/stream"
	"reflect"
	"strings"
	"unsafe"
)

type WalkFunc func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type) error

type ParamsWalkFunc[T any] func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type, params T) error

// 查找strut point的所有元素，包括子struct，将其中所有的field都调用WalkFunc进行处理
func _walk(v any, walkFn WalkFunc, rootTypes []reflect.Type) error {
	val := reflect.ValueOf(v)

	if reflect.Ptr == val.Kind() {
		val = val.Elem()
	}

	valType := val.Type()
	rootTypes = append(rootTypes, valType)

	for i := 0; i < valType.NumField(); i++ {
		fieldValue := val.Field(i)
		structField := valType.Field(i)

		ff := fieldValue
		// 尝试遍历指针型struct, 为空则不用遍历
		if reflect.Ptr == fieldValue.Kind() && !fieldValue.IsNil() {
			ff = fieldValue.Elem()
		}

		// 尝试遍历内嵌型struct
		// TODO 尝试遍历非导出结构体field
		if reflect.Struct == ff.Kind() && ff.CanAddr() && ff.CanInterface() {
			if err := _walk(ff.Addr().Interface(), walkFn, rootTypes); err != nil {
				return err
			}
		}

		if reflect.Slice == ff.Kind() {
			// 处理slice下的point struct
			if reflect.Ptr == ff.Type().Elem().Kind() && reflect.Struct == ff.Type().Elem().Elem().Kind() {
				for i := 0; i < ff.Len(); i++ {
					if ff.Index(i).CanAddr() && !ff.Index(i).IsNil() && ff.Index(i).CanInterface() {
						if err := _walk(ff.Index(i).Interface(), walkFn, rootTypes); err != nil {
							return err
						}
					}
				}
			}

			// 获取slice下的元素类型是否是struct
			if reflect.Struct == ff.Type().Elem().Kind() {
				for i := 0; i < ff.Len(); i++ {
					if ff.Index(i).CanAddr() && ff.Index(i).CanInterface() {
						if err := _walk(ff.Index(i).Addr().Interface(), walkFn, rootTypes); err != nil {
							return err
						}
					}
				}
			}
		}

		if err := walkFn(fieldValue, structField, rootTypes); err != nil {
			return err
		}
	}

	return nil
}

func WalkField(v any, walkFn WalkFunc) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer")
	}

	val = val.Elem()
	if !val.CanAddr() {
		return errors.New("result must be addressable (a pointer)")
	}

	if val.Kind() != reflect.Struct {
		return errors.New("result must be a struct")
	}

	return _walk(v, walkFn, nil)
}

func WalkWithTagNames(v any, tagNames []string, walkFn ParamsWalkFunc[map[string]string]) error {
	return WalkField(v, func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type) error {
		tags := map[string]string{}
		for _, tagName := range tagNames {
			if tagValue, ok := structField.Tag.Lookup(tagName); ok {
				tags[tagName] = tagValue
			}
		}

		if len(tags) > 0 {
			if err := walkFn(fieldValue, structField, rootTypes, tags); err != nil {
				return err
			}
		}

		return nil
	})
}

func SetField(fieldValue reflect.Value, v any) error {

	if !fieldValue.CanAddr() {
		return errors.New("fieldValue is not addressable")
	}

	valueType := reflect.TypeOf(v)
	if !valueType.ConvertibleTo(fieldValue.Type()) {
		return errors.New("value is not assignable to the field type")
	}

	if fieldValue.CanSet() {
		fieldValue.Set(reflect.ValueOf(v))
	} else {
		// 通过unsafe包绕过CanSet的限制
		rf := reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
		rf.Set(reflect.ValueOf(v))
	}

	return nil
}

func GetFieldPath(structField reflect.StructField, rootTypes []reflect.Type) string {
	var results []string
	if len(rootTypes) > 0 {
		results = append(results, rootTypes[0].PkgPath())
		names := stream.Must(stream.Map(stream.Of(rootTypes), func(in reflect.Type) (string, error) {
			return in.Name(), nil
		}).ToSlice())
		results = append(results, names...)
	}

	result := ""
	if len(results) > 0 {
		result = strings.Join(results, "/") + "."
	}

	result += fmt.Sprintf("%s(%s)", structField.Name, structField.Type.String())

	return result
}
