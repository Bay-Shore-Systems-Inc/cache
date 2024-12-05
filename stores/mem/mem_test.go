package mem_test

import (
	"crypto/rand"
	"encoding/base32"
	"testing"
	"time"

	"github.com/Bay-Shore-Systems-Inc/cache"
	"github.com/Bay-Shore-Systems-Inc/cache/stores/mem"
	"github.com/stretchr/testify/assert"
)

// TestMemStore tests the in-memory storage
func TestMemStore(t *testing.T) {
	a := assert.New(t)

	memStore, err := mem.New(&mem.Store{
		MaxAge: 4,
	})
	a.NoError(err)
	a.NotNil(memStore)

	// Create a New local cache
	c := cache.New(&cache.Options{
		TrimTime: 5,
		Stores:   cache.MakeStores(memStore),
	})
	a.NotZero(len(c.Stores))

	m := mem.Get(memStore)
	a.NotNil(m)
	a.NoError(c.Start())

	// Make some random key-value strings to store
	s := 10
	key := make([]string, s)
	value := make([]string, s)
	for i := 0; i < len(key); i++ {
		key[i], _ = randString(16)
	}
	for i := 0; i < len(value); i++ {
		value[i], _ = randString(16)
	}

	// Test adding key-values and reading them from memory store
	for i := 0; i < len(key); i++ {
		err = m.Write(key[i], []byte(value[i]), false)
		a.NoError(err)

		v, err := m.Read(key[i])
		if a.NoError(err) {
			a.Equal(value[i], string(v))
		}
	}

	// Test deleting items from memory store
	for i := 0; i < len(key); i++ {
		oldVal := value[i]
		m.Remove(key[i])
		v, err := m.Read(key[i])
		if a.Error(err) {
			a.NotEqual(oldVal, string(v))
		}
	}

	// Test Trim method
	k, _ := randString(16)
	v, _ := randString(16)
	m.Write(k, []byte(v), true)
	time.Sleep(time.Second * 6)
	_, err = m.Read(k)
	a.Error(err)

	a.NoError(cache.StopCacheInstance(c.CacheNum))
}

// create random strings for testing
func randString(length int) (string, error) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(randBytes)[:length], nil
}
