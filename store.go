package cache

import (
	"fmt"
	"time"
)

type (
	// Store must be implemented for each store
	// This is used to integrate into the cache and provide access
	// to starting, stopping, and internal cache maintenance.
	// a writer will still need to be implemented for each cache type.
	Store interface {
		Type() string
		Trim()
		Purge() error
	}

	// Stores are used to access all the available stores
	Stores map[string]Store

	// MaxAge Defines the maximum allowed age of items stored in the cache measured in seconds.
	// MaxAge = 1800 would set the MaxAge for all files to 30 minutes.
	// This needs to be implemented for each store and used with the store specific Purge() method
	MaxAge time.Duration
)

// DefaultMaxAge is used to define the MaxAge if one is not set for an individual store.
// This needs to be implemented at the store level on initialization.
var DefaultMaxAge = MaxAge(1800)

// MakeStores makes the stores after they have been initialized and are ready to be added to the cache
func MakeStores(stores ...Store) Stores {
	s := make(Stores)
	for _, store := range stores {
		s[store.Type()] = store
	}
	return s
}

// getStore is an internal method that should return a specific store for use
func getStore(sType string, Stores Stores) (Store, error) {
	store := Stores[sType]
	if store == nil {
		return nil, fmt.Errorf("store of type %v not found", sType)
	}
	return store, nil
}
