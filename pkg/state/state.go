package state

import (
	"sync"
)

type ChildVisibleSet struct {
	mu  sync.RWMutex
	set map[string]struct{}
}

func New() *ChildVisibleSet {
	return &ChildVisibleSet{
		set: make(map[string]struct{}),
	}
}

func (s *ChildVisibleSet) Contains(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.set[name]
	return ok
}

func (s *ChildVisibleSet) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set = make(map[string]struct{})
}

func (s *ChildVisibleSet) Toggle(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.set[name]
	if ok {
		delete(s.set, name)
	} else {
		s.set[name] = struct{}{}
	}
}
