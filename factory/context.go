package factory

import (
	"fmt"
	"github.com/peace0phmind/bud/util"
	"reflect"
)

var _context = &context{}

type Getter func() any

type context struct {
	defaultMustBuilderCache util.Cache[reflect.Type, *contextCachedItem] // package:name -> must builder
	namedMustBuilderCache   util.Cache[string, *contextCachedItem]       // name -> must builder
	wiringCache             util.Cache[reflect.Type, bool]
	option                  Option
}

type contextCachedItem struct {
	_type  reflect.Type
	getter Getter
}

func UseConstructor(useConstructor bool) {
	_context.option.UseConstructor(useConstructor)
}

func InitMethodName(initMethodName string) {
	_context.option.InitMethodName(initMethodName)
}

func Get[T any]() *T {
	vt := reflect.TypeOf((*T)(nil))

	result := _context._get(vt)

	resultType := reflect.TypeOf(result)
	if resultType.Kind() == reflect.Ptr && resultType.ConvertibleTo(vt) {
		return result.(*T)
	}

	// panic
	panic(fmt.Sprintf("Invalid type: need %v, get %v", vt, resultType))
}

func GetByName[T any](name string) *T {
	vt := reflect.TypeOf((*T)(nil))

	result := _context._getByName(name)

	resultType := reflect.TypeOf(result)
	if resultType.Kind() == reflect.Ptr && resultType.ConvertibleTo(vt) {
		return result.(*T)
	}

	// panic
	panic(fmt.Sprintf("Invalid type: need %v, get %v", vt, resultType))
}

func Range[T any](interfaceFunc func(T) bool, structFunc func(*T) bool) {
	vT := reflect.TypeOf((*T)(nil))

	isInterface := false
	if vT.Elem().Kind() == reflect.Interface {
		vT = vT.Elem()
		isInterface = true
	}

	if isInterface {
		if interfaceFunc == nil {
			panic("T is an interface, interfaceFunc must not be nil")
		}
	} else {
		if structFunc == nil {
			panic("T is a struct, structFunc must not be nil")
		}
	}

	_context.defaultMustBuilderCache.Range(func(k reflect.Type, v *contextCachedItem) bool {
		if k.ConvertibleTo(vT) {
			if isInterface {
				return interfaceFunc(v.getter().(T))
			} else {
				return structFunc(v.getter().(*T))
			}
		}
		return true
	})
}

func (c *context) _getByNameOrType(name string, vt reflect.Type) any {
	mb, ok := c.namedMustBuilderCache.Get(name)

	if ok {
		result := mb.getter()
		rt := reflect.TypeOf(result)
		if vt.ConvertibleTo(rt) {
			return result
		}
	}

	return c._get(vt)
}

func (c *context) _get(vt reflect.Type) any {
	mb, ok := c.defaultMustBuilderCache.Get(vt)

	if ok {
		return mb.getter()
	}

	convertibleList := c.defaultMustBuilderCache.Filter(func(k reflect.Type, v *contextCachedItem) bool {
		return k.ConvertibleTo(vt)
	})

	convertibleListSize := convertibleList.Size()

	if convertibleListSize > 1 {
		panic(fmt.Sprintf("Multiple default builders found for type: %v, please use named singleton", vt))
	}

	if convertibleListSize == 1 {
		convertibleList.Range(func(k reflect.Type, v *contextCachedItem) bool {
			mb = v
			ok = true
			return false
		})

		if ok {
			return mb.getter()
		}
	}

	svt := vt
	if svt.Kind() == reflect.Ptr {
		svt = svt.Elem()
	}

	panic(fmt.Sprintf("use type to get Getter, %s:%s not found.", svt.PkgPath(), svt.Name()))
}

func (c *context) _set(vt reflect.Type, builder Getter) {
	_, getOk := c.defaultMustBuilderCache.GetOrStore(vt, &contextCachedItem{_type: vt, getter: builder})
	if getOk {
		panic(fmt.Sprintf("Default builder allready exist: %s", vt.String()))
	}
}

func (c *context) _getByName(name string) any {
	mb, ok := c.namedMustBuilderCache.Get(name)

	if ok {
		return mb.getter()
	}

	panic(fmt.Sprintf("Named builder %s not found.", name))
}

func (c *context) _setByName(name string, vt reflect.Type, builder Getter) {
	_, getOk := c.namedMustBuilderCache.GetOrStore(name, &contextCachedItem{_type: vt, getter: builder})
	if getOk {
		panic(fmt.Sprintf("Named builder allready exist: %s", name))
	}
}

func (c *context) _addWiring(vt reflect.Type) {
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}

	_, ok := c.wiringCache.GetOrStore(vt, true)
	if ok {
		panic(fmt.Sprintf("%s:%s is wiring, possible circular reference exists.", vt.PkgPath(), vt.Name()))
	}
}

func (c *context) _deleteWiring(vt reflect.Type) {
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}

	c.wiringCache.Delete(vt)
}
