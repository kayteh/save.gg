package dummycache

import (
	"sync"
	"time"
)

type DummyCache struct {
	hash map[string]record
	mu   sync.RWMutex
}

type record struct {
	Value []byte
	TTL   time.Time
}

func NewDummyCache() (c *DummyCache, err error) {
	return &DummyCache{
		hash: make(map[string]record),
		mu:   sync.RWMutex{},
	}, err
}

func (c *DummyCache) Get(key string) (s []byte, err error) {

	c.mu.RLock()
	d := c.hash[key]
	c.mu.RUnlock()

	if time.Now().After(d.TTL) {
		c.Del(key)
		return nil, nil
	}

	return d.Value, nil
}

func (c *DummyCache) Set(key string, data []byte, ttl time.Duration) (err error) {

	ttlTime := time.Now().Add(ttl)

	c.mu.Lock()
	c.hash[key] = record{
		Value: data,
		TTL:   ttlTime,
	}
	c.mu.Unlock()

	return nil
}

func (c *DummyCache) Del(key string) (err error) {
	c.mu.Lock()
	delete(c.hash, key)
	c.mu.Unlock()

	return nil
}
