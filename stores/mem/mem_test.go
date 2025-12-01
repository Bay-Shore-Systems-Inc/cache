package mem_test

import (
	"crypto/rand"
	"encoding/base32"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tmstorm/cache"
	"github.com/tmstorm/cache/stores/mem"
)

// TestMemStore tests the in-memory storage
func TestMemStore(t *testing.T) {
	a := assert.New(t)

	memStore := mem.New(&mem.Store{
		MaxAge: 4,
	})
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
	s := 2000
	key := make([]string, s)
	value := make([]string, s)
	for i := 0; i < len(key); i++ {
		key[i], _ = randString(16)
	}
	for i := 0; i < len(value); i++ {
		value[i], _ = randString(16)
	}

	// Test adding key-values and reading them from memory store
	//
	// WARNING: When testing race conditions in multiple go routines
	// you can't use the assert package or %v to format lines.
	// the assert package uses %v and will cause a data race.
	// See https://github.com/stretchr/testify/pull/1598
	// %v causes race conditions because it uses the GoStringer interface
	// which is not safe for concurrency
	var wg sync.WaitGroup
	for i := 0; i < len(key); i++ {
		wg.Add(1)
		go func() {
			err := m.Write(key[i], []byte(value[i]), false)
			if err != nil {
				t.Fail()
				t.Log(err)
			}

			v, err := m.Read(key[i])
			if err != nil {
				t.Fail()
				t.Log(err)
			} else {
				if value[i] != string(v) {
					t.Fail()
					t.Logf("%s doesn't match %s", value[i], string(v))
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Test deleting items from memory store
	for i := 0; i < len(key); i++ {
		wg.Add(1)
		go func() {
			oldVal := value[i]
			err := m.Remove(key[i])
			if err != nil {
				t.Fail()
				t.Log(err)

			}
			v, err := m.Read(key[i])
			if err != nil {
				if oldVal == string(v) {
					t.Fail()
					t.Logf("%s should not match %s", oldVal, string(v))
				}
			} else {
				t.Fail()
				t.Log("expected error but got nil")
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Test Trim method
	k, _ := randString(16)
	v, _ := randString(16)
	m.Write(k, []byte(v), true)
	time.Sleep(time.Second * 6)
	_, err := m.Read(k)
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
