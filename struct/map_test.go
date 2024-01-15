package _struct

import (
	"reflect"
	"testing"
)

type Foo struct {
	Name string
}

type Bar struct {
	Name string
}

type Baz struct {
	ID int
}

func TestAddTypeAliasMap(t *testing.T) {
	tests := []struct {
		name      string
		runFunc   func()
		wantPanic bool
	}{
		{
			name:      "Adds valid alias successfully",
			runFunc:   func() { AddTypeAliasMap[Foo, Bar]() },
			wantPanic: false,
		},
		{
			name:      "Adding the same alias twice triggers panic",
			runFunc:   func() { AddTypeAliasMap[Foo, Bar]() },
			wantPanic: true,
		},
		{
			name:      "Add different aliases to the same type",
			runFunc:   func() { AddTypeAliasMap[Foo, Bar]() },
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("AddTypeAliasMap should panic when adding the same alias twice")
					}
				}()
			} else {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Got unexpected panic: %v", r)
					}
				}()
			}
			tt.runFunc()
		})
	}
}

func TestConvert(t *testing.T) {
	testCases := []struct {
		name          string
		fn            func(any) (any, error)
		from          any
		to            any
		expectedPanic string
		expectError   bool
	}{
		{
			name:          "int to string",
			fn:            func(from any) (any, error) { return ConvertTo[string](from) },
			from:          123,
			to:            "123",
			expectedPanic: "",
			expectError:   false,
		},
		{
			name:          "string to int",
			fn:            func(from any) (any, error) { return ConvertTo[int](from) },
			from:          123,
			to:            "123",
			expectedPanic: "",
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tc.expectError {
					t.Errorf("panic = %v, expectError %v", r, tc.expectError)
					return
				}
				if tc.expectError && r != nil {
					if r != tc.expectedPanic {
						t.Errorf("Got panic = %v, want %v", r, tc.expectedPanic)
					}
				}
			}()

			got, err := tc.fn(tc.from)
			if err != nil {
				t.Errorf("convert `%s` error: %v", tc.name, err)
			}
			if reflect.DeepEqual(got, tc.to) {
				t.Errorf("convert err, want: %v  got: %+v", tc.to, got)
			}
		})
	}
}
