package rache

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("Test New Cache", func(t *testing.T) {
		for i := range 10 {
			cache := NewCache[string, string](i)
			if cache.limit != i {
				t.Errorf("wanted %q got %q", i, cache.limit)
			}
		}
	})

	t.Run("Test Read and Write", func(t *testing.T) {
		tests := []struct {
			key, value string
		}{
			{
				"hello",
				"world",
			},
			{
				"Sleep",
				"Awake",
			},
			{
				"Sit",
				"Stand",
			},
		}
		cache := NewCache[string, string](3)
		for _, tt := range tests {
			cache.Put(tt.key, tt.value)
			val, _ := cache.Get(tt.key)
			if *val != tt.value {
				t.Errorf("wanted %q, got %q", tt.value, *val)
			}
		}
	})

	t.Run("Test entry limit", func(t *testing.T) {
		tests := []struct {
			entryLimit, testLimit int
		}{
			{10, 20},
			{5, 10},
			{20, 10},
			{0, 10},
		}
		for _, tt := range tests {
			cache := NewCache[int, string](tt.entryLimit)
			for i := range tt.testLimit {
				cache.Put(i, fmt.Sprintf("data-%d", i))
			}

			if len(cache.entries) > tt.entryLimit {
				t.Errorf("wanted %d, got %d", tt.entryLimit, len(cache.entries))
			}
		}
	})

}

func TestCacheConcurrent(t *testing.T) {
	t.Run("Test concurrent writes", func(t *testing.T) {
		var wg sync.WaitGroup
		cache := NewCache[int, string](10)
		for i := range 10 {
			wg.Add(1)
			go func() {
				cache.Put(i, fmt.Sprintf("data-%d", i))
				wg.Done()

			}()
		}
		wg.Wait()

		var reads int
		for i := range 10 {
			if _, ok := cache.Get(i); ok {
				reads++
			}
		}

		if reads != 10 {
			t.Errorf("wanted 10 got %d", reads)
		}

	})

	t.Run("Test concurrent reads", func(t *testing.T) {
		var wg sync.WaitGroup
		cache := NewCache[int, string](10)
		for i := range 10 {
			wg.Add(1)
			go func() {
				cache.Put(i, fmt.Sprintf("data-%d", i))
				wg.Done()

			}()
		}
		wg.Wait()

		var reads atomic.Int32
		for i := range 10 {
			wg.Add(1)
			go func() {
				if _, ok := cache.Get(i); ok {
					reads.Add(1)
				}
				wg.Done()

			}()
		}
		wg.Wait()
		if reads.Load() != int32(cache.successfulReads) {
			t.Errorf("wanted %d got %d", cache.successfulReads, reads.Load())
		}

	})
}

func TestEvictionAlgorithms(t *testing.T) {
	t.Run("Test LRU algorithm", func(t *testing.T) {
		cache := NewCache[string, string](2)
		cache.Put("hello", "world")
		cache.Put("foo", "bar")
		cache.Put("johnny", "cage")
		if len(cache.entries) > cache.limit {
			t.Errorf("wanted %d, got %d", len(cache.entries), cache.limit)
		}
	})

}

func BenchmarkLRUEviction(b *testing.B) {
	cache := NewCache[int, string](5)
	for range b.N {
		for i := range 100 {
			val := fmt.Sprintf("data %d", i)
			cache.Put(i, val)
		}
	}
}
