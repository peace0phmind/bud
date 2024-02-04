package enum

import (
	"errors"
	"fmt"
)

const (
	// SuffixTestSomeItem is a SuffixTest of type some_item.
	SuffixTestSomeItem SuffixTest = "some_item"
)

var ErrInvalidSuffixTest = errors.New("not a valid SuffixTest")

// Name is the attribute of SuffixTest.
func (x SuffixTest) Name() string {
	if v, ok := _SuffixTestNameMap[string(x)]; ok {
		return string(v)
	}
	panic(ErrInvalidSuffixTest)
}

// Value is the attribute of SuffixTest.
func (x SuffixTest) Value() string {
	if x.IsValid() {
		return string(x)
	}
	panic(ErrInvalidSuffixTest)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x SuffixTest) IsValid() bool {
	_, ok := _SuffixTestNameMap[string(x)]
	return ok
}

// String implements the Stringer interface.
func (x SuffixTest) String() string {
	return x.Name()
}

var _SuffixTestNameMap = map[string]SuffixTest{
	"some_item": SuffixTestSomeItem,
}

// ParseSuffixTest converts a string to a SuffixTest.
func ParseSuffixTest(value string) (SuffixTest, error) {
	if x, ok := _SuffixTestNameMap[value]; ok {
		return x, nil
	}
	return "", fmt.Errorf("%s is %w", value, ErrInvalidSuffixTest)
}
