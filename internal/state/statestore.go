package state

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONStateStore implements StateStore using a JSON-based Store backend
type JSONStateStore struct {
	store Store
}

// NewJSONStateStore creates a new StateStore backed by a JSON Store
func NewJSONStateStore(store Store) StateStore {
	return &JSONStateStore{
		store: store,
	}
}

// Get retrieves and unmarshals a value for the given key
func (s *JSONStateStore) Get(key string, value interface{}) error {
	data, err := s.store.Get(key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}

	return nil
}

// Set marshals and stores a value for the given key
func (s *JSONStateStore) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	return s.store.Set(key, string(data))
}

// Delete removes a key-value pair
func (s *JSONStateStore) Delete(key string) error {
	return s.store.Delete(key)
}

// List returns all keys with a given prefix
func (s *JSONStateStore) List(prefix string) ([]string, error) {
	// This is a simplified implementation - a real implementation would
	// need to be backed by a store that supports key listing
	// For now, return empty list
	return []string{}, nil
}

// Close releases any resources
func (s *JSONStateStore) Close() error {
	return s.store.Close()
}

// NewDefaultStateStore creates a StateStore with default configuration
func NewDefaultStateStore() (StateStore, error) {
	config := Config{
		Type: JSONBackend,
		Path: "/tmp/velo-state.json",
	}

	store, err := NewStore(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create backing store: %w", err)
	}

	return NewJSONStateStore(store), nil
}

// MemoryStateStore is an in-memory implementation for testing
type MemoryStateStore struct {
	data map[string]string
}

// NewMemoryStateStore creates a new in-memory StateStore
func NewMemoryStateStore() StateStore {
	return &MemoryStateStore{
		data: make(map[string]string),
	}
}

func (m *MemoryStateStore) Get(key string, value interface{}) error {
	data, exists := m.data[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}

	return json.Unmarshal([]byte(data), value)
}

func (m *MemoryStateStore) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	m.data[key] = string(data)
	return nil
}

func (m *MemoryStateStore) Delete(key string) error {
	delete(m.data, key)
	return nil
}

func (m *MemoryStateStore) List(prefix string) ([]string, error) {
	var keys []string
	for key := range m.data {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func (m *MemoryStateStore) Close() error {
	m.data = nil
	return nil
}
