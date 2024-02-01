package enum

import (
	"errors"
	"fmt"
)

const (
	// SuffixGen is a Suffix of type gen.
	SuffixGen Suffix = "gen"
)

var ErrInvalidSuffix = errors.New("not a valid Suffix")

// Name is the attribute of Suffix.
func (x Suffix) Name() string {
	if v, ok := _SuffixNameMap[string(x)]; ok {
		return string(v)
	}
	panic(ErrInvalidSuffix)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Suffix) IsValid() bool {
	_, ok := _SuffixNameMap[string(x)]
	return ok
}

// String implements the Stringer interface.
func (x Suffix) String() string {
	if v, ok := _SuffixNameMap[string(x)]; ok {
		return string(v)
	}
	return fmt.Sprintf("Suffix(%s)", string(x))
}

var _SuffixNameMap = map[string]Suffix{
	"gen": SuffixGen,
}

// ParseSuffix converts a string to a Suffix.
func ParseSuffix(value string) (Suffix, error) {
	if x, ok := _SuffixNameMap[value]; ok {
		return x, nil
	}
	return "", fmt.Errorf("%s is %w", value, ErrInvalidSuffix)
}
