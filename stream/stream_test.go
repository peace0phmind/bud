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

func TestStreamFlatMap(t *testing.T) {
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
				// try once again
				shuffled = stream.Shuffle()
				if reflect.DeepEqual(stream, shuffled) {
					t.Errorf("Shuffle() must return new stream object")
				}
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

func Test_Stream_Limit(t *testing.T) {
	testCases := []struct {
		name      string
		input     Stream[int]
		inputN    int
		want      Stream[int]
		wantError bool
	}{
		{
			name:      "LimitNegative",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    -2,
			want:      Of([]int{}),
			wantError: false,
		},
		{
			name:      "LimitZero",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    0,
			want:      Of([]int{}),
			wantError: false,
		},
		{
			name:      "LimitInRange",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    3,
			want:      Of([]int{1, 2, 3}),
			wantError: false,
		},
		{
			name:      "LimitUpperRange",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    5,
			want:      Of([]int{1, 2, 3, 4, 5}),
			wantError: false,
		},
		{
			name:      "LimitOverflow",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    10,
			want:      Of([]int{1, 2, 3, 4, 5}),
			wantError: false,
		},
		{
			name:      "LimitError",
			input:     Stream[int]{err: errors.New("example error")},
			inputN:    2,
			want:      Stream[int]{err: errors.New("example error")},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Limit(tc.inputN)

			// Check error cases
			if (got.err != nil) != tc.wantError {
				t.Errorf("Stream[int].Limit() error = %v, wantError %v", got.err, tc.wantError)
				return
			}
			if got.err != nil && tc.wantError {
				if got.err.Error() != tc.want.err.Error() {
					t.Errorf("Stream[int].Limit() error = %v, wantError %v", got.err, tc.want.err)
				}
				return
			}

			// Check valid cases
			gotSlice, _ := got.ToSlice()
			wantSlice, _ := tc.want.ToSlice()
			if !reflect.DeepEqual(gotSlice, wantSlice) {
				t.Errorf("Stream[int].Limit() = %v, want %v", gotSlice, wantSlice)
			}
		})
	}
}

func TestStream_Skip(t *testing.T) {
	testCases := []struct {
		name      string
		input     Stream[int]
		inputN    int
		want      Stream[int]
		wantError bool
	}{
		{
			name:      "skip 0",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    0,
			want:      Of([]int{1, 2, 3, 4, 5}),
			wantError: false,
		},
		{
			name:      "skip all elements",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    5,
			want:      Of([]int{}),
			wantError: false,
		},
		{
			name:      "skip 3 elements",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    3,
			want:      Of([]int{4, 5}),
			wantError: false,
		},
		{
			name:      "skip more than length",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    10,
			want:      Of([]int{}),
			wantError: false,
		},
		{
			name:      "negative skip",
			input:     Of([]int{1, 2, 3, 4, 5}),
			inputN:    -1,
			want:      Of([]int{1, 2, 3, 4, 5}),
			wantError: false,
		},
		{
			name:      "skip with error",
			input:     Stream[int]{elems: []int{1, 2, 3, 4, 5}, err: errors.New("skip error")},
			inputN:    2,
			want:      Stream[int]{err: errors.New("skip error")},
			wantError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Skip(tc.inputN)

			// Check error cases
			if (got.err != nil) != tc.wantError {
				t.Errorf("Stream[int].Skip() error = %v, wantError %v", got.err, tc.wantError)
				return
			}
			if got.err != nil && tc.wantError {
				if got.err.Error() != tc.want.err.Error() {
					t.Errorf("Stream[int].Skip() error = %v, wantError %v", got.err, tc.want.err)
				}
				return
			}

			// Check valid cases
			gotSlice, _ := got.ToSlice()
			wantSlice, _ := tc.want.ToSlice()
			if !reflect.DeepEqual(gotSlice, wantSlice) {
				t.Errorf("Stream[int].Skip() = %v, want %v", gotSlice, wantSlice)
			}
		})
	}
}

