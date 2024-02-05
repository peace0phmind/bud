package enum

import (
	"errors"
	"fmt"
)

const (
	// TestOnlyEnumAbcdx is a TestOnlyEnum of type ABCD (x).
	TestOnlyEnumAbcdx TestOnlyEnum = "ABCD (x)"
	// TestOnlyEnumEfghy is a TestOnlyEnum of type EFGH (y).
	TestOnlyEnumEfghy TestOnlyEnum = "EFGH (y)"
)

var ErrInvalidTestOnlyEnum = errors.New("not a valid TestOnlyEnum")

var _TestOnlyEnumNameMap = map[string]TestOnlyEnum{
	"ABCD (x)": TestOnlyEnumAbcdx,
	"EFGH (y)": TestOnlyEnumEfghy,
}

// Name is the attribute of TestOnlyEnum.
func (x TestOnlyEnum) Name() string {
	if v, ok := _TestOnlyEnumNameMap[string(x)]; ok {
		return string(v)
	}
	panic(ErrInvalidTestOnlyEnum)
}

// Val is the attribute of TestOnlyEnum.
func (x TestOnlyEnum) Val() string {
	if x.IsValid() {
		return string(x)
	}
	panic(ErrInvalidTestOnlyEnum)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x TestOnlyEnum) IsValid() bool {
	_, ok := _TestOnlyEnumNameMap[string(x)]
	return ok
}

// String implements the Stringer interface.
func (x TestOnlyEnum) String() string {
	return x.Name()
}

// ParseTestOnlyEnum converts a string to a TestOnlyEnum.
func ParseTestOnlyEnum(value string) (TestOnlyEnum, error) {
	if x, ok := _TestOnlyEnumNameMap[value]; ok {
		return x, nil
	}
	return "", fmt.Errorf("%s is %w", value, ErrInvalidTestOnlyEnum)
}
