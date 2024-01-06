package stream

// Stream is a generic type representing a stream of elems of type T.
// It contains a slice of elems of type T.
type Stream[T any] struct {
	elems []T
	err   error
}

// Of creates a new stream with the given elems.
func Of[T any](values []T) *Stream[T] {
	return &Stream[T]{elems: values}
}

// Filter filters the stream by applying the keep function to each element.
// It creates a new stream containing only the elements for which keep returns true.
// The original stream is not modified.
// The elements of the new stream are in the same order as in the original stream.
// The keep function should return true for elements that should be kept in the new stream, and false for elements that should be excluded.
// The new stream is returned as a pointer to Stream[T].
func (s *Stream[T]) Filter(keep func(T) bool) *Stream[T] {
	var result Stream[T]

	for _, v := range s.elems {
		if keep(v) {
			result.elems = append(result.elems, v)
		}
	}

	return &result
}

// Append appends the given values to the stream.
// It modifies the stream by adding the values to the end of its elems slice.
// The original stream is modified.
// The elements are appended in the order they are supplied.
// The appended values can be of any type specified by T in the stream declaration.
// The modified stream is returned as a pointer to Stream[T].
func (s *Stream[T]) Append(values ...T) *Stream[T] {
	s.elems = append(s.elems, values...)
	return s
}

// AllMatch returns true if all elements in the stream satisfy the given predicate function.
//
// It iterates through each element in the stream and applies the predicate function to determine if the element satisfies the condition.
// If any element fails the condition, the function immediately returns false.
// If all elements pass the condition, the function returns true.
//
// The original stream is not modified.
// The predicate function should return true for elements that satisfy the condition, and false for elements that do not.
// The stream is a pointer to Stream[T] type.
//
// Example usage:
//
//	stream := Of([]T{1, 2, 3, 4, 5})
//	result := stream.AllMatch(func(elem T) bool {
//	  return elem > 0
//	})
//	// result is true, since all elements in the stream are greater than 0
//
// Note: The elements of the stream should be of the same type as the type specified for Stream[T].
// For example, if the Stream[T] is created with Stream[int], the elements should be of type int.
// The behavior of the method is undefined if this condition is violated.
func (s *Stream[T]) AllMatch(predicate func(T) bool) bool {
	for _, elem := range s.elems {
		if !predicate(elem) {
			return false
		}
	}

	return true
}

func (s *Stream[T]) ToSlice() []T {
	return s.elems
}

func (s *Stream[T]) ToAny() []any {
	var result []any

	for _, v := range s.elems {
		result = append(result, any(v))
	}

	return result
}

func (s *Stream[T]) Err() error {
	return s.err
}

func (s *Stream[T]) Range(forEach func(T) bool) {
	for _, elem := range s.elems {
		if !forEach(elem) {
			break
		}
	}
}

// GroupBy groups the elements of the input stream based on the provided key function.
// It returns a map where each key corresponds to a group, and the value is a stream
// containing the elements that belong to that group.
// The key function determines the grouping criteria by extracting a key from each element.
// If two elements have the same key, they will belong to the same group.
//
// Example:
//
//	s := Of([]int{1, 2, 3, 4, 5})
//	groups := GroupBy(s, func(num int) string {
//	    if num%2 == 0 {
//	        return "even"
//	    }
//	    return "odd"
//	})
//
//	// The resulting groups map will be:
//	// {
//	//    "even": {elems: [2, 4]},
//	//    "odd": {elems: [1, 3, 5]},
//	// }
//
// If the input stream is empty, the result will be an empty map.
//
// Parameters:
// - s: The input stream to group.
// - getKey: A function that extracts the key from each element in the stream.
//
// Returns:
//   - A map where each key corresponds to a group, and the value is a stream
//     containing the elements that belong to that group.
func GroupBy[T any, K comparable](s *Stream[T], getKey func(T) K) map[K]*Stream[T] {
	result := make(map[K]*Stream[T])

	for _, v := range s.elems {
		key := getKey(v)
		if _, ok := result[key]; !ok {
			result[key] = Of([]T{v})
		} else {
			result[key].Append(v)
		}
	}

	return result
}

func Map[In any, Out any](s *Stream[In], f func(In) (Out, error)) *Stream[Out] {
	var result Stream[Out]

	for _, v := range s.elems {
		elem, err := f(v)
		if err != nil {
			result.err = err
			return &result
		}
		result.elems = append(result.elems, elem)
	}

	return &result
}
