package stream

import (
	"errors"
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	testCases := []struct {
		name     string
		initial  Stream[int]
		keepFunc func(int) (bool, error)
		want     Stream[int]
	}{
		{
			"EmptyStream",
			Stream[int]{elems: []int{}},
			func(v int) (bool, error) { return v%2 == 0, nil },
			Stream[int]{elems: []int{}},
		},
		{
			"AllElementsFilteredOut",
			Stream[int]{elems: []int{1, 3, 5, 7, 9}},
			func(v int) (bool, error) { return v%2 == 0, nil },
			Stream[int]{elems: []int{}},
		},
		{
			"SomeElementsFilteredOut",
			Stream[int]{elems: []int{1, 2, 3, 4, 5}},
			func(v int) (bool, error) { return v%2 == 0, nil },
			Stream[int]{elems: []int{2, 4}},
		},
		{
			"NoElementsFilteredOut",
			Stream[int]{elems: []int{2, 4, 6, 8, 10}},
			func(v int) (bool, error) { return v%2 == 0, nil },
			Stream[int]{elems: []int{2, 4, 6, 8, 10}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.initial.Filter(tc.keepFunc)

			for i, v := range got.elems {
				if v != tc.want.elems[i] {
					t.Errorf("Filter failed for case %s, want %#v, got %#v", tc.name, tc.want, got)
				}
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	testCases := []struct {
		name    string
		stream  Stream[int]
		flatFun func(int) Stream[int]
		want    Stream[int]
		wantErr bool
	}{
		{
			name: "empty stream",
			stream: Stream[int]{
				elems: []int{},
			},
			flatFun: func(n int) Stream[int] {
				return Stream[int]{elems: []int{n, n}}
			},
			want:    Stream[int]{elems: []int{}},
			wantErr: false,
		},
		{
			name: "stream with error",
			stream: Stream[int]{
				elems: []int{1, 2, 3},
				err:   errors.New("previous error"),
			},
			flatFun: func(n int) Stream[int] {
				return Stream[int]{elems: []int{n, n}}
			},
			want:    Stream[int]{elems: []int{}, err: errors.New("previous error")},
			wantErr: true,
		},
		{
			name: "flat function returns error",
			stream: Stream[int]{
				elems: []int{1, 2, 3},
			},
			flatFun: func(n int) Stream[int] {
				return Stream[int]{elems: []int{n, n}, err: errors.New("flat function error")}
			},
			want:    Stream[int]{elems: []int{}, err: errors.New("flat function error")},
			wantErr: true,
		},
		{
			name: "normal case",
			stream: Stream[int]{
				elems: []int{1, 2, 3},
			},
			flatFun: func(n int) Stream[int] {
				return Stream[int]{elems: []int{n, n}}
			},
			want:    Stream[int]{elems: []int{1, 1, 2, 2, 3, 3}},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.stream.FlatMap(tc.flatFun)
			for i, v := range got.elems {
				if v != tc.want.elems[i] {
					t.Errorf("Filter failed for case %s, want %#v, got %#v", tc.name, tc.want, got)
				}
			}

			if (got.Err() != nil) != tc.wantErr {
				t.Errorf("FlatMap() error = %v, wantErr %v", got.Err(), tc.wantErr)
			}
		})
	}
}

func TestStream_Shuffle(t *testing.T) {
	// Function for generating sample data
	generateData := func(n int) []int {
		data := make([]int, n)
		for i := 0; i < n; i++ {
			data[i] = i
		}
		return data
	}

	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stream := Of(tc.elems)
			shuffled := stream.Shuffle()

			// Check if shuffle returns a new stream object
			if reflect.DeepEqual(stream, shuffled) {
				t.Errorf("Shuffle() must return new stream object")
			}

			shuffledElems := shuffled.MustToSlice()

			// Check the number of elements is same in the original and shuffled stream
			if got, want := len(tc.elems), len(shuffledElems); got != want {
				t.Errorf("len(shuffled) got %v, want %v", got, want)
			}

			if len(tc.elems) < 10 {
				t.Logf("shuffled got %v, origin %v", shuffledElems, tc.elems)
			}

			// Check at least one element is in a different position
			var found bool
			for i, v := range tc.elems {
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

	testCases := []struct {
		name    string
		stream  Stream[int]
		funcMap func(int) (int, error)
		want    Stream[int]
		err     error
	}{
		{
			name:    "increment",
			stream:  Of[int]([]int{1, 2, 3}),
			funcMap: increment,
			want:    Of[int]([]int{2, 3, 4}),
		},
		{
			name:    "empty",
			stream:  Of[int]([]int{}),
			funcMap: increment,
			want:    Of[int]([]int{}),
		},
		{
			name:    "errorFunction",
			stream:  Of[int]([]int{1, 2, 3}),
			funcMap: errFunc,
			err:     errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Map[int, int](tc.stream, tc.funcMap)

			if tc.err != nil {
				if result.err.Error() != tc.err.Error() {
					t.Errorf("got error %q, want %q", result.err, tc.err)
				}
				return
			}

			if len(result.elems) != len(tc.want.elems) {
				t.Errorf("got length %d, want %d", len(result.elems), len(tc.want.elems))
				return
			}

			for i, v := range result.elems {
				if v != tc.want.elems[i] {
					t.Errorf("at index %d: got %v, want %v", i, v, tc.want.elems[i])
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

	testCases := []struct {
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
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := GroupBy(tc.args.s, tc.args.getKey); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("GroupBy() = %v, want %v", got, tc.want)
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

func TestStreamSort(t *testing.T) {
	// Comparing function.
	compare := func(a, b int) int {
		return a - b
	}
	testCases := []struct {
		name  string
		input []int
		want  []int
	}{
		{"Empty Slice", []int{}, []int{}},
		{"Single Element", []int{1}, []int{1}},
		{"Two Elements Sorted", []int{1, 2}, []int{1, 2}},
		{"Two Elements Unsorted", []int{2, 1}, []int{1, 2}},
		{"Multiple Elements", []int{3, 1, 2}, []int{1, 2, 3}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stream := Of(tc.input).Sort(compare)
			result := stream.MustToSlice()

			if len(tc.want) != len(result) {
				t.Fatalf("want length %v but got %v", len(tc.want), len(result))
			}

			for i, v := range tc.want {
				if v != result[i] {
					t.Fatalf("at index %d, want %v but got %v", i, v, result[i])
				}
			}
		})
	}
}

func TestMustReduceWithInit(t *testing.T) {
	accumulator := func(preItem, nextItem int) (int, error) {
		return preItem + nextItem, nil
	}

	testCases := []struct {
		name  string
		elems []int
		init  int
		want  int
	}{
		{
			name:  "EmptyStream",
			elems: []int{},
			init:  0,
			want:  0,
		},
		{
			name:  "SingleElement",
			elems: []int{5},
			init:  0,
			want:  5,
		},
		{
			name:  "MultipleElements",
			elems: []int{1, 2, 3, 4, 5},
			init:  0,
			want:  15,
		},
		{
			name:  "MultipleElementsWithInit",
			elems: []int{1, 2, 3, 4, 5},
			init:  10,
			want:  25,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := Of(tc.elems)
			got := s.MustReduceWithInit(tc.init, accumulator)
			if got != tc.want {
				t.Errorf("MustReduceWithInit() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDistinct(t *testing.T) {
	equalInt := func(preItem, nextItem int) (bool, error) {
		return preItem == nextItem, nil
	}

	testCases := []struct {
		name     string
		elems    []int
		equalFun func(preItem, nextItem int) (bool, error)
		want     []int
	}{
		{
			name:     "EmptySlice",
			elems:    []int{},
			equalFun: equalInt,
			want:     []int{},
		},
		{
			name:     "NoDuplicates",
			elems:    []int{1, 2, 3, 4},
			equalFun: equalInt,
			want:     []int{1, 2, 3, 4},
		},
		{
			name:     "AllDuplicates",
			elems:    []int{2, 2, 2, 2, 2},
			equalFun: equalInt,
			want:     []int{2},
		},
		{
			name:     "SomeDuplicates",
			elems:    []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			equalFun: equalInt,
			want:     []int{1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := Of(tc.elems)
			distinctS := s.Distinct(tc.equalFun)
			got, _ := distinctS.ToSlice()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Stream.Distinct() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFindFirst(t *testing.T) {
	testCases := []struct {
		name    string
		stream  Stream[int]
		keep    func(int) (bool, error)
		want    int
		wantErr bool
	}{
		{
			name:   "FindFirstPositiveNumberInIntStream",
			stream: Of([]int{0, -1, -3, 10, -2, 100}),
			keep: func(i int) (bool, error) {
				return i > 0, nil
			},
			want:    10,
			wantErr: false,
		},
		{
			name:   "FindFirstError",
			stream: Stream[int]{elems: []int{0, -1, -3, 10, -2, 100}, err: errors.New("test error")},
			keep: func(i int) (bool, error) {
				return i > 0, nil
			},
			wantErr: true,
		},
		{
			name:   "FindFirstInEmptyStream",
			stream: Stream[int]{elems: []int{}},
			keep: func(i int) (bool, error) {
				return i > 0, nil
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.stream.FindFirst(tc.keep)
			if (err != nil) != tc.wantErr {
				t.Errorf("FindFirst() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && got != tc.want {
				t.Errorf("FindFirst() got = %v, want %v", got, tc.want)
			}
		})
	}
}
