package structures

import "sync"

// StringSet the set of Strings
type StringSet struct {
	items map[string]bool
	lock  sync.RWMutex
}

// Add adds a new element to the Set. Returns a pointer to the Set.
func (s *StringSet) Add(t string) *StringSet {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items == nil {
		s.items = make(map[string]bool)
	}
	_, ok := s.items[t]
	if !ok {
		s.items[t] = true
	}
	return s
}

// Clear removes all elements from the Set
func (s *StringSet) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = make(map[string]bool)
}

// Delete removes the string from the Set and returns Has(string)
func (s *StringSet) Delete(item string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.items[item]
	if ok {
		delete(s.items, item)
	}
	return ok
}

// Has returns true if the Set contains the string
func (s *StringSet) Has(item string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.items[item]
	return ok
}

// Strings returns the string(s) stored
func (s *StringSet) Strings() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	items := []string{}
	for i := range s.items {
		items = append(items, i)
	}
	return items
}

// Size returns the size of the set
func (s *StringSet) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

// Package set creates an IntSet data structure for the int type

// IntSet the set of Ints
type IntSet struct {
	items map[int]bool
	lock  sync.RWMutex
}

// Add adds a new element to the Set. Returns a pointer to the Set.
func (s *IntSet) Add(t int) *IntSet {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items == nil {
		s.items = make(map[int]bool)
	}
	_, ok := s.items[t]
	if !ok {
		s.items[t] = true
	}
	return s
}

// Clear removes all elements from the Set
func (s *IntSet) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = make(map[int]bool)
}

// Delete removes the int from the Set and returns Has(int)
func (s *IntSet) Delete(item int) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.items[item]
	if ok {
		delete(s.items, item)
	}
	return ok
}

// Has returns true if the Set contains the int
func (s *IntSet) Has(item int) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.items[item]
	return ok
}

// Ints returns the int(s) stored
func (s *IntSet) Ints() []int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	items := []int{}
	for i := range s.items {
		items = append(items, i)
	}
	return items
}

// Size returns the size of the set
func (s *IntSet) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}
