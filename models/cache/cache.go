// Abstraction layer on top of the cache system. Lets us change it
// without needing to change app code.
package cache

import (
	"encoding/json"
	"errors"
	"time"
)

type Cache struct {
	backend CacheBackend
}

type CacheBackend interface {
	Get(string) ([]byte, error)
	Set(string, []byte, time.Duration) error
	Del(string) error
}

func NewCache(b CacheBackend) (c *Cache, err error) {

	return &Cache{backend: b}, nil
}

func (c *Cache) Get(key string, v interface{}) (err error) {

	s, err := c.backend.Get(key)
	if err != nil {
		return err
	}

	if len(s) == 0 {
		err = c.Del(key)
		if err != nil {
			return err
		}

		return errors.New("cache miss")
	}

	err = json.Unmarshal(s, v)
	return err

}

func (c *Cache) Set(key string, data interface{}, ttl time.Duration) (err error) {

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = c.backend.Set(key, d, ttl)
	return err

}

func (c *Cache) Del(key string) (err error) {

	err = c.backend.Del(key)
	return err

}
