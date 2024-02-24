package app

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	shardedmap "github.com/zutto/shardedmap"
)

type TrueTimeMockValidator struct{}

func (validator *TrueTimeMockValidator) IsValid(cache *EndpointCache) bool {
	return true
}

type FalseTimeMockValidator struct{}

func (validator *FalseTimeMockValidator) IsValid(cache *EndpointCache) bool {
	return false
}

func makeTrueCacheMock() Cache {
	return &TimedCache{
		cacheMap:      shardedmap.NewShardMap(2),
		cacheTimeout:  10 * time.Second,
		estimatedSize: atomic.Uint64{},
		validator:     &TrueTimeMockValidator{},
	}
}

func makeFalseCacheMock() Cache {
	return &TimedCache{
		cacheMap:      shardedmap.NewShardMap(2),
		cacheTimeout:  10 * time.Second,
		estimatedSize: atomic.Uint64{},
		validator:     &FalseTimeMockValidator{},
	}
}

func TestCacheAddition(t *testing.T) {

	test_data := []struct {
		name     string
		contents []byte
	}{
		{"first", []byte("hello")},
		{"second", []byte("the quick brown fox does some weird stuff")},
		{"veryLonNameThatIsProbablyTooLong", []byte("Hello there my friends")},
		{"nPe9Rkff6ER6EzAxPUIpxc8UBBLm71hhq2MO9hkQWisrfihUqv", []byte("oA7Hv1A7vOuZSKrPT4ZN5DGKNSHZqpLEvUA5hu54CMyIt8c78u")},
	}

	cache := makeTrueCacheMock()

	rolling_size := uint64(0)
	for _, test_case := range test_data {

		rolling_size += uint64(len(test_case.contents))
		err := cache.Store(test_case.name, test_case.contents)
		assert.Nil(t, err)
		assert.Equal(t, cache.Size(), rolling_size)

		endpoint_cache, err := cache.Get(test_case.name)
		assert.Nil(t, err)
		assert.Equal(t, endpoint_cache.contents, test_case.contents)
	}
}

func TestCacheFailure(t *testing.T) {

	test_data := []struct {
		name     string
		contents []byte
	}{
		{"first", []byte("hello")},
		{"second", []byte("the quick brown fox does some weird stuff")},
		{"veryLonNameThatIsProbablyTooLong", []byte("Hello there my friends")},
		{"nPe9Rkff6ER6EzAxPUIpxc8UBBLm71hhq2MO9hkQWisrfihUqv", []byte("oA7Hv1A7vOuZSKrPT4ZN5DGKNSHZqpLEvUA5hu54CMyIt8c78u")},
	}

	cache := makeFalseCacheMock()

	rolling_size := uint64(0)
	for _, test_case := range test_data {

		rolling_size += uint64(len(test_case.contents))
		err := cache.Store(test_case.name, test_case.contents)
		assert.Nil(t, err)
		assert.Equal(t, cache.Size(), rolling_size)

		_, err = cache.Get(test_case.name)
		assert.NotNil(t, err)
	}
}
