package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type JsonStore struct {
	filePath string
	mu       sync.RWMutex
	data     map[string]string
}

// NewJSONStore creates and initializes a new JSON file-backed store.
func NewJSONStore(filePath string) (*JsonStore, error) {
	store := &JsonStore{
		filePath: filePath,
		data:     make(map[string]string),
	}
	if err := store.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		// Ignore "file not found" on initial load, but report other errors.
		return nil, fmt.Errorf("failed to load initial json config from %q: %w", filePath, err)
	}
	return store, nil
}

// load reads the JSON file into the in-memory map. MUST be called with write lock held or during init.
func (s *JsonStore) load() error {
	content, err := os.ReadFile(s.filePath)
	if err != nil {
		// Let caller handle os.ErrNotExist if necessary
		return err
	}
	// If file is empty, just return (treat as empty map)
	if len(content) == 0 {
		s.data = make(map[string]string) // Ensure it's empty map, not nil
		return nil
	}

	newData := make(map[string]string)
	if err := json.Unmarshal(content, &newData); err != nil {
		return fmt.Errorf("failed to unmarshal json from %q: %w", s.filePath, err)
	}
	s.data = newData
	return nil
}

// save writes the in-memory map back to the JSON file. MUST be called with write lock held.
func (s *JsonStore) save() error {
	jsonData, err := json.MarshalIndent(s.data, "", "  ") // Pretty print JSON
	if err != nil {
		return fmt.Errorf("failed to marshal json data: %w", err)
	}

	// Write atomically if possible (rename over existing)
	tempFile := s.filePath + ".tmp"
	err = os.WriteFile(tempFile, jsonData, 0640) // Sensible permissions
	if err != nil {
		return fmt.Errorf("failed to write temp json file %q: %w", tempFile, err)
	}

	// Attempt atomic rename
	if err = os.Rename(tempFile, s.filePath); err != nil {
		// Cleanup temp file on rename failure
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temp json file to %q: %w", s.filePath, err)
	}

	return nil
}

func (s *JsonStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return value, nil
}

func (s *JsonStore) Set(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure data map is initialized if load failed initially with ErrNotExist
	if s.data == nil {
		s.data = make(map[string]string)
	}

	s.data[key] = value
	if err := s.save(); err != nil {
		// Note: In-memory data is updated even if save fails.
		// Could potentially revert, but keeping it simple for now.
		return fmt.Errorf("json Set succeeded in memory but failed to save: %w", err)
	}
	return nil
}

func (s *JsonStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure data map is initialized
	if s.data == nil {
		return nil // Key definitely doesn't exist if map is nil
	}

	_, exists := s.data[key]
	if !exists {
		return nil // Deleting non-existent key is not an error
	}

	delete(s.data, key)
	if err := s.save(); err != nil {
		// Note: In-memory data is updated even if save fails.
		return fmt.Errorf("json Delete succeeded in memory but failed to save: %w", err)
	}
	return nil
}

func (s *JsonStore) Close() error {
	// No explicit resources to close for the JSON store (file handles are managed per operation).
	// Could potentially force a save here if needed, but current design saves on modify.
	s.mu.Lock() // Ensure no operations are ongoing if we decide to add cleanup later
	defer s.mu.Unlock()
	s.data = nil // Clear map to release memory if desired
	return nil
}
