package util

import "sync"

type SafeBoolArray struct {
	mu   *sync.Mutex
	vals []bool
}

func NewSafeBoolArray(size int) *SafeBoolArray {
	return &SafeBoolArray{
		mu:   &sync.Mutex{},
		vals: make([]bool, size),
	}
}

func (s *SafeBoolArray) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.vals {
		s.vals[i] = false
	}
}

func (s *SafeBoolArray) Get(idx int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.vals[idx]
}

func (s *SafeBoolArray) Set(idx int, val bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vals[idx] = val
}

func (s *SafeBoolArray) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.vals)
}
