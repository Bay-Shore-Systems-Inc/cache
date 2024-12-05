package cache_test

import (
	"testing"

	"github.com/Bay-Shore-Systems-Inc/cache"
	"github.com/Bay-Shore-Systems-Inc/cache/stores/sham"
	"github.com/stretchr/testify/assert"
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
