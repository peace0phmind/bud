package enum

import (
	"errors"
	"fmt"
)

const (
	// ForceUpperTypeDataSwap is a ForceUpperType of type DATASWAP.
	ForceUpperTypeDataSwap ForceUpperType = iota
	// ForceUpperTypeBootNode is a ForceUpperType of type BOOTNODE.
	ForceUpperTypeBootNode
)

var ErrInvalidForceUpperType = errors.New("not a valid ForceUpperType")

var _ForceUpperTypeName = "DATASWAPBOOTNODE"

var _ForceUpperTypeMapName = map[ForceUpperType]string{
	ForceUpperTypeDataSwap: _ForceUpperTypeName[0:8],
	ForceUpperTypeBootNode: _ForceUpperTypeName[8:16],
}

// Name is the attribute of ForceUpperType.
func (x ForceUpperType) Name() string {
	if v, ok := _ForceUpperTypeMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("ForceUpperType(%d).Name", x)
}

// Val is the attribute of ForceUpperType.
func (x ForceUpperType) Val() int {
	return int(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ForceUpperType) IsValid() bool {
	_, ok := _ForceUpperTypeMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x ForceUpperType) String() string {
	return x.Name()
}

var _ForceUpperTypeNameMap = map[string]ForceUpperType{
	_ForceUpperTypeName[0:8]:  ForceUpperTypeDataSwap,
	_ForceUpperTypeName[8:16]: ForceUpperTypeBootNode,
}

// ParseForceUpperType converts a string to a ForceUpperType.
func ParseForceUpperType(value string) (ForceUpperType, error) {
	if x, ok := _ForceUpperTypeNameMap[value]; ok {
		return x, nil
	}
	return ForceUpperType(0), fmt.Errorf("%s is %w", value, ErrInvalidForceUpperType)
}
