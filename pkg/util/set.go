package util

// Set is a generic set implementation.
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet creates a new set with the given items.
func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{
		m: make(map[T]struct{}),
	}

	s.Add(items...)
	return s
}

// Add adds items to the set.
func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

// Remove removes items from the set.
func (s *Set[T]) Remove(items ...T) {
	for _, item := range items {
		delete(s.m, item)
	}
}

// Contains returns true if the set contains the given item.
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

// Len returns the number of items in the set.
func (s *Set[T]) Len() int {
	return len(s.m)
}

// Items returns the items in the set.
func (s *Set[T]) Items() []T {
	items := make([]T, 0, len(s.m))
	for item := range s.m {
		items = append(items, item)
	}

	return items
}

// ForEach calls the given function for each item in the set until the function returns false or the set is exhausted.
func (s *Set[T]) ForEach(f func(item T) (shouldContinue bool)) {
	for item := range s.m {
		if !f(item) {
			return
		}
	}
}

// Equal returns true if the set contains the same items as the given set.
func (s *Set[T]) Equal(other *Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}

	for item := range s.m {
		if !other.Contains(item) {
			return false
		}
	}

	return true
}

// HasOnly returns true if the set contains only the given items.
func (s *Set[T]) HasOnly(items ...T) bool {
	if len(items) != s.Len() {
		return false
	}

	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}

	return true
}

// HasAll returns true if the set contains all of the given items.
func (s *Set[T]) HasAll(items ...T) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}

	return true
}

// HasAny returns true if the set contains any of the given items.
func (s *Set[T]) HasAny(items ...T) bool {
	for _, item := range items {
		if s.Contains(item) {
			return true
		}
	}

	return false
}
