package factory

import (
	"reflect"
	"strings"
	"sync"
)

type singleton[T any] struct {
	once     sync.Once
	obj      *T
	initFunc func() *T
	option   Option

	_name   string
	_type   reflect.Type
	_getter Getter
	lock    sync.Mutex
}

func _singleton[T any]() *singleton[T] {
	result := &singleton[T]{
		option: Option{
			doSetDefault: true,
			doAutoWire:   true,
		},
	}

	result._type = reflect.TypeOf((*T)(nil))

	result._getter = func() any {
		return result.Get()
	}

	return result
}

func Singleton[T any]() *singleton[T] {
	result := _singleton[T]()

	_context._set(result._type, result._getter)

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
		_context._setByName(name, s._type, s._getter)
		s._name = name
	} else {
		panic("name already set")
	}

	return s
}

func (s *singleton[T]) SetInitFunc(initFunc func() *T) *singleton[T] {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.initFunc = initFunc
	return s
}

func (s *singleton[T]) AutoWire(autoWire bool) *singleton[T] {
	s.option.AutoWire(autoWire)
	return s
}

func (s *singleton[T]) SetDefault(setDefault bool) *singleton[T] {
	s.option.SetDefault(setDefault)
	return s
}

func (s *singleton[T]) UseConstructor(useConstructor bool) *singleton[T] {
	s.option.UseConstructor(useConstructor)
	return s
}

func (s *singleton[T]) InitMethodName(initMethodName string) *singleton[T] {
	s.option.InitMethodName(initMethodName)
	return s
}

func (s *singleton[T]) Get() *T {
	s.once.Do(func() {
		if s.initFunc != nil {
			s.obj = s.initFunc()
		} else {
			// 如果设置了全局参数，singleton未设置，则使用全局参数覆盖配置
			if _context.option.useInitMethod() && !s.option.useInitMethod() {
				s.option.customInitMethod = _context.option.customInitMethod
			}
			s.obj = NewWithOption[T](&s.option)
		}
	})

	return s.obj
}

func (s *singleton[T]) Getter() func() *T {
	return func() *T {
		return s.Get()
	}
}
