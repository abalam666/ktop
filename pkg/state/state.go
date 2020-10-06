package state

import (
	"sync"
)

type TableVisibleSet struct {
	mu  sync.RWMutex
	set map[string]struct{}
}

func NewTableVisibleSet() *TableVisibleSet {
	return &TableVisibleSet{
		set: make(map[string]struct{}),
	}
}

func (s *TableVisibleSet) Contains(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.set[name]
	return ok
}

func (s *TableVisibleSet) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set = make(map[string]struct{})
}

func (s *TableVisibleSet) Toggle(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.set[name]
	if ok {
		delete(s.set, name)
	} else {
		s.set[name] = struct{}{}
	}
}
