package goCache

import (
	"errors"
	"sync"
	"time"
)

type ProviderInfo struct {
	SystemName string `json:"systemName"`
	Address    string `json:"address"`
	Port       int    `json:"port"`
}

type CachedProvider struct {
	ProviderInfo
	expireAtTimestamp int64
}

type LocalCache struct {
	Stop chan struct{}

	Wg            sync.WaitGroup
	Mu            sync.RWMutex
	SystemService map[string]CachedProvider
}

func NewLocalCache(cleanupInterval time.Duration) *LocalCache {
	lc := &LocalCache{
		SystemService: make(map[string]CachedProvider),
		Stop:          make(chan struct{}),
	}

	lc.Wg.Add(1)
	go func(cleanupInterval time.Duration) {
		defer lc.Wg.Done()
		lc.cleanupLoop(cleanupInterval)
	}(cleanupInterval)

	return lc
}

func (lc *LocalCache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-lc.Stop:
			return
		case <-t.C:
			lc.Mu.Lock()
			for uid, cu := range lc.SystemService {
				if cu.expireAtTimestamp <= time.Now().Unix() {
					delete(lc.SystemService, uid)
				}
			}
			lc.Mu.Unlock()
		}
	}
}

func (lc *LocalCache) stopCleanup() {
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
	errUserNotInCache = errors.New("the user isn't in cache")
)

func (lc *LocalCache) Read(serviceDefinition string) (ProviderInfo, error) {
	lc.Mu.RLock()
	defer lc.Mu.RUnlock()

	cu, ok := lc.SystemService[serviceDefinition]
	if !ok {
		return ProviderInfo{}, errUserNotInCache
	}

	return cu.ProviderInfo, nil
}

func (lc *LocalCache) delete(serviceDefinition string) {
	lc.Mu.Lock()
	defer lc.Mu.Unlock()

	delete(lc.SystemService, serviceDefinition)
}
