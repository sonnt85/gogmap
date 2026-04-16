package gogmap

import (
	"hash/fnv"
	"sync"
)

const defaultShards = 16

type shard[T any] struct {
	mu   sync.RWMutex
	data map[string]T
}

// GlobalMap is a sharded concurrent map for reduced lock contention.
type GlobalMap[T any] struct {
	shards [defaultShards]shard[T]
}

func NewGlobalMap[T any]() *GlobalMap[T] {
	gm := &GlobalMap[T]{}
	for i := range gm.shards {
		gm.shards[i].data = make(map[string]T)
	}
	return gm
}

// NewGlobalMapWithCapacity creates a GlobalMap pre-allocated with the given capacity spread across shards.
func NewGlobalMapWithCapacity[T any](capacity int) *GlobalMap[T] {
	perShard := capacity/defaultShards + 1
	gm := &GlobalMap[T]{}
	for i := range gm.shards {
		gm.shards[i].data = make(map[string]T, perShard)
	}
	return gm
}

//go:inline
func (gm *GlobalMap[T]) getShard(key string) *shard[T] {
	h := fnv.New32a()
	h.Write([]byte(key))
	return &gm.shards[h.Sum32()%defaultShards]
}

func (gm *GlobalMap[T]) GetVal(key string) (T, bool) {
	s := gm.getShard(key)
	s.mu.RLock()
	value, ok := s.data[key]
	s.mu.RUnlock()
	if ok {
		return value, true
	}
	var zero T
	return zero, false
}

func (gm *GlobalMap[T]) Get(key string) T {
	s := gm.getShard(key)
	s.mu.RLock()
	value, ok := s.data[key]
	s.mu.RUnlock()
	if ok {
		return value
	}
	var zero T
	return zero
}

func (gm *GlobalMap[T]) Set(key string, value T) {
	s := gm.getShard(key)
	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()
}

func (gm *GlobalMap[T]) Del(key string) {
	s := gm.getShard(key)
	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()
}

func (gm *GlobalMap[T]) Map() map[string]T {
	cp := make(map[string]T)
	for i := range gm.shards {
		s := &gm.shards[i]
		s.mu.RLock()
		for k, v := range s.data {
			cp[k] = v
		}
		s.mu.RUnlock()
	}
	return cp
}

// Range calls f for each key-value pair. If f returns false, iteration stops.
func (gm *GlobalMap[T]) Range(f func(key string, value T) bool) {
	for i := range gm.shards {
		s := &gm.shards[i]
		s.mu.RLock()
		for k, v := range s.data {
			if !f(k, v) {
				s.mu.RUnlock()
				return
			}
		}
		s.mu.RUnlock()
	}
}

// Len returns the number of entries in the map.
func (gm *GlobalMap[T]) Len() int {
	n := 0
	for i := range gm.shards {
		s := &gm.shards[i]
		s.mu.RLock()
		n += len(s.data)
		s.mu.RUnlock()
	}
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
