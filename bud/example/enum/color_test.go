package enum

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	ColorX Color `json:"color"`
}

func TestColorString(t *testing.T) {
	x := Color(109)
	assert.Equal(t, "Color(109)Name", x.String())

	assert.Equal(t, Color(33), ColorGreen)
	assert.Equal(t, Color(34), ColorBlue)
	assert.Equal(t, &x, Color(109).Ptr())
}

func TestColorMustParse(t *testing.T) {
	x := `avocadogreen`

	assert.PanicsWithError(t, x+" is not a valid Color", func() { MustParseColor(x) })
	assert.NotPanics(t, func() { MustParseColor(ColorGreen.String()) })
	assert.NotPanics(t, func() { MustParseColor("BLack") })
}

func TestColorUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		output        *testData
		errorExpected bool
		err           error
	}{
		{
			name:          "black",
			input:         `{"color":"Black"}`,
			output:        &testData{ColorX: ColorBlack},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "blacklower",
			input:         `{"color":"black"}`,
			output:        &testData{ColorX: ColorBlack},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "white",
			input:         `{"color":"White"}`,
			output:        &testData{ColorX: ColorWhite},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "whitelower",
			input:         `{"color":"white"}`,
			output:        &testData{ColorX: ColorWhite},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "red",
			input:         `{"color":"Red"}`,
			output:        &testData{ColorX: ColorRed},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "redlower",
			input:         `{"color":"red"}`,
			output:        &testData{ColorX: ColorRed},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "green",
			input:         `{"color":"Green"}`,
			output:        &testData{ColorX: ColorGreen},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "greenlower",
			input:         `{"color":"green"}`,
			output:        &testData{ColorX: ColorGreen},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "blue",
			input:         `{"color":"Blue"}`,
			output:        &testData{ColorX: ColorBlue},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "bluelower",
			input:         `{"color":"blue"}`,
			output:        &testData{ColorX: ColorBlue},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "grey",
			input:         `{"color":"grey"}`,
			output:        &testData{ColorX: ColorGrey},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "greylower",
			input:         `{"color":"grey"}`,
			output:        &testData{ColorX: ColorGrey},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "yellow",
			input:         `{"color":"yellow"}`,
			output:        &testData{ColorX: ColorYellow},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "yellowlower",
			input:         `{"color":"yellow"}`,
			output:        &testData{ColorX: ColorYellow},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "yellow_green",
			input:         `{"color":"yellow_green"}`,
			output:        &testData{ColorX: ColorYellowGreen},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "magenta",
			input:         `{"color":"Magenta"}`,
			output:        &testData{ColorX: ColorYellow},
			errorExpected: true,
			err:           errors.New("Magenta is not a valid Color"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			x := &testData{}
			err := json.Unmarshal([]byte(test.input), x)
			if !test.errorExpected {
				require.NoError(tt, err, "failed unmarshalling the json.")
				assert.Equal(tt, test.output.ColorX, x.ColorX)
			} else {
				require.Error(tt, err)
				assert.EqualError(tt, err, test.err.Error())
			}
		})
	}
}

func TestColorMarshal(t *testing.T) {
	tests := []struct {
		name          string
		input         *testData
		output        string
		errorExpected bool
		err           error
	}{
		{
			name:          "black",
			output:        `{"color":"Black"}`,
			input:         &testData{ColorX: ColorBlack},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "white",
			output:        `{"color":"White"}`,
			input:         &testData{ColorX: ColorWhite},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "red",
			output:        `{"color":"Red"}`,
			input:         &testData{ColorX: ColorRed},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "green",
			output:        `{"color":"Green"}`,
			input:         &testData{ColorX: ColorGreen},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "blue",
			output:        `{"color":"Blue"}`,
			input:         &testData{ColorX: ColorBlue},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "grey",
			output:        `{"color":"grey"}`,
			input:         &testData{ColorX: ColorGrey},
			errorExpected: false,
			err:           nil,
		},
		{
			name:          "yellow",
			output:        `{"color":"yellow"}`,
			input:         &testData{ColorX: ColorYellow},
			errorExpected: false,
			err:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			raw, err := json.Marshal(test.input)
			require.NoError(tt, err, "failed marshalling to json")
			assert.JSONEq(tt, test.output, string(raw))
		})
	}
}

func BenchmarkColorParse(b *testing.B) {
	knownItems := []string{
		ColorRedOrangeBlue.String(),
		strings.ToLower(ColorRedOrangeBlue.String()),
		"BLack",
		// "2",  Leave this in to add an int as string parsing option in future.
	}

	var err error
	for _, item := range knownItems {
		b.Run(item, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err = ParseColor(item)
				assert.NoError(b, err)
			}
		})
	}
}

func TestColor_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    Color
		want bool
	}{
		{
			name: "TestIsValidWithValidColor",
			c:    Color(1), // assuming 1 is a valid Color, change as needed
			want: true,
		},
		{
			name: "TestIsValidWithInvalidColor",
			c:    Color(999), // assuming 999 is a not valid Color, change as needed
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("Color.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorName(t *testing.T) {
	testCases := []struct {
		desc string
		in   Color
		out  string
	}{
		{
			desc: "Valid Color",
			in:   ColorBlack,
			out:  "Black",
		},
		{
			desc: "Invalid Color",
			in:   Color(100),
			out:  "",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && tC.out != "" {
					t.Errorf("The code panicked %v", r)
				}
			}()
			got := tC.in.Name()
			if got != tC.out {
				t.Errorf("Expected Color name %v, got %v", tC.out, got)
			}
		})
	}
}
