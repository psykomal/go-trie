package main

import "sync"

type MapStore struct {
	mu    sync.RWMutex
	kvMap map[string]string
}

func NewMapStore() *MapStore {
	return &MapStore{
		kvMap: make(map[string]string),
	}
}

func (s *MapStore) Set(key, value string) error {
	if key == "" {
		return ErrKeyIsEmpty
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.kvMap[key] = value
	return nil
}

func (s *MapStore) Get(key string) (string, error) {
	if key == "" {
		return "", ErrKeyIsEmpty
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.kvMap[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	return value, nil
}

func (s *MapStore) Delete(key string) error {
	if key == "" {
		return ErrKeyIsEmpty
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.kvMap, key)
	return nil
}
