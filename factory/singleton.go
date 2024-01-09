package factory

import (
	"reflect"
	"strings"
	"sync"
)

// singleton is a generic type that implements the singleton design pattern.
// It ensures that only one instance of an object is created and provides a global point
// of access to that instance.
//
// Usage:
//
//	singleton[T] is a type parameterized by the concrete type T.
//	To create a singleton instance, use the Singleton[T] function.
//	To get the instance, call the GetInstance method on the singleton object.
//	You can set an initialization function using the SetInitOnce method or a must
//	initialization function using the SetMustInitOnce method.
//	The singleton object provides builder functions that return a closure
//	for easy initialization of dependency injection containers.
//
//	Initialization can be done through the initialization method of singleton, or the corresponding interface
//	implemented by T type. If you call Singleton and set the initialization methods InitOnce and
//	MustInitOnce at the same time, initialization will only execute the InitOnce method. If T implements
//	the ISingleton and MustInitOnce interfaces at the same time, initialization will only execute the
//	InitOnce interface. If both the initialization method of Singleton is set and T also implements the
//	corresponding interface, only the method of Singleton is executed during initialization,
//	and the interface implemented by T is not executed.
//
// Example:
//
//	type cat struct {
//	}
//	func (c* cat) Meow() {
//	}
//
//	var Cat = Singleton[cat]().MustBuilder()
//
// Call method:
//
//	Cat().Meow()
/*
  通过 initOnce func 和 mustInitOnce func进行初始化的类，不设置default和进行auto wire，如有需要可以手动调用
  当使用接口或没有实现任何接口进行类的初始化时， 默认先调用set default, 再进行auto wire，最后调用这些实现额接口
*/
type singleton[T any] struct {
	once         sync.Once
	obj          *T
	err          error
	initOnce     func() (*T, error)
	mustInitOnce func() *T
	autoWire     bool
	setDefault   bool

	_name        string
	_type        reflect.Type
	_mustBuilder func() *T
	lock         sync.Mutex
}

type ISingleton interface {
	InitOnce() error
}

type IMustSingleton interface {
	MustInitOnce()
}

func _singleton[T any]() *singleton[T] {
	result := &singleton[T]{
		autoWire:   true,
		setDefault: true,
	}

	var t T
	result._type = reflect.TypeOf(&t)

	result._mustBuilder = func() *T {
		if instance, err := result.GetInstance(); err != nil {
			panic(err)
		} else {
			return instance
		}
	}

	return result
}

func Singleton[T any]() *singleton[T] {
	result := _singleton[T]()

	_context._set(result._type, func() any {
		return result._mustBuilder()
	})

	return result
}

func NamedSingleton[T any](name string) *singleton[T] {
	result := _singleton[T]()

	return result.Name(name)
}

func (s *singleton[T]) Name(name string) *singleton[T] {
	s.lock.Lock()
	defer s.lock.Unlock()

	name = strings.TrimSpace(name)
	if len(name) == 0 {
		panic("name must not be empty")
	}

	if len(s._name) == 0 {
		_context._setByName(name, s._type, func() any {
			return s._mustBuilder()
		})
		s._name = name
	}

	return s
}

func (s *singleton[T]) AutoWire(autoWire bool) *singleton[T] {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.autoWire = autoWire
	return s
}

func (s *singleton[T]) SetDefault(setDefault bool) *singleton[T] {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.setDefault = setDefault
	return s
}

func (s *singleton[T]) CreateOnly() *singleton[T] {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.autoWire = false
	s.setDefault = false
	return s
}

func (s *singleton[T]) GetInstance() (*T, error) {
	s.once.Do(func() {
		if s.initOnce != nil {
			s.obj, s.err = s.initOnce()
		} else if s.mustInitOnce != nil {
			s.obj = s.mustInitOnce()
		} else {
			s.obj = new(T)

			if s.setDefault {
				s.err = SetDefault(s.obj)
				if s.err != nil {
					return
				}
			}

			if s.autoWire {
				s.err = AutoWire(s.obj)
				if s.err != nil {
					return
				}
			}

			if initializer, ok := any(s.obj).(ISingleton); ok {
				s.err = initializer.InitOnce()
			} else if mustInitializer, mustOk := any(s.obj).(IMustSingleton); mustOk {
				mustInitializer.MustInitOnce()
			}
		}
	})

	return s.obj, s.err
}

func (s *singleton[T]) SetInitOnce(initOnce func() (*T, error)) *singleton[T] {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.initOnce = initOnce
	return s
}

func (s *singleton[T]) SetMustInitOnce(mustInitOnce func() *T) *singleton[T] {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.mustInitOnce = mustInitOnce
	return s
}

func (s *singleton[T]) Builder() func() (*T, error) {
	return func() (*T, error) {
		return s.GetInstance()
	}
}

func (s *singleton[T]) MustBuilder() func() *T {
	return s._mustBuilder
}

func New[T any]() *T {
	t := new(T)

	// set default first
	if err := SetDefault(t); err != nil {
		panic(err)
	}

	// call auto wire
	if err := AutoWire(t); err != nil {
		panic(err)
	}

	return t
}
