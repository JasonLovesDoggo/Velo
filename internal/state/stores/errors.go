package stores

import "errors"

// ErrNotFound is returned when a key is not found in the store.
var ErrNotFound = errors.New("stores: key not found")
