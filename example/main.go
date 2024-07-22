package main

import (
	"fmt"

	"github.com/quamejnr/rache/rache"
)

func main() {
	// using default cache policy
	cache := rache.NewCache[int, string](10)
	cache.Put(1, "world")
	val, ok := cache.Get(1)
	if ok {
		fmt.Println(val)
	}
	// You can change the default cache policy by using
	lruTime := rache.NewLRUTimePolicy[int, string]()
	cache.Policy = lruTime
  cache.Put(2, "Love")
	val, ok = cache.Get(2)
	if ok {
		fmt.Println(val)
	}
}
