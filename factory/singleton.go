package factory

import "sync"

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
//	The singleton object also provides builder functions that return a closure
//	for easy initialization of dependency injection containers.
//
// singleton Initialization can be done through the initialization method of singleton, or the corresponding interface implemented by T type
//
// If you call Singleton and set the initialization methods InitOnce and MustInitOnce at the same time
// Initialization will only execute the InitOnce method
//
// If T implements the ISingleton and MustInitOnce interfaces at the same time
// Initialization will only execute the InitOnce interface
//
// If both the initialization method of Singleton is set and T also implements the corresponding interface
// Only the method of Singleton is executed during initialization, and the interface implemented by T is not executed
//
// Example:
// type cat struct {
// }
// func (c* cat) Meow() {
// }
//
//	var Cat = Singleton[cat]().MustBuilder()
//
// call method:
// Cat().Meow()
type singleton[T any] struct {
	once         sync.Once
	obj          *T
	err          error
	initOnce     func() (*T, error)
	mustInitOnce func() *T
}

type ISingleton interface {
	InitOnce() error
}

type IMustSingleton interface {
	MustInitOnce()
}

func Singleton[T any]() *singleton[T] {
	return &singleton[T]{}
}

func (s *singleton[T]) GetInstance() (*T, error) {
	s.once.Do(func() {
		if s.initOnce != nil {
			s.obj, s.err = s.initOnce()
		} else if s.mustInitOnce != nil {
			s.obj = s.mustInitOnce()
		} else {
			s.obj = new(T)
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
	s.initOnce = initOnce
	return s
}

func (s *singleton[T]) SetMustInitOnce(mustInitOnce func() *T) *singleton[T] {
	s.mustInitOnce = mustInitOnce
	return s
}

func (s *singleton[T]) Builder() func() (*T, error) {
	return func() (*T, error) {
		return s.GetInstance()
	}
}

func (s *singleton[T]) MustBuilder() func() *T {
	return func() *T {
		if instance, err := s.GetInstance(); err != nil {
			panic(err)
		} else {
			return instance
		}
	}
}
