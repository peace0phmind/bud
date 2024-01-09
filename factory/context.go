package factory

import (
	"fmt"
	"github.com/peace0phmind/bud/util"
	"reflect"
)

var _context = Singleton[context]().MustBuilder()

type MustBuilder func() any

type contextKey struct {
	Package string
	Name    string
}

func getContextKeyFromType(vt reflect.Type) contextKey {
	return contextKey{
		Package: vt.PkgPath(),
		Name:    vt.Name(),
	}
}

type mustBuilder struct {
	build MustBuilder
}

type context struct {
	defaultMustBuilderCache util.Cache[contextKey, *mustBuilder] // package:name -> must builder
	namedMustBuilderCache   util.Cache[string, *mustBuilder]     // name -> must builder
	wiringCache             util.Cache[contextKey, bool]
}

func Get[T any]() *T {
	var t T
	vt := reflect.TypeOf(&t)

	result := _context()._get(vt)

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

	result := _context()._getByName(name)

	resultType := reflect.TypeOf(result)
	if resultType.Kind() == reflect.Ptr && resultType.ConvertibleTo(vt) {
		return result.(*T)
	}

	// panic
	panic(fmt.Sprintf("Invalid type: need %v, get %v", vt, resultType))
}

func (c *context) _get(vt reflect.Type) any {
	ck := getContextKeyFromType(vt)

	mb, ok := c.defaultMustBuilderCache.Get(ck)

	if ok {
		return mb.build()
	}

	panic(fmt.Sprintf("Default Builder %s:%s  not found.", ck.Package, ck.Name))
}

func (c *context) _getByName(name string) any {
	mb, ok := c.namedMustBuilderCache.Get(name)

	if ok {
		return mb.build()
	}

	panic(fmt.Sprintf("Named builder %s not found.", name))
}

func (c *context) _setByName(name string, builder MustBuilder) {
	_, getOk := c.namedMustBuilderCache.GetOrStore(name, &mustBuilder{build: builder})
	if !getOk {
		panic(fmt.Sprintf("Named builder all ready exist: %s", name))
	}
}
