package app

import (
	"fmt"
	"sync/atomic"
	"time"

	shardedmap "github.com/zutto/shardedmap"
)

// Do not store over this amount
// of MBs in the cache
const MAX_CACHE_SIZE_MB = 10;

type EndpointCache struct {
	name string
	contents []byte
	validUntil time.Time
}

func emptyEndpointCache() EndpointCache {
	return EndpointCache{"", []byte{}, time.Now()}
}

type Cache struct {
	cacheMap shardedmap.ShardMap
	cacheTimeout time.Duration
	estimatedSize atomic.Uint64 // in bytes
}

func (self *Cache) Store(name string, buffer []byte) error {
	// Only store to the cache if we have enough space left
	afterSizeMB := float64(self.estimatedSize.Load() + uint64(len(buffer))) / 1000000;
	if afterSizeMB > MAX_CACHE_SIZE_MB {
		return fmt.Errorf("maximum size reached")
	}

	var cache_entry interface{} = EndpointCache{
		name: name,
		contents: buffer,
		validUntil: time.Now().Add(self.cacheTimeout),
	}
	self.cacheMap.Set(name, &cache_entry)
	self.estimatedSize.Add(uint64(len(buffer)))
	return nil
}

func (self *Cache) Get(name string) (EndpointCache, error) {
	// if the endpoint is cached
	cached_entry := self.cacheMap.Get(name)
	if cached_entry != nil {
		cache_contents := (*cached_entry).(EndpointCache)

		// We only return the cache if it's still valid
		if cache_contents.validUntil.After(time.Now()) {
			return cache_contents, nil
		} else {
				self.cacheMap.Delete(name)
				return emptyEndpointCache(), fmt.Errorf("cached endpoint had expired")
		}
	}

	return emptyEndpointCache(), fmt.Errorf("cache does not contain key")
}

func (self *Cache) Size() uint64 {
	return self.estimatedSize.Load()
}

func makeCache(n_shards int, expiry_duration time.Duration) Cache {
	return Cache{
		cacheMap: shardedmap.NewShardMap(n_shards),
		cacheTimeout: expiry_duration,
		estimatedSize: atomic.Uint64{},
	}
}
