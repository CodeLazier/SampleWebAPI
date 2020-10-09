package tests

import (
	"sync"
	"testing"

	"test/cache"
)

func Benchmark_Add_Get(b *testing.B) {
	c := cache.GetInstance()

	wg := sync.WaitGroup{}
	w := func(key int, ca chan int) <-chan int {
		ca <- key
		if err := c.Add(key, cache.CacheItem{}); err == cache.ErrCacheExist {
			b.Fatal(err)
		}
		return ca
	}

	r := func(ca <-chan int) {
		if _, err := c.Get(<-ca); err == cache.ErrCacheNotFound {
			b.Fatal(err)
		}
	}

	wg.Add(b.N)

	ch := make(chan int, 1)
	for i := 0; i < b.N; i++ {
		go func(i int) {
			defer wg.Done()
			r(w(i, ch))
		}(i)
	}

	wg.Wait()

	c.Clear()
}
