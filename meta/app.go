package meta

import (
	log "github.com/Sirupsen/logrus"
	redis "gopkg.in/redis.v3"

	"save.gg/sgg/models/cache"
	dummyCache "save.gg/sgg/models/cache/dummy"
	redisCache "save.gg/sgg/models/cache/redis"
)

type Application struct {
	Cache *cache.Cache
	Conf  Config
	Log   *log.Entry

	Env string
}

var App *Application

func SetupApp() (a *Application, err error) {

	a.Log = log.New().WithFields(log.Fields{})

	configLoc := ResolveConfigLocation()
	a.Conf = NewConfig(configLoc)

	a.Env = a.Conf.Self.Env

	err = a.MountCache()
	if err != nil {
		return a, err
	}

	// err = a.MountPg()
	// if err != nil {
	// 	return a, err
	// }

	return a, err
}

func (a *Application) MountCache() (err error) {
	cacheBackendType := "bolt"
	var cacheBackend cache.CacheBackend

	switch cacheBackendType {
	case "redis":
		cacheBackend, err = redisCache.NewRedisCache(&redis.Options{})
	case "dummy":
		cacheBackend, err = dummyCache.NewDummyCache()
	}

	if err != nil {
		return err
	}

	a.Cache, err = cache.NewCache(cacheBackend)
	if err != nil {
		return err
	}

	return nil
}
