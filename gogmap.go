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

func (gm *GlobalMap[T]) Get(key string) T {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	if value, ok := gm.data[key]; ok {
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
	gm.mu.Lock()
	defer gm.mu.Unlock()
	return gm.data
}

var GMap = NewGlobalMap[string]()

func Get(key string) string {
	return GMap.Get(key)
}

func Set(key, value string) {
	GMap.Set(key, value)
}

func Del(key string) {
	GMap.Del(key)
}
func Map() map[string]string {
	return GMap.data
}
