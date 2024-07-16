package main

import (
	"fmt"

	"github.com/quamejnr/rache/rache"
)

func main() {
	cache := rache.NewCache[int, string](10)
	cache.Put(1, "world")
	val, ok := cache.Get(1)
	if ok {
		fmt.Println(val)
	}
}
