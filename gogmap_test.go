package gogmap

import (
	"fmt"
	"sync"
	"testing"
)

func TestBasicOps(t *testing.T) {
	m := NewGlobalMap[int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	if v := m.Get("a"); v != 1 {
		t.Errorf("Get(a) = %d, want 1", v)
	}
	if v, ok := m.GetVal("b"); !ok || v != 2 {
		t.Errorf("GetVal(b) = %d,%v, want 2,true", v, ok)
	}
	if _, ok := m.GetVal("missing"); ok {
		t.Error("GetVal(missing) should return false")
	}
	if m.Len() != 3 {
		t.Errorf("Len() = %d, want 3", m.Len())
	}

	m.Del("b")
	if m.Len() != 2 {
		t.Errorf("Len() after Del = %d, want 2", m.Len())
	}

	cp := m.Map()
	if len(cp) != 2 {
		t.Errorf("Map() len = %d, want 2", len(cp))
	}

	count := 0
	m.Range(func(k string, v int) bool {
		count++
		return true
	})
	if count != 2 {
		t.Errorf("Range count = %d, want 2", count)
	}
}

func TestRangeEarlyStop(t *testing.T) {
	m := NewGlobalMap[int]()
	for i := 0; i < 100; i++ {
		m.Set(fmt.Sprintf("key%d", i), i)
	}
	count := 0
	m.Range(func(k string, v int) bool {
		count++
		return count < 5
	})
	if count != 5 {
		t.Errorf("Range early stop: count = %d, want 5", count)
	}
}

func TestConcurrent(t *testing.T) {
	m := NewGlobalMap[int]()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("k%d", n)
			m.Set(key, n)
			m.Get(key)
			m.GetVal(key)
			m.Del(key)
		}(i)
	}
	wg.Wait()
}

// benchKeys returns a pre-computed slice of n keys to avoid
// measuring fmt.Sprintf cost inside benchmark hot loops.
func benchKeys(n int) []string {
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}
	return keys
}

func BenchmarkGetParallel(b *testing.B) {
	const n = 1000
	keys := benchKeys(n)
	m := NewGlobalMap[string]()
	for _, k := range keys {
		m.Set(k, "value")
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get(keys[i%n])
			i++
		}
	})
}

func BenchmarkSetParallel(b *testing.B) {
	const n = 1000
	keys := benchKeys(n)
	m := NewGlobalMap[string]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(keys[i%n], "value")
			i++
		}
	})
}

func BenchmarkMixedParallel(b *testing.B) {
	const n = 1000
	keys := benchKeys(n)
	m := NewGlobalMap[string]()
	for _, k := range keys {
		m.Set(k, "value")
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			k := keys[i%n]
			if i%10 == 0 {
				m.Set(k, "new")
			} else {
				m.Get(k)
			}
			i++
		}
	})
}
