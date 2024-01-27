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

func AnimalValues() []Animal {
	return []Animal{
		AnimalCat,
		AnimalDog,
		AnimalFish,
		AnimalFishPlusPlus,
		AnimalFishSharp,
	}
}

var _AnimalName = "CatDogFishFish++Fish#"

var _AnimalMapName = map[Animal]string{
	AnimalCat:          _AnimalName[0:3],
	AnimalDog:          _AnimalName[3:6],
	AnimalFish:         _AnimalName[6:10],
	AnimalFishPlusPlus: _AnimalName[10:16],
	AnimalFishSharp:    _AnimalName[16:21],
}

func (x Animal) Name() string {
	if result, ok := _AnimalMapName[x]; ok {
		return result
	}
	panic(ErrInvalidAnimal)
}

func (x Animal) String() string {
	if str, ok := _AnimalMapName[x]; ok {
		return str
	}
	return fmt.Sprintf("Animal(%d)", x)
}

func (x Animal) IsValid() bool {
	_, ok := _AnimalMapName[x]
	return ok
}

var _AnimalNameMap = map[string]Animal{
	_AnimalName[0:3]:   AnimalCat,
	_AnimalName[3:6]:   AnimalDog,
	_AnimalName[6:10]:  AnimalFish,
	_AnimalName[10:16]: AnimalFishPlusPlus,
	_AnimalName[16:21]: AnimalFishSharp,
}

func ParseAnimal(name string) (Animal, error) {
	if x, ok := _AnimalNameMap[name]; ok {
		return x, nil
	}
	return Animal(0), fmt.Errorf("%s is %w", name, ErrInvalidAnimal)
}
