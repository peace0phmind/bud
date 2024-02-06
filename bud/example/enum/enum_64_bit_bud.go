package enum

import (
	"errors"
	"fmt"
)

const (
	// Enum64bitUnkno is an Enum64bit of type Unkno.
	Enum64bitUnkno Enum64bit = 0
	// Enum64bitE2P15 is an Enum64bit of type E2P15.
	Enum64bitE2P15 Enum64bit = 32768
	// Enum64bitE2P16 is an Enum64bit of type E2P16.
	Enum64bitE2P16 Enum64bit = 65536
	// Enum64bitE2P17 is an Enum64bit of type E2P17.
	Enum64bitE2P17 Enum64bit = 131072
	// Enum64bitE2P18 is an Enum64bit of type E2P18.
	Enum64bitE2P18 Enum64bit = 262144
	// Enum64bitE2P19 is an Enum64bit of type E2P19.
	Enum64bitE2P19 Enum64bit = 524288
	// Enum64bitE2P20 is an Enum64bit of type E2P20.
	Enum64bitE2P20 Enum64bit = 1048576
	// Enum64bitE2P21 is an Enum64bit of type E2P21.
	Enum64bitE2P21 Enum64bit = 2097152
	// Enum64bitE2P22 is an Enum64bit of type E2P22.
	Enum64bitE2P22 Enum64bit = 33554432
	// Enum64bitE2P23 is an Enum64bit of type E2P23.
	Enum64bitE2P23 Enum64bit = 67108864
	// Enum64bitE2P28 is an Enum64bit of type E2P28.
	Enum64bitE2P28 Enum64bit = 536870912
	// Enum64bitE2P30 is an Enum64bit of type E2P30.
	Enum64bitE2P30 Enum64bit = 1073741824
	// Enum64bitE2P31 is an Enum64bit of type E2P31.
	Enum64bitE2P31 Enum64bit = 2147483648
	// Enum64bitE2P32 is an Enum64bit of type E2P32.
	Enum64bitE2P32 Enum64bit = 4294967296
	// Enum64bitE2P33 is an Enum64bit of type E2P33.
	Enum64bitE2P33 Enum64bit = 8454967296
	// Enum64bitE2P63 is an Enum64bit of type E2P63.
	Enum64bitE2P63 Enum64bit = 18446744073709551615
)

var ErrInvalidEnum64bit = errors.New("not a valid Enum64bit")

var _Enum64bitName = "UnknoE2P15E2P16E2P17E2P18E2P19E2P20E2P21E2P22E2P23E2P28E2P30E2P31E2P32E2P33E2P63"

var _Enum64bitMapName = map[Enum64bit]string{
	Enum64bitUnkno: _Enum64bitName[0:5],
	Enum64bitE2P15: _Enum64bitName[5:10],
	Enum64bitE2P16: _Enum64bitName[10:15],
	Enum64bitE2P17: _Enum64bitName[15:20],
	Enum64bitE2P18: _Enum64bitName[20:25],
	Enum64bitE2P19: _Enum64bitName[25:30],
	Enum64bitE2P20: _Enum64bitName[30:35],
	Enum64bitE2P21: _Enum64bitName[35:40],
	Enum64bitE2P22: _Enum64bitName[40:45],
	Enum64bitE2P23: _Enum64bitName[45:50],
	Enum64bitE2P28: _Enum64bitName[50:55],
	Enum64bitE2P30: _Enum64bitName[55:60],
	Enum64bitE2P31: _Enum64bitName[60:65],
	Enum64bitE2P32: _Enum64bitName[65:70],
	Enum64bitE2P33: _Enum64bitName[70:75],
	Enum64bitE2P63: _Enum64bitName[75:80],
}

// Name is the attribute of Enum64bit.
func (x Enum64bit) Name() string {
	if v, ok := _Enum64bitMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("Enum64bit(%d).Name", x)
}

// Val is the attribute of Enum64bit.
func (x Enum64bit) Val() uint64 {
	return uint64(x)
}

var _Enum64bitValues = []Enum64bit{
	Enum64bitUnkno,
	Enum64bitE2P15,
	Enum64bitE2P16,
	Enum64bitE2P17,
	Enum64bitE2P18,
	Enum64bitE2P19,
	Enum64bitE2P20,
	Enum64bitE2P21,
	Enum64bitE2P22,
	Enum64bitE2P23,
	Enum64bitE2P28,
	Enum64bitE2P30,
	Enum64bitE2P31,
	Enum64bitE2P32,
	Enum64bitE2P33,
	Enum64bitE2P63,
}

// Enum64bitValues returns a list of the values of Enum64bit
func Enum64bitValues() []Enum64bit {
	return _Enum64bitValues
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Enum64bit) IsValid() bool {
	_, ok := _Enum64bitMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x Enum64bit) String() string {
	return x.Name()
}

var _Enum64bitNameMap = map[string]Enum64bit{
	_Enum64bitName[0:5]:   Enum64bitUnkno,
	_Enum64bitName[5:10]:  Enum64bitE2P15,
	_Enum64bitName[10:15]: Enum64bitE2P16,
	_Enum64bitName[15:20]: Enum64bitE2P17,
	_Enum64bitName[20:25]: Enum64bitE2P18,
	_Enum64bitName[25:30]: Enum64bitE2P19,
	_Enum64bitName[30:35]: Enum64bitE2P20,
	_Enum64bitName[35:40]: Enum64bitE2P21,
	_Enum64bitName[40:45]: Enum64bitE2P22,
	_Enum64bitName[45:50]: Enum64bitE2P23,
	_Enum64bitName[50:55]: Enum64bitE2P28,
	_Enum64bitName[55:60]: Enum64bitE2P30,
	_Enum64bitName[60:65]: Enum64bitE2P31,
	_Enum64bitName[65:70]: Enum64bitE2P32,
	_Enum64bitName[70:75]: Enum64bitE2P33,
	_Enum64bitName[75:80]: Enum64bitE2P63,
}

// ParseEnum64bit converts a string to an Enum64bit.
func ParseEnum64bit(value string) (Enum64bit, error) {
	if x, ok := _Enum64bitNameMap[value]; ok {
		return x, nil
	}
	return Enum64bit(0), fmt.Errorf("%s is %w", value, ErrInvalidEnum64bit)
}
