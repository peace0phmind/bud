package factory

import (
	"fmt"
	"github.com/expr-lang/expr"
	"reflect"
	"testing"
)

type DoInf interface {
	Hello() string
}

var _baseStruct = Singleton[BaseStruct]().Getter()

type BaseStruct struct {
	self DoInf  `wire:"self"`
	name string `wire:"value:cfg.Name"`
	cfg  *Cfg   `wire:"auto"`
}

func (b *BaseStruct) Greet() string {
	return b.self.Hello()
}

func (b *BaseStruct) Hello() string {
	return fmt.Sprintf("Hello(%s): base struct", b.name)
}

var _extStruct = Singleton[ExtStruct]().Getter()

type ExtStruct struct {
	BaseStruct
}

var _cfg = Singleton[Cfg]().Name("cfg").Getter()

type Cfg struct {
	Name string
}

func (c *Cfg) Init() {
	c.Name = "py"
}

func (e *ExtStruct) Hello() string {
	return fmt.Sprintf("Hello(%s): ext struct", e.name)
}

func TestUpdateSelf(t *testing.T) {
	testCases := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "base struct greet",
			got:  New[BaseStruct]().Greet(),
			want: "Hello(py): base struct",
		},
		{
			name: "ext struct greet",
			got:  New[ExtStruct]().Greet(),
			want: "Hello(py): ext struct",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("AutoWire = %v, want %v", tc.got, tc.want)
			}
		})
	}
}

func TestExpr(t *testing.T) {
	env := _context.getByName("env")

	envMap := map[string]any{}
	envMap["env"] = env

	program, err := expr.Compile("'abc'", expr.Env(envMap), expr.AsAny())
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	output, err := expr.Run(program, env)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%v\n", reflect.ValueOf(output).Kind())

	//out, err := expr.Eval("env.PWD1 ?? 123", envMap,)
	//if err != nil {
	//	t.Errorf("eval err: %v", err)
	//} else {
	//	t.Errorf("output: %v", reflect.ValueOf(out).Kind())
	//}
}
