// Package bench compares gogmap against ecosystem alternatives.
//
// Run: go test -bench=. -benchmem -run=^$
package bench

import (
	"fmt"
	"sync"
	"testing"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/sonnt85/gogmap"
)

const benchKeySet = 1000

// precomputed keys isolate library cost from fmt.Sprintf.
var benchKeys = func() []string {
	keys := make([]string, benchKeySet)
	for i := 0; i < benchKeySet; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}
	return keys
}()

// --- Get parallel (read-only) ---

func BenchmarkGet_Gogmap(b *testing.B) {
	m := gogmap.NewGlobalMap[int]()
	for i, k := range benchKeys {
		m.Set(k, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get(benchKeys[i%benchKeySet])
			i++
		}
	})
}

func BenchmarkGet_SyncMap(b *testing.B) {
	var m sync.Map
	for i, k := range benchKeys {
		m.Store(k, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Load(benchKeys[i%benchKeySet])
			i++
		}
	})
}

func BenchmarkGet_Orcaman(b *testing.B) {
	m := cmap.New[int]()
	for i, k := range benchKeys {
		m.Set(k, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get(benchKeys[i%benchKeySet])
			i++
		}
	})
}

func BenchmarkGet_Xsync(b *testing.B) {
	m := xsync.NewMapOf[string, int]()
	for i, k := range benchKeys {
		m.Store(k, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Load(benchKeys[i%benchKeySet])
			i++
		}
	})
}

// --- Set parallel (write-only) ---

func BenchmarkSet_Gogmap(b *testing.B) {
	m := gogmap.NewGlobalMap[int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(benchKeys[i%benchKeySet], i)
			i++
		}
	})
}

func BenchmarkSet_SyncMap(b *testing.B) {
	var m sync.Map
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Store(benchKeys[i%benchKeySet], i)
			i++
		}
	})
}

func BenchmarkSet_Orcaman(b *testing.B) {
	m := cmap.New[int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(benchKeys[i%benchKeySet], i)
			i++
		}
	})
}

func BenchmarkSet_Xsync(b *testing.B) {
	m := xsync.NewMapOf[string, int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Store(benchKeys[i%benchKeySet], i)
			i++
		}
	})
}

// --- Mixed parallel (90% read, 10% write) ---

func BenchmarkMixed_Gogmap(b *testing.B) {
	m := gogmap.NewGlobalMap[int]()
	for i, k := range benchKeys {
		m.Set(k, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			k := benchKeys[i%benchKeySet]
			if i%10 == 0 {
				m.Set(k, i)
			} else {
				m.Get(k)
			}
			i++
		}
	})
}

func BenchmarkMixed_SyncMap(b *testing.B) {
	var m sync.Map
	for i, k := range benchKeys {
		m.Store(k, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			k := benchKeys[i%benchKeySet]
			if i%10 == 0 {
				m.Store(k, i)
			} else {
				m.Load(k)
			}
			i++
		}
	})
}

func BenchmarkMixed_Orcaman(b *testing.B) {
	m := cmap.New[int]()
	for i, k := range benchKeys {
		m.Set(k, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			k := benchKeys[i%benchKeySet]
			if i%10 == 0 {
				m.Set(k, i)
			} else {
				m.Get(k)
			}
			i++
		}
	})
}

func BenchmarkMixed_Xsync(b *testing.B) {
	m := xsync.NewMapOf[string, int]()
	for i, k := range benchKeys {
		m.Store(k, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			k := benchKeys[i%benchKeySet]
			if i%10 == 0 {
				m.Store(k, i)
			} else {
				m.Load(k)
			}
			i++
		}
	})
}
