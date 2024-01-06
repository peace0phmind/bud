package stream

import (
	"math/rand"
	"time"
)

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
	if s.err != nil {
		return s
	}

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
	if s.err != nil {
		return s
	}

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
func (s *Stream[T]) AllMatch(predicate func(T) bool) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for _, elem := range s.elems {
		if !predicate(elem) {
			return false, nil
		}
	}

	return true, nil
}

// AnyMatch checks if any element in the stream satisfies the given predicate.
// It iterates over each element in the stream and applies the predicate function to it.
// If the predicate returns true for any element, the method returns true.
// If the predicate returns false for all elements, the method returns false.
// The original stream is not modified.
// The predicate function should return true for elements that satisfy the condition and false otherwise.
// Returns true if any element in the stream satisfies the predicate, false otherwise.
func (s *Stream[T]) AnyMatch(predicate func(T) bool) (bool, error) {
	if s.err != nil {
		return false, s.err
	}

	for _, elem := range s.elems {
		if predicate(elem) {
			return true, nil
		}
	}
	return false, nil
}

// Shuffle randomly rearranges the elements in the stream.
// It creates a new stream with the same elements as the original stream, but in a random order.
// The original stream remains unchanged.
// The new stream is returned as a pointer to Stream[T].
// The shuffle algorithm used is the Fisher-Yates shuffle.
// The seed for the random number generator is set using the current time.
// Example usage:
//
//	stream := Of([]int{1, 2, 3, 4, 5})
//	shuffled := stream.Shuffle()
//	shuffledElems := shuffled.ToSlice()
//	fmt.Println(shuffledElems)  // Output: [4 3 1 2 5]
func (s *Stream[T]) Shuffle() *Stream[T] {
	if s.err != nil {
		return s
	}

	//Create a new Stream and copy the data from the original Stream over
	newStream := &Stream[T]{elems: append([]T(nil), s.elems...)}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for n := len(newStream.elems); n > 0; n-- {
		randIndex := r.Intn(n)
		newStream.elems[n-1], newStream.elems[randIndex] = newStream.elems[randIndex], newStream.elems[n-1]
	}

	return newStream
}

// ToSlice returns a slice containing all the elements of the stream.
// The original stream is not modified.
// The elements in the returned slice are in the same order as in the original stream.
// The returned slice has the type []T, where T is the type of elements in the stream.
// Example usage:
//
//		stream := Of([]int{1, 2, 3})
//	 result := stream.ToSlice() // result is []int{1, 2, 3}
func (s *Stream[T]) ToSlice() []T {
	if s.err != nil {
		var ret []T
		return ret
	}

	return s.elems
}

// ToAny converts the elements of the stream to the `any` type and returns them as a slice.
// It creates a new slice and appends the converted elements of the stream to it.
// The original stream is not modified.
// The elements in the resulting slice follow the same order as in the original stream.
// The resulting slice is returned as a value of type `[]any`.
func (s *Stream[T]) ToAny() []any {
	if s.err != nil {
		var ret []any
		return ret
	}

	var result []any

	for _, v := range s.elems {
		result = append(result, any(v))
	}

	return result
}

// Err returns the error associated with the stream.
// It retrieves the error value stored in the 'err' field of the Stream struct.
// This method can be used to check if an error occurred during stream processing.
// If no error occurred, it returns nil.
// The error value is returned as an instance of the 'error' interface.
// Example usage:
//
//	stream := &Stream[T]{}
//	err := stream.Err()
//	if err != nil {
//	    fmt.Println("An error occurred:", err.Error())
//	}
func (s *Stream[T]) Err() error {
	return s.err
}

// Range iterates over each element in the stream and applies the forEach function to it.
// If the forEach function returns false for any element, the iteration is stopped.
// The forEach function should return true for elements that need to be processed, and false for elements that can be skipped.
// This method does not modify the original stream.
// The elements are iterated in the same order as in the stream.
// This method does not return any value.
func (s *Stream[T]) Range(forEach func(T) bool) error {
	if s.err != nil {
		return s.err
	}

	for _, elem := range s.elems {
		if !forEach(elem) {
			break
		}
	}

	return nil
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

// Map applies the provided function `f` to each element of the input stream `s`
// and returns a new stream containing the resulting elements. If an error occurs during
// the mapping process, the resulting stream will have the corresponding error value.
//
// Example:
//
//	str := Of([]int{1, 2, 3, 4})
//	double := func(i int) (int, error) {
//	    return i * 2, nil
//	}
//	doubledStr := Map(str, double)
//	fmt.Println(doubledStr.ToSlice()) // Output: [2, 4, 6, 8]
//
//	str2 := Of([]string{"hello", "world"})
//	length := func(s string) (int, error) {
//	    return len(s), nil
//	}
//	lengthStr := Map(str2, length)
//	fmt.Println(lengthStr.ToSlice()) // Output: [5, 5]
func Map[In any, Out any](s *Stream[In], f func(In) (Out, error)) *Stream[Out] {
	var result Stream[Out]
	if s.err != nil {
		result.err = s.err
		return &result
	}

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
