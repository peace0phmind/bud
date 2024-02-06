package enum

import (
	"fmt"
	"strings"
)

const (
	// MakeToyota is a Make of type Toyota.
	MakeToyota Make = iota
	// Skipped value.
	_
	// MakeChevy is a Make of type Chevy.
	MakeChevy
	// Skipped value.
	_
	// MakeFord is a Make of type Ford.
	MakeFord
	// Skipped value.
	_
	// MakeTesla is a Make of type Tesla.
	MakeTesla
	// Skipped value.
	_
	// MakeHyundai is a Make of type Hyundai.
	MakeHyundai
	// Skipped value.
	_
	// MakeNissan is a Make of type Nissan.
	MakeNissan
	// Skipped value.
	_
	// MakeJaguar is a Make of type Jaguar.
	MakeJaguar
	// Skipped value.
	_
	// MakeAudi is a Make of type Audi.
	MakeAudi
	// Skipped value.
	_
	// MakeBmw is a Make of type BMW.
	MakeBmw
	// Skipped value.
	_
	// MakeMercedesBenz is a Make of type Mercedes_Benz.
	MakeMercedesBenz
	// Skipped value.
	_
	// MakeVolkswagon is a Make of type Volkswagon.
	MakeVolkswagon
)

const (
	// NoZerosStart is a NoZeros of type start.
	NoZerosStart NoZeros = 20
	// NoZerosMiddle is a NoZeros of type middle.
	NoZerosMiddle NoZeros = 21
	// NoZerosEnd is a NoZeros of type end.
	NoZerosEnd NoZeros = 22
	// NoZerosPs is a NoZeros of type ps.
	NoZerosPs NoZeros = 23
	// NoZerosPps is a NoZeros of type pps.
	NoZerosPps NoZeros = 24
	// NoZerosPpps is a NoZeros of type ppps.
	NoZerosPpps NoZeros = 25
)

var ErrInvalidMake = fmt.Errorf("not a valid Make, try [%s]", strings.Join(_MakeNames, ", "))

var _MakeName = "ToyotaChevyFordTeslaHyundaiNissanJaguarAudiBMWMercedes_BenzVolkswagon"

var _MakeMapName = map[Make]string{
	MakeToyota:       _MakeName[0:6],
	MakeChevy:        _MakeName[6:11],
	MakeFord:         _MakeName[11:15],
	MakeTesla:        _MakeName[15:20],
	MakeHyundai:      _MakeName[20:27],
	MakeNissan:       _MakeName[27:33],
	MakeJaguar:       _MakeName[33:39],
	MakeAudi:         _MakeName[39:43],
	MakeBmw:          _MakeName[43:46],
	MakeMercedesBenz: _MakeName[46:59],
	MakeVolkswagon:   _MakeName[59:69],
}

// Name is the attribute of Make.
func (x Make) Name() string {
	if v, ok := _MakeMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("Make(%d).Name", x)
}

// Val is the attribute of Make.
func (x Make) Val() int32 {
	return int32(x)
}

var _MakeValues = []Make{
	MakeToyota,
	MakeChevy,
	MakeFord,
	MakeTesla,
	MakeHyundai,
	MakeNissan,
	MakeJaguar,
	MakeAudi,
	MakeBmw,
	MakeMercedesBenz,
	MakeVolkswagon,
}

// MakeValues returns a list of the values of Make
func MakeValues() []Make {
	return _MakeValues
}

var _MakeNames = []string{
	_MakeName[0:6],
	_MakeName[6:11],
	_MakeName[11:15],
	_MakeName[15:20],
	_MakeName[20:27],
	_MakeName[27:33],
	_MakeName[33:39],
	_MakeName[39:43],
	_MakeName[43:46],
	_MakeName[46:59],
	_MakeName[59:69],
}

// MakeNames returns a list of the names of Make
func MakeNames() []string {
	return _MakeNames
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Make) IsValid() bool {
	_, ok := _MakeMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x Make) String() string {
	return x.Name()
}

