package enum

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// ColorBlack is a Color of type Black.
	ColorBlack Color = 0
	// ColorWhite is a Color of type White.
	ColorWhite Color = 1
	// ColorRed is a Color of type Red.
	ColorRed Color = 2
	// ColorGreen is a Color of type Green.
	ColorGreen Color = 33 // Green starts with 33
	// ColorBlue is a Color of type Blue.
	ColorBlue Color = 34
	// ColorGrey is a Color of type grey.
	ColorGrey Color = 45
	// Skipped value.
	_
	// Skipped value.
	_
	// ColorYellow is a Color of type yellow.
	ColorYellow Color = 48
	// ColorBlueGreen is a Color of type blue-green.
	ColorBlueGreen Color = 49
	// ColorRedOrange is a Color of type red-orange.
	ColorRedOrange Color = 50
	// ColorYellowGreen is a Color of type yellow_green.
	ColorYellowGreen Color = 51
	// ColorRedOrangeBlue is a Color of type red-orange-blue.
	ColorRedOrangeBlue Color = 52
)

var ErrInvalidColor = errors.New("not a valid Color")

var _ColorName = "BlackWhiteRedGreenBluegreyyellowblue-greenred-orangeyellow_greenred-orange-blue"

var _ColorMapName = map[Color]string{
	ColorBlack:         _ColorName[0:5],
	ColorWhite:         _ColorName[5:10],
	ColorRed:           _ColorName[10:13],
	ColorGreen:         _ColorName[13:18],
	ColorBlue:          _ColorName[18:22],
	ColorGrey:          _ColorName[22:26],
	ColorYellow:        _ColorName[26:32],
	ColorBlueGreen:     _ColorName[32:42],
	ColorRedOrange:     _ColorName[42:52],
	ColorYellowGreen:   _ColorName[52:64],
	ColorRedOrangeBlue: _ColorName[64:79],
}

// Name is the attribute of Color.
func (x Color) Name() string {
	if v, ok := _ColorMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("Color(%d).Name", x)
}

// Val is the attribute of Color.
func (x Color) Val() int {
	return int(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Color) IsValid() bool {
	_, ok := _ColorMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x Color) String() string {
	return x.Name()
}

var _ColorNameMap = map[string]Color{
	_ColorName[0:5]:                    ColorBlack,
	strings.ToLower(_ColorName[0:5]):   ColorBlack,
	_ColorName[5:10]:                   ColorWhite,
	strings.ToLower(_ColorName[5:10]):  ColorWhite,
	_ColorName[10:13]:                  ColorRed,
	strings.ToLower(_ColorName[10:13]): ColorRed,
	_ColorName[13:18]:                  ColorGreen,
	strings.ToLower(_ColorName[13:18]): ColorGreen,
	_ColorName[18:22]:                  ColorBlue,
	strings.ToLower(_ColorName[18:22]): ColorBlue,
	_ColorName[22:26]:                  ColorGrey,
	_ColorName[26:32]:                  ColorYellow,
	_ColorName[32:42]:                  ColorBlueGreen,
	_ColorName[42:52]:                  ColorRedOrange,
	_ColorName[52:64]:                  ColorYellowGreen,
	_ColorName[64:79]:                  ColorRedOrangeBlue,
}

// ParseColor converts a string to a Color.
func ParseColor(value string) (Color, error) {
	if x, ok := _ColorNameMap[value]; ok {
		return x, nil
	}
	if x, ok := _ColorNameMap[strings.ToLower(value)]; ok {
		return x, nil
	}
	return Color(0), fmt.Errorf("%s is %w", value, ErrInvalidColor)
}

// MustParseColor converts a string to a Color, and panics if is not valid.
func MustParseColor(value string) Color {
	val, err := ParseColor(value)
	if err != nil {
		panic(err)
	}
	return val
}

func (x Color) Ptr() *Color {
	return &x
}

// MarshalText implements the text marshaller method.
func (x Color) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Color) UnmarshalText(text []byte) error {
	val, err := ParseColor(string(text))
	if err != nil {
		return err
	}
	*x = val
	return nil
}
