package cache

import "sync"

// InMemory struct
type InMemory struct {
	hash map[string][]byte
	mu   sync.RWMutex
}

// NewInMemory initializes in-memory structure
func NewInMemory() *InMemory {
	return &InMemory{hash: make(map[string][]byte)}
}

// Get cache value
func (im *InMemory) Get(key string) ([]byte, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	if v, found := im.hash[key]; found {
		return v, nil
	}

	return []byte{}, errNotFound
}

// Set cache value
func (im *InMemory) Set(key string, value []byte) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, found := im.hash[key]; found {
		return errAlreadyExists
	}

	im.hash[key] = value
	return nil
}

// Remove cache value
func (im *InMemory) Remove(key string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, found := im.hash[key]; found {
		delete(im.hash, key)
		return nil
	}

	return errNotFound
}
