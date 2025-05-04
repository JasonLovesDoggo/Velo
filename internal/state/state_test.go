package state

import (
	"errors"
	"fmt"
	stores "github.com/jasonlovesdoggo/velo/internal/state/stores"
	_ "g
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runStoreImplementationTests(t *testing.T, store Store) {
	t.Helper()
	key1, val1 := "testKey1", "testValue1"
	key2, val2 := "testKey2", "testValue2"
	updatedVal1 := "testValue1-updated"

	t.Run("InitialGetNotFound", func(t *testing.T) {
		_, err := store.Get(key1)
		if !errors.Is(err, stores.ErrNotFound) {
			t.Errorf("Expected ErrNotFound for initial Get(%q), got %v", key1, err)
		}
	})

	t.Run("SetAndGetFirstKey", func(t *testing.T) {
		err := store.Set(key1, val1)
		if err != nil {
			t.Fatalf("Set(%q, %q) failed: %v", key1, val1, err)
		}
		retVal, err := store.Get(key1)
		if err != nil {
			t.Fatalf("Get(%q) after Set failed: %v", key1, err)
		}
		if retVal != val1 {
			t.Errorf("Get(%q) returned %q, want %q", key1, retVal, val1)
		}
	})

	t.Run("SetAndGetSecondKey", func(t *testing.T) {
		err := store.Set(key2, val2)
		if err != nil {
			t.Fatalf("Set(%q, %q) failed: %v", key2, val2, err)
		}
		retVal, err := store.Get(key2)
		if err != nil {
			t.Fatalf("Get(%q) after Set failed: %v", key2, err)
		}
		if retVal != val2 {
			t.Errorf("Get(%q) returned %q, want %q", key2, retVal, val2)
		}
	})

	t.Run("UpdateAndGetFirstKey", func(t *testing.T) {
		// Ensure key1 exists from previous subtest (tests run sequentially within t.Run)
		_, err := store.Get(key1)
		if err != nil {
			t.Skipf("Skipping update test as key %q doesn't exist", key1)
		}

		err = store.Set(key1, updatedVal1)
		if err != nil {
			t.Fatalf("Update Set(%q, %q) failed: %v", key1, updatedVal1, err)
		}
		retVal, err := store.Get(key1)
		if err != nil {
			t.Fatalf("Get(%q) after update failed: %v", key1, err)
		}
		if retVal != updatedVal1 {
			t.Errorf("Get(%q) after update returned %q, want %q", key1, retVal, updatedVal1)
		}
	})

	t.Run("DeleteAndVerifyFirstKey", func(t *testing.T) {
		err := store.Delete(key1)
		if err != nil {
			t.Fatalf("Delete(%q) failed: %v", key1, err)
		}
		_, err = store.Get(key1)
		if !errors.Is(err, stores.ErrNotFound) {
			t.Errorf("Expected ErrNotFound after Delete(%q), got %v", key1, err)
		}
	})

	t.Run("VerifySecondKeyAfterDelete", func(t *testing.T) {
		// Ensure key2 exists from previous subtest
		retVal, err := store.Get(key2)
		if err != nil {
			t.Fatalf("Get(%q) after deleting another key failed: %v", key2, err)
		}
		if retVal != val2 {
			t.Errorf("Get(%q) after deleting another key returned %q, want %q", key2, retVal, val2)
		}
	})

	t.Run("DeleteNonExistentKey", func(t *testing.T) {
		err := store.Delete("non_existent_key_for_delete")
		if err != nil {
			t.Errorf("Delete() of non-existent key failed: %v", err)
		}
	})
}

func runStorePersistenceTest(t *testing.T, cfg Config) {
	t.Helper()
	key, val := fmt.Sprintf("persistKey_%s", cfg.Type), fmt.Sprintf("persistValue_%s", cfg.Type)

	store1, err := NewStore(cfg)
	if err != nil {
		t.Fatalf("Phase 1: Failed to create store: %v", err)
	}
	err = store1.Set(key, val)
	if err != nil {
		store1.Close() // Attempt close even on error
		t.Fatalf("Phase 1: Failed to set key %q: %v", key, err)
	}
	err = store1.Close()
	if err != nil {
		// Log warning for JSON as close is less critical, fail for SQLite
		if cfg.Type == SQLiteBackend {
			t.Fatalf("Phase 1: Failed to close store: %v", err)
		} else {
			t.Logf("Phase 1: Warning - closing store failed: %v", err)
		}
	}

	store2, err := NewStore(cfg)
	if err != nil {
		t.Fatalf("Phase 2: Failed to reopen store: %v", err)
	}
	defer store2.Close()

	retVal, err := store2.Get(key)
	if err != nil {
		t.Fatalf("Phase 2: Failed to get key %q after reopen: %v", key, err)
	}
	if retVal != val {
		t.Errorf("Phase 2: Got value %q after reopen, want %q", retVal, val)
	}

	err = store2.Delete(key)
	if err != nil {
		t.Logf("Phase 3: Warning - failed to delete persisted key %q: %v", key, err)
	}
}

func TestStoreImplementations(t *testing.T) {
	testCases := []struct {
		name        string
		backendType BackendType
		setupPath   func(t *testing.T) string
	}{
		{
			name:        "SQLite",
			backendType: SQLiteBackend,
			setupPath: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "test_matrix.db")
			},
		},
		{
			name:        "JSON",
			backendType: JSONBackend,
			setupPath: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "test_matrix.json")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path := tc.setupPath(t)
			cfg := Config{
				Type: tc.backendType,
				Path: path,
			}

			t.Run("CoreInterface", func(t *testing.T) {
				store, err := NewStore(cfg)
				if err != nil {
					t.Fatalf("Failed to create store for backend %q: %v", tc.backendType, err)
				}
				defer store.Close()
				runStoreImplementationTests(t, store)
			})

			t.Run("Persistence", func(t *testing.T) {
				runStorePersistenceTest(t, cfg)
			})

			t.Run("CloseIdempotency", func(t *testing.T) {
				store, err := NewStore(cfg)
				if err != nil {
					t.Fatalf("Failed to create store for close test: %v", err)
				}
				err = store.Close()
				if err != nil {
					t.Errorf("First Close() failed: %v", err)
				}
				err = store.Close()
				if err != nil {
					t.Errorf("Second Close() failed: %v", err)
				}
			})
		})
	}
}

