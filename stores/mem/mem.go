package mem

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Bay-Shore-Systems-Inc/cache"
)

type (
	// Store implements cache.Store
	Store struct {
		storeType string
		mtx       sync.Mutex
		data      memStoreMap

		// MaxAge is the implementation of cache.MaxAge for use during trimming old key-value pairs
		MaxAge cache.MaxAge
	}

	// writer is used to read, write, and remove key-value pairs
	writer struct {
		Store *Store
	}

	// value is used to save in a key value store.
	// Using the value passed and setting time saved.
	valueStore struct {
		value     []byte
		timeStamp time.Time
	}

	// memStoreMap type is used for the in-memory store data
	memStoreMap map[string]valueStore
)

// New will initialize a new in-memory store
// This is not persistent through reboots and should be used for short lived tasks.
func New(s *Store) (*Store, error) {
	// Set store type and make data map for writing
	s.storeType = "mem"
	s.data = make(memStoreMap)

	// Check if MaxAge is set.
	// If not set to the default value.
	if s.MaxAge == 0 {
		s.MaxAge = cache.DefaultMaxAge
	}
	return s, nil
}

// Type is the implementation from cache to get the store type
func (s *Store) Type() string {
	return s.storeType
}

// Mem is used to retrieve the in-memory store created at startup
func Get(s *Store) *writer {
	var w writer
	w.Store = s
	return &w
}

// Write adds a new key-value pair in memory
// If overwrite = true data will be overwriten if it alreay exists
func (w *writer) Write(key string, value []byte, overwrite bool) error {
	w.Store.mtx.Lock()
	defer w.Store.mtx.Unlock()
	if _, ok := w.Store.data[key]; ok {
		err := errors.New("key already exists in memory store")
		return err
	}
	save := valueStore{
		value:     value,
		timeStamp: time.Now(),
	}
	w.Store.data[key] = save
	return nil
}

// Read gets key-value pair from the in memory store and return is as a byte slice
func (w *writer) Read(key string) ([]byte, error) {
	w.Store.mtx.Lock()
	defer w.Store.mtx.Unlock()
	if _, ok := w.Store.data[key]; !ok {
		err := errors.New("key not found in memory store")
		return []byte{}, err
	}
	value := w.Store.data[key]
	return value.value, nil
}

// Remove deletes a key-value pair from in memory store
// Map delete() is no-op if map is nil or there is no matching key
func (w *writer) Remove(key string) error {
	w.Store.mtx.Lock()
	defer w.Store.mtx.Unlock()
	delete(w.Store.data, key)
	return nil
}

// Trim is used for trimming keys older then the MaxAge.
// It is called by the caches trim worker.
// This can be called directly if needed.
func (s *Store) Trim() {
	log.Println("Starting file store trimming...")
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for key, stored := range s.data {
		age := stored.timeStamp.Add(time.Second * time.Duration(s.MaxAge))
		if time.Now().Local().After(age) {
			delete(s.data, key)
		}
	}

	log.Println("File store trimming complete")
}

// Purge clears the entire in-memory store.
// This function should only be used when stopping the service.
// If you need to flush the store without stopping it you can
// call this method directly.
// This will never return an error.
func (s *Store) Purge() error {
	log.Println("In-memory store is being purged...")
	s.mtx.Lock()
	defer s.mtx.Unlock()

	newMap := make(memStoreMap)
	s.data = newMap

	log.Println("In-memory store purge complete")
	return nil
}
