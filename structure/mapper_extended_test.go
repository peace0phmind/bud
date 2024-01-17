package structure

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

type DummyTextUnmarshaler struct {
	value string
}

func (dtu *DummyTextUnmarshaler) UnmarshalText(text []byte) error {
	if len(text) > 0 {
		dtu.value = "Hello: " + string(text)
	}
	return nil
}

func TestString2TextUnmarshalerMapper(t *testing.T) {
	cases := []struct {
		name     string
		from     string
		expected string
		isError  bool
	}{
		{name: "Unmarshal successful", from: "test", expected: "Hello: test", isError: false},
		{name: "Unmarshal empty string", from: "", expected: "", isError: false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			dtu := &DummyTextUnmarshaler{}
			from := reflect.ValueOf(testCase.from)
			to := reflect.ValueOf(dtu)

			err := string2TextUnmarshalerMapper(from, to)

			if (err != nil) != testCase.isError {
				t.Errorf("Unexpected error state: got %v, but expected error to be %v", err != nil, testCase.isError)
			}

			if err == nil && dtu.value != testCase.expected {
				t.Errorf("Unmarshalled value mismatch: got %v, but expected %v", dtu.value, testCase.expected)
			}
		})
	}
}

type binaryUnmarshalerStub struct {
	data []byte
}

func (b *binaryUnmarshalerStub) UnmarshalBinary(data []byte) error {
	b.data = make([]byte, len(data))
	copy(b.data, data)
	return nil
}

func TestString2TBinaryUnmarshalerMapper(t *testing.T) {
	testCases := []struct {
		name    string
		from    string
		to      interface{}
		want    []byte
		wantErr bool
	}{
		{
			name: "valid binary unmarshaler implementation",
			from: "testdata",
			to:   &binaryUnmarshalerStub{},
			want: []byte("testdata"),
		},
		{
			name:    "invalid type without binary unmarshaler",
			from:    "testdata",
			to:      bytes.NewBuffer(nil),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := string2BinaryUnmarshalerMapper(reflect.ValueOf(tt.from), reflect.ValueOf(tt.to))
			if err != nil {
				if !tt.wantErr {
					t.Errorf("string2BinaryUnmarshalerMapper() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if got := tt.to.(*binaryUnmarshalerStub).data; !bytes.Equal(got, tt.want) {
					t.Errorf("string2BinaryUnmarshalerMapper() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestStringToDurationMapper(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Duration
		err   bool
	}{
		{
			name:  "ValidDuration",
			input: "1h",
			want:  time.Hour,
			err:   false,
		},
		{
			name:  "ValidDuration",
			input: "1m",
			want:  time.Minute,
			err:   false,
		},
		{
			name:  "ValidDuration",
			input: "1s",
			want:  time.Second,
			err:   false,
		},
		{
			name:  "InvalidDuration",
			input: "not a duration",
			want:  0,
			err:   true,
		},
		{
			name:  "EmptyDuration",
			input: "",
			want:  0,
			err:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from := reflect.ValueOf(tt.input)
			to := reflect.New(reflect.TypeOf((*time.Duration)(nil)).Elem()).Elem()
			err := string2durationMapper(from, to)
			if (err != nil) != tt.err {
				t.Errorf("string2durationMapper() error = %v, expected error? %v", err, tt.err)
			}
			if !tt.err && to.Interface() != tt.want {
				t.Errorf("string2durationMapper() = %v, want %v", to.Interface(), tt.want)
			}
		})
	}
}
