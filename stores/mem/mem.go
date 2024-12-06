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
		data      *sync.Map

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
)

// New will initialize a new in-memory store
// This is not persistent through reboots and should be used for short lived tasks.
func New(s *Store) (*Store, error) {
	// Set store type and make data map for writing
	s.storeType = "mem"
	s.data = new(sync.Map)

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
	if !overwrite {
		_, ok := w.Store.data.LoadOrStore(key, &valueStore{
			value:     value,
			timeStamp: time.Now(),
		})
		if ok {
			err := errors.New("key already exists in memory store")
			return err
		}
	} else {
		w.Store.data.Store(key, &valueStore{
			value:     value,
			timeStamp: time.Now(),
		})
	}
	return nil
}

// Read gets key-value pair from the in memory store and return is as a byte slice
func (w *writer) Read(key string) ([]byte, error) {
	value, ok := w.Store.data.Load(key)
	if !ok {
		err := errors.New("key not found in memory store")
		return []byte{}, err
	}
	return value.(*valueStore).value, nil
}

// Remove deletes a key-value pair from in memory store
// Map delete() is no-op if map is nil or there is no matching key
func (w *writer) Remove(key string) error {
	w.Store.data.Delete(key)
	_, ok := w.Store.data.Load(key)
	if !ok {
		return errors.New("key was not removed")
	}
	return nil
}

// Trim is used for trimming keys older then the MaxAge.
// It is called by the caches trim worker.
// This can be called directly if needed.
func (s *Store) Trim() {
	log.Println("Starting file store trimming...")

	s.data.Range(func(key interface{}, stored interface{}) bool {
		age := stored.(*valueStore).timeStamp.Add(time.Second * time.Duration(s.MaxAge))
		if time.Now().Local().After(age) {
			s.data.Delete(key)
		}
		return true
	})

	log.Println("File store trimming complete")
}

// Purge clears the entire in-memory store.
// This function should only be used when stopping the service.
// If you need to flush the store without stopping it you can
// call this method directly.
// This will never return an error.
func (s *Store) Purge() error {
	log.Println("In-memory store is being purged...")

	s.data.Range(func(key interface{}, stored interface{}) bool {
		s.data.Delete(key)
		return true
	})

	log.Println("In-memory store purge complete")
	return nil
}
