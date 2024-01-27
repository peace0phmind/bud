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
	// Green starts with 33
	ColorGreen Color = 33
	// ColorBlue is a Color of type Blue.
	ColorBlue Color = 34
	// ColorGrey is a Color of type grey.
	ColorGrey Color = 45
	// ColorYellow is a Color of type yellow.
	ColorYellow Color = 46
	// ColorBlueGreen is a Color of type blue-green.
	ColorBlueGreen Color = 47
	// ColorRedOrange is a Color of type red-orange.
	ColorRedOrange Color = 48
	// ColorYellowGreen is a Color of type yellow_green.
	ColorYellowGreen Color = 49
	// ColorRedOrangeBlue is a Color of type red-orange-blue.
	ColorRedOrangeBlue Color = 50
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

func (x Color) Name() string {
	if result, ok := _ColorMapName[x]; ok {
		return result
	}
	panic(ErrInvalidColor)
}

func (x Color) String() string {
	if str, ok := _ColorMapName[x]; ok {
		return str
	}
	return fmt.Sprintf("Color(%d)", x)
}

func (x Color) IsValid() bool {
	_, ok := _ColorMapName[x]
	return ok
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
	strings.ToLower(_ColorName[22:26]): ColorGrey,
	_ColorName[26:32]:                  ColorYellow,
	strings.ToLower(_ColorName[26:32]): ColorYellow,
	_ColorName[32:42]:                  ColorBlueGreen,
	strings.ToLower(_ColorName[32:42]): ColorBlueGreen,
	_ColorName[42:52]:                  ColorRedOrange,
	strings.ToLower(_ColorName[42:52]): ColorRedOrange,
	_ColorName[52:64]:                  ColorYellowGreen,
	strings.ToLower(_ColorName[52:64]): ColorYellowGreen,
	_ColorName[64:79]:                  ColorRedOrangeBlue,
	strings.ToLower(_ColorName[64:79]): ColorRedOrangeBlue,
}

func ParseColor(name string) (Color, error) {
	if x, ok := _ColorNameMap[name]; ok {
		return x, nil
	}
	if x, ok := _ColorNameMap[strings.ToLower(name)]; ok {
		return x, nil
	}
	return Color(0), fmt.Errorf("%s is %w", name, ErrInvalidColor)
}

func MustParseColor(name string) Color {
	val, err := ParseColor(name)
	if err != nil {
		panic(err)
	}
	return val
}

func (x Color) Ptr() *Color {
	return &x
}

func (x Color) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

func (x *Color) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseColor(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
