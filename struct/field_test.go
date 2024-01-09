package _struct

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

type unmarshaler struct {
	time.Duration
}

// TextUnmarshaler implements encoding.TextUnmarshaler.
func (d *unmarshaler) UnmarshalText(data []byte) (err error) {
	if len(data) != 0 {
		d.Duration, err = time.ParseDuration(string(data))
	} else {
		d.Duration = 0
	}
	return err
}

type PointStruct struct {
	NestedNonDefined struct {
		NonDefined struct {
			String string `env:"STR"`
		} `envPrefix:"NONDEFINED_"`
	} `envPrefix:"PRF_"`
}

type WalkStruct struct {
	String          string    `env:"STRING"`
	StringPtr       *string   `env:"STRING"`
	Strings         []string  `env:"STRINGS"`
	StringPtrs      []*string `env:"STRINGS"`
	StringPtrInited *string   `env:"STRINGInited"`

	Bool     bool    `env:"BOOL"`
	BoolPtr  *bool   `env:"BOOL"`
	Bools    []bool  `env:"BOOLS"`
	BoolPtrs []*bool `env:"BOOLS"`

	Int     int    `env:"INT"`
	IntPtr  *int   `env:"INT"`
	Ints    []int  `env:"INTS"`
	IntPtrs []*int `env:"INTS"`

	Int8     int8    `env:"INT8"`
	Int8Ptr  *int8   `env:"INT8"`
	Int8s    []int8  `env:"INT8S"`
	Int8Ptrs []*int8 `env:"INT8S"`

	Int16     int16    `env:"INT16"`
	Int16Ptr  *int16   `env:"INT16"`
	Int16s    []int16  `env:"INT16S"`
	Int16Ptrs []*int16 `env:"INT16S"`

	Int32     int32    `env:"INT32"`
	Int32Ptr  *int32   `env:"INT32"`
	Int32s    []int32  `env:"INT32S"`
	Int32Ptrs []*int32 `env:"INT32S"`

	Int64     int64    `env:"INT64"`
	Int64Ptr  *int64   `env:"INT64"`
	Int64s    []int64  `env:"INT64S"`
	Int64Ptrs []*int64 `env:"INT64S"`

	Uint     uint    `env:"UINT"`
	UintPtr  *uint   `env:"UINT"`
	Uints    []uint  `env:"UINTS"`
	UintPtrs []*uint `env:"UINTS"`

	Uint8     uint8    `env:"UINT8"`
	Uint8Ptr  *uint8   `env:"UINT8"`
	Uint8s    []uint8  `env:"UINT8S"`
	Uint8Ptrs []*uint8 `env:"UINT8S"`

	Uint16     uint16    `env:"UINT16"`
	Uint16Ptr  *uint16   `env:"UINT16"`
	Uint16s    []uint16  `env:"UINT16S"`
	Uint16Ptrs []*uint16 `env:"UINT16S"`

	Uint32     uint32    `env:"UINT32"`
	Uint32Ptr  *uint32   `env:"UINT32"`
	Uint32s    []uint32  `env:"UINT32S"`
	Uint32Ptrs []*uint32 `env:"UINT32S"`

	Uint64     uint64    `env:"UINT64"`
	Uint64Ptr  *uint64   `env:"UINT64"`
	Uint64s    []uint64  `env:"UINT64S"`
	Uint64Ptrs []*uint64 `env:"UINT64S"`

	Float32     float32    `env:"FLOAT32"`
	Float32Ptr  *float32   `env:"FLOAT32"`
	Float32s    []float32  `env:"FLOAT32S"`
	Float32Ptrs []*float32 `env:"FLOAT32S"`

	Float64     float64    `env:"FLOAT64"`
	Float64Ptr  *float64   `env:"FLOAT64"`
	Float64s    []float64  `env:"FLOAT64S"`
	Float64Ptrs []*float64 `env:"FLOAT64S"`

	Duration     time.Duration    `env:"DURATION"`
	Durations    []time.Duration  `env:"DURATIONS"`
	DurationPtr  *time.Duration   `env:"DURATION"`
	DurationPtrs []*time.Duration `env:"DURATIONS"`

	Unmarshaler     unmarshaler    `env:"UNMARSHALER"`
	UnmarshalerPtr  *unmarshaler   `env:"UNMARSHALER"`
	Unmarshalers    []unmarshaler  `env:"UNMARSHALERS"`
	UnmarshalerPtrs []*unmarshaler `env:"UNMARSHALERS"`

	URL     url.URL    `env:"URL"`
	URLPtr  *url.URL   `env:"URL"`
	URLs    []url.URL  `env:"URLS"`
	URLPtrs []*url.URL `env:"URLS"`

	StringWithdefault string `env:"DATABASE_URL" envDefault:"postgres://localhost:5432/db"`

	CustomSeparator []string `env:"SEPSTRINGS" envSeparator:":"`

	NonDefined struct {
		String string `env:"NONDEFINED_STR"`
	}

	NestedNonDefined struct {
		NonDefined struct {
			String string `env:"STR"`
		} `envPrefix:"NONDEFINED_"`
	} `envPrefix:"PRF_"`

	NotAnEnv   string
	unexported string `env:"FOO"`

	PointStruct
	PS         PointStruct
	PPS        *PointStruct
	PPSNotInit *PointStruct

	SlicePointStruct   []PointStruct
	SlicePPointStruct  []*PointStruct
	PSlicePointStruct  *[]PointStruct
	PSlicePPointStruct *[]PointStruct
}

