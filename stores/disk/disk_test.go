package disk_test

import (
	"encoding/json"
	"testing"

	"github.com/Bay-Shore-Systems-Inc/cache"
	"github.com/Bay-Shore-Systems-Inc/cache/stores/disk"
	"github.com/stretchr/testify/assert"
)

type User struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

var (
	u = User{
		FirstName: "John",
		LastName:  "Doe",
	}
	path = "json"
	file = "somejson.json"
	root = "./testcache"
)

// Test writing, reading, and removing in the on disk store
func TestDiskStore(t *testing.T) {
	a := assert.New(t)

	// Create a disk store
	diskStore := disk.New(&disk.Store{
		RootDir: root,
		MaxAge:  20,
	})
	a.NotNil(diskStore)

	// Create a new cache
	c := cache.New(&cache.Options{
		TrimTime: 10,
		Stores:   cache.MakeStores(diskStore),
	})
	a.NoError(c.Start())

	d := disk.Get(diskStore)
	a.NotNil(d)

	data, _ := json.MarshalIndent(u, " ", " ")

	// Save to file
	err := d.Write(path, file, data, false)
	a.NoError(err)

	// Attempt to overwrite file with overwrite = false
	err = d.Write(path, file, data, false)
	a.Error(err)

	// Attempt to overwrite the same file at the same time with overwrite = true
	for i := 0; i < 2000; i++ {
		go a.NoError(d.Write(path, file, data, true))
	}

	var readUser User
	b, _ := d.Read(path, file)

	err = json.Unmarshal(b, &readUser)
	assert.NoError(t, err)

	if assert.NotNil(t, readUser) {
		a.Equal(readUser.FirstName, u.FirstName)
		a.Equal(readUser.LastName, u.LastName)
	}

	// Remove the file
	a.NoError(d.Remove(path, file))

	// Check if the file have been removed
	_, err = d.Read(path, file)
	a.Error(err)

	cache.StopCacheInstance(c.CacheNum)
}
