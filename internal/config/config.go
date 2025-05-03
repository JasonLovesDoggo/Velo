package config

import "sync"

// temp in memory store

var store = struct {
	data map[string]string
	sync.RWMutex
}{data: make(map[string]string)}

func Set(key, value string) {
	store.Lock()
	defer store.Unlock()
	store.data[key] = value
}

func Get(key string) string {
	store.RLock()
	defer store.RUnlock()
	return store.data[key]
}
