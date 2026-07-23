package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"go.uber.org/zap"
)

// setupTestCache creates a miniredis server and returns a ready-to-use Cache.
func setupTestCache(t *testing.T) (*Cache, *miniredis.Miniredis) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	logger := mustNewLogger(t)

	cfg := Config{
		RedisAddr:     mr.Addr(),
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    5 * time.Minute,
	}

	cache, err := New(cfg, logger)
	if err != nil {
		mr.Close()
		t.Fatalf("failed to create cache: %v", err)
	}

	t.Cleanup(func() {
		cache.Close()
		mr.Close()
	})

	return cache, mr
}

// mustNewLogger creates a no-op logger for testing.
func mustNewLogger(t *testing.T) *zap.Logger {
	t.Helper()
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	return logger
}

// =========================================================================
// Test: Basic Set & Get
// =========================================================================

func TestCache_SetAndGet(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()
	key := "test:key1"
	value := []byte("hello-world")

	err := cache.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, ok := cache.Get(ctx, key)
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if string(got) != string(value) {
		t.Fatalf("got %q, want %q", string(got), string(value))
	}
}

func TestCache_Get_Miss(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()

	got, ok := cache.Get(ctx, "nonexistent:key")
	if ok {
		t.Fatal("expected cache miss, got hit")
	}
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestCache_Get_EmptyValue(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()
	key := "test:emptyvalue"

	err := cache.Set(ctx, key, []byte{}, 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, ok := cache.Get(ctx, key)
	if !ok {
		t.Fatal("expected cache hit for empty value, got miss")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

// =========================================================================
// Test: SetJSON
// =========================================================================

type testData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestCache_SetJSON(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()
	key := "test:json"
	data := testData{Name: "test-user", Value: 42}

	err := cache.SetJSON(ctx, key, data, 0)
	if err != nil {
		t.Fatalf("SetJSON failed: %v", err)
	}

	raw, ok := cache.Get(ctx, key)
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}

	var decoded testData
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}
	if decoded.Name != data.Name || decoded.Value != data.Value {
		t.Fatalf("got %+v, want %+v", decoded, data)
	}
}

func TestCache_SetJSON_NilData(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()

	err := cache.SetJSON(ctx, "test:niljson", nil, 0)
	if err != nil {
		t.Fatalf("SetJSON with nil should not error: %v", err)
	}
}

// =========================================================================
// Test: Invalidate
// =========================================================================

func TestCache_Invalidate(t *testing.T) {
	cache, mr := setupTestCache(t)
	ctx := context.Background()
	key := "test:invalidate"

	// Set value
	err := cache.Set(ctx, key, []byte("data"), 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it's in Redis
	if !mr.Exists(key) {
		t.Fatal("expected key to exist in Redis after Set")
	}

	// Invalidate
	err = cache.Invalidate(ctx, key)
	if err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}

	// Verify it's gone from Redis
	if mr.Exists(key) {
		t.Fatal("expected key to be deleted from Redis after Invalidate")
	}

	// Verify cache miss
	_, ok := cache.Get(ctx, key)
	if ok {
		t.Fatal("expected cache miss after Invalidate")
	}
}

func TestCache_Invalidate_NonExistentKey(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()

	// Invalidating a non-existent key should not error
	err := cache.Invalidate(ctx, "nonexistent:key")
	if err != nil {
		t.Fatalf("Invalidate on non-existent key should not error: %v", err)
	}
}

// =========================================================================
// Test: InvalidatePrefix
// =========================================================================

func TestCache_InvalidatePrefix(t *testing.T) {
	cache, mr := setupTestCache(t)
	ctx := context.Background()

	keys := []string{"prefix:a", "prefix:b", "prefix:c", "other:z"}
	values := [][]byte{[]byte("1"), []byte("2"), []byte("3"), []byte("4")}

	for i, k := range keys {
		if err := cache.Set(ctx, k, values[i], 0); err != nil {
			t.Fatalf("Set %s failed: %v", k, err)
		}
	}

	// Invalidate all keys with "prefix:" prefix
	err := cache.InvalidatePrefix(ctx, "prefix:")
	if err != nil {
		t.Fatalf("InvalidatePrefix failed: %v", err)
	}

	// Check prefix keys are gone
	for _, k := range keys[:3] {
		if mr.Exists(k) {
			t.Fatalf("expected key %s to be deleted", k)
		}
		_, ok := cache.Get(ctx, k)
		if ok {
			t.Fatalf("expected cache miss for %s", k)
		}
	}

	// Check other key still exists
	if !mr.Exists("other:z") {
		t.Fatal("expected 'other:z' to still exist")
	}
	_, ok := cache.Get(ctx, "other:z")
	if !ok {
		t.Fatal("expected 'other:z' to still be cached")
	}
}

func TestCache_InvalidatePrefix_NoMatch(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()

	// Set one key
	err := cache.Set(ctx, "test:abc", []byte("data"), 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Invalidate a prefix that doesn't match anything
	err = cache.InvalidatePrefix(ctx, "zzz:")
	if err != nil {
		t.Fatalf("InvalidatePrefix on non-matching prefix should not error: %v", err)
	}

	// Original key should still exist
	_, ok := cache.Get(ctx, "test:abc")
	if !ok {
		t.Fatal("expected key to still exist")
	}
}

// =========================================================================
// Test: Local Cache Behavior
// =========================================================================

func TestCache_LocalCache_Hit(t *testing.T) {
	// Verify that local cache is populated after a Set and that
	// subsequent Get hits local cache (even if Redis is not queried).
	cache, mr := setupTestCache(t)
	ctx := context.Background()
	key := "test:local"
	value := []byte("local-data")

	err := cache.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Delete key from Redis directly to simulate local-only cache
	mr.Del(key)

	// Get should still work because it's in local cache
	got, ok := cache.Get(ctx, key)
	if !ok {
		t.Fatal("expected local cache hit")
	}
	if string(got) != string(value) {
		t.Fatalf("got %q, want %q", string(got), string(value))
	}
}

func TestCache_LocalCache_Expiry(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()
	key := "test:localexpiry"
	value := []byte("expiring-data")

	// Set with very short TTL (only 1ms) so local cache expires quickly
	err := cache.Set(ctx, key, value, time.Millisecond)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it's in local cache initially
	_, ok := cache.getLocal(key)
	if !ok {
		t.Fatal("expected local cache hit immediately after Set")
	}

	// Wait for TTL to expire
	time.Sleep(5 * time.Millisecond)

	// Local cache should have expired
	_, ok = cache.getLocal(key)
	if ok {
		t.Fatal("expected local cache to expire")
	}
}

// =========================================================================
// Test: Custom TTL
// =========================================================================

func TestCache_Set_WithCustomTTL(t *testing.T) {
	cache, mr := setupTestCache(t)
	ctx := context.Background()
	key := "test:customttl"
	value := []byte("ttl-data")
	shortTTL := 100 * time.Millisecond

	err := cache.Set(ctx, key, value, shortTTL)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify in Redis with TTL
	redisTTL := mr.TTL(key)
	if redisTTL <= 0 {
		t.Fatal("expected Redis key to have TTL > 0")
	}
	if redisTTL > shortTTL {
		t.Fatalf("expected Redis TTL <= %v, got %v", shortTTL, redisTTL)
	}

	// Advance miniredis clock past TTL
	mr.FastForward(shortTTL + 50*time.Millisecond)

	// Key should have expired in Redis
	if mr.Exists(key) {
		t.Fatal("expected key to expire in Redis after TTL")
	}
}

func TestCache_DefaultTTL_Applied(t *testing.T) {
	cache, mr := setupTestCache(t)
	ctx := context.Background()
	key := "test:defaultttl"
	value := []byte("default-ttl")

	// Set with zero TTL (should use default: 5 minutes)
	err := cache.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	redisTTL := mr.TTL(key)
	if redisTTL <= 0 {
		t.Fatal("expected Redis key to have TTL > 0")
	}
	if redisTTL < 4*time.Minute {
		t.Fatalf("expected TTL around 5m (default), got %v", redisTTL)
	}
}

// =========================================================================
// Test: Ping (Health Check)
// =========================================================================

func TestCache_Ping(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()

	err := cache.Ping(ctx)
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}

// =========================================================================
// Test: Pub/Sub Distributed Invalidation
// =========================================================================

func TestCache_PubSub_Invalidation(t *testing.T) {
	// Create two cache instances sharing the same miniredis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewLogger(t)

	cfg := Config{
		RedisAddr:     mr.Addr(),
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    5 * time.Minute,
	}

	cache1, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create cache1: %v", err)
	}
	defer cache1.Close()

	cache2, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create cache2: %v", err)
	}
	defer cache2.Close()

	ctx := context.Background()
	key := "test:pubsub"
	value := []byte("shared-data")

	// Set on cache1 (stored in Redis + local of cache1)
	err = cache1.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Set on cache1 failed: %v", err)
	}

	// Verify cache2 can get it from Redis (not local yet)
	got, ok := cache2.Get(ctx, key)
	if !ok {
		t.Fatal("expected cache2 to get value from Redis")
	}
	if string(got) != string(value) {
		t.Fatalf("cache2 got %q, want %q", string(got), string(value))
	}

	// Now both caches have it in their local cache.
	// Verify by removing from Redis and checking local cached is used.
	mr.Del(key)
	_, ok = cache1.Get(ctx, key)
	if !ok {
		t.Fatal("expected cache1 local cache hit after Redis deletion")
	}
	_, ok = cache2.Get(ctx, key)
	if !ok {
		t.Fatal("expected cache2 local cache hit after Redis deletion")
	}

	// Invalidate from cache2 — triggers Pub/Sub message
	err = cache2.Invalidate(ctx, key)
	if err != nil {
		t.Fatalf("Invalidate from cache2 failed: %v", err)
	}

	// Wait for Pub/Sub message to propagate (retry loop for determinism)
	deadline := time.Now().Add(500 * time.Millisecond)
	for {
		_, ok1 := cache1.Get(ctx, key)
		_, ok2 := cache2.Get(ctx, key)
		if !ok1 && !ok2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for Pub/Sub invalidation")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestCache_PubSub_InvalidatePrefix(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewLogger(t)

	cfg := Config{
		RedisAddr:     mr.Addr(),
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    5 * time.Minute,
	}

	cache1, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create cache1: %v", err)
	}
	defer cache1.Close()

	cache2, err := New(cfg, logger)
	if err != nil {
		t.Fatalf("failed to create cache2: %v", err)
	}
	defer cache2.Close()

	ctx := context.Background()
	keys := []string{"pubsub:a", "pubsub:b", "other:c"}

	for _, k := range keys {
		if err := cache1.Set(ctx, k, []byte(k), 0); err != nil {
			t.Fatalf("Set %s failed: %v", k, err)
		}
		// Populate cache2's local cache too
		cache2.Get(ctx, k)
	}

	// Invalidate prefix "pubsub:" from cache2
	err = cache2.InvalidatePrefix(ctx, "pubsub:")
	if err != nil {
		t.Fatalf("InvalidatePrefix failed: %v", err)
	}

	// Wait for Pub/Sub propagation (retry loop for determinism)
	deadline := time.Now().Add(500 * time.Millisecond)
outer1:
	for {
		for _, k := range keys[:2] {
			_, ok := cache1.Get(ctx, k)
			if ok {
				if time.Now().After(deadline) {
					t.Fatalf("timed out waiting for local cache eviction of %s", k)
				}
				time.Sleep(10 * time.Millisecond)
				continue outer1
			}
		}
		break
	}
	for _, k := range keys[:2] {
		_, ok := cache1.Get(ctx, k)
		if ok {
			t.Fatalf("expected cache1 local cache to evict %s via Pub/Sub", k)
		}
	}

	// other:c should still be cached locally
	_, ok := cache1.Get(ctx, "other:c")
	if !ok {
		t.Fatal("expected 'other:c' to still be in cache1 local cache")
	}
}

