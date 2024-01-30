package enum

import (
	"errors"
	"fmt"
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

var ErrInvalidMake = errors.New("not a valid Make")

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

func (x Make) IsValid() bool {
	_, ok := _MakeMapName[x]
	return ok
}

func (x Make) Name() string {
	if v, ok := _MakeMapName[x]; ok {
		return v
	}
	panic(ErrInvalidMake)
}

func (x Make) String() string {
	if v, ok := _MakeMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("Make(%d)Name", x)
}

var _MakeNameMap = map[string]Make{
	_MakeName[0:6]:   MakeToyota,
	_MakeName[6:11]:  MakeChevy,
	_MakeName[11:15]: MakeFord,
	_MakeName[15:20]: MakeTesla,
	_MakeName[20:27]: MakeHyundai,
	_MakeName[27:33]: MakeNissan,
	_MakeName[33:39]: MakeJaguar,
	_MakeName[39:43]: MakeAudi,
	_MakeName[43:46]: MakeBmw,
	_MakeName[46:59]: MakeMercedesBenz,
	_MakeName[59:69]: MakeVolkswagon,
}

func ParseMake(value string) (Make, error) {
	if x, ok := _MakeNameMap[value]; ok {
		return x, nil
	}

	return Make(0), fmt.Errorf("%s is %w", value, ErrInvalidMake)
}

var ErrInvalidNoZeros = errors.New("not a valid NoZeros")

var _NoZerosName = "startmiddleendpsppsppps"

var _NoZerosMapName = map[NoZeros]string{
	NoZerosStart:  _NoZerosName[0:5],
	NoZerosMiddle: _NoZerosName[5:11],
	NoZerosEnd:    _NoZerosName[11:14],
	NoZerosPs:     _NoZerosName[14:16],
	NoZerosPps:    _NoZerosName[16:19],
	NoZerosPpps:   _NoZerosName[19:23],
}

func (x NoZeros) IsValid() bool {
	_, ok := _NoZerosMapName[x]
	return ok
}

func (x NoZeros) Name() string {
	if v, ok := _NoZerosMapName[x]; ok {
		return v
	}
	panic(ErrInvalidNoZeros)
}

func (x NoZeros) String() string {
	if v, ok := _NoZerosMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("NoZeros(%d)Name", x)
}

var _NoZerosNameMap = map[string]NoZeros{
	_NoZerosName[0:5]:   NoZerosStart,
	_NoZerosName[5:11]:  NoZerosMiddle,
	_NoZerosName[11:14]: NoZerosEnd,
	_NoZerosName[14:16]: NoZerosPs,
	_NoZerosName[16:19]: NoZerosPps,
	_NoZerosName[19:23]: NoZerosPpps,
}

func ParseNoZeros(value string) (NoZeros, error) {
	if x, ok := _NoZerosNameMap[value]; ok {
		return x, nil
	}

	return NoZeros(0), fmt.Errorf("%s is %w", value, ErrInvalidNoZeros)
}
