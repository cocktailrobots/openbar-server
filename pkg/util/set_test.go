package util

import "testing"

func TestSet_Add(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("expected set to have 3 items, got %d", s.Len())
	}
}

func TestSet_Remove(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)
	s.Remove(2)
	if s.Len() != 2 {
		t.Errorf("expected set to have 2 items, got %d", s.Len())
	}

	s.Remove(8)
	if s.Len() != 2 {
		t.Errorf("expected set to have 2 items, got %d", s.Len())
	}
}

func TestSet_Contains(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)
	if !s.Contains(2) {
		t.Errorf("expected set to contain 2")
	}

	if s.Contains(4) {
		t.Errorf("expected set to not contain 4")
	}
}

func TestSet_Len(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("expected set to have 3 items, got %d", s.Len())
	}
}

func TestSet_Items(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)
	items := s.Items()
	if len(items) != 3 {
		t.Errorf("expected set to have 3 items, got %d", len(items))
	}
}

func TestSet_ForEach(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)
	sum := 0
	s.ForEach(func(item int) bool {
		sum += item
		return true
	})
	if sum != 6 {
		t.Errorf("expected sum to be 6, got %d", sum)
	}
}


