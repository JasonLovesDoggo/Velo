package state

import (
	"fmt"
	stores "github.com/jasonlovesdoggo/velo/internal/state/stores"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

// Store defines the interface for key-value configuration storage.
type Store interface {
	// Get retrieves the value for a given key.
	// Returns ErrNotFound if the key does not exist.
	Get(key string) (string, error)

	// Set stores a key-value pair.
	// If the key already exists, its value is updated.
	Set(key string, value string) error

	// Delete removes a key-value pair.
	// It's not an error to delete a non-existent key.
	Delete(key string) error

	// Close releases any resources used by the store (like database connections).
	Close() error
}

// StateStore defines the interface for storing structured data with JSON serialization.
type StateStore interface {
	// Get retrieves and unmarshals a value for a given key into the provided interface.
	Get(key string, value interface{}) error

	// Set marshals and stores a value for a given key.
	Set(key string, value interface{}) error

	// Delete removes a key-value pair.
	Delete(key string) error

	// List returns all keys with a given prefix.
	List(prefix string) ([]string, error)

	// Close releases any resources used by the store.
	Close() error
}

// BackendType defines the type of storage backend.
type BackendType string

const (
	SQLiteBackend BackendType = "sqlite"
	JSONBackend   BackendType = "json"
)

// Config holds configuration for creating a new store.
type Config struct {
	// Type specifies the backend to use (e.g., SQLiteBackend, JSONBackend).
	Type BackendType
	// Path specifies the file path for the storage (db file for SQLite, json file for JSON).
	Path string
}

// NewStore creates a new configuration store based on the provided config.
func NewStore(cfg Config) (Store, error) {
	// Ensure the directory for the path exists.
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create directory %q: %w", dir, err)
	}

	switch cfg.Type {
	case SQLiteBackend:
		return stores.NewSQLiteStore(cfg.Path)
	case JSONBackend:
		return stores.NewJSONStore(cfg.Path)
	default:
		return nil, fmt.Errorf("unsupported config store type: %s", cfg.Type)
	}
}