func TestAllMatch(t *testing.T) {
	testCases := []struct {
		name  string
		elems []int
		match func(int) (bool, error)
		want  bool
	}{
		{
			"empty slice",
			[]int{},
			func(i int) (bool, error) { return i > 0, nil },
			true,
		},
		{
			"all match",
			[]int{1, 2, 3},
			func(i int) (bool, error) { return i > 0, nil },
			true,
		},
		{
			"some not match",
			[]int{1, -1, 3},
			func(i int) (bool, error) { return i > 0, nil },
			false,
		},
		{
			"none match",
			[]int{-1, -2, -3},
			func(i int) (bool, error) { return i > 0, nil },
			false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := Stream[int]{elems: tc.elems}
			got, err := s.AllMatch(tc.match)
			if err == nil && got != tc.want {
				t.Errorf("AllMatch got = %v, want = %v", got, tc.want)
			}
		})
	}
}

func TestStream_AnyMatch(t *testing.T) {
	type testCase struct {
		name           string
		stream         Stream[int]
		matchFunc      func(int) (bool, error)
		expectedResult bool
		expectError    bool
	}

	testCases := []testCase{
		{
			name:           "Match Found",
			stream:         Stream[int]{elems: []int{1, 2, 3, 4, 5}},
			matchFunc:      func(n int) (bool, error) { return n > 3, nil },
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "Match Not Found",
			stream:         Stream[int]{elems: []int{1, 2, 3, 4, 5}},
			matchFunc:      func(n int) (bool, error) { return n > 5, nil },
			expectedResult: false,
			expectError:    false,
		},
		{
			name:           "Error In Match Function",
			stream:         Stream[int]{elems: []int{1, 2, 3, 4, 5}},
			matchFunc:      func(n int) (bool, error) { return false, errors.New("Unknown Error") },
			expectedResult: false,
			expectError:    true,
		},
		{
			name:           "Empty Stream",
			stream:         Stream[int]{elems: []int{}},
			matchFunc:      func(n int) (bool, error) { return n > 3, nil },
			expectedResult: false,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.stream.AnyMatch(tc.matchFunc)
			if err != nil && !tc.expectError {
				t.Fatalf("want no error, but got: %v", err)
			}
			if err == nil && tc.expectError {
				t.Fatalf("want error, but got none")
			}
			if result != tc.expectedResult {
				t.Fatalf("want %v, but got %v", tc.expectedResult, result)
			}
		})
	}
}

func compareInts(x, y int) (int, error) {
	return x - y, nil
}

