package mock_test

import (
	"testing"
	"time"

	"github.com/kolosys/neuron"
	"github.com/kolosys/neuron/mock"
)

func TestMockCache_Get_Miss(t *testing.T) {
	cache := mock.NewMockCache(nil)

	entry, hit := cache.Get("nonexistent")

	if hit {
		t.Error("expected miss")
	}
	if entry != nil {
		t.Error("expected nil entry")
	}
	if cache.Misses() != 1 {
		t.Errorf("expected 1 miss, got %d", cache.Misses())
	}
}

func TestMockCache_Set_Get_Hit(t *testing.T) {
	cache := mock.NewMockCache(nil)

	testEntry := neuron.CacheEntry{
		Data:       []byte(`test`),
		StatusCode: 200,
		Timestamp:  time.Now(),
	}

	cache.Set("key1", testEntry)
	entry, hit := cache.Get("key1")

	if !hit {
		t.Error("expected hit")
	}
	if entry == nil {
		t.Fatal("expected non-nil entry")
	}
	if string(entry.Data) != "test" {
		t.Errorf("expected 'test', got %s", entry.Data)
	}
	if cache.Hits() != 1 {
		t.Errorf("expected 1 hit, got %d", cache.Hits())
	}
}

func TestMockCache_Statistics(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})
	cache.Set("key2", neuron.CacheEntry{})

	cache.Get("key1") // hit
	cache.Get("key2") // hit
	cache.Get("key3") // miss

	if cache.Sets() != 2 {
		t.Errorf("expected 2 sets, got %d", cache.Sets())
	}
	if cache.Hits() != 2 {
		t.Errorf("expected 2 hits, got %d", cache.Hits())
	}
	if cache.Misses() != 1 {
		t.Errorf("expected 1 miss, got %d", cache.Misses())
	}
}

func TestMockCache_Delete(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})
	cache.Delete("key1")

	if cache.Deletes() != 1 {
		t.Errorf("expected 1 delete, got %d", cache.Deletes())
	}

	_, hit := cache.Get("key1")
	if hit {
		t.Error("expected miss after delete")
	}
}

func TestMockCache_Clear(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})
	cache.Set("key2", neuron.CacheEntry{})

	cache.Clear()

	if cache.Clears() != 1 {
		t.Errorf("expected 1 clear, got %d", cache.Clears())
	}

	if cache.Size() != 0 {
		t.Errorf("expected size 0, got %d", cache.Size())
	}
}

func TestMockCache_Size(t *testing.T) {
	cache := mock.NewMockCache(nil)

	if cache.Size() != 0 {
		t.Error("expected size 0 initially")
	}

	cache.Set("key1", neuron.CacheEntry{})
	cache.Set("key2", neuron.CacheEntry{})

	if cache.Size() != 2 {
		t.Errorf("expected size 2, got %d", cache.Size())
	}
}

func TestMockCache_Contains(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})

	if !cache.Contains("key1") {
		t.Error("expected key1 to exist")
	}

	if cache.Contains("key2") {
		t.Error("expected key2 to not exist")
	}
}

func TestMockCache_Keys(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})
	cache.Set("key2", neuron.CacheEntry{})
	cache.Set("key3", neuron.CacheEntry{})

	keys := cache.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}

	// Check that all keys are present
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}

	for _, expected := range []string{"key1", "key2", "key3"} {
		if !keyMap[expected] {
			t.Errorf("expected key %s to be present", expected)
		}
	}
}

func TestMockCache_HitRate(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})

	cache.Get("key1") // hit
	cache.Get("key1") // hit
	cache.Get("key2") // miss

	rate := cache.HitRate()
	expected := 2.0 / 3.0

	if rate < expected-0.01 || rate > expected+0.01 {
		t.Errorf("expected hit rate ~%.2f, got %.2f", expected, rate)
	}
}

func TestMockCache_HitRateEmpty(t *testing.T) {
	cache := mock.NewMockCache(nil)

	if cache.HitRate() != 0 {
		t.Error("expected hit rate 0 with no requests")
	}
}

func TestMockCache_InitialData(t *testing.T) {
	initial := map[string]neuron.CacheEntry{
		"key1": {Data: []byte("value1")},
		"key2": {Data: []byte("value2")},
	}

	cache := mock.NewMockCache(&mock.MockCacheOptions{
		InitialData: initial,
	})

	if cache.Size() != 2 {
		t.Errorf("expected initial size 2, got %d", cache.Size())
	}

	entry, hit := cache.Get("key1")
	if !hit {
		t.Error("expected key1 to exist")
	}
	if string(entry.Data) != "value1" {
		t.Errorf("expected 'value1', got %s", entry.Data)
	}
}

