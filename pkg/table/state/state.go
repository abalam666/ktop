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

// name is for any one of node, pod or container.
func (s *ChildVisibleSet) Switch(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.set[name]
	if ok {
		delete(s.set, name)
	} else {
		s.set[name] = struct{}{}
	}
}

func (s *ChildVisibleSet) Contains(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.set[name]
	return ok
}
