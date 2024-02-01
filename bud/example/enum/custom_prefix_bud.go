package enum

import (
	"errors"
	"fmt"
)

const (
	// AcmeIncProductAnvil is a Product of type Anvil.
	AcmeIncProductAnvil Product = iota
	// AcmeIncProductDynamite is a Product of type Dynamite.
	AcmeIncProductDynamite
	// AcmeIncProductGlue is a Product of type Glue.
	AcmeIncProductGlue
)

var ErrInvalidProduct = errors.New("not a valid Product")

var _ProductName = "AnvilDynamiteGlue"

var _ProductMapName = map[Product]string{
	AcmeIncProductAnvil:    _ProductName[0:5],
	AcmeIncProductDynamite: _ProductName[5:13],
	AcmeIncProductGlue:     _ProductName[13:17],
}

// Name is the attribute of Product.
func (x Product) Name() string {
	if v, ok := _ProductMapName[x]; ok {
		return v
	}
	panic(ErrInvalidProduct)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Product) IsValid() bool {
	_, ok := _ProductMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x Product) String() string {
	if v, ok := _ProductMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("Product(%d)", x)
}

var _ProductNameMap = map[string]Product{
	_ProductName[0:5]:   AcmeIncProductAnvil,
	_ProductName[5:13]:  AcmeIncProductDynamite,
	_ProductName[13:17]: AcmeIncProductGlue,
}

// ParseProduct converts a string to a Product.
func ParseProduct(value string) (Product, error) {
	if x, ok := _ProductNameMap[value]; ok {
		return x, nil
	}
	return Product(0), fmt.Errorf("%s is %w", value, ErrInvalidProduct)
}
