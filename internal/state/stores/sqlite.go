package stores

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // Driver import
)

const (
	createTableSQL = `
	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	// Index for potential performance improvement on Get
	createKeyIndexSQL = `CREATE INDEX IF NOT EXISTS idx_config_key ON config(key);`

	getSQL    = `SELECT value FROM config WHERE key = ?;`
	setSQL    = `INSERT INTO config (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP;`
	deleteSQL = `DELETE FROM config WHERE key = ?;`
)

type SqliteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates and initializes a new SQLite-backed store.
func NewSQLiteStore(dbPath string) (*SqliteStore, error) {
	// Add necessary parameters for better performance and WAL mode (good for Litestream if we use that down the line)
	// _busy_timeout helps with concurrent writes.
	// _journal_mode=WAL allows concurrent reads and writes.
	// _synchronous=NORMAL is a good balance for WAL mode.
	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL", dbPath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database %q: %w", dbPath, err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // Important for SQLite to avoid concurrency issues at the connection level
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour) // Optional: recycle connections periodically

	// Ping to ensure connection is valid
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close() // Close the potentially invalid DB handle
		return nil, fmt.Errorf("failed to ping sqlite database %q: %w", dbPath, err)
	}

	// Ensure table and index exist
	if _, err := db.ExecContext(ctx, createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create config table: %w", err)
	}
	if _, err := db.ExecContext(ctx, createKeyIndexSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create config key index: %w", err)
	}

	return &SqliteStore{db: db}, nil
}

func (s *SqliteStore) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var value string
	err := s.db.QueryRowContext(ctx, getSQL, key).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("sqlite Get failed for key %q: %w", key, err)
	}
	return value, nil
}

func (s *SqliteStore) Set(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, setSQL, key, value)
	if err != nil {
		return fmt.Errorf("sqlite Set failed for key %q: %w", key, err)
	}
	return nil
}

func (s *SqliteStore) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, deleteSQL, key)
	if err != nil {
		return fmt.Errorf("sqlite Delete failed for key %q: %w", key, err)
	}
	// We don't check RowsAffected, as deleting a non-existent key is not an error per interface.
	return nil
}

func (s *SqliteStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
