package goCache

import (
	"errors"
	"sync"
	"time"
)

// Provider info struct for storing provider info
type ProviderInfo struct {
	SystemName string `json:"systemName"`
	Address    string `json:"address"`
	Port       int    `json:"port"`
}

// CachedProvider is the value in the map
type CachedProvider struct {
	ProviderInfo
	expireAtTimestamp int64
}

// This will be the cache-object, i.e. it will contain
type LocalCache struct {
	Stop chan struct{} //Sending s non-cost signal to stop any chaneling operation

	Wg            sync.WaitGroup            // Waitgroup -> kind of similar to semaphores
	Mu            sync.RWMutex              // Mutex lock
	SystemService map[string]CachedProvider // A map to store all service and respective provider information
}

func NewLocalCache(cleanupInterval time.Duration) *LocalCache {
	lc := &LocalCache{
		SystemService: make(map[string]CachedProvider),
		Stop:          make(chan struct{}),
	}

	lc.Wg.Add(1)
	// Anonymous function which is run by a gorotuine (thread)
	go func(cleanupInterval time.Duration) {
		// defer decrementation on the waitgroup once the anonymous function is done executing
		defer lc.Wg.Done()
		lc.CleanupLoop(cleanupInterval)
	}(cleanupInterval)

	return lc
}

func (lc *LocalCache) CleanupLoop(interval time.Duration) {

	// After each interval the ticker will send the time on the channel given to the ticker, this interval is used to do cleanup.
	// Each interval a message is sent throught the ticker channel (t.C)
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		//blocks (busy-wait) until there is data available on lc.Stop or t.C
		select {
		case <-lc.Stop:
			return
		case <-t.C:
			lc.Mu.Lock()
			for service, cp := range lc.SystemService {
				//Clearing the cache after each interval
				if cp.expireAtTimestamp <= time.Now().Unix() {
					delete(lc.SystemService, service)
				}
			}
			lc.Mu.Unlock()
		}
	}
}

func (lc *LocalCache) StopCleanup() {
	//Closes the channel Stop
	close(lc.Stop)
	lc.Wg.Wait()
}

func (lc *LocalCache) Update(serviceDefinition string, p ProviderInfo, expireAtTimestamp int64) {
	lc.Mu.Lock()
	defer lc.Mu.Unlock()

	lc.SystemService[serviceDefinition] = CachedProvider{
		ProviderInfo:      p,
		expireAtTimestamp: expireAtTimestamp,
	}
}

var (
	ErrServiceNotInCache = errors.New("the Service isn't in cache")
)

func (lc *LocalCache) Read(serviceDefinition string) (ProviderInfo, error) {
	lc.Mu.RLock()
	defer lc.Mu.RUnlock()

	cu, ok := lc.SystemService[serviceDefinition]
	if !ok {
		return ProviderInfo{}, ErrServiceNotInCache
	}

	return cu.ProviderInfo, nil
}

func (lc *LocalCache) delete(serviceDefinition string) {
	lc.Mu.Lock()
	defer lc.Mu.Unlock()

	delete(lc.SystemService, serviceDefinition)
}
