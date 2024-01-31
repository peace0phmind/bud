package factory

import (
	"errors"
	"fmt"
	"strings"
)

const (
	WireValueSelf  WireValue = "self"
	WireValueAuto  WireValue = "auto"
	WireValueType  WireValue = "type"
	WireValueName  WireValue = "name"
	WireValueValue WireValue = "value"
)

var ErrInvalidWireValue = errors.New("not a valid WireValue")

var _WireValueNameMap = map[string]WireValue{
	"self":  WireValueSelf,
	"auto":  WireValueAuto,
	"type":  WireValueType,
	"name":  WireValueName,
	"value": WireValueValue,
}

func (x WireValue) IsValid() bool {
	_, ok := _WireValueNameMap[string(x)]
	return ok
}

func WireValueValues() []WireValue {
	return []WireValue{
		WireValueSelf,
		WireValueAuto,
		WireValueType,
		WireValueName,
		WireValueValue,
	}
}

func (x WireValue) Name() string {
	if v, ok := _WireValueNameMap[string(x)]; ok {
		return string(v)
	}
	panic(ErrInvalidWireValue)
}

func (x WireValue) String() string {
	if v, ok := _WireValueNameMap[string(x)]; ok {
		return string(v)
	}
	return fmt.Sprintf("WireValue(%s)", string(x))
}

func ParseWireValue(value string) (WireValue, error) {
	if x, ok := _WireValueNameMap[value]; ok {
		return x, nil
	}
	if x, ok := _WireValueNameMap[strings.ToLower(value)]; ok {
		return x, nil
	}
	return "", fmt.Errorf("%s is %w", value, ErrInvalidWireValue)
}

func (x WireValue) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

func (x *WireValue) UnmarshalText(text []byte) error {
	val, err := ParseWireValue(string(text))
	if err != nil {
		return err
	}
	*x = val
	return nil
}
