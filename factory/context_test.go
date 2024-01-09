package factory

import "testing"

func TestContextRange(t *testing.T) {
	Range[DoInf](func(d DoInf) bool {
		println(d.Hello())
		return true
	}, nil)

	Range[ExtStruct](nil, func(e *ExtStruct) bool {
		println(e.Hello())
		return true
	})
}
