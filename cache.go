package cache

import (
	"fmt"
	"log"
	"time"
)

type (
	// Options sets an individual caches options
	Options struct {
		// CacheNum is used to call which cache the user wants to use.
		// This is set when the cache instance is started and should be saved to access the cache later.
		CacheNum int

		// Stores is used to access the current caches active stores
		Stores Stores

		// TrimTime Defines the interval at which trim() calls the store specified Trim() method
		// to check the MaxAge of items the store.
		// TrimTime = 900 would set TrimTime to 15 minutes
		TrimTime time.Duration

		// w adds the Workers struct to the current cache instance
		w workers
	}

	// Workers are used to control the trim() worker in the current cache
	workers struct {
		// stop is an internal channel used to close the trim() worker when stopping the cache.
		stop chan bool

		// done is used to signal the StopAll() or StopCacheInstance() methods that the
		// trim() worker has closed and it is safe to purge the cache
		done chan bool
	}
)

var (
	// defaultOpts sets the default options of the file cache.
	defaultOpts = Options{
		TrimTime: 900,
		w: workers{
			stop: make(chan bool),
			done: make(chan bool),
		},
	}

	// caches is an internal map created to access an individual cache
	// When a new cache is created its options and cache number are added to this map.
	caches = make(map[int]*Options)
)

// New returns a new cache instance for use using the options struct.
// If no options are passed in or any are omitted the defaults will be applied.
func New(o *Options) *Options {
	// Check if options are set.
	// If not set to the default value.
	if o.TrimTime == 0 {
		o.TrimTime = defaultOpts.TrimTime
	}

	// Set the cache number
	cacheNum := len(caches)
	if cacheNum == 0 {
		cacheNum = 1
	} else {
		cacheNum++
	}

	// Create and initiates new cache options
	newCache := &Options{
		CacheNum: cacheNum,
		Stores:   o.Stores,
		TrimTime: o.TrimTime,
		w:        defaultOpts.w,
	}

	caches[newCache.CacheNum] = newCache

	return caches[newCache.CacheNum]
}

// Start initiates a new cache instance.
// It needs to be initialized by New()
// It should start trim() as a goroutine for each store to maintain the cache size.
func (o *Options) Start() error {
	// check if a stores have been set
	if len(o.Stores) == 0 {
		return fmt.Errorf("no store has been provided")
	}

	log.Println("Starting local cache...")

	// Start a gorouting for trimming the cache
	for store := range o.Stores {
		s, err := getStore(store, o.Stores)
		if err != nil {
			return err
		}
		go o.trim(s)
	}

	return nil
}

// StopCacheInstance gracefully closes the specified cache instances.
// It should stop the trim() worker.
// Purge the cache by calling purge().
// Remove the cache from the active Caches map.
func StopCacheInstance(cn int) error {
	// Stop the trim worker
	caches[cn].w.stop <- true

	// Wait for the trim worker to stop
	<-caches[cn].w.done

	// After the trim worker has closed purge the cache
	for store := range caches[cn].Stores {
		s, err := getStore(store, caches[cn].Stores)
		if err != nil {
			return err
		}

		err = caches[cn].purge(s)
		if err != nil {
			return err
		}
	}

	// Remove the cache from the map and return
	delete(caches, cn)

	return nil
}

// StopAll gracefully closes the all cache instances.
// It should stop the trim() worker.
// Purge the cache and remove the root directory.
// Reset caches to an empty map.
func StopAll() error {
	for index, cache := range caches {
		// Stop the trim worker
		cache.w.stop <- true

		// Wait for the trim worker to stop
		<-cache.w.done

		// After the trim worker has closed purge the cache
		for store := range cache.Stores {
			s, err := getStore(store, cache.Stores)
			if err != nil {
				return err
			}

			err = cache.purge(s)
			if err != nil {
				return err
			}
		}

		// Remove the cache from the map and return
		delete(caches, index)
	}

	return nil
}

// trim calls the Trim method for each store instance in a cache.
// It uses the Workers channel to check is a store is currently
// trimming which is started based on the TrimTime provided when the cache
// is initialized.
func (o *Options) trim(store Store) {
	for {
		select {
		case <-o.w.stop:
			log.Println("Trim worker closing")
			o.w.done <- true
		case <-time.After(time.Second * o.TrimTime):
			store.Trim()
		}
	}
}

// purge calls the Purge method for each store in the cache when the
// cache is signaled to stop.
func (o *Options) purge(store Store) error {
	err := store.Purge()
	if err != nil {
		return err
	}
	return nil
}
