package rache

import (
	"sync"
	"time"
)

// Interface for cache eviction policies
type Policy[T comparable, V any] interface {
	Evict(data map[T]*cacheEntry[V]) (T, bool)
	Insert(val T)
	Update(val T)
}

type cacheEntryStats struct {
	reads        int
	writes       int
	lastAccessed time.Time
}

type cacheEntry[V any] struct {
	value V
	stats cacheEntryStats
}

type cache[T comparable, V any] struct {
	sync.Mutex
	limit           int
	totalReads      int
	successfulReads int
	totalWrites     int
	entries         map[T]*cacheEntry[V]
	policy          Policy[T, V]
}

func NewCache[T comparable, V any](entryLimit int) *cache[T, V] {
	LRU := NewLRUPolicy[T, V]()
	return &cache[T, V]{
		limit:   entryLimit,
		entries: make(map[T]*cacheEntry[V], entryLimit),
		policy:  LRU,
	}
}

// Put entries in cache and returns true if entry already existed
// If entry already present, replace entry.
func (c *cache[T, V]) Put(key T, val V) bool {
	if c.limit == 0 {
		return false
	}

	c.Lock()
	defer c.Unlock()

	e, present := c.entries[key]

	if !present {
		// remove least recently used if entries is filled to the limit
		if len(c.entries) == c.limit {
			k, ok := c.policy.Evict(c.entries)
			if ok {
				delete(c.entries, k)
			}
		}
		e = &cacheEntry[V]{}
		c.entries[key] = e
		c.policy.Insert(key)
	} else {
		c.policy.Update(key)
	}
	e.value = val
	e.stats.lastAccessed = time.Now()
	e.stats.writes++
	c.totalWrites++
	return present
}

// Get returns the value of the key passed and a boolean that shows whether the entry existed
// If the entry did not exist, the boolean is false and value is nil.
func (c *cache[T, V]) Get(key T) (V, bool) {
	c.Lock()
	defer c.Unlock()
	var zero V
	e, ok := c.entries[key]
	c.totalReads++
	if !ok {
		return zero, false
	}
	c.policy.Update(key)
	c.successfulReads++
	e.stats.reads++
	e.stats.lastAccessed = time.Now()
	return e.value, ok
}

// This is a simple algorithm using a doubly linked list where the least recently accessed entries are
// removed when the cache is full.
//
// The idea is anytime an entry is accessed, it's moved to the head of the list.
// Following that trend, the least recently accessed entry will always be at the end of the list.
//
// Example: We have these entries in our list
// 1 -> 2 -> 3. If we access 2, then it will be updated as;
// 2 -> 1 -> 3. If we access 3, then it will be updated as;
// 3 -> 2 -> 1. If we proceed to access 2 again, then it will be updated as;
// 2 -> 3 -> 1. Looking at the trend, 1 is the least recently accessed and thus at the tail of the list.
// Thus 1 will be removed if there needs to be an eviction
type LRUPolicy[T comparable, V any] struct {
	list *DLL[T]
}

func NewLRUPolicy[T comparable, V any]() *LRUPolicy[T, V] {
	dll := &DLL[T]{}
	return &LRUPolicy[T, V]{list: dll}
}

// Evicts the least recently used(LRU) data
// The LRU data will be at the tail of the list.
func (p *LRUPolicy[T, V]) Evict(data map[T]*cacheEntry[V]) (T, bool) {
	return p.list.deleteBack()
}

// Insert new values to cache
func (p *LRUPolicy[T, V]) Insert(val T) {
	p.list.insertFront(val)
}

// Update values in cache.
// This moves a value in the linked list to the front of the list
func (p *LRUPolicy[T, V]) Update(val T) {
	ok := p.list.remove(val)
	if ok {
		p.list.insertFront(val)
	}
}
