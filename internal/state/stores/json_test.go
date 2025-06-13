package stores

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// --- JSON Specific Load Tests

func TestJSONStore_LoadSpecifics(t *testing.T) {
	t.Run("LoadEmptyFile", func(t *testing.T) {
		tempDir := t.TempDir()
		jsonPath := filepath.Join(tempDir, "empty.json")

		f, err := os.Create(jsonPath)
		if err != nil {
			t.Fatalf("Failed to create empty file %q: %v", jsonPath, err)
		}
		f.Close()

		store, err := NewJSONStore(jsonPath)
		if err != nil {
			t.Fatalf("Failed to create store from empty file %q: %v", jsonPath, err)
		}
		defer store.Close()

		_, err = store.Get("anyKey")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Expected ErrNotFound when getting from store loaded from empty file, got %v", err)
		}

		err = store.Set("newKey", "newValue")
		if err != nil {
			t.Errorf("Failed to set key on store loaded from empty file: %v", err)
		}
	})

	t.Run("LoadMalformedFile", func(t *testing.T) {
		tempDir := t.TempDir()
		jsonPath := filepath.Join(tempDir, "malformed.json")

		malformedData := []byte(`{"key": "value",`)
		err := os.WriteFile(jsonPath, malformedData, 0644)
		if err != nil {
			t.Fatalf("Failed to write malformed file %q: %v", jsonPath, err)
		}

		_, err = NewJSONStore(jsonPath)
		if err == nil {
			t.Fatalf("Expected error when creating store from malformed file %q, got nil", jsonPath)
		}
		t.Logf("Got expected error for malformed JSON: %v", err)
	})
}
