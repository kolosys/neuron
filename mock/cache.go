package mock

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/kolosys/neuron"
)

// MockCache is a mock implementation of neuron.Cache for testing.
type MockCache struct {
	data       map[string]neuron.CacheEntry
	hits       atomic.Int64
	misses     atomic.Int64
	sets       atomic.Int64
	deletes    atomic.Int64
	clears     atomic.Int64
	operations []CacheOperation
	recordOps  atomic.Bool
	mu         sync.RWMutex
}

// MockCacheOptions configures mock cache behavior.
type MockCacheOptions struct {
	RecordOperations bool
	InitialData      map[string]neuron.CacheEntry
}

// NewMockCache creates a new mock cache.
func NewMockCache(opts *MockCacheOptions) *MockCache {
	mc := &MockCache{
		data:       make(map[string]neuron.CacheEntry),
		operations: make([]CacheOperation, 0),
	}

	if opts == nil {
		mc.recordOps.Store(true)
		return mc
	}

	mc.recordOps.Store(opts.RecordOperations)

	if opts.InitialData != nil {
		mc.mu.Lock()
		for k, v := range opts.InitialData {
			mc.data[k] = v
		}
		mc.mu.Unlock()
	}

	return mc
}

// Get retrieves an entry from the cache.
func (mc *MockCache) Get(key string) (*neuron.CacheEntry, bool) {
	mc.mu.RLock()
	entry, exists := mc.data[key]
	mc.mu.RUnlock()

	if exists {
		mc.hits.Add(1)
	} else {
		mc.misses.Add(1)
	}

	if mc.recordOps.Load() {
		op := CacheOperation{
			Op:    "get",
			Key:   key,
			Time:  time.Now().UnixNano(),
			Hit:   exists,
			Entry: nil,
		}
		if exists {
			op.Entry = &entry
		}

		mc.mu.Lock()
		mc.operations = append(mc.operations, op)
		mc.mu.Unlock()
	}

	if exists {
		return &entry, true
	}
	return nil, false
}

// Set stores an entry in the cache.
func (mc *MockCache) Set(key string, entry neuron.CacheEntry) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.data[key] = entry
	mc.sets.Add(1)

	if mc.recordOps.Load() {
		op := CacheOperation{
			Op:    "set",
			Key:   key,
			Time:  time.Now().UnixNano(),
			Entry: entry,
		}
		mc.operations = append(mc.operations, op)
	}
}

// Delete removes an entry from the cache.
func (mc *MockCache) Delete(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.data, key)
	mc.deletes.Add(1)

	if mc.recordOps.Load() {
		op := CacheOperation{
			Op:   "delete",
			Key:  key,
			Time: time.Now().UnixNano(),
		}
		mc.operations = append(mc.operations, op)
	}
}

// Clear removes all entries from the cache.
func (mc *MockCache) Clear() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.data = make(map[string]neuron.CacheEntry)
	mc.clears.Add(1)

	if mc.recordOps.Load() {
		op := CacheOperation{
			Op:   "clear",
			Time: time.Now().UnixNano(),
		}
		mc.operations = append(mc.operations, op)
	}
}

// Hits returns the number of cache hits.
func (mc *MockCache) Hits() int64 {
	return mc.hits.Load()
}

// Misses returns the number of cache misses.
func (mc *MockCache) Misses() int64 {
	return mc.misses.Load()
}

// Sets returns the number of cache sets.
func (mc *MockCache) Sets() int64 {
	return mc.sets.Load()
}

// Deletes returns the number of cache deletes.
func (mc *MockCache) Deletes() int64 {
	return mc.deletes.Load()
}

// Clears returns the number of cache clears.
func (mc *MockCache) Clears() int64 {
	return mc.clears.Load()
}

// Size returns the number of entries in the cache.
func (mc *MockCache) Size() int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return len(mc.data)
}

// Contains checks if a key exists in the cache.
func (mc *MockCache) Contains(key string) bool {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	_, exists := mc.data[key]
	return exists
}

// Keys returns a copy of all keys in the cache.
func (mc *MockCache) Keys() []string {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	keys := make([]string, 0, len(mc.data))
	for k := range mc.data {
		keys = append(keys, k)
	}
	return keys
}

// HitRate returns the ratio of hits to total requests.
func (mc *MockCache) HitRate() float64 {
	hits := mc.hits.Load()
	misses := mc.misses.Load()
	total := hits + misses

	if total == 0 {
		return 0
	}

	return float64(hits) / float64(total)
}

// Operations returns a copy of all recorded operations.
func (mc *MockCache) Operations() []CacheOperation {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	ops := make([]CacheOperation, len(mc.operations))
	copy(ops, mc.operations)
	return ops
}

// ClearRecorded clears all recorded operations and statistics.
func (mc *MockCache) ClearRecorded() {
	mc.mu.Lock()
	mc.operations = mc.operations[:0]
	mc.mu.Unlock()

	mc.hits.Store(0)
	mc.misses.Store(0)
	mc.sets.Store(0)
	mc.deletes.Store(0)
	mc.clears.Store(0)
}

// Reset clears all cache state, statistics, and recorded operations.
func (mc *MockCache) Reset() {
	mc.mu.Lock()
	mc.data = make(map[string]neuron.CacheEntry)
	mc.operations = mc.operations[:0]
	mc.mu.Unlock()

	mc.hits.Store(0)
	mc.misses.Store(0)
	mc.sets.Store(0)
	mc.deletes.Store(0)
	mc.clears.Store(0)
}

// GetEntry returns a copy of an entry without updating hit counts.
func (mc *MockCache) GetEntry(key string) (*neuron.CacheEntry, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, exists := mc.data[key]
	if exists {
		return &entry, true
	}
	return nil, false
}

// EnableRecording enables or disables operation recording.
func (mc *MockCache) EnableRecording(enabled bool) {
	mc.recordOps.Store(enabled)
}
