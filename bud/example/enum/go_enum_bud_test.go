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

func (x TestOnlyEnum) IsValid() bool {
	_, ok := _TestOnlyEnumNameMap[string(x)]
	return ok
}

func (x TestOnlyEnum) Name() string {
	if v, ok := _TestOnlyEnumNameMap[string(x)]; ok {
		return string(v)
	}
	panic(ErrInvalidTestOnlyEnum)
}

func (x TestOnlyEnum) String() string {
	if v, ok := _TestOnlyEnumNameMap[string(x)]; ok {
		return string(v)
	}
	return fmt.Sprintf("TestOnlyEnum(%s)", string(x))
}

func ParseTestOnlyEnum(value string) (TestOnlyEnum, error) {
	if x, ok := _TestOnlyEnumNameMap[value]; ok {
		return x, nil
	}
	return "", fmt.Errorf("%s is %w", value, ErrInvalidTestOnlyEnum)
}
