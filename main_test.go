package main

import (
	"fmt"
	"sync"
	"testing"

	"logiq.ai/cache/cache"
)

func TestCache(t *testing.T) {
	c := cache.New(5)

	c.Store("key", true)
	value, err := c.Retrieve("key")
	if err != nil {
		t.Error("Should have found the key")
	}

	if value != true {
		t.Errorf("Expected value : true, found : %v", value)
	}

	ok := c.Delete("key")
	if !ok {
		t.Error("Should have deleted")
	}

	val, err := c.Retrieve("key")
	if err == nil {
		t.Errorf("Expected error, found value : %v at %s", val, "key")
	}
}

func TestCacheEviction(t *testing.T) {
	c := cache.New(5)

	// Populate the cache
	for i := 1; i <= 5; i++ {
		c.Store(fmt.Sprintf("key%d", i), i)
	}

	for i := 2; i <= 5; i++ {
		// revisit all except the first one aka key1
		c.Retrieve(fmt.Sprintf("key%d", i))
	}

	// this should trigger our eviction policy
	// key1 should be booted off.
	c.Store("key6", 6)

	_, err := c.Retrieve(fmt.Sprintf("key1"))
	if err == nil {
		t.Errorf("eviction policy failed")
	}

}

func BenchmarkCache(b *testing.B) {
	c := cache.New(100)

	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Store(fmt.Sprintf("key%d", i), i)
		}(i)
	}

	wg.Wait()
}