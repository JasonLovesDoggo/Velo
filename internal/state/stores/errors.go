package config

import "errors"

// ErrNotFound is returned when a key is not found in the store.
var ErrNotFound = errors.New("config: key not found")
