package factory

import (
	"encoding"
	"errors"
	"github.com/peace0phmind/bud/stream"
	"strings"
)

//go:generate go-enum --marshal --values --nocomments --nocase

// WireTag is a constant that defines the annotation string used for wire injection in Go code.
const WireTag = "wire"

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
