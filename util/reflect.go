package util

import (
	"reflect"
)

// Function to parse object's tag for `default`
func ParseTag(obj interface{}) (defaultValues map[string]string, err error) {
	defaultValues = make(map[string]string)
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		defaultValue := field.Tag.Get("default")
		if defaultValue != "" {
			defaultValues[field.Name] = defaultValue
		}
	}

	return defaultValues, nil
}
