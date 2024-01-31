package enum

import (
	"errors"
	"fmt"
)

const (
	// ForceLowerTypeDataSwap is a ForceLowerType of type dataswap.
	ForceLowerTypeDataSwap ForceLowerType = iota
	// ForceLowerTypeBootNode is a ForceLowerType of type bootnode.
	ForceLowerTypeBootNode
)

var ErrInvalidForceLowerType = errors.New("not a valid ForceLowerType")

var _ForceLowerTypeName = "dataswapbootnode"

var _ForceLowerTypeMapName = map[ForceLowerType]string{
	ForceLowerTypeDataSwap: _ForceLowerTypeName[0:8],
	ForceLowerTypeBootNode: _ForceLowerTypeName[8:16],
}

func (x ForceLowerType) IsValid() bool {
	_, ok := _ForceLowerTypeMapName[x]
	return ok
}

func (x ForceLowerType) Name() string {
	if v, ok := _ForceLowerTypeMapName[x]; ok {
		return v
	}
	panic(ErrInvalidForceLowerType)
}

func (x ForceLowerType) String() string {
	if v, ok := _ForceLowerTypeMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("ForceLowerType(%d)", x)
}

var _ForceLowerTypeNameMap = map[string]ForceLowerType{
	_ForceLowerTypeName[0:8]:  ForceLowerTypeDataSwap,
	_ForceLowerTypeName[8:16]: ForceLowerTypeBootNode,
}

func ParseForceLowerType(value string) (ForceLowerType, error) {
	if x, ok := _ForceLowerTypeNameMap[value]; ok {
		return x, nil
	}
	return ForceLowerType(0), fmt.Errorf("%s is %w", value, ErrInvalidForceLowerType)
}