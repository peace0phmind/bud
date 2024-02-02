package opt

import (
	"encoding/json"
	"strings"
)

const JsonNullValue = "null"

type JsonNull[T any] struct {
	JsonV T
	Set   bool
}

func NewJsonNull[T any](t T) *JsonNull[T] {
	return &JsonNull[T]{JsonV: t, Set: true}
}

func (jn *JsonNull[T]) MarshalJSON() ([]byte, error) {
	if !jn.Set {
		return []byte(JsonNullValue), nil
	}
	return json.Marshal(jn.JsonV)
}

func (jn *JsonNull[T]) UnmarshalJSON(b []byte) error {
	jn.Set = true
	if strings.EqualFold(string(b), JsonNullValue) {
		return nil
	}
	return json.Unmarshal(b, &jn.JsonV)
}
