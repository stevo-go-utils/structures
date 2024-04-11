package structures

import "sync"

type Set[T comparable] struct {
	items map[T]bool
	lock  sync.RWMutex
}

func NewSet[T comparable](items ...T) *Set[T] {
	m := map[T]bool{}
	for _, item := range items {
		m[item] = true
	}
	return &Set[T]{
		items: m,
	}
}

// Add adds a new item(s) to the Set. Returns a pointer to the Set.
func (s *Set[T]) Add(items ...T) *Set[T] {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items == nil {
		s.items = make(map[T]bool)
	}
	for _, item := range items {
		_, ok := s.items[item]
		if !ok {
			s.items[item] = true
		}
	}
	return s
}

// Clear removes all items from the Set
func (s *Set[T]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = make(map[T]bool)
}

// Delete removes the item from the Set and returns Has(item)
func (s *Set[T]) Delete(item T) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.items[item]
	if ok {
		delete(s.items, item)
	}
	return ok
}

// Has returns true if the Set contains the item
func (s *Set[T]) Has(item T) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.items[item]
	return ok
}

// Vals returns the val(s) stored
func (s *Set[T]) Vals() []T {
	s.lock.RLock()
	defer s.lock.RUnlock()
	items := []T{}
	for i := range s.items {
		items = append(items, i)
	}
	return items
}

// Size returns the size of the set
func (s *Set[T]) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}
