package structure

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

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

func TestString2URLMapper(t *testing.T) {

	mustParse := func(strU string) url.URL {
		u, err := url.Parse(strU)
		if err != nil {
			panic(err)
		}
		return *u
	}

	tests := []struct {
		name      string
		input     string
		want      url.URL
		wantError bool
	}{
		{
			name:      "ValidURL",
			input:     "http://example.com/",
			want:      mustParse("http://example.com/"),
			wantError: false,
		},
		{
			name:      "InvalidURL",
			input:     ":/example.com",
			want:      url.URL{},
			wantError: true,
		},
		{
			name:      "EmptyString",
			input:     "",
			want:      url.URL{},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fromValue := reflect.ValueOf(tt.input)
			toValue := reflect.New(reflect.TypeOf(url.URL{})).Elem()
			gotError := string2urlMapper(fromValue, toValue)

			if !tt.wantError && gotError != nil {
				t.Errorf("%s: did not expect an error but got: %v", tt.name, gotError)
			}

			if tt.wantError && gotError == nil {
				t.Errorf("%s: expected an error but not got any", tt.name)
			}

			toUrl := toValue.Interface().(url.URL)
			if !tt.wantError && !(toUrl.String() == tt.want.String()) {
				t.Errorf("%s: wrong result url, want: %v, got: %v", tt.name, tt.want.String(), toUrl.String())
			}
		})
	}
}
