package app

import (
	"fmt"
	"sync/atomic"
	"time"

	shardedmap "github.com/zutto/shardedmap"
)

// Do not store over this amount
// of MBs in the cache
const MAX_CACHE_SIZE_MB = 10

type EndpointCache struct {
	name       string
	contents   []byte
	validUntil time.Time
}

func emptyEndpointCache() EndpointCache {
	return EndpointCache{"", []byte{}, time.Now()}
}

type Cache interface {
	Get(name string) (EndpointCache, error)
	Store(name string, buffer []byte) error
	Size() uint64
}

type CacheValidator interface {
	IsValid(cache *EndpointCache) bool
}

type TimeValidator struct{}

func (validator *TimeValidator) IsValid(cache *EndpointCache) bool {
	// We only return the cache if it's still valid
	return cache.validUntil.After(time.Now())
}

type TimedCache struct {
	cacheMap      shardedmap.ShardMap
	cacheTimeout  time.Duration
	estimatedSize atomic.Uint64 // in bytes
	validator     CacheValidator
}

func (cache *TimedCache) Store(name string, buffer []byte) error {
	// Only store to the cache if we have enough space left
	afterSizeMB := float64(cache.estimatedSize.Load()+uint64(len(buffer))) / 1000000
	if afterSizeMB > MAX_CACHE_SIZE_MB {
		return fmt.Errorf("maximum size reached")
	}

	var cache_entry interface{} = EndpointCache{
		name:       name,
		contents:   buffer,
		validUntil: time.Now().Add(cache.cacheTimeout),
	}
	cache.cacheMap.Set(name, &cache_entry)
	cache.estimatedSize.Add(uint64(len(buffer)))
	return nil
}

func (cache *TimedCache) Get(name string) (EndpointCache, error) {
	// if the endpoint is cached
	cached_entry := cache.cacheMap.Get(name)
	if cached_entry != nil {
		cache_contents := (*cached_entry).(EndpointCache)

		// We only return the cache if it's still valid
		if cache.validator.IsValid(&cache_contents) {
			return cache_contents, nil
		} else {
			cache.cacheMap.Delete(name)
			return emptyEndpointCache(), fmt.Errorf("cached endpoint had expired")
		}
	}

	return emptyEndpointCache(), fmt.Errorf("cache does not contain key")
}

func (cache *TimedCache) Size() uint64 {
	return cache.estimatedSize.Load()
}

func makeCache(n_shards int, expiry_duration time.Duration) Cache {
	return &TimedCache{
		cacheMap:      shardedmap.NewShardMap(n_shards),
		cacheTimeout:  expiry_duration,
		estimatedSize: atomic.Uint64{},
		validator:     &TimeValidator{},
	}
}
