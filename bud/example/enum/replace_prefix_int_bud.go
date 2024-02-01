package enum

import (
	"fmt"
	"strings"
)

const (
	// AcmeInt_SOME_PLACE_AWESOME is an IntShop of type SOME_PLACE_AWESOME.
	AcmeInt_SOME_PLACE_AWESOME IntShop = iota
	// AcmeInt_SomewhereElse is an IntShop of type SomewhereElse.
	AcmeInt_SomewhereElse
	// AcmeInt_LocationUnknown is an IntShop of type LocationUnknown.
	AcmeInt_LocationUnknown
)

var ErrInvalidIntShop = fmt.Errorf("not a valid IntShop, try [%s]", strings.Join(_IntShopNames, ", "))

var _IntShopName = "SOME_PLACE_AWESOMESomewhereElseLocationUnknown"

var _IntShopMapName = map[IntShop]string{
	AcmeInt_SOME_PLACE_AWESOME: _IntShopName[0:18],
	AcmeInt_SomewhereElse:      _IntShopName[18:31],
	AcmeInt_LocationUnknown:    _IntShopName[31:46],
}

// Name is the attribute of IntShop.
func (x IntShop) Name() string {
	if v, ok := _IntShopMapName[x]; ok {
		return v
	}
	panic(ErrInvalidIntShop)
}

var _IntShopNames = []string{
	_IntShopName[0:18],
	_IntShopName[18:31],
	_IntShopName[31:46],
}

// IntShopNames returns a list of the names of IntShop
func IntShopNames() []string {
	return _IntShopNames
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x IntShop) IsValid() bool {
	_, ok := _IntShopMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x IntShop) String() string {
	if v, ok := _IntShopMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("IntShop(%d)", x)
}

var _IntShopNameMap = map[string]IntShop{
	_IntShopName[0:18]:  AcmeInt_SOME_PLACE_AWESOME,
	_IntShopName[18:31]: AcmeInt_SomewhereElse,
	_IntShopName[31:46]: AcmeInt_LocationUnknown,
}

// ParseIntShop converts a string to an IntShop.
func ParseIntShop(value string) (IntShop, error) {
	if x, ok := _IntShopNameMap[value]; ok {
		return x, nil
	}
	return IntShop(0), fmt.Errorf("%s is %w", value, ErrInvalidIntShop)
}

// MarshalText implements the text marshaller method.
func (x IntShop) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *IntShop) UnmarshalText(text []byte) error {
	val, err := ParseIntShop(string(text))
	if err != nil {
		return err
	}
	*x = val
	return nil
}
