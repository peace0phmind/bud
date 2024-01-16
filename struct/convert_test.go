package _struct

import (
	"reflect"
	"testing"
)

func TestMustConvertTo(t *testing.T) {
	tests := []struct {
		name  string
		fn    func() any
		want  any
		panic bool
	}{
		// test bool convert to int, uint, float, string
		{"Bool true to int", func() any { return MustConvertTo[int](true) }, 1, false},
		{"Bool false to int", func() any { return MustConvertTo[int](false) }, 0, false},
		{"Bool true to uint", func() any { return MustConvertTo[uint](true) }, uint(1), false},
		{"Bool false to uint", func() any { return MustConvertTo[uint](false) }, uint(0), false},
		{"Bool true float", func() any { return MustConvertTo[float32](true) }, float32(1.0), false},
		{"Bool false float", func() any { return MustConvertTo[float64](false) }, 0.0, false},
		{"Bool true String", func() any { return MustConvertTo[string](true) }, "true", false},
		{"Bool false String", func() any { return MustConvertTo[string](false) }, "false", false},

		// test int convert to bool, uint, float, string
		{"Int to Bool true", func() any { return MustConvertTo[bool](-1) }, true, false},
		{"Int to Bool false", func() any { return MustConvertTo[bool](0) }, false, false},
		{"Int to uint", func() any { return MustConvertTo[uint](456) }, uint(456), false},
		{"Int to float", func() any { return MustConvertTo[float64](456) }, 456.0, false},
		{"Int to String", func() any { return MustConvertTo[string](123) }, "123", false},

		// test uint convert to bool, int, float, string
		{"Uint to Bool true", func() any { return MustConvertTo[bool](uint(1)) }, true, false},
		{"Uint to Bool false", func() any { return MustConvertTo[bool](uint(0)) }, false, false},
		{"Uint to int", func() any { return MustConvertTo[int](uint(456)) }, 456, false},
		{"Uint to float", func() any { return MustConvertTo[float64](uint(456)) }, 456.0, false},
		{"Uint to String", func() any { return MustConvertTo[string](uint(123)) }, "123", false},

		// test float convert to bool, int, uint, string
		{"Float to Bool true", func() any { return MustConvertTo[bool](1.23) }, true, false},
		{"Float to Bool false", func() any { return MustConvertTo[bool](0.0) }, false, false},
		{"Float to int", func() any { return MustConvertTo[int](123.45) }, 123, false},
		{"Float to uint", func() any { return MustConvertTo[uint](123.45) }, uint(123), false},
		{"Float to String", func() any { return MustConvertTo[string](123.45) }, "123.45", false},

		// test string convert to bool, int, uint, float
		{"String to Bool true", func() any { return MustConvertTo[bool]("true") }, true, false},
		{"String to Bool false", func() any { return MustConvertTo[bool]("false") }, false, false},
		{"String to Bool error case", func() any { return MustConvertTo[bool]("notaboolean") }, nil, true},
		{"String to int", func() any { return MustConvertTo[int]("123") }, 123, false},
		{"String to uint", func() any { return MustConvertTo[uint]("123") }, uint(123), false},
		{"String to float", func() any { return MustConvertTo[float64]("1.23") }, 1.23, false},

		// test int convert to int8, int16, int32, int64
		{"Int to Int8", func() any { return MustConvertTo[int8](123) }, int8(123), false},
		{"Int to Int16", func() any { return MustConvertTo[int16](123) }, int16(123), false},
		{"Int to Int32", func() any { return MustConvertTo[int32](123) }, int32(123), false},
		{"Int to Int64", func() any { return MustConvertTo[int64](123) }, int64(123), false},

		// boundary test int convert to int8, int16, int32, int64
		{"Int to Int8 boundary", func() any { return MustConvertTo[int8](127) }, int8(127), false},
		{"Int to Int8 boundary", func() any { return MustConvertTo[int8](-127) }, int8(-127), false},
		{"Int to Int8 boundary", func() any { return MustConvertTo[int8](129) }, int8(-127), false},
		{"Int to Int16 boundary", func() any { return MustConvertTo[int16](32767) }, int16(32767), false},
		{"Int to Int16 boundary", func() any { return MustConvertTo[int16](-32767) }, int16(-32767), false},
		{"Int to Int16 boundary", func() any { return MustConvertTo[int16](32769) }, int16(-32767), false},
		{"Int to Int32 boundary", func() any { return MustConvertTo[int32](2147483647) }, int32(2147483647), false},
		{"Int to Int32 boundary", func() any { return MustConvertTo[int32](-2147483647) }, int32(-2147483647), false},
		{"Int to Int32 boundary", func() any { return MustConvertTo[int32](2147483649) }, int32(-2147483647), false},
		{"Int to Int64 boundary", func() any { return MustConvertTo[int64](9223372036854775807) }, int64(9223372036854775807), false},

		// test uint convert to uint8, uint16, uint32, uint64
		{"Uint to Uint8", func() any { return MustConvertTo[uint8](uint(123)) }, uint8(123), false},
		{"Uint to Uint16", func() any { return MustConvertTo[uint16](uint(123)) }, uint16(123), false},
		{"Uint to Uint32", func() any { return MustConvertTo[uint32](uint(123)) }, uint32(123), false},
		{"Uint to Uint64", func() any { return MustConvertTo[uint64](uint(123)) }, uint64(123), false},

		// boundary test uint convert to uint8, uint16, uint32, uint64
		{"Uint to Uint8 boundary", func() any { return MustConvertTo[uint8](uint(255)) }, uint8(255), false},
		{"Uint to Uint8 boundary", func() any { return MustConvertTo[uint8](uint(256)) }, uint8(0), false},
		{"Uint to Uint8 boundary", func() any { return MustConvertTo[uint8](uint(257)) }, uint8(1), false},
		{"Uint to Uint16 boundary", func() any { return MustConvertTo[uint16](uint(65535)) }, uint16(65535), false},
		{"Uint to Uint16 boundary", func() any { return MustConvertTo[uint16](uint(65536)) }, uint16(0), false},
		{"Uint to Uint16 boundary", func() any { return MustConvertTo[uint16](uint(65537)) }, uint16(1), false},
		{"Uint to Uint32 boundary", func() any { return MustConvertTo[uint32](uint(4294967295)) }, uint32(4294967295), false},
		{"Uint to Uint32 boundary", func() any { return MustConvertTo[uint32](uint(4294967296)) }, uint32(0), false},
		{"Uint to Uint32 boundary", func() any { return MustConvertTo[uint32](uint(4294967297)) }, uint32(1), false},
		{"Uint to Uint64 boundary", func() any { return MustConvertTo[uint64](uint(18446744073709551615)) }, uint64(18446744073709551615), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.panic {
						t.Errorf("The code panicked")
					}
				} else {
					if tc.panic {
						t.Errorf("The code did not panic")
					}
				}
			}()

			// Convert 'fn' to the same type as 'to'.
			// This may panic, depending on the test case.
			result := tc.fn()

			if !tc.panic {
				// Check the result.
				if !reflect.DeepEqual(result, tc.want) {
					t.Errorf("Want %v, but got %v", tc.want, result)
				}
			}
		})
	}
}
