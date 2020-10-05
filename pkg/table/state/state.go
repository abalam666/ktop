package state

import (
	"sync"
)

type VisibleSet struct {
	mu           sync.RWMutex
	childVisible map[string]struct{}
}

// name is for any one of node, pod or container.
func (s *VisibleSet) Add(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.childVisible[name] = struct{}{}
}

func (s *VisibleSet) Remove(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.childVisible, name)
}

func (s *VisibleSet) Contains(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.childVisible[name]
	return ok
}
