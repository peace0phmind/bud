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

func (x Commented) Name() string {
	if result, ok := _CommentedMapName[x]; ok {
		return result
	}
	panic(ErrInvalidCommented)
}

func (x Commented) IsValid() bool {
	_, ok := _CommentedMapName[x]
	return ok
}

func (x Commented) String() string {
	if str, ok := _CommentedMapName[x]; ok {
		return str
	}
	return fmt.Sprintf("Commented(%d)", x)
}

var _CommentedNameMap = map[string]Commented{
	_CommentedName[0:6]:   CommentedValue1,
	_CommentedName[6:12]:  CommentedValue2,
	_CommentedName[12:18]: CommentedValue3,
}

func ParseCommented(name string) (Commented, error) {
	if x, ok := _CommentedNameMap[name]; ok {
		return x, nil
	}
	if x, ok := _CommentedNameMap[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Commented(0), fmt.Errorf("%s is %w", name, ErrInvalidCommented)
}

func (x Commented) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

func (x *Commented) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseCommented(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var ErrInvalidComplexCommented = errors.New("not a valid ComplexCommented")

var _ComplexCommentedName = "value1value2value3"

var _ComplexCommentedMapName = map[ComplexCommented]string{
	ComplexCommentedValue1: _ComplexCommentedName[0:6],
	ComplexCommentedValue2: _ComplexCommentedName[6:12],
	ComplexCommentedValue3: _ComplexCommentedName[12:18],
}

func (x ComplexCommented) Name() string {
	if result, ok := _ComplexCommentedMapName[x]; ok {
		return result
	}
	panic(ErrInvalidComplexCommented)
}

func (x ComplexCommented) IsValid() bool {
	_, ok := _ComplexCommentedMapName[x]
	return ok
}

func (x ComplexCommented) String() string {
	if str, ok := _ComplexCommentedMapName[x]; ok {
		return str
	}
	return fmt.Sprintf("ComplexCommented(%d)", x)
}

var _ComplexCommentedNameMap = map[string]ComplexCommented{
	_ComplexCommentedName[0:6]:   ComplexCommentedValue1,
	_ComplexCommentedName[6:12]:  ComplexCommentedValue2,
	_ComplexCommentedName[12:18]: ComplexCommentedValue3,
}

func ParseComplexCommented(name string) (ComplexCommented, error) {
	if x, ok := _ComplexCommentedNameMap[name]; ok {
		return x, nil
	}
	if x, ok := _ComplexCommentedNameMap[strings.ToLower(name)]; ok {
		return x, nil
	}
	return ComplexCommented(0), fmt.Errorf("%s is %w", name, ErrInvalidComplexCommented)
}

func (x ComplexCommented) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

func (x *ComplexCommented) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseComplexCommented(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
