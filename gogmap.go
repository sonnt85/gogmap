package gogmap

import (
	"sync"
)

type GlobalMap[T any] struct {
	data map[string]T
	mu   sync.RWMutex
}

func NewGlobalMap[T any]() *GlobalMap[T] {
	return &GlobalMap[T]{
		data: make(map[string]T),
	}
}

// NewGlobalMapWithCapacity creates a GlobalMap pre-allocated with the given capacity.
func NewGlobalMapWithCapacity[T any](capacity int) *GlobalMap[T] {
	return &GlobalMap[T]{data: make(map[string]T, capacity)}
}

func (gm *GlobalMap[T]) GetVal(key string) (T, bool) {
	gm.mu.RLock()
	value, ok := gm.data[key]
	gm.mu.RUnlock()
	if ok {
		return value, true
	}
	var defaultValue T
	return defaultValue, false
}

func (gm *GlobalMap[T]) Get(key string) T {
	gm.mu.RLock()
	value, ok := gm.data[key]
	gm.mu.RUnlock()
	if ok {
		return value
	}
	var defaultValue T
	return defaultValue
}

func (gm *GlobalMap[T]) Set(key string, value T) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.data[key] = value
}

func (gm *GlobalMap[T]) Del(key string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	delete(gm.data, key)
}

func (gm *GlobalMap[T]) Map() map[string]T {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	cp := make(map[string]T, len(gm.data))
	for k, v := range gm.data {
		cp[k] = v
	}
	return cp
}

// Range calls f for each key-value pair. If f returns false, iteration stops.
func (gm *GlobalMap[T]) Range(f func(key string, value T) bool) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	for k, v := range gm.data {
		if !f(k, v) {
			return
		}
	}
}

// Len returns the number of entries in the map.
func (gm *GlobalMap[T]) Len() int {
	gm.mu.RLock()
	n := len(gm.data)
	gm.mu.RUnlock()
	return n
}

var GMap = NewGlobalMap[string]()

func Get(key string) string {
	return GMap.Get(key)
}

func GetVal(key string) (string, bool) {
	return GMap.GetVal(key)
}

func Set(key, value string) {
	GMap.Set(key, value)
}

func Del(key string) {
	GMap.Del(key)
}

func Map() map[string]string {
	return GMap.Map()
}

// Range calls f for each key-value pair in the global GMap. If f returns false, iteration stops.
func Range(f func(key string, value string) bool) {
	GMap.Range(f)
}

// Len returns the number of entries in the global GMap.
func Len() int {
	return GMap.Len()
}