type WalkSliceStruct struct {
	PointStruct
	SlicePointStructInit   []PointStruct
	SlicePPointStructInit  []*PointStruct
	SlicePointStruct       []PointStruct
	SlicePPointStruct      []*PointStruct
	PSlicePointStructInit  *[]PointStruct
	PSlicePPointStructInit *[]*PointStruct
	PSlicePointStruct      *[]PointStruct
	PSlicePPointStruct     *[]*PointStruct
}

func TestWalkSliceStruct(t *testing.T) {

	init1 := make([]PointStruct, 2)
	init2 := make([]*PointStruct, 2)
	init2 = append(init2, &PointStruct{})
	wk := &WalkSliceStruct{
		SlicePointStructInit:   make([]PointStruct, 2),
		SlicePPointStructInit:  make([]*PointStruct, 2),
		PSlicePointStructInit:  &init1,
		PSlicePPointStructInit: &init2,
	}
	wk.SlicePPointStructInit = append(wk.SlicePPointStructInit, &PointStruct{})

	count := 0
	err := WalkField(wk, func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type) error {
		count++
		println(count, GetFieldPath(structField, rootTypes), fieldValue.Kind().String())
		return nil
	})

	if err != nil {
		t.Error(err)
	}
}

func TestWalk(t *testing.T) {

	count := 0
	var hello = "abc"
	wk := &WalkStruct{
		StringPtrInited: &hello,
		PPS:             &PointStruct{},
	}
	err := WalkField(wk, func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type) error {
		count++
		println(count, GetFieldPath(structField, rootTypes), fieldValue.Kind().String())
		return nil
	})

	if err != nil {
		t.Error(err)
	}
}

func TestWalkWithTag(t *testing.T) {

	count := 0
	var hello = "abc"
	wk := &WalkStruct{
		StringPtrInited: &hello,
		PPS:             &PointStruct{},
	}

	err := WalkWithTagName(wk, "env", func(fieldValue reflect.Value, structField reflect.StructField, rootTypes []reflect.Type, tagValue string) error {
		count++
		println(count, GetFieldPath(structField, rootTypes), fieldValue.Kind().String(), tagValue)
		return nil
	})

	if err != nil {
		t.Error(err)
	}
}
