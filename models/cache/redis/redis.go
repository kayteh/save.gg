package rediscache

import (
	redis "gopkg.in/redis.v3"
	"time"
)

type RedisCache struct {
	Redis *redis.Client
}

func NewRedisCache(o *redis.Options) (c *RedisCache, err error) {
	//return nil, errors.New("not implemented")

	r := redis.NewClient(o)
	return &RedisCache{
		Redis: r,
	}, err
}

func (c *RedisCache) Get(key string) (s []byte, err error) {
	o, err := c.Redis.Get(key).Bytes()

	if err != nil && err.Error() != "redis: nil" {
		return s, err
	}

	return o, nil
}

func (c *RedisCache) Set(key string, data []byte, ttl time.Duration) (err error) {
	_, err = c.Redis.Set(key, string(data), ttl).Result()
	return err
}

func (c *RedisCache) Del(key string) (err error) {
	_, err = c.Redis.Del(key).Result()
	return err
}
