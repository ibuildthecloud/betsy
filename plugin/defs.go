package plugin

import "sync"

// Store manages the plugin inventory in memory and on-disk
type Store struct {
	sync.RWMutex
}

// NewStore creates a Store.
func NewStore(libRoot string) *Store {
	return &Store{}
}
