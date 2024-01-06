package stream

import (
	"errors"
	"reflect"
	"testing"
)

// Function for generating sample data
func generateData(n int) []int {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = i
	}
	return data
}

func TestStream_Shuffle(t *testing.T) {
	tests := []struct {
		name  string
		elems []int
	}{
		{
			name:  "Three Elements",
			elems: []int{1, 2, 3},
		},
		{
			name:  "Four Elements",
			elems: []int{1, 2, 3, 4},
		},
		{
			name:  "Multiple Elements",
			elems: generateData(1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := Of(tt.elems)
			shuffled := stream.Shuffle()

			// Check if shuffle returns a new stream object
			if reflect.DeepEqual(stream, shuffled) {
				t.Errorf("Shuffle() must return new stream object")
			}

			shuffledElems := shuffled.MustToSlice()

			// Check the number of elements is same in the original and shuffled stream
			if got, want := len(tt.elems), len(shuffledElems); got != want {
				t.Errorf("len(shuffled) got %v, want %v", got, want)
			}

			if len(tt.elems) < 10 {
				t.Logf("shuffled got %v, origin %v", shuffledElems, tt.elems)
			}

			// Check at least one element is in a different position
			var found bool
			for i, v := range tt.elems {
				if v != shuffledElems[i] {
					found = true
					break
				}
			}

			if !found {
				t.Error("Shuffle() should alter the order of the elements")
			}
		})
	}
}

func TestMap(t *testing.T) {
	increment := func(n int) (int, error) {
		return n + 1, nil
	}

	errFunc := func(n int) (int, error) {
		return 0, errors.New("Error")
	}

	tests := []struct {
		name     string
		stream   Stream[int]
		funcMap  func(int) (int, error)
		expected Stream[int]
		err      error
	}{
		{
			name:     "increment",
			stream:   Of[int]([]int{1, 2, 3}),
			funcMap:  increment,
			expected: Of[int]([]int{2, 3, 4}),
		},
		{
			name:     "empty",
			stream:   Of[int]([]int{}),
			funcMap:  increment,
			expected: Of[int]([]int{}),
		},
		{
			name:    "errorFunction",
			stream:  Of[int]([]int{1, 2, 3}),
			funcMap: errFunc,
			err:     errors.New("Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map[int, int](tt.stream, tt.funcMap)

			if tt.err != nil {
				if result.err.Error() != tt.err.Error() {
					t.Errorf("got error %q, want %q", result.err, tt.err)
				}
				return
			}

			if len(result.elems) != len(tt.expected.elems) {
				t.Errorf("got length %d, want %d", len(result.elems), len(tt.expected.elems))
				return
			}

			for i, v := range result.elems {
				if v != tt.expected.elems[i] {
					t.Errorf("at index %d: got %v, want %v", i, v, tt.expected.elems[i])
				}
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	type args struct {
		s      Stream[int]
		getKey func(int) int
	}
	tests := []struct {
		name string
		args args
		want map[int]Stream[int]
	}{
		{
			name: "Test with empty stream",
			args: args{
				s: Of([]int{}),
				getKey: func(i int) int {
					return i % 2
				},
			},
			want: make(map[int]Stream[int]),
		},
		{
			name: "Test with non empty stream",
			args: args{
				s: Of([]int{1, 2, 3, 4, 5}),
				getKey: func(i int) int {
					return i % 2
				},
			},
			want: map[int]Stream[int]{
				// Assuming that Of and Append works properly.
				0: Of([]int{2, 4}),
				1: Of([]int{1, 3, 5}),
			},
		},
		{
			name: "Test with all same values",
			args: args{
				s: Of([]int{1, 1, 1, 1}),
				getKey: func(i int) int {
					return i % 2
				},
			},
			want: map[int]Stream[int]{
				// Assuming that Of and Append works properly.
				1: Of([]int{1, 1, 1, 1}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupBy(tt.args.s, tt.args.getKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupByGetNotExistKey(t *testing.T) {
	ret := GroupBy(Of([]int{1, 2, 3, 4, 5}), func(i int) int {
		return i % 2
	})

	println(ret[3].Size())
}
