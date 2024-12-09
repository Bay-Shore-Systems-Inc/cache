![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/bay-shore-systems-inc/cache)
[![Go Reference](https://pkg.go.dev/badge/github.com/Bay-Shore-Systems-Inc/cache.svg)](https://pkg.go.dev/github.com/Bay-Shore-Systems-Inc/cache)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bay-Shore-Systems-Inc/cache)](https://goreportcard.com/report/github.com/Bay-Shore-Systems-Inc/cache)
[![Go](https://github.com/Bay-Shore-Systems-Inc/cache/actions/workflows/go.yml/badge.svg)](https://github.com/Bay-Shore-Systems-Inc/cache/actions/workflows/go.yml)
# Cache
Cache is a simple dependency free cache package with the ability to add new storage methods by implementing a new store.

This cache will run as many stores as needed and can run more than one instance at a time. You can access and maintain the stores individually per cache instance. 
Each cache will run a trimming function which calls the trim method implemented by each store to remove old items. The cache will also shut down gracefully. You can shut down all cache instances or a single instance.

This package supports the last 2 stable go versions.

- [Install](#install)
- [Stores](#stores)
- [Implementing](#implementing)
  - [Cache Options](#cache-options)
  - [Accessing Stores](#accessing-stores)
  - [Direct Store Maintenance](#direct-store-maintenance)
  - [Stopping](#stopping)

## Install
```
$ go get github.com/bay-shore-systems-inc/cache
```
## Stores
* Mem  (In-memory)
* disk (On disk file store)

## Implementing
Select and initiate each store you want to run in the cache.
> [!NOTE]
> The Options for each store will vary depending on the store.
```go
import "github.com/bay-shore-systems-inc/stores/mem"

func main() {
  ...
  store := mem.New(&mem.Store{
    // Every store should implement a MaxAge to be used for trimming old cache items.
    MaxAge: 1800,
  })
  ...
}
```
> [!NOTE]
> MaxAge is measured in seconds so a value of 1800 would be 30 minutes

Initialize a new cache.
```go
import (
  "github.com/bay-shore-systems-inc/mem"
  "github.com/bay-shore-systems-inc/stores/mem"
)
func main() {
  ...
  c := cache.New(&cache.Options{
    TrimTime: 900,
    Stores: cache.MakeStores(store),
  })
  ...
}
```

Start the cache.
```go
import (
  "log"
  "github.com/bay-shore-systems-inc/mem"
  "github.com/bay-shore-systems-inc/stores/mem"
)
func main() {
  ...
  err := c.Start()
  if err != nil {
    log.Panicf("Cache instance could not be started: %v", err)
  }
  ...
}
```
## Cache Options
| Option | Default | Description |
|---     |---      | --- |
| TrimTime | 1800s (30m) | Used by the internal trim method to decide when the trim method for each store should be ran.|
| Stores | nil | Used be the current cache to access the store or stores used to save items |

## Accessing Stores
Once you have the cache started you can access the store. You will need to call the store to get the store's writer.
```go
func main() {
  ...
  m := mem.Get(store)
  ...
}
```
Now you can access the store's writer methods.
```go
func main() {
  ...
  key := "foo"
  value := "bar"
  overwrite := false

  // Get in-memory store
  m := mem.Get(store)

  // Write key-value pair
  err := m.Write(key, []byte(value), overwrite)
  if err != nil {
    fmt.Println(err)
  }

  // Read value with given key
  v, err := m.Read(key)
  if err != nil {
    fmt.Println(err)
  }

  // Delete key-value pair
  m.Remove(key)
  ...
}
```
## Direct Store Maintenance
If you would like to trim or purge a store directly you can do so by calling their methods directly.
```go
func main() {
  ...
  // Runs the stores trim method for removing old items
  m.Trim()

  // Runs the stores purge method to remove all items
  m.Purge()
  ...
}
```

## Stopping
When you want to stop the cache there are two ways to do this. You can either stop all instances or a single cache instance.

To stop a single cache instance.
```go
func main() {
  ...
  err := cache.StopCacheInstance(c.CacheNum)
  if err != nil {
    log.Println(err)
  }
  ...
}
```
To stop all running cache instances.
```go
func main() {
  ...
  err := cache.StopAll()
  if err != nil {
    log.Println(err)
  }
  ...
}
```
