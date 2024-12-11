package cache

import (
    "sync"
    "time"
)

type VolumeInfo struct {
    Symbol      string `json:"Symbol"`
    Volume24h   string `json:"24h Volume"`
    LastUpdated string `json:"Last Updated"`
}

type VolumeCache struct {
    Data        VolumeInfo
    LastUpdated time.Time
}

type VolumeCacheManager struct {
    cache      map[string]VolumeCache
    mutex      sync.RWMutex
    expiration time.Duration
}

var (
    volumeCacheManager *VolumeCacheManager
    once              sync.Once
)

// GetVolumeCacheManager returns a singleton instance of VolumeCacheManager
func GetVolumeCacheManager() *VolumeCacheManager {
    once.Do(func() {
        volumeCacheManager = &VolumeCacheManager{
            cache:      make(map[string]VolumeCache),
            expiration: 15 * time.Minute,
        }
    })
    return volumeCacheManager
}

// Get retrieves cached volume data for a symbol
func (m *VolumeCacheManager) Get(symbol string) (VolumeInfo, bool) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()

    if cache, exists := m.cache[symbol]; exists {
        if time.Since(cache.LastUpdated) < m.expiration {
            return cache.Data, true
        }
    }
    return VolumeInfo{}, false
}

// Set stores volume data in cache
func (m *VolumeCacheManager) Set(symbol string, data VolumeInfo) {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    m.cache[symbol] = VolumeCache{
        Data:        data,
        LastUpdated: time.Now(),
    }
}

// Clear removes expired entries from cache
func (m *VolumeCacheManager) Clear() {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    for symbol, cache := range m.cache {
        if time.Since(cache.LastUpdated) >= m.expiration {
            delete(m.cache, symbol)
        }
    }
}

// StartCleanupRoutine starts a goroutine to periodically clean up expired cache entries
func (m *VolumeCacheManager) StartCleanupRoutine() {
    go func() {
        ticker := time.NewTicker(m.expiration)
        for range ticker.C {
            m.Clear()
        }
    }()
}