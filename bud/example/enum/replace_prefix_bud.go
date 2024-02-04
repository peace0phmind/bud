package enum

import (
	"fmt"
	"strings"
)

const (
	// AcmeInc_SOME_PLACE_AWESOME is a Shop of type SOME_PLACE_AWESOME.
	AcmeInc_SOME_PLACE_AWESOME Shop = "SOME_PLACE_AWESOME"
	// AcmeInc_SomewhereElse is a Shop of type SomewhereElse.
	AcmeInc_SomewhereElse Shop = "SomewhereElse"
	// AcmeInc_LocationUnknown is a Shop of type LocationUnknown.
	AcmeInc_LocationUnknown Shop = "LocationUnknown"
)

var ErrInvalidShop = fmt.Errorf("not a valid Shop, try [%s]", strings.Join(_ShopNames, ", "))

// Name is the attribute of Shop.
func (x Shop) Name() string {
	if v, ok := _ShopNameMap[string(x)]; ok {
		return string(v)
	}
	panic(ErrInvalidShop)
}

// Value is the attribute of Shop.
func (x Shop) Value() string {
	if x.IsValid() {
		return string(x)
	}
	panic(ErrInvalidShop)
}

var _ShopNames = []string{
	"SOME_PLACE_AWESOME",
	"SomewhereElse",
	"LocationUnknown",
}

// ShopNames returns a list of the names of Shop
func ShopNames() []string {
	return _ShopNames
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Shop) IsValid() bool {
	_, ok := _ShopNameMap[string(x)]
	return ok
}

// String implements the Stringer interface.
func (x Shop) String() string {
	return x.Name()
}

var _ShopNameMap = map[string]Shop{
	"SOME_PLACE_AWESOME": AcmeInc_SOME_PLACE_AWESOME,
	"SomewhereElse":      AcmeInc_SomewhereElse,
	"LocationUnknown":    AcmeInc_LocationUnknown,
}

// ParseShop converts a string to a Shop.
func ParseShop(value string) (Shop, error) {
	if x, ok := _ShopNameMap[value]; ok {
		return x, nil
	}
	return "", fmt.Errorf("%s is %w", value, ErrInvalidShop)
}

// MarshalText implements the text marshaller method.
func (x Shop) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Shop) UnmarshalText(text []byte) error {
	val, err := ParseShop(string(text))
	if err != nil {
		return err
	}
	*x = val
	return nil
}