func TestNewStore_UnsupportedType(t *testing.T) {
	_, err := NewStore(Config{Type: "invalid-backend", Path: "dummy"})
	if err == nil {
		t.Error("Expected an error for unsupported backend type, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("Expected error message to contain 'unsupported', got: %v", err)
	}
}

func TestNewStore_BadPath(t *testing.T) {
	badPath := filepath.Join(string(os.PathSeparator), "dev", "null", "non_existent_subdir", "test.db")
	parentDir := filepath.Dir(badPath)

	_, errSqlite := NewStore(Config{Type: SQLiteBackend, Path: badPath})
	if errSqlite == nil {
		t.Errorf("Expected an error for SQLite with unwritable path %q, got nil", badPath)
		store, _ := NewStore(Config{Type: SQLiteBackend, Path: badPath})
		if store != nil {
			store.Close()
		}
		os.Remove(badPath)
		os.Remove(parentDir)
	} else {
		t.Logf("Got expected error for SQLite bad path: %v", errSqlite)
	}

	_, errJson := NewStore(Config{Type: JSONBackend, Path: badPath})
	if errJson == nil {
		t.Errorf("Expected an error for JSON with unwritable path %q, got nil", badPath)
		store, _ := NewStore(Config{Type: JSONBackend, Path: badPath})
		if store != nil {
			store.Close()
		}
		os.Remove(badPath)
		os.Remove(parentDir)
	} else {
		t.Logf("Got expected error for JSON bad path: %v", errJson)
	}
}
