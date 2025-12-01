package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tmstorm/cache"
	"github.com/tmstorm/cache/stores/sham"
)

func TestMakeStores(t *testing.T) {
	a := assert.New(t)
	s := sham.New(&sham.Store{})
	stores := cache.MakeStores(s)

	c := cache.New(&cache.Options{
		Stores: stores,
	})
	a.NotNil(c.Stores)

	c = cache.New(&cache.Options{})
	a.Nil(c.Stores)
}
