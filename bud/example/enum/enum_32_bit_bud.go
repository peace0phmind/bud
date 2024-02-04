package enum

import (
	"errors"
	"fmt"
)

const (
	// Enum32bitUnkno is an Enum32bit of type Unkno.
	Enum32bitUnkno Enum32bit = 0
	// Enum32bitE2P15 is an Enum32bit of type E2P15.
	Enum32bitE2P15 Enum32bit = 32768
	// Enum32bitE2P16 is an Enum32bit of type E2P16.
	Enum32bitE2P16 Enum32bit = 65536
	// Enum32bitE2P17 is an Enum32bit of type E2P17.
	Enum32bitE2P17 Enum32bit = 131072
	// Enum32bitE2P18 is an Enum32bit of type E2P18.
	Enum32bitE2P18 Enum32bit = 262144
	// Enum32bitE2P19 is an Enum32bit of type E2P19.
	Enum32bitE2P19 Enum32bit = 524288
	// Enum32bitE2P20 is an Enum32bit of type E2P20.
	Enum32bitE2P20 Enum32bit = 1048576
	// Enum32bitE2P21 is an Enum32bit of type E2P21.
	Enum32bitE2P21 Enum32bit = 2097152
	// Enum32bitE2P22 is an Enum32bit of type E2P22.
	Enum32bitE2P22 Enum32bit = 33554432
	// Enum32bitE2P23 is an Enum32bit of type E2P23.
	Enum32bitE2P23 Enum32bit = 67108864
	// Enum32bitE2P28 is an Enum32bit of type E2P28.
	Enum32bitE2P28 Enum32bit = 536870912
	// Enum32bitE2P30 is an Enum32bit of type E2P30.
	Enum32bitE2P30 Enum32bit = 1073741824
)

var ErrInvalidEnum32bit = errors.New("not a valid Enum32bit")

var _Enum32bitName = "UnknoE2P15E2P16E2P17E2P18E2P19E2P20E2P21E2P22E2P23E2P28E2P30"

var _Enum32bitMapName = map[Enum32bit]string{
	Enum32bitUnkno: _Enum32bitName[0:5],
	Enum32bitE2P15: _Enum32bitName[5:10],
	Enum32bitE2P16: _Enum32bitName[10:15],
	Enum32bitE2P17: _Enum32bitName[15:20],
	Enum32bitE2P18: _Enum32bitName[20:25],
	Enum32bitE2P19: _Enum32bitName[25:30],
	Enum32bitE2P20: _Enum32bitName[30:35],
	Enum32bitE2P21: _Enum32bitName[35:40],
	Enum32bitE2P22: _Enum32bitName[40:45],
	Enum32bitE2P23: _Enum32bitName[45:50],
	Enum32bitE2P28: _Enum32bitName[50:55],
	Enum32bitE2P30: _Enum32bitName[55:60],
}

// Name is the attribute of Enum32bit.
func (x Enum32bit) Name() string {
	if v, ok := _Enum32bitMapName[x]; ok {
		return v
	}
	panic(ErrInvalidEnum32bit)
}

// Value is the attribute of Enum32bit.
func (x Enum32bit) Value() uint32 {
	if x.IsValid() {
		return uint32(x)
	}
	panic(ErrInvalidEnum32bit)
}

var _Enum32bitValues = []Enum32bit{
	Enum32bitUnkno,
	Enum32bitE2P15,
	Enum32bitE2P16,
	Enum32bitE2P17,
	Enum32bitE2P18,
	Enum32bitE2P19,
	Enum32bitE2P20,
	Enum32bitE2P21,
	Enum32bitE2P22,
	Enum32bitE2P23,
	Enum32bitE2P28,
	Enum32bitE2P30,
}

// Enum32bitValues returns a list of the values of Enum32bit
func Enum32bitValues() []Enum32bit {
	return _Enum32bitValues
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Enum32bit) IsValid() bool {
	_, ok := _Enum32bitMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x Enum32bit) String() string {
	return x.Name()
}

var _Enum32bitNameMap = map[string]Enum32bit{
	_Enum32bitName[0:5]:   Enum32bitUnkno,
	_Enum32bitName[5:10]:  Enum32bitE2P15,
	_Enum32bitName[10:15]: Enum32bitE2P16,
	_Enum32bitName[15:20]: Enum32bitE2P17,
	_Enum32bitName[20:25]: Enum32bitE2P18,
	_Enum32bitName[25:30]: Enum32bitE2P19,
	_Enum32bitName[30:35]: Enum32bitE2P20,
	_Enum32bitName[35:40]: Enum32bitE2P21,
	_Enum32bitName[40:45]: Enum32bitE2P22,
	_Enum32bitName[45:50]: Enum32bitE2P23,
	_Enum32bitName[50:55]: Enum32bitE2P28,
	_Enum32bitName[55:60]: Enum32bitE2P30,
}

// ParseEnum32bit converts a string to an Enum32bit.
func ParseEnum32bit(value string) (Enum32bit, error) {
	if x, ok := _Enum32bitNameMap[value]; ok {
		return x, nil
	}
	return Enum32bit(0), fmt.Errorf("%s is %w", value, ErrInvalidEnum32bit)
}
