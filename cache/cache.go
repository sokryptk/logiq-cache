package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	DEFAULT_EXPIRATION = time.Minute * 5
)

type Object struct {
	Value      any
	Visited    int
	Expiration int64 // epoch time for this key to expire
}

type Cache struct {
	capacity   int
	expiration time.Duration // default expiration is 5minutes
	objects    map[string]*Object
	mutex      sync.RWMutex
	Logging    bool
}

func New(capacity int) *Cache {
	objects := make(map[string]*Object)

	cache := Cache{
		capacity:   capacity,
		expiration: DEFAULT_EXPIRATION,
		objects:    objects,
	}

	go cache.cleanup()

	return &cache
}

func (c *Cache) debugLogf(format string, a ...any) {
	if c.Logging {
		log.Println(fmt.Sprintf(format, a...))
	}
}

// Takes place for newer objects
// Old objects remain unaltered
func (c *Cache) SetExpiration(duration time.Duration) {
	c.expiration = duration
}

// Values() retrieves back all values currently in Cache.
func (c *Cache) Values() map[string]*Object {
	return c.objects
}

func (c *Cache) Retrieve(key string) (value any, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	result, ok := c.objects[key]
	if !ok {
		return nil, fmt.Errorf("key not found in cache")
	}

	result.Visited += 1
	c.debugLogf("Retrieved key :%s", key)
	return result.Value, nil
}

// Stores the object while being under capacity.
func (c *Cache) Store(key string, value any) {
	if len(c.objects) >= c.capacity {
		c.removeLeastVisited()
		// deletes the key which has been visited the least
	}

	c.mutex.Lock()
	c.objects[key] = &Object{
		Value:      value,
		Expiration: time.Now().Add(c.expiration).Unix(),
	}
	c.mutex.Unlock()
	c.debugLogf("Stored key : %s with value : %v", key, value)
}

// Delete the object at key.
// If the key doesn't exist, it is a no-op.
func (c *Cache) Delete(key string) bool {
	c.mutex.Lock()
	delete(c.objects, key)
	c.mutex.Unlock()
	return true
}

func (c *Cache) cleanup() {
	timer := time.NewTicker(c.expiration / 10)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			c.mutex.RLock()
			for k, v := range c.objects {
				if v.Expiration < time.Now().Unix() {
					c.debugLogf("Cleaning up key : %s", k)

					// No need of using c.Delete() since we're already locking through mutex
					delete(c.objects, k)
				}
			}
			c.mutex.RUnlock()
		}
	}
}

func (c *Cache) removeLeastVisited() {
	// Highest number
	leastVisited := 2 << 31
	var leastVisitedKey string

	c.mutex.RLock()
	for k, v := range c.objects {
		if v.Visited < leastVisited {
			leastVisited = v.Visited
			leastVisitedKey = k
		}
	}

	c.mutex.RUnlock()
	c.debugLogf("LRU Key : %s, proceeding for deletion", leastVisitedKey)
	c.Delete(leastVisitedKey)
}