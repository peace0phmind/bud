package factory

import (
	"fmt"
	"reflect"
	"sync"
)

type IInit interface {
	Init()
}

type customInitMethod struct {
	useConstructor bool
	initMethodName string
}

type Option struct {
	customInitMethod
	doSetDefault bool
	doAutoWire   bool
	lock         sync.Mutex
}

func NewOption() *Option {
	return &Option{
		doSetDefault: true,
		doAutoWire:   true,
	}
}

func (o *Option) SetDefault(setDefault bool) *Option {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.doSetDefault = setDefault
	return o
}

func (o *Option) AutoWire(autoWire bool) *Option {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.doAutoWire = autoWire
	return o
}

func (o *Option) UseConstructor(useConstructor bool) *Option {
	o.lock.Lock()
	defer o.lock.Unlock()

	if useConstructor && len(o.initMethodName) > 0 {
		panic("initMethodName must be empty when UseConstructor is true")
	}

	o.useConstructor = useConstructor
	return o
}

func (o *Option) InitMethodName(initMethodName string) *Option {
	o.lock.Lock()
	defer o.lock.Unlock()

	if len(initMethodName) > 0 && o.useConstructor {
		panic("useConstructor must be false when initMethodName is set")
	}

	o.initMethodName = initMethodName
	return o
}

func (o *Option) useInitMethod() bool {
	return o.useConstructor || len(o.initMethodName) > 0
}

var newDefaultOption = NewOption()

func GetNewDefaultOption() *Option {
	return newDefaultOption
}

func New[T any]() *T {
	return NewWithOption[T](newDefaultOption)
}

func NewWithOption[T any](option *Option) *T {
	t := new(T)

	if option.useInitMethod() {
		vt := reflect.TypeOf(t)

		if vt.Kind() == reflect.Ptr && vt.Elem().Kind() == reflect.Struct {
			vt = vt.Elem()

			methodName := option.initMethodName
			if option.useConstructor {
				methodName = vt.Name()
			}

			initMethod, ok := vt.MethodByName(methodName)
			if ok {
				if initMethod.Type.NumOut() > 0 {
					panic(fmt.Sprintf("Init method '%s' must not have return values", methodName))
				}

				defer initMethod.Func.Call([]reflect.Value{reflect.ValueOf(t)})
			}
		} else {
			panic("T must be a struct type")
		}
	} else {
		if init, ok := any(t).(IInit); ok {
			// 最后执行
			defer init.Init()
		}
	}

	// do set default
	if option.doSetDefault {
		if err := SetDefault(t); err != nil {
			panic(err)
		}
	}

	// do auto wire
	if option.doAutoWire {
		if err := AutoWire(t); err != nil {
			panic(err)
		}
	}

	return t
}
