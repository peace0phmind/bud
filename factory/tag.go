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

// ENUM(self, auto, type, name, value)
type WireValue string

func splitAndTrimValue(value, sep string) []string {
	return stream.Of(strings.Split(strings.TrimSpace(value), sep)).
		Map(func(s string) string { return strings.TrimSpace(s) }).
		Filter(func(s string) (bool, error) { return len(s) > 0, nil }).MustToSlice()
}

func ParseTagValue[T any](tagValue string) (t T, v string, err error) {
	values := splitAndTrimValue(tagValue, ":")
	if len(values) == 0 {
		return t, v, errors.New("tag value is empty")
	}

	if len(values) > 2 {
		return t, v, errors.New("tag value contains multiple `:`")
	}

	if unmarshaler, ok := any(&t).(encoding.TextUnmarshaler); ok {
		if err = unmarshaler.UnmarshalText([]byte(values[0])); err != nil {
			return t, v, err
		} else {
			if len(values) == 2 {
				v = values[1]
			}
			return t, v, nil
		}
	} else {
		panic("parse type muse implements encoding.TextUnmarshaler")
	}
}