var _MakeNameMap = map[string]Make{
	_MakeName[0:6]:                    MakeToyota,
	strings.ToLower(_MakeName[0:6]):   MakeToyota,
	_MakeName[6:11]:                   MakeChevy,
	strings.ToLower(_MakeName[6:11]):  MakeChevy,
	_MakeName[11:15]:                  MakeFord,
	strings.ToLower(_MakeName[11:15]): MakeFord,
	_MakeName[15:20]:                  MakeTesla,
	strings.ToLower(_MakeName[15:20]): MakeTesla,
	_MakeName[20:27]:                  MakeHyundai,
	strings.ToLower(_MakeName[20:27]): MakeHyundai,
	_MakeName[27:33]:                  MakeNissan,
	strings.ToLower(_MakeName[27:33]): MakeNissan,
	_MakeName[33:39]:                  MakeJaguar,
	strings.ToLower(_MakeName[33:39]): MakeJaguar,
	_MakeName[39:43]:                  MakeAudi,
	strings.ToLower(_MakeName[39:43]): MakeAudi,
	_MakeName[43:46]:                  MakeBmw,
	strings.ToLower(_MakeName[43:46]): MakeBmw,
	_MakeName[46:59]:                  MakeMercedesBenz,
	strings.ToLower(_MakeName[46:59]): MakeMercedesBenz,
	_MakeName[59:69]:                  MakeVolkswagon,
	strings.ToLower(_MakeName[59:69]): MakeVolkswagon,
}

// ParseMake converts a string to a Make.
func ParseMake(value string) (Make, error) {
	if x, ok := _MakeNameMap[value]; ok {
		return x, nil
	}
	if x, ok := _MakeNameMap[strings.ToLower(value)]; ok {
		return x, nil
	}
	return Make(0), fmt.Errorf("%s is %w", value, ErrInvalidMake)
}

// Set implements the Golang flag.Value interface func.
func (x *Make) Set(value string) error {
	v, err := ParseMake(value)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x Make) Get() any {
	return x
}

// MarshalText implements the text marshaller method.
func (x Make) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Make) UnmarshalText(text []byte) error {
	val, err := ParseMake(string(text))
	if err != nil {
		return err
	}
	*x = val
	return nil
}

var ErrInvalidNoZeros = fmt.Errorf("not a valid NoZeros, try [%s]", strings.Join(_NoZerosNames, ", "))

var _NoZerosName = "startmiddleendpsppsppps"

var _NoZerosMapName = map[NoZeros]string{
	NoZerosStart:  _NoZerosName[0:5],
	NoZerosMiddle: _NoZerosName[5:11],
	NoZerosEnd:    _NoZerosName[11:14],
	NoZerosPs:     _NoZerosName[14:16],
	NoZerosPps:    _NoZerosName[16:19],
	NoZerosPpps:   _NoZerosName[19:23],
}

// Name is the attribute of NoZeros.
func (x NoZeros) Name() string {
	if v, ok := _NoZerosMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("NoZeros(%d).Name", x)
}

// Val is the attribute of NoZeros.
func (x NoZeros) Val() int32 {
	return int32(x)
}

var _NoZerosValues = []NoZeros{
	NoZerosStart,
	NoZerosMiddle,
	NoZerosEnd,
	NoZerosPs,
	NoZerosPps,
	NoZerosPpps,
}

// NoZerosValues returns a list of the values of NoZeros
func NoZerosValues() []NoZeros {
	return _NoZerosValues
}

var _NoZerosNames = []string{
	_NoZerosName[0:5],
	_NoZerosName[5:11],
	_NoZerosName[11:14],
	_NoZerosName[14:16],
	_NoZerosName[16:19],
	_NoZerosName[19:23],
}

// NoZerosNames returns a list of the names of NoZeros
func NoZerosNames() []string {
	return _NoZerosNames
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x NoZeros) IsValid() bool {
	_, ok := _NoZerosMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x NoZeros) String() string {
	return x.Name()
}

var _NoZerosNameMap = map[string]NoZeros{
	_NoZerosName[0:5]:   NoZerosStart,
	_NoZerosName[5:11]:  NoZerosMiddle,
	_NoZerosName[11:14]: NoZerosEnd,
	_NoZerosName[14:16]: NoZerosPs,
	_NoZerosName[16:19]: NoZerosPps,
	_NoZerosName[19:23]: NoZerosPpps,
}

// ParseNoZeros converts a string to a NoZeros.
func ParseNoZeros(value string) (NoZeros, error) {
	if x, ok := _NoZerosNameMap[value]; ok {
		return x, nil
	}
	if x, ok := _NoZerosNameMap[strings.ToLower(value)]; ok {
		return x, nil
	}
	return NoZeros(0), fmt.Errorf("%s is %w", value, ErrInvalidNoZeros)
}

// Set implements the Golang flag.Value interface func.
func (x *NoZeros) Set(value string) error {
	v, err := ParseNoZeros(value)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x NoZeros) Get() any {
	return x
}

// MarshalText implements the text marshaller method.
func (x NoZeros) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *NoZeros) UnmarshalText(text []byte) error {
	val, err := ParseNoZeros(string(text))
	if err != nil {
		return err
	}
	*x = val
	return nil
}
