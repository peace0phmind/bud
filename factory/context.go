package factory

import (
	"fmt"
	"github.com/peace0phmind/bud/util"
	"reflect"
)

var _context = &context{}

type MustBuilder func() any

type context struct {
	defaultMustBuilderCache util.Cache[reflect.Type, *contextCachedItem] // package:name -> must builder
	namedMustBuilderCache   util.Cache[string, *contextCachedItem]       // name -> must builder
	wiringCache             util.Cache[reflect.Type, bool]
}

type contextCachedItem struct {
	_type reflect.Type
	build MustBuilder
}

func Get[T any]() *T {
	var t T
	vt := reflect.TypeOf(&t)

	result := _context._get(vt)

	resultType := reflect.TypeOf(result)
	if resultType.Kind() == reflect.Ptr && resultType.ConvertibleTo(vt) {
		return result.(*T)
	}

	// panic
	panic(fmt.Sprintf("Invalid type: need %v, get %v", vt, resultType))
}

func GetByName[T any](name string) *T {
	var t T
	vt := reflect.TypeOf(&t)

	result := _context._getByName(name)

	resultType := reflect.TypeOf(result)
	if resultType.Kind() == reflect.Ptr && resultType.ConvertibleTo(vt) {
		return result.(*T)
	}

	// panic
	panic(fmt.Sprintf("Invalid type: need %v, get %v", vt, resultType))
}

func (c *context) _get(vt reflect.Type) any {
	mb, ok := c.defaultMustBuilderCache.Get(vt)

	if ok {
		return mb.build()
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
			return mb.build()
		}
	}

	panic(fmt.Sprintf("Default Builder %s:%s  not found.", vt.PkgPath(), vt.Name()))
}

func (c *context) _set(vt reflect.Type, builder MustBuilder) {
	_, getOk := c.defaultMustBuilderCache.GetOrStore(vt, &contextCachedItem{_type: vt, build: builder})
	if getOk {
		panic(fmt.Sprintf("Default builder allready exist: %v", vt))
	}
}

func (c *context) _getByName(name string) any {
	mb, ok := c.namedMustBuilderCache.Get(name)

	if ok {
		return mb.build()
	}

	panic(fmt.Sprintf("Named builder %s not found.", name))
}

func (c *context) _setByName(name string, vt reflect.Type, builder MustBuilder) {
	_, getOk := c.namedMustBuilderCache.GetOrStore(name, &contextCachedItem{_type: vt, build: builder})
	if getOk {
		panic(fmt.Sprintf("Named builder allready exist: %s", name))
	}
}

func (c *context) _addWire(vt reflect.Type) {
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}

	_, ok := c.wiringCache.GetOrStore(vt, true)
	if ok {
		panic(fmt.Sprintf("%s:%s is wiring, possible circular reference exists.", vt.PkgPath(), vt.Name()))
	}
}

func (c *context) _deleteWire(vt reflect.Type) {
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}

	c.wiringCache.Delete(vt)
}
