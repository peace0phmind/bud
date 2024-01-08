package factory

import (
	"testing"
)

type DoInf interface {
	Hello() string
}

type BaseStruct struct {
	self DoInf `wire:"self"`
}

func (b *BaseStruct) Greet() string {
	return b.self.Hello()
}

func (b *BaseStruct) Hello() string {
	return "Hello: base struct"
}

type ExtStruct struct {
	BaseStruct
}

func (e *ExtStruct) Hello() string {
	return "Hello: ext struct"
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
			want: "Hello: base struct",
		},
		{
			name: "ext struct greet",
			got:  New[ExtStruct]().Greet(),
			want: "Hello: ext struct",
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
