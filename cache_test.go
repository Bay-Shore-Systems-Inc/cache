package cache_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tmstorm/cache"
	"github.com/tmstorm/cache/stores/sham"
)

// Test trimming and purging
func TestCache(t *testing.T) {
	a := assert.New(t)
	// Create a test store
	store := sham.New(&sham.Store{
		MaxAge: 3,
	})

	// Attempt to make a cache with zero stores
	badCache := cache.New(&cache.Options{})
	a.Error(badCache.Start())

	// Create a New local cache
	c := cache.New(&cache.Options{
		TrimTime: 10,
		Stores:   cache.MakeStores(store),
	})
	a.NoError(c.Start())

	s := sham.Get(store)
	a.NotNil(s)

	// Create a second cache
	cSecond := cache.New(&cache.Options{
		TrimTime: 2,
		Stores:   cache.MakeStores(store),
	})
	a.NoError(cSecond.Start())

	// Stop individual cache
	a.NoError(cache.StopCacheInstance(c.CacheNum))

	// Test if the second cache has its own identifier
	a.NotEqual(c.CacheNum, cSecond.CacheNum)

	// Check give trim time to run and stop all caches
	cSecond.Start()
	time.Sleep(time.Second * 4)
	a.NoError(cache.StopAll())
}
