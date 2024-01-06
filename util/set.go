package util

var null struct{}

type Set[K comparable] struct {
	values Cache[K, any]
}

func NewSet[K comparable]() *Set[K] {
	ret := &Set[K]{}

	return ret
}

func (s *Set[K]) AddValue(v K) bool {
	if IsNil(v) {
		return false
	}

	_, loaded := s.values.GetOrStore(v, null)
	return !loaded
}

func (s *Set[K]) Contain(v K) bool {
	if IsNil(v) {
		return false
	}

	_, ok := s.values.Get(v)
	return ok
}

func (s *Set[K]) DeleteValue(v K) bool {
	if IsNil(v) {
		return false
	}

	_, loaded := s.values.GetAndDelete(v)
	return loaded
}

func (s *Set[K]) Size() int {
	return s.values.Size()
}

func (s *Set[K]) DoEach(f func(v K)) {
	s.values.Range(func(v K, Any any) bool {
		f(v)
		return true
	})
}

func (s *Set[K]) DoEachAsync(f func(v K)) {
	s.values.Range(func(v K, Any any) bool {
		go f(v)
		return true
	})
}
