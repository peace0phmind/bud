package enum

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// CommentedValue1 is a Commented of type value1.
	CommentedValue1 Commented = iota // Commented value 1
	// CommentedValue2 is a Commented of type value2.
	CommentedValue2
	// CommentedValue3 is a Commented of type value3.
	CommentedValue3 // Commented value 3
)

const (
	// Skipped value.
	_ ComplexCommented = iota // Placeholder with a ','  in it. (for harder testing)
	// ComplexCommentedValue1 is a ComplexCommented of type value1.
	ComplexCommentedValue1 // Commented value 1
	// ComplexCommentedValue2 is a ComplexCommented of type value2.
	ComplexCommentedValue2
	// Skipped value.
	_
	// Skipped value.
	_
	// ComplexCommentedValue3 is a ComplexCommented of type value3.
	ComplexCommentedValue3 // Commented value 3
)

var ErrInvalidCommented = errors.New("not a valid Commented")

var _CommentedName = "value1value2value3"

var _CommentedMapName = map[Commented]string{
	CommentedValue1: _CommentedName[0:6],
	CommentedValue2: _CommentedName[6:12],
	CommentedValue3: _CommentedName[12:18],
}

// Name is the attribute of Commented.
func (x Commented) Name() string {
	if v, ok := _CommentedMapName[x]; ok {
		return v
	}
	panic(ErrInvalidCommented)
}

// Val is the attribute of Commented.
func (x Commented) Val() int {
	if x.IsValid() {
		return int(x)
	}
	panic(ErrInvalidCommented)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Commented) IsValid() bool {
	_, ok := _CommentedMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x Commented) String() string {
	return x.Name()
}

var _CommentedNameMap = map[string]Commented{
	_CommentedName[0:6]:   CommentedValue1,
	_CommentedName[6:12]:  CommentedValue2,
	_CommentedName[12:18]: CommentedValue3,
}

// ParseCommented converts a string to a Commented.
func ParseCommented(value string) (Commented, error) {
	if x, ok := _CommentedNameMap[value]; ok {
		return x, nil
	}
	if x, ok := _CommentedNameMap[strings.ToLower(value)]; ok {
		return x, nil
	}
	return Commented(0), fmt.Errorf("%s is %w", value, ErrInvalidCommented)
}

// MarshalText implements the text marshaller method.
func (x Commented) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Commented) UnmarshalText(text []byte) error {
	val, err := ParseCommented(string(text))
	if err != nil {
		return err
	}
	*x = val
	return nil
}

var ErrInvalidComplexCommented = errors.New("not a valid ComplexCommented")

var _ComplexCommentedName = "value1value2value3"

var _ComplexCommentedMapName = map[ComplexCommented]string{
	ComplexCommentedValue1: _ComplexCommentedName[0:6],
	ComplexCommentedValue2: _ComplexCommentedName[6:12],
	ComplexCommentedValue3: _ComplexCommentedName[12:18],
}

// Name is the attribute of ComplexCommented.
func (x ComplexCommented) Name() string {
	if v, ok := _ComplexCommentedMapName[x]; ok {
		return v
	}
	panic(ErrInvalidComplexCommented)
}

// Val is the attribute of ComplexCommented.
func (x ComplexCommented) Val() int {
	if x.IsValid() {
		return int(x)
	}
	panic(ErrInvalidComplexCommented)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ComplexCommented) IsValid() bool {
	_, ok := _ComplexCommentedMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x ComplexCommented) String() string {
	return x.Name()
}

var _ComplexCommentedNameMap = map[string]ComplexCommented{
	_ComplexCommentedName[0:6]:   ComplexCommentedValue1,
	_ComplexCommentedName[6:12]:  ComplexCommentedValue2,
	_ComplexCommentedName[12:18]: ComplexCommentedValue3,
}

// ParseComplexCommented converts a string to a ComplexCommented.
func ParseComplexCommented(value string) (ComplexCommented, error) {
	if x, ok := _ComplexCommentedNameMap[value]; ok {
		return x, nil
	}
	if x, ok := _ComplexCommentedNameMap[strings.ToLower(value)]; ok {
		return x, nil
	}
	return ComplexCommented(0), fmt.Errorf("%s is %w", value, ErrInvalidComplexCommented)
}

// MarshalText implements the text marshaller method.
func (x ComplexCommented) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ComplexCommented) UnmarshalText(text []byte) error {
	val, err := ParseComplexCommented(string(text))
	if err != nil {
		return err
	}
	*x = val
	return nil
}
