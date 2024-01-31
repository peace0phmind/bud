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

func (x ForceUpperType) IsValid() bool {
	_, ok := _ForceUpperTypeMapName[x]
	return ok
}

func (x ForceUpperType) Name() string {
	if v, ok := _ForceUpperTypeMapName[x]; ok {
		return v
	}
	panic(ErrInvalidForceUpperType)
}

func (x ForceUpperType) String() string {
	if v, ok := _ForceUpperTypeMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("ForceUpperType(%d)", x)
}

var _ForceUpperTypeNameMap = map[string]ForceUpperType{
	_ForceUpperTypeName[0:8]:  ForceUpperTypeDataSwap,
	_ForceUpperTypeName[8:16]: ForceUpperTypeBootNode,
}

func ParseForceUpperType(value string) (ForceUpperType, error) {
	if x, ok := _ForceUpperTypeNameMap[value]; ok {
		return x, nil
	}
	return ForceUpperType(0), fmt.Errorf("%s is %w", value, ErrInvalidForceUpperType)
}
