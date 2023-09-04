package main

import (
	"sync"
)

type trieStore struct {
	rootMu sync.RWMutex // used when swapping old root with new
	cowMu  sync.Mutex   // used for COW operations - Set and Delete
	root   *COWNode     // COW Trie Node
}

func NewTrieStore() *trieStore {
	return &trieStore{
		root: NewCOWNode(),
	}
}

func (s *trieStore) Set(key string, value string) error {
	if key == "" {
		return ErrKeyIsEmpty
	}

	s.cowMu.Lock()
	defer s.cowMu.Unlock()

	new_root, err := s.root.Set(key, value)
	if err != nil {
		return err
	}

	s.rootMu.Lock()
	defer s.rootMu.Unlock()

	s.root = new_root
	return nil
}

func (s *trieStore) Get(key string) (string, error) {
	if key == "" {
		return "", ErrKeyIsEmpty
	}

	s.rootMu.RLock()
	defer s.rootMu.RUnlock()

	return s.root.Get(key)
}

func (s *trieStore) Delete(key string) error {
	if key == "" {
		return ErrKeyIsEmpty
	}

	s.cowMu.Lock()
	defer s.cowMu.Unlock()

	new_root, err := s.root.Delete(key)
	if err != nil {
		return err
	}

	s.rootMu.Lock()
	defer s.rootMu.Unlock()

	s.root = new_root
	return nil
}
