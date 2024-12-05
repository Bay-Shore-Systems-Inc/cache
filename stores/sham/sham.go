package sham

import (
	"fmt"
	"sync"
	"time"

	"github.com/Bay-Shore-Systems-Inc/cache"
)

// Test implementation of the store interface
type Store struct {
	storeType string
	mtx       sync.RWMutex
	MaxAge    cache.MaxAge
}

type writer struct {
	Store *Store
}

func New(s *Store) *Store {
	s.storeType = "testStore"
	return s
}

func Get(s *Store) *writer {
	var w writer
	w.Store = s
	return &w
}

func (s *Store) Type() string {
	return s.storeType
}

func (s *Store) Trim() {
	fmt.Println("running trim function")
	s.mtx.Lock()
	time.Sleep(time.Second * 4)
	defer s.mtx.Unlock()
}

func (s *Store) Purge() error {
	fmt.Println("running purge function")
	s.mtx.Lock()
	time.Sleep(time.Second * 4)
	defer s.mtx.Unlock()
	return nil
}
