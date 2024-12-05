package disk

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Bay-Shore-Systems-Inc/cache"
)

type (
	// Store implements cache.Store
	Store struct {
		storeType string
		mtx       sync.RWMutex

		// RootDir defines the root file path of the cache
		RootDir string

		// MaxAge is the implementation of cache.MaxAge for use during trimming old files
		MaxAge cache.MaxAge
	}

	// writer is used to implement the store for read, write, and remove
	writer struct {
		Store *Store
	}
)

// New initializes a new instance of the disk store to be added to the current cache.
// If MaxAge or RootDir are not provided the will be set to their default values.
func New(s *Store) *Store {
	// define the store type
	s.storeType = "disk"

	// Check if RootDir and MaxAge are set.
	// If not set to the default value.
	if s.RootDir == "" {
		fmt.Println("cache: a root directory for the cache has not been provided. It will now be set to cache/ in the applications root directory.")
		s.RootDir = filepath.Clean("./cache")
	} else {
		s.RootDir = filepath.Clean(s.RootDir)
	}

	if s.MaxAge == 0 {
		s.MaxAge = cache.DefaultMaxAge
	}

	// create the cache root directory
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, err := os.Stat(s.RootDir)
	if err != nil {
		err = os.MkdirAll(s.RootDir, 0o750)
		if err != nil {
			fmt.Printf("cannot make root directory: %v", err)
			return nil
		}
	}
	return s
}

// Type returns the stores type as a string
func (s *Store) Type() string {
	return s.storeType
}

// Get returns a writer for the current disk store
func Get(s *Store) *writer {
	var w writer
	w.Store = s
	return &w
}

// Write saves the data passed with the fileName give to the given directory.
// The given directory is joined with the RootDir path set when the store was created.
// If overwrite = true the file will be overwriten if it already exists
func (w *writer) Write(path string, fileName string, data []byte, overwrite bool) error {
	w.Store.mtx.Lock()
	defer w.Store.mtx.Unlock()

	// Check if the save directory exists.
	// If not create it.
	saveDir := w.Store.buildPath(path)
	_, err := os.Stat(saveDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(saveDir, 0o750)
		if err != nil {
			return err
		}
	}

	fullPath := w.Store.buildPath(saveDir, fileName)

	// Check if ok to overwrite an already existing file.
	if !overwrite {
		_, err = os.Stat(fullPath)
		if err == nil {
			return fmt.Errorf("file already exists in store: %v", err)
		}
	}

	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0o600) //#nosec G304
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// Remove deletes the file passed in at the given path from the store.
func (w *writer) Remove(path string, fileName string) error {
	w.Store.mtx.Lock()
	defer w.Store.mtx.Unlock()

	fullPath := w.Store.buildPath(path, fileName)
	file, err := os.Open(fullPath) //#nosec G304
	if err != nil {
		return err
	}
	defer file.Close()

	err = os.Remove(fullPath)
	if err != nil {
		return err
	}
	return nil
}

// Read reads the file passed in from the store in the given path,
// and return it as a byte slice.
func (w *writer) Read(path string, fileName string) ([]byte, error) {
	w.Store.mtx.Lock()
	defer w.Store.mtx.Unlock()

	fullPath := w.Store.buildPath(path, fileName)
	file, err := os.Open(fullPath) //#nosec G304
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return []byte{}, err
	}

	b := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(b)
	if err != nil && err != io.EOF {
		return []byte{}, nil
	}

	return b, nil
}

// buildPath is an internal method used for building a complete cleaned file path
// It will join the RootDir if it is not already present.
func (s *Store) buildPath(elem ...string) string {
	// Return empty string if there a no elements
	if len(elem) == 0 {
		return ""
	}

	// Join all elements then join with RootDir
	path := filepath.Join(elem...)

	// If the path doesn't include the RootDir add it
	if !strings.HasPrefix(path, s.RootDir) {
		fullPath := filepath.Join(s.RootDir, path)
		return fullPath
	}

	return path
}

// Purge will clear the entire cache and remove the RootDir.
// This function should only be used when stopping the service.
// If you need to flush the store without stopping it you can
// call this method directly.
func (s *Store) Purge() error {
	log.Println("File store is being purged...")
	s.mtx.Lock()
	defer s.mtx.Unlock()

	err := os.RemoveAll(s.RootDir)
	if err != nil {
		return err
	}

	log.Println("File store purge complete")
	return nil
}

// Trim is used for trimming files older then the MaxAge.
// It is called by the caches trim worker.
// This can be called directly if needed.
func (s *Store) Trim() {
	log.Println("Starting file store trimming...")
	s.mtx.Lock()
	defer s.mtx.Unlock()

	err := filepath.Walk(s.RootDir, s.walk)
	if err != nil {
		log.Printf("unable to read path: %v", err)
	}

	log.Println("File store trimming complete")
}

// walk is an internal method used for filepath.WalkFunc
// to check a files MaxAge and remove it if to old.
// It will also check for empty directories and remove them.
func (s *Store) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	cleanPath := filepath.Clean(path)
	f, err := os.Open(cleanPath)
	if err != nil {
		return err
	}
	defer f.Close()

	switch info.IsDir() {
	// If the path is not a directory check if it has reached the MaxAge.
	// If so delete the file.
	case false:
		stat, _ := os.Stat(cleanPath)
		age := stat.ModTime().Add(time.Second * time.Duration(s.MaxAge))
		if time.Now().Local().After(age) {
			err := os.RemoveAll(cleanPath)
			if err != nil {
				return err
			}
		}
	// If the path is a directory check if it is empty.
	// If so remove the empty directory.
	case true:
		if path != s.RootDir {
			_, err := f.Readdirnames(1)
			if err == io.EOF {
				err = os.Remove(cleanPath)
				if err != nil {
					return err
				}
			}
		}
	default:
		return errors.New("unknown file type")
	}

	return nil
}