func TestMockCache_Operations_Recording(t *testing.T) {
	cache := mock.NewMockCache(&mock.MockCacheOptions{
		RecordOperations: true,
	})

	cache.Set("key1", neuron.CacheEntry{Data: []byte("test")})
	cache.Get("key1")
	cache.Get("key2")
	cache.Delete("key1")
	cache.Clear()

	ops := cache.Operations()
	if len(ops) != 5 {
		t.Fatalf("expected 5 operations, got %d", len(ops))
	}

	if ops[0].Op != "set" {
		t.Error("expected first op to be set")
	}
	if ops[1].Op != "get" || !ops[1].Hit {
		t.Error("expected second op to be get hit")
	}
	if ops[2].Op != "get" || ops[2].Hit {
		t.Error("expected third op to be get miss")
	}
	if ops[3].Op != "delete" {
		t.Error("expected fourth op to be delete")
	}
	if ops[4].Op != "clear" {
		t.Error("expected fifth op to be clear")
	}
}

func TestMockCache_Operations_Disabled(t *testing.T) {
	cache := mock.NewMockCache(&mock.MockCacheOptions{
		RecordOperations: false,
	})

	cache.Set("key1", neuron.CacheEntry{})
	cache.Get("key1")

	ops := cache.Operations()
	if len(ops) != 0 {
		t.Errorf("expected 0 operations when recording disabled, got %d", len(ops))
	}
}

func TestMockCache_ClearRecorded(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})
	cache.Get("key1")

	if cache.Sets() != 1 || cache.Hits() != 1 {
		t.Error("expected statistics before clear")
	}

	cache.ClearRecorded()

	if cache.Sets() != 0 || cache.Hits() != 0 {
		t.Error("expected statistics to be cleared")
	}

	if len(cache.Operations()) != 0 {
		t.Error("expected operations to be cleared")
	}

	// Data should still be there
	if cache.Size() != 1 {
		t.Error("expected cache data to be preserved")
	}
}

func TestMockCache_Reset(t *testing.T) {
	cache := mock.NewMockCache(nil)

	cache.Set("key1", neuron.CacheEntry{})
	cache.Get("key1")

	cache.Reset()

	if cache.Size() != 0 {
		t.Error("reset should clear cache data")
	}

	if cache.Sets() != 0 || cache.Hits() != 0 {
		t.Error("reset should clear statistics")
	}

	if len(cache.Operations()) != 0 {
		t.Error("reset should clear operations")
	}
}

func TestMockCache_GetEntry(t *testing.T) {
	cache := mock.NewMockCache(nil)

	testEntry := neuron.CacheEntry{Data: []byte("test")}
	cache.Set("key1", testEntry)

	// GetEntry should not affect statistics
	entry, hit := cache.GetEntry("key1")

	if !hit {
		t.Error("expected hit")
	}

	if cache.Hits() != 0 {
		t.Error("GetEntry should not update hit count")
	}

	if string(entry.Data) != "test" {
		t.Errorf("expected 'test', got %s", entry.Data)
	}
}

func TestMockCache_EnableRecording(t *testing.T) {
	cache := mock.NewMockCache(&mock.MockCacheOptions{
		RecordOperations: true,
	})

	cache.Set("key1", neuron.CacheEntry{})

	if len(cache.Operations()) == 0 {
		t.Error("expected operations to be recorded")
	}

	cache.EnableRecording(false)
	cache.Set("key2", neuron.CacheEntry{})

	if len(cache.Operations()) != 1 {
		t.Error("expected recording to be disabled")
	}
}

func TestMockCache_ThreadSafety(t *testing.T) {
	cache := mock.NewMockCache(&mock.MockCacheOptions{
		RecordOperations: true,
	})

	// Concurrent operations
	for i := 0; i < 10; i++ {
		go func(idx int) {
			key := "key" + string(rune(idx))
			cache.Set(key, neuron.CacheEntry{Data: []byte("test")})

			for j := 0; j < 10; j++ {
				cache.Get(key)
				cache.Get("missing")
			}
		}(i)
	}

	time.Sleep(100 * time.Millisecond)

	if cache.Sets() != 10 {
		t.Errorf("expected 10 sets, got %d", cache.Sets())
	}

	if cache.Size() != 10 {
		t.Errorf("expected size 10, got %d", cache.Size())
	}
}

func TestMockCache_MultipleTypes(t *testing.T) {
	cache := mock.NewMockCache(nil)

	entries := []neuron.CacheEntry{
		{Data: []byte("json"), StatusCode: 200},
		{Data: []byte("html"), StatusCode: 200},
		{Data: []byte("error"), StatusCode: 404},
	}

	for i, entry := range entries {
		key := "key" + string('0'+byte(i))
		cache.Set(key, entry)
	}

	if cache.Size() != 3 {
		t.Errorf("expected size 3, got %d", cache.Size())
	}

	entry, hit := cache.Get("key0")
	if !hit {
		t.Fatal("expected key0 to exist")
	}
	if entry == nil {
		t.Fatal("expected non-nil entry")
	}
	if entry.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", entry.StatusCode)
	}
}
