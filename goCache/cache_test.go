package goCache

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestCacheCleanuploop(t *testing.T) {

	newCache := &LocalCache{
		SystemService: make(map[string]CachedProvider),
		Stop:          make(chan struct{}),
	}

	// Empty cache
	want := make(map[string]CachedProvider)

	for i := 0; i < 10; i++ {

		newCache.SystemService[strconv.Itoa(i)] = CachedProvider{
			expireAtTimestamp: time.Now().Unix(),
		}
	}
	// Creating a ticker with a given time duration
	timeDuration := time.Second * 5
	ticker := time.NewTicker(timeDuration)

	// Creating a goroutine with the newCache CleanupLoop function and passing
	// the same time duaration as the one given to the ticker above
	go newCache.CleanupLoop(timeDuration)

	for range ticker.C {
		newCache.StopCleanup()
		break
	}

	// By the end of the time duration, the cache (newCache) should be cleared or else
	// an error is thrown
	if !reflect.DeepEqual(newCache.SystemService, want) {
		t.Error()
	}

}

func TestCacheUpdate(t *testing.T) {

	newCache := &LocalCache{
		SystemService: make(map[string]CachedProvider),
		Stop:          make(chan struct{}),
	}

	// Empty cache
	want := make(map[string]CachedProvider)

	for i := 0; i < 10; i++ {

		want[strconv.Itoa(i)] = CachedProvider{
			expireAtTimestamp: time.Now().Unix(),
		}

		// Testing update
		newCache.Update(strconv.Itoa(i), ProviderInfo{}, time.Now().Unix())
	}

	// Comparing each element in both maps
	for i, x := range newCache.SystemService {

		if !reflect.DeepEqual(x, want[i]) {
			t.Error()
		}

	}

}

func TestCacheRead(t *testing.T) {

	newCache := &LocalCache{
		SystemService: make(map[string]CachedProvider),
		Stop:          make(chan struct{}),
	}

	var test int
	for i := 0; i < 10; i++ {

		newCache.SystemService[strconv.Itoa(i)] = CachedProvider{
			expireAtTimestamp: time.Now().Unix(),
		}
		test = i
	}
	// The element test++ should not exists in newCache
	test++

	if _, err := newCache.Read(strconv.Itoa(test)); err == nil {
		t.Error()
	}

}
