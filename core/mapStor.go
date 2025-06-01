package mapstore

import (
	"fmt"
	"sync"
)

// MapStore is an interface that abstracts a simple key-value store.
// It defines methods to load, store, and delete values by key.
type MapStore interface {
	Load(key string) (interface{}, bool) // Retrieves the value for the given key, if present
	Store(key string, value interface{}) // Stores a value under the given key
	Delete(key string)                   // Removes the key-value pair from the store
}

// muMapStore is a concrete implementation of MapStore using a standard Go map.
type muMapStore struct {
	mapStore map[string]interface{}
}

// muInstance holds the singleton instance of muMapStore.
// This prevents multiple instantiations of the map store.
var muInstance *muMapStore

// NewMuMapStore initializes and returns a singleton instance of muMapStore.
// It ensures that the map store is created only once, using the provided mutex for thread safety during instantiation.
func NewMuMapStore(lock *sync.Mutex) MapStore {
	if muInstance == nil {
		fmt.Println("Application is using Balanced RW Map Mechanism")
		lock.Lock()
		defer lock.Unlock()

		// Double-check inside lock to avoid race in concurrent calls
		if muInstance == nil {
			muInstance = &muMapStore{
				mapStore: make(map[string]interface{}),
			}
		}
	}
	return muInstance
}

// Store saves the given value associated with the specified key in the map.
func (r *muMapStore) Store(key string, value interface{}) {
	r.mapStore[key] = value
}

// Load retrieves the value for a given key. Returns the value and a boolean indicating if the key exists.
func (r *muMapStore) Load(key string) (interface{}, bool) {
	val, ok := r.mapStore[key]
	return val, ok
}

// Delete removes the key and its associated value from the map.
func (r *muMapStore) Delete(key string) {
	delete(r.mapStore, key)
}
