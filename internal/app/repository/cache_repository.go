package repository

import (
	"time"

	"github.com/allegro/bigcache"

)

type CacheRepository interface {
	Get(key string) ([]byte, error)
	Set(key string, entry []byte) error
	Delete(key string) error
}

func NewCacheRepository() CacheRepository {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic("failed to connect cache")
	}
	return cache
}
