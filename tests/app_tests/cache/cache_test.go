package cache_tests

import (
	"testing"
	"time"

	"github.com/matheusgomes28/urchin/app"
	"github.com/stretchr/testify/assert"
)

type TrueTimeMockValidator struct{}

func (validator *TrueTimeMockValidator) IsValid(cache *app.EndpointCache) bool {
	return true
}

type FalseTimeMockValidator struct{}

func (validator *FalseTimeMockValidator) IsValid(cache *app.EndpointCache) bool {
	return false
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

	cache := app.MakeCache(1, 10*time.Second, &TrueTimeMockValidator{})

	rolling_size := uint64(0)
	for _, test_case := range test_data {

		rolling_size += uint64(len(test_case.contents))
		err := cache.Store(test_case.name, test_case.contents)
		assert.Nil(t, err)
		assert.Equal(t, cache.Size(), rolling_size)

		endpoint_cache, err := cache.Get(test_case.name)
		assert.Nil(t, err)
		assert.Equal(t, endpoint_cache.Contents, test_case.contents)
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

	cache := app.MakeCache(1, 10*time.Second, &FalseTimeMockValidator{})

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

// Tests that storing over 10MB fails
func TestCacheStoreMaxBytes(t *testing.T) {
	cache := app.MakeCache(1, 10*time.Second, &FalseTimeMockValidator{})

	err := cache.Store("fatty", make([]byte, 10000000))
	assert.Nil(t, err)

	err = cache.Store("slim", make([]byte, 1000))
	assert.NotNil(t, err)
}
