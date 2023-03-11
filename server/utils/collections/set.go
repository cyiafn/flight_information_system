package collections

import (
	"sync"

	"github.com/pkg/errors"
)

// Set is a generic concurrent-safe implementation
type Set[T comparable] struct {
	sync.RWMutex
	data map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}),
	}
}

func (s *Set[T]) Has(a T) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.data[a]
	return ok
}

func (s *Set[T]) Add(a T) error {
	if s.Has(a) {
		return errors.Errorf("element already exists")
	}
	s.MustAdd(a)
	return nil
}

// MustAdd simply overwrites the previous data and does not throw an error even if key exists while add does
func (s *Set[T]) MustAdd(a T) {
	s.Lock()
	defer s.Unlock()
	s.data[a] = struct{}{}
}

func (s *Set[T]) Remove(a T) error {
	if !s.Has(a) {
		return errors.Errorf("element does not exist")
	}
	s.MustRemove(a)
	return nil
}

// MustRemove will remove whether the key is present or not without throwing an error.
func (s *Set[T]) MustRemove(a T) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, a)
}

// ToList converts the set to an array
func (s *Set[T]) ToList() []T {
	if len(s.data) == 0 {
		return nil
	}
	s.RLock()
	defer s.RUnlock()

	output := make([]T, len(s.data))

	ptr := 0
	for key := range s.data {
		output[ptr] = key
		ptr += 1
	}

	return output
}

func (s *Set[T]) Len() int {
	return len(s.data)
}
