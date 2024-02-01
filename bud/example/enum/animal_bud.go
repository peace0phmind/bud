package enum

import (
	"errors"
	"fmt"
)

const (
	// AnimalCat is an Animal of type Cat.
	AnimalCat Animal = iota
	// AnimalDog is an Animal of type Dog.
	AnimalDog
	// AnimalFish is an Animal of type Fish.
	AnimalFish
	// AnimalFishPlusPlus is an Animal of type Fish++.
	AnimalFishPlusPlus
	// AnimalFishSharp is an Animal of type Fish#.
	AnimalFishSharp
)

var ErrInvalidAnimal = errors.New("not a valid Animal")

var _AnimalName = "CatDogFishFish++Fish#"

var _AnimalMapName = map[Animal]string{
	AnimalCat:          _AnimalName[0:3],
	AnimalDog:          _AnimalName[3:6],
	AnimalFish:         _AnimalName[6:10],
	AnimalFishPlusPlus: _AnimalName[10:16],
	AnimalFishSharp:    _AnimalName[16:21],
}

// Name is the attribute of Animal.
func (x Animal) Name() string {
	if v, ok := _AnimalMapName[x]; ok {
		return v
	}
	panic(ErrInvalidAnimal)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Animal) IsValid() bool {
	_, ok := _AnimalMapName[x]
	return ok
}

// String implements the Stringer interface.
func (x Animal) String() string {
	if v, ok := _AnimalMapName[x]; ok {
		return v
	}
	return fmt.Sprintf("Animal(%d)", x)
}

var _AnimalNameMap = map[string]Animal{
	_AnimalName[0:3]:   AnimalCat,
	_AnimalName[3:6]:   AnimalDog,
	_AnimalName[6:10]:  AnimalFish,
	_AnimalName[10:16]: AnimalFishPlusPlus,
	_AnimalName[16:21]: AnimalFishSharp,
}

// ParseAnimal converts a string to an Animal.
func ParseAnimal(value string) (Animal, error) {
	if x, ok := _AnimalNameMap[value]; ok {
		return x, nil
	}
	return Animal(0), fmt.Errorf("%s is %w", value, ErrInvalidAnimal)
}
