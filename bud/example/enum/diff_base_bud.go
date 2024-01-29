package enum

import (
	"errors"
	"fmt"
)

const (
	// DiffBaseB3 is a DiffBase of type b3.
	DiffBaseB3 DiffBase = 3
	// DiffBaseB4 is a DiffBase of type b4.
	DiffBaseB4 DiffBase = 4
	// DiffBaseB5 is a DiffBase of type b5.
	DiffBaseB5 DiffBase = 5
	// DiffBaseB6 is a DiffBase of type b6.
	DiffBaseB6 DiffBase = 6
	// DiffBaseB7 is a DiffBase of type b7.
	DiffBaseB7 DiffBase = 7
	// DiffBaseB8 is a DiffBase of type b8.
	DiffBaseB8 DiffBase = 8
	// DiffBaseB9 is a DiffBase of type b9.
	DiffBaseB9 DiffBase = 9
	// DiffBaseB10 is a DiffBase of type b10.
	DiffBaseB10 DiffBase = 11
	// DiffBaseB11 is a DiffBase of type b11.
	DiffBaseB11 DiffBase = 43
)

var ErrInvalidDiffBase = errors.New("not a valid DiffBase")

var _DiffBaseName = "b3b4b5b6b7b8b9b10b11"

var _DiffBaseMapName = map[DiffBase]string{
	DiffBaseB3:  _DiffBaseName[0:2],
	DiffBaseB4:  _DiffBaseName[2:4],
	DiffBaseB5:  _DiffBaseName[4:6],
	DiffBaseB6:  _DiffBaseName[6:8],
	DiffBaseB7:  _DiffBaseName[8:10],
	DiffBaseB8:  _DiffBaseName[10:12],
	DiffBaseB9:  _DiffBaseName[12:14],
	DiffBaseB10: _DiffBaseName[14:17],
	DiffBaseB11: _DiffBaseName[17:20],
}

func (x DiffBase) Name() string {
	if result, ok := _DiffBaseMapName[x]; ok {
		return result
	}
	panic(ErrInvalidDiffBase)
}

func (x DiffBase) IsValid() bool {
	_, ok := _DiffBaseMapName[x]
	return ok
}

func (x DiffBase) String() string {
	if str, ok := _DiffBaseMapName[x]; ok {
		return str
	}
	return fmt.Sprintf("DiffBase(%d)", x)
}

var _DiffBaseNameMap = map[string]DiffBase{
	_DiffBaseName[0:2]:   DiffBaseB3,
	_DiffBaseName[2:4]:   DiffBaseB4,
	_DiffBaseName[4:6]:   DiffBaseB5,
	_DiffBaseName[6:8]:   DiffBaseB6,
	_DiffBaseName[8:10]:  DiffBaseB7,
	_DiffBaseName[10:12]: DiffBaseB8,
	_DiffBaseName[12:14]: DiffBaseB9,
	_DiffBaseName[14:17]: DiffBaseB10,
	_DiffBaseName[17:20]: DiffBaseB11,
}

func ParseDiffBase(name string) (DiffBase, error) {
	if x, ok := _DiffBaseNameMap[name]; ok {
		return x, nil
	}
	return DiffBase(0), fmt.Errorf("%s is %w", name, ErrInvalidDiffBase)
}
