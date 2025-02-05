package generic

import "sync"

type Set[T comparable] struct {
	rw       sync.RWMutex
	instance map[T]struct{}
}

func (s *Set[T]) Add(items ...T) {
	s.rw.Lock()
	defer s.rw.Unlock()
	for _, it := range items {
		s.instance[it] = struct{}{}
	}
}

func (s *Set[T]) Delete(items ...T) {
	s.rw.Lock()
	defer s.rw.Unlock()
	for _, it := range items {
		delete(s.instance, it)
	}
}

func (s *Set[T]) Has(item T) bool {
	s.rw.RLock()
	defer s.rw.RUnlock()
	_, ok := s.instance[item]
	return ok
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		instance: make(map[T]struct{}),
	}
}
