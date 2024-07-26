package main

// A memory based storage engine
type MemoryStorageEngine struct {
	Store map[string]any
}

// Create a new MemoryStorageEngine
func NewMemoryStorageEngine() MemoryStorageEngine {
	return MemoryStorageEngine{
		Store: make(map[string]any),
	}
}

// Set a value in this storage engine
func (s MemoryStorageEngine) SetValue(key string, value any) {
	s.Store[key] = value
}

// Get a value from the storage engine
func (s MemoryStorageEngine) GetValue(key string) any {
	val, ok := s.Store[key]

	if ok {
		return val
	}

	return nil
}