// =========================================================================
// Test: Concurrent Access Safety
// =========================================================================

func TestCache_ConcurrentAccess(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()
	var wg sync.WaitGroup
	n := 50

	// Concurrent writes
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("test:concurrent:w:%d", i)
			err := cache.Set(ctx, key, []byte{byte(i)}, 0)
			if err != nil {
				t.Errorf("concurrent Set failed: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("test:concurrent:w:%d", i)
			_, ok := cache.Get(ctx, key)
			if !ok {
				t.Errorf("concurrent Get miss for %s", key)
			}
		}(i)
	}
	wg.Wait()

	// Concurrent invalidates
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("test:concurrent:w:%d", i)
			err := cache.Invalidate(ctx, key)
			if err != nil {
				t.Errorf("concurrent Invalidate failed: %v", err)
			}
		}(i)
	}
	wg.Wait()
}

// =========================================================================
// Test: New with Default Config
// =========================================================================

func TestCache_New_DefaultTTL(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	cfg := Config{
		RedisAddr:     mr.Addr(),
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    0,
	}

	cache, err := New(cfg, mustNewLogger(t))
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer cache.Close()

	if cache.defaultTTL != 5*time.Minute {
		t.Fatalf("expected defaultTTL to be 5m, got %v", cache.defaultTTL)
	}
}