func TestMustMax(t *testing.T) {
	testCases := []struct {
		name string
		init Stream[int]
		want int
		err  bool
	}{
		{
			name: "empty",
			init: Stream[int]{elems: []int{}},
			want: 0,
			err:  true,
		},
		{
			name: "one element",
			init: Stream[int]{elems: []int{20}},
			want: 20,
			err:  false,
		},
		{
			name: "duplicates",
			init: Stream[int]{elems: []int{4, 4, 4, 4}},
			want: 4,
			err:  false,
		},
		{
			name: "order",
			init: Stream[int]{elems: []int{10, 5, 15, 20}},
			want: 20,
			err:  false,
		},
		{
			name: "negative numbers",
			init: Stream[int]{elems: []int{-1, -2, -3, -4}},
			want: -1,
			err:  false,
		},
		{
			name: "mixed numbers",
			init: Stream[int]{elems: []int{-1, 2, -3, 4}},
			want: 4,
			err:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := recover(); (err != nil) != tc.err {
					t.Errorf("MustMax() recovered panic = %v, wantErr %v", err, tc.err)
				}
			}()
			if got := tc.init.MustMax(compareInts); got != tc.want {
				t.Errorf("MustMax() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestMustMin(t *testing.T) {
	testCases := []struct {
		name        string
		stream      Stream[int]
		compareFunc func(x, y int) (int, error)
		expected    int
		expectPanic bool
	}{
		{
			name:        "ValidMin",
			stream:      Of([]int{5, 2, 9, 11, 3}),
			compareFunc: func(x, y int) (int, error) { return x - y, nil },
			expected:    2,
			expectPanic: false,
		},
		{
			name:        "EmptyStream",
			stream:      Of([]int{}),
			compareFunc: func(x, y int) (int, error) { return x - y, nil },
			expected:    0,
			expectPanic: true,
		},
		{
			name:   "CompareFuncError",
			stream: Of([]int{5, 2, 9, 11, 3}),
			compareFunc: func(x, y int) (int, error) {
				return 0, errors.New("error in comparison")
			},
			expected:    2,
			expectPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil && !tc.expectPanic {
					t.Errorf("The code panicked, %+v", r)
				} else if r == nil && tc.expectPanic {
					t.Error("The code did not panic")
				}
			}()

			if result := tc.stream.MustMin(tc.compareFunc); result != tc.expected && !tc.expectPanic {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestMustFirst(t *testing.T) {
	testCases := []struct {
		name string
		data []int
		want int
		err  error
	}{
		{
			name: "Empty",
			data: []int{},
			want: 0,
			err:  errors.New("stream is empty"),
		},
		{
			name: "Single Element",
			data: []int{1},
			want: 1,
			err:  nil,
		},
		{
			name: "Multiple Elements",
			data: []int{5, 10, 15},
			want: 5,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := Of(tc.data[:])
			defer func() {
				if r := recover(); r != nil {
					if tc.err == nil {
						t.Errorf("MustFirst() panic = %v, no panic want", r)
					} else if r.(error).Error() != tc.err.Error() {
						t.Errorf("MustFirst() panic = %v, want panic = %v", r, tc.err.Error())
					}
				}
			}()
			got := s.MustFirst()
			if got != tc.want {
				t.Errorf("MustFirst() = %v, want = %v", got, tc.want)
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

func TestFlatMap(t *testing.T) {
	testCases := []struct {
		name    string
		input   Stream[int]
		flatMap func(int) Stream[int]
		want    Stream[int]
	}{
		{
			name:    "empty stream",
			input:   Of([]int{}),
			flatMap: func(x int) Stream[int] { return Of([]int{x, x * 2}) },
			want:    Of([]int{}),
		},
		{
			name:    "non-empty stream",
			input:   Of([]int{1, 2, 3}),
			flatMap: func(x int) Stream[int] { return Of([]int{x, x * 2}) },
			want:    Of([]int{1, 2, 2, 4, 3, 6}),
		},
		{
			name:    "mapToInt error",
			input:   Stream[int]{err: errors.New("stream error")},
			flatMap: func(x int) Stream[int] { return Of([]int{x, x * 2}) },
			want:    Stream[int]{err: errors.New("stream error")},
		},
		{
			name:  "flatMap error",
			input: Of([]int{1, 2, 3}),
			flatMap: func(x int) Stream[int] {
				if x == 2 {
					return Stream[int]{err: errors.New("flatMap error")}
				}
				return Of([]int{x, x * 2})
			},
			want: Stream[int]{err: errors.New("flatMap error")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := FlatMap(tc.input, tc.flatMap)

			if tc.want.err != nil {
				if got.err.Error() != tc.want.err.Error() {
					t.Errorf("FlatMap error: want: %v, got: %v", tc.want.err, got.err)
				} else {
					return
				}
			}

			for i, v := range got.elems {
				if v != tc.want.elems[i] {
					t.Errorf("Filter failed for case %s, want %#v, got %#v", tc.name, tc.want, got)
				}
			}
		})
	}
}

func TestToMap(t *testing.T) {
	// Simple util function to compare maps. This is necessary since the
	// use of generics means the map can't be compared directly
	compareMaps := func(map1, map2 map[int]string) bool {
		if len(map1) != len(map2) {
			return false
		}
		for k, v := range map1 {
			if map2[k] != v {
				return false
			}
		}
		return true
	}

	testCases := []struct {
		name      string
		stream    Stream[string]
		mapFunc   func(string) (int, string, error)
		want      map[int]string
		expectErr bool
	}{
		{
			name: "Valid Map 1",
			stream: Stream[string]{
				elems: []string{"one", "four", "three"},
			},
			mapFunc: func(s string) (int, string, error) {
				return len(s), s, nil
			},
			want: map[int]string{
				3: "one",
				4: "four",
				5: "three",
			},
			expectErr: false,
		},
		{
			name: "Valid Map 2",
			stream: Stream[string]{
				elems: []string{"one", "two", "three"},
			},
			mapFunc: func(s string) (int, string, error) {
				return len(s), s, nil
			},
			want: map[int]string{
				3: "two",
				5: "three",
			},
			expectErr: false,
		},
		{
			name: "Stream With Error",
			stream: Stream[string]{
				err: errors.New("stream test error"),
			},
			expectErr: true,
		},
		{
			name: "transform function Error",
			stream: Stream[string]{
				elems: []string{"one", "two", "three"},
			},
			mapFunc: func(s string) (int, string, error) {
				return 0, "", errors.New("map func test error")
			},
			expectErr: true,
		},
		{
			name:   "Empty Stream",
			stream: Stream[string]{},
			mapFunc: func(s string) (int, string, error) {
				return len(s), s, nil
			},
			want:      map[int]string{},
			expectErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gotMap, err := ToMap(tt.stream, tt.mapFunc)
			if (err != nil) != tt.expectErr {
				t.Errorf("ToMap() error = %v, expectErr %v", err, tt.expectErr)
			}
			if err == nil && !compareMaps(gotMap, tt.want) {
				t.Errorf("ToMap() gotMap = %v, want %v", gotMap, tt.want)
			}
		})
	}
}

func TestRange(t *testing.T) {

	testCases := []struct {
		name      string
		input     Stream[int]
		want      []int
		shouldErr bool
	}{
		{"empty stream", Stream[int]{elems: []int{}}, []int{}, false},
		{"one element", Stream[int]{elems: []int{1}}, []int{1}, false},
		{"multiple elements", Stream[int]{elems: []int{1, 2, 3}}, []int{1, 2, 3}, false},
		{"stream with error", Stream[int]{elems: []int{1, 2, 3}, err: errors.New("tc error")}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got = make([]int, 0)
			err := tc.input.Range(func(i int) error {
				got = append(got, i)
				return nil
			})
			if (err != nil) != tc.shouldErr {
				t.Errorf("got error = %v, want %v", err, tc.shouldErr)
			} else if err == nil && !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStream_Reverse(t *testing.T) {
	type fields struct {
		elems []interface{}
		err   error
	}
	testCases := []struct {
		name   string
		fields fields
		want   []interface{}
	}{
		{
			name: "Test Reverse string slice 3",
			fields: fields{
				elems: []interface{}{"one", "two", "three"},
				err:   nil,
			},
			want: []interface{}{"three", "two", "one"},
		},
		{
			name: "Test Reverse string slice 4",
			fields: fields{
				elems: []interface{}{"one", "two", "three", "four"},
				err:   nil,
			},
			want: []interface{}{"four", "three", "two", "one"},
		},
		{
			name: "Test Reverse int slice",
			fields: fields{
				elems: []interface{}{1, 2, 3},
				err:   nil,
			},
			want: []interface{}{3, 2, 1},
		},
		{
			name: "Test Reverse empty slice",
			fields: fields{
				elems: []interface{}{},
				err:   nil,
			},
			want: []interface{}{},
		},
		{
			name: "Test Reverse slice with error",
			fields: fields{
				elems: []interface{}{"one", "two", "three"},
				err:   errors.New("test error"),
			},
			want: nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := Stream[interface{}]{
				elems: tt.fields.elems,
				err:   tt.fields.err,
			}
			got := s.Reverse()
			if !reflect.DeepEqual(got.elems, tt.want) {
				t.Errorf("Stream.Reverse() = %v, want %v", got.elems, tt.want)
			}
		})
	}
}