// =========================================================================
// Test: Close
// =========================================================================

func TestCache_Close(t *testing.T) {
	cache, _ := setupTestCache(t)

	// Close should not panic
	err := cache.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

// =========================================================================
// Test: New with Redis Connection Failure (Error Path)
// =========================================================================

// TestCache_New_InvalidAddress verifies that New() returns an error when
// the Redis server is unreachable.
func TestCache_New_InvalidAddress(t *testing.T) {
	logger := mustNewLogger(t)

	cfg := Config{
		RedisAddr:     "127.0.0.1:1", // Port 1 — unlikely to have Redis
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    5 * time.Minute,
	}

	_, err := New(cfg, logger)
	if err == nil {
		t.Fatal("expected error when Redis is unreachable, got nil")
	}

	// Verify the error message is descriptive (contains the address)
	if !strings.Contains(err.Error(), "127.0.0.1:1") {
		t.Fatalf("error should contain the Redis address, got: %v", err)
	}
	if !strings.Contains(err.Error(), "failed to connect to Redis") {
		t.Fatalf("error should mention 'failed to connect to Redis', got: %v", err)
	}
}

// TestCache_New_EmptyAddress verifies that New() returns an error when
// Redis address is empty.
func TestCache_New_EmptyAddress(t *testing.T) {
	logger := mustNewLogger(t)

	cfg := Config{
		RedisAddr:     "127.0.0.1:0", // Port 0 is invalid for TCP
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    5 * time.Minute,
	}

	_, err := New(cfg, logger)
	if err == nil {
		t.Fatal("expected error when Redis address has invalid port, got nil")
	}
}

// TestCache_New_ClosedRedis verifies that New() returns an error when
// the provided Redis client cannot be pinged (simulated via stopped server).
func TestCache_New_ClosedRedis(t *testing.T) {
	// Start and immediately stop a miniredis to get a port that was
	// briefly active but is now closed.
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	addr := mr.Addr()
	mr.Close() // Stop the server — port is now closed

	logger := mustNewLogger(t)

	cfg := Config{
		RedisAddr:     addr,
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    5 * time.Minute,
	}

	// New() should fail because the server is already closed
	_, err = New(cfg, logger)
	if err == nil {
		t.Fatal("expected error when connecting to closed Redis server, got nil")
	}
}

// TestCache_New_ErrorMessage_Format verifies the error message format
// contains all expected context for debugging.
func TestCache_New_ErrorMessage_Format(t *testing.T) {
	logger := mustNewLogger(t)

	cfg := Config{
		RedisAddr:     "192.0.2.1:6379", // TEST-NET address (RFC 5737)
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    5 * time.Minute,
	}

	_, err := New(cfg, logger)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	errMsg := err.Error()
	// Should contain the address for debugging
	if !strings.Contains(errMsg, "192.0.2.1:6379") {
		t.Fatalf("error should contain Redis address, got: %s", errMsg)
	}
	t.Logf("connection error message: %s", errMsg)
}



// =========================================================================
// Test: EvictLocal (called by Pub/Sub)
// =========================================================================

func TestCache_evictLocal(t *testing.T) {
	cache, _ := setupTestCache(t)
	ctx := context.Background()

	// Populate local cache via Set
	err := cache.Set(ctx, "test:evict:a", []byte("a"), 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	err = cache.Set(ctx, "test:evict:b", []byte("b"), 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify both in local cache
	_, ok := cache.getLocal("test:evict:a")
	if !ok {
		t.Fatal("expected local cache hit for a")
	}
	_, ok = cache.getLocal("test:evict:b")
	if !ok {
		t.Fatal("expected local cache hit for b")
	}

	// evictLocal should delete both from local cache
	cache.evictLocal("test:evict:a", "test:evict:b")

	_, ok = cache.getLocal("test:evict:a")
	if ok {
		t.Fatal("expected a to be evicted from local cache")
	}
	_, ok = cache.getLocal("test:evict:b")
	if ok {
		t.Fatal("expected b to be evicted from local cache")
	}
}
