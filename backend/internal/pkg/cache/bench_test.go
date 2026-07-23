package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// =========================================================================
// Benchmark Helpers
// =========================================================================

// benchLogger creates a no-op zap logger for benchmarks.
func benchLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

// benchSetupCache creates a miniredis + Cache for benchmarking.
// Returns the cache, miniredis, and a cleanup function.
func benchSetupCache(b *testing.B) (*Cache, *miniredis.Miniredis, func()) {
	b.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		b.Fatalf("failed to start miniredis: %v", err)
	}

	logger := benchLogger()

	cfg := Config{
		RedisAddr:     mr.Addr(),
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    5 * time.Minute,
	}

	cache, err := New(cfg, logger)
	if err != nil {
		mr.Close()
		b.Fatalf("failed to create cache: %v", err)
	}

	cleanup := func() {
		cache.Close()
		mr.Close()
	}

	return cache, mr, cleanup
}

// benchSetupPubSub creates a miniredis + PubSub for benchmarking.
// Returns the pubsub, client, miniredis, and a cleanup function.
func benchSetupPubSub(b *testing.B) (*PubSub, *redis.Client, *miniredis.Miniredis, func()) {
	b.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		b.Fatalf("failed to start miniredis: %v", err)
	}

	logger := benchLogger()
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	callback := func(keys ...string) {}

	ps, err := NewPubSub(client, logger, callback)
	if err != nil {
		client.Close()
		mr.Close()
		b.Fatalf("failed to create PubSub: %v", err)
	}

	cleanup := func() {
		ps.Close()
		client.Close()
		mr.Close()
	}

	return ps, client, mr, cleanup
}

// =========================================================================
// Cache Benchmarks: Set
// =========================================================================

// BenchmarkCache_Set_SmallData benchmarks Set with a small value (100 bytes).
func BenchmarkCache_Set_SmallData(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()
	value := make([]byte, 100)
	for i := range value {
		value[i] = byte(i)
	}

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:set:small:%d", i)
		if err := cache.Set(ctx, key, value, 5*time.Minute); err != nil {
			b.Fatalf("Set failed: %v", err)
		}
	}
}

// BenchmarkCache_Set_LargeData benchmarks Set with a large value (10KB).
func BenchmarkCache_Set_LargeData(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()
	value := make([]byte, 10_240)
	for i := range value {
		value[i] = byte(i % 256)
	}

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:set:large:%d", i)
		if err := cache.Set(ctx, key, value, 5*time.Minute); err != nil {
			b.Fatalf("Set failed: %v", err)
		}
	}
}

// =========================================================================
// Cache Benchmarks: Get
// =========================================================================

// BenchmarkCache_Get_Hot benchmarks Get when data is already in local cache.
func BenchmarkCache_Get_Hot(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()
	key := "bench:get:hot"
	value := []byte("hot-cache-value")

	if err := cache.Set(ctx, key, value, 5*time.Minute); err != nil {
		b.Fatalf("Set failed: %v", err)
	}

	// First Get to ensure it's in local cache
	cache.Get(ctx, key)

	b.ResetTimer()
	for range b.N {
		_, ok := cache.Get(ctx, key)
		if !ok {
			b.Fatal("expected cache hit")
		}
	}
}

// BenchmarkCache_Get_Cold benchmarks Get when data must be fetched from Redis
// (not in local cache). Uses a fixed key pool with cycling to handle large b.N.
func BenchmarkCache_Get_Cold(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()

	const poolSize = 1000
	// Pre-populate Redis with a fixed pool of keys
	for i := 0; i < poolSize; i++ {
		key := fmt.Sprintf("bench:get:cold:%d", i)
		if err := cache.Set(ctx, key, []byte("cold-value"), 5*time.Minute); err != nil {
			b.Fatalf("pre-Set failed: %v", err)
		}
		cache.local.Delete(key)
	}

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:get:cold:%d", i%poolSize)
		cache.local.Delete(key) // force Redis read each time
		_, ok := cache.Get(ctx, key)
		if !ok {
			b.Fatalf("expected cache hit for %s", key)
		}
	}
}

// BenchmarkCache_Get_Miss benchmarks Get on non-existent keys.
func BenchmarkCache_Get_Miss(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:get:miss:%d", i)
		_, ok := cache.Get(ctx, key)
		if ok {
			b.Fatal("expected cache miss")
		}
	}
}

// =========================================================================
// Cache Benchmarks: SetJSON
// =========================================================================

// benchmarkPayload is a moderately sized struct for JSON benchmarks.
type benchmarkPayload struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Active    bool              `json:"active"`
	Tags      []string          `json:"tags"`
	Metadata  map[string]string `json:"metadata"`
}

// BenchmarkCache_SetJSON benchmarks JSON serialization + cache set.
func BenchmarkCache_SetJSON(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()
	payload := benchmarkPayload{
		ID:     42,
		Name:   "Benchmark User",
		Email:  "benchmark@example.com",
		Active: true,
		Tags:   []string{"go", "redis", "cache", "benchmark"},
		Metadata: map[string]string{
			"department": "engineering",
			"level":      "senior",
		},
	}

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:setjson:%d", i)
		if err := cache.SetJSON(ctx, key, payload, 5*time.Minute); err != nil {
			b.Fatalf("SetJSON failed: %v", err)
		}
	}
}

// =========================================================================
// Cache Benchmarks: Invalidate
// =========================================================================

// BenchmarkCache_Invalidate benchmarks single-key invalidation.
// Uses a fixed key pool with cycling to handle large b.N.
func BenchmarkCache_Invalidate(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()

	const poolSize = 1000
	// Pre-populate Redis with a fixed pool of keys
	for i := 0; i < poolSize; i++ {
		key := fmt.Sprintf("bench:invalidate:%d", i)
		if err := cache.Set(ctx, key, []byte("data"), 5*time.Minute); err != nil {
			b.Fatalf("pre-Set failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:invalidate:%d", i%poolSize)
		if err := cache.Invalidate(ctx, key); err != nil {
			b.Fatalf("Invalidate failed: %v", err)
		}
	}
}

// BenchmarkCache_InvalidatePrefix benchmarks prefix-based invalidation.
// Uses StopTimer/StartTimer to exclude insert phase from measurement.
func BenchmarkCache_InvalidatePrefix(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()
	const batchSize = 10

	b.ResetTimer()
	for i := range b.N {
		prefix := fmt.Sprintf("bench:prefix:%d:", i)

		// Insert batch — excluded from timing
		b.StopTimer()
		for j := 0; j < batchSize; j++ {
			key := fmt.Sprintf("%s%d", prefix, j)
			if err := cache.Set(ctx, key, []byte("data"), 5*time.Minute); err != nil {
				b.Fatalf("Set failed: %v", err)
			}
		}
		b.StartTimer()

		// Invalidate prefix — measured
		if err := cache.InvalidatePrefix(ctx, prefix); err != nil {
			b.Fatalf("InvalidatePrefix failed: %v", err)
		}
	}
}

// =========================================================================
// Cache Benchmarks: Concurrent Operations
// =========================================================================

// BenchmarkCache_ConcurrentSet benchmarks concurrent Set operations.
func BenchmarkCache_ConcurrentSet(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()
	value := []byte("concurrent-value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			key := fmt.Sprintf("bench:concurrent:set:%d", i)
			if err := cache.Set(ctx, key, value, 5*time.Minute); err != nil {
				b.Errorf("Set failed: %v", err)
			}
			i++
		}
	})
}

// BenchmarkCache_ConcurrentGetHot benchmarks concurrent Get on hot (local) cache.
func BenchmarkCache_ConcurrentGetHot(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()

	// Pre-populate shared keys
	const numKeys = 100
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("bench:concurrent:get:hot:%d", i)
		if err := cache.Set(ctx, key, []byte("hot-value"), 5*time.Minute); err != nil {
			b.Fatalf("pre-Set failed: %v", err)
		}
		// Warm up local cache
		cache.Get(ctx, key)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			key := fmt.Sprintf("bench:concurrent:get:hot:%d", i%numKeys)
			_, ok := cache.Get(ctx, key)
			if !ok {
				b.Error("expected cache hit")
			}
			i++
		}
	})
}

// BenchmarkCache_ConcurrentGetCold benchmarks concurrent Get on cold cache.
func BenchmarkCache_ConcurrentGetCold(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()

	// Pre-populate Redis with many keys (clear local cache after)
	const numKeys = 1000
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("bench:concurrent:get:cold:%d", i)
		if err := cache.Set(ctx, key, []byte("cold-value"), 5*time.Minute); err != nil {
			b.Fatalf("pre-Set failed: %v", err)
		}
		cache.local.Delete(key) // ensure cold start
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			key := fmt.Sprintf("bench:concurrent:get:cold:%d", i%numKeys)
			_, ok := cache.Get(ctx, key)
			if !ok {
				b.Error("expected cache hit")
			}
			i++
		}
	})
}

// BenchmarkCache_ConcurrentMixed benchmarks a mix of Set/Get/Invalidate.
func BenchmarkCache_ConcurrentMixed(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()
	value := []byte("mixed-value")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			key := fmt.Sprintf("bench:mixed:%d", i)
			switch i % 3 {
			case 0:
				_ = cache.Set(ctx, key, value, 5*time.Minute)
			case 1:
				cache.Get(ctx, key)
			case 2:
				_ = cache.Invalidate(ctx, key)
			}
			i++
		}
	})
}

// =========================================================================
// Cache Benchmarks: Local vs Redis Latency Comparison
// =========================================================================

// BenchmarkCache_LocalOnly benchmarks local cache hits (bypasses Redis).
func BenchmarkCache_LocalOnly(b *testing.B) {
	cache, mr, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()
	key := "bench:localonly"
	value := []byte("local-only-value")

	if err := cache.Set(ctx, key, value, 5*time.Minute); err != nil {
		b.Fatalf("Set failed: %v", err)
	}

	// Ensure it's in local cache
	cache.Get(ctx, key)

	// Delete from Redis to force local-only access
	mr.Del(key)

	b.ResetTimer()
	for range b.N {
		_, ok := cache.Get(ctx, key)
		if !ok {
			b.Fatal("expected local cache hit")
		}
	}
}

// BenchmarkCache_RedisOnly benchmarks Redis-only access (no local cache).
func BenchmarkCache_RedisOnly(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()

	// Pre-populate keys and clear local cache each time
	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:redisonly:%d", i)
		if err := cache.Set(ctx, key, []byte("redis-value"), 5*time.Minute); err != nil {
			b.Fatalf("Set failed: %v", err)
		}
		// Force Redis read by clearing local cache
		cache.local.Delete(key)
		_, ok := cache.Get(ctx, key)
		if !ok {
			b.Fatalf("expected Redis hit for %s", key)
		}
	}
}

// =========================================================================
// PubSub Benchmarks
// =========================================================================

// BenchmarkPubSub_PublishInvalidate_SingleKey benchmarks publishing a single key.
func BenchmarkPubSub_PublishInvalidate_SingleKey(b *testing.B) {
	ps, _, _, cleanup := benchSetupPubSub(b)
	defer cleanup()

	ctx := context.Background()

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:pubsub:single:%d", i)
		if err := ps.PublishInvalidate(ctx, key); err != nil {
			b.Fatalf("PublishInvalidate failed: %v", err)
		}
	}
}

// BenchmarkPubSub_PublishInvalidate_MultipleKeys benchmarks publishing 10 keys at once.
func BenchmarkPubSub_PublishInvalidate_MultipleKeys(b *testing.B) {
	ps, _, _, cleanup := benchSetupPubSub(b)
	defer cleanup()

	ctx := context.Background()

	b.ResetTimer()
	for i := range b.N {
		keys := make([]string, 10)
		for j := 0; j < 10; j++ {
			keys[j] = fmt.Sprintf("bench:pubsub:multi:%d:%d", i, j)
		}
		if err := ps.PublishInvalidate(ctx, keys...); err != nil {
			b.Fatalf("PublishInvalidate failed: %v", err)
		}
	}
}

// BenchmarkPubSub_PublishInvalidate_Empty benchmarks publishing empty keys.
func BenchmarkPubSub_PublishInvalidate_Empty(b *testing.B) {
	ps, _, _, cleanup := benchSetupPubSub(b)
	defer cleanup()

	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		if err := ps.PublishInvalidate(ctx); err != nil {
			b.Fatalf("PublishInvalidate empty failed: %v", err)
		}
	}
}

// BenchmarkPubSub_HandleMessage_Valid benchmarks direct handleMessage call with valid JSON.
func BenchmarkPubSub_HandleMessage_Valid(b *testing.B) {
	ps, _, _, cleanup := benchSetupPubSub(b)
	defer cleanup()

	payload := `{"keys":["bench:key1","bench:key2","bench:key3"]}`

	b.ResetTimer()
	for range b.N {
		ps.handleMessage(payload)
	}
}

// BenchmarkPubSub_HandleMessage_Invalid benchmarks direct handleMessage call with invalid JSON.
func BenchmarkPubSub_HandleMessage_Invalid(b *testing.B) {
	ps, _, _, cleanup := benchSetupPubSub(b)
	defer cleanup()

	payload := "{invalid json}"

	b.ResetTimer()
	for range b.N {
		ps.handleMessage(payload)
	}
}

// BenchmarkPubSub_HandleMessage_LargeKeys benchmarks handleMessage with many keys.
func BenchmarkPubSub_HandleMessage_LargeKeys(b *testing.B) {
	ps, _, _, cleanup := benchSetupPubSub(b)
	defer cleanup()

	// Build payload with 100 keys
	keys := make([]string, 100)
	for i := range keys {
		keys[i] = fmt.Sprintf("bench:large:key:%d", i)
	}
	msg := invalidationMessage{Keys: keys}
	data, _ := json.Marshal(msg)
	payload := string(data)

	b.ResetTimer()
	for range b.N {
		ps.handleMessage(payload)
	}
}

// =========================================================================
// PubSub Benchmarks: Concurrent
// =========================================================================

// BenchmarkPubSub_ConcurrentPublish benchmarks concurrent publish operations.
func BenchmarkPubSub_ConcurrentPublish(b *testing.B) {
	ps, _, _, cleanup := benchSetupPubSub(b)
	defer cleanup()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			key := fmt.Sprintf("bench:concurrent:pubsub:%d", i)
			if err := ps.PublishInvalidate(ctx, key); err != nil {
				b.Errorf("PublishInvalidate failed: %v", err)
			}
			i++
		}
	})
}

// =========================================================================
// Data Size Comparison Benchmark
// =========================================================================

// BenchmarkCache_Set_DataSize compares Set performance across data sizes.
func BenchmarkCache_Set_DataSize(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"64B", 64},
		{"1KB", 1_024},
		{"10KB", 10_240},
		{"100KB", 102_400},
		{"1MB", 1_048_576},
	}

	for _, s := range sizes {
		b.Run(s.name, func(b *testing.B) {
			cache, _, cleanup := benchSetupCache(b)
			defer cleanup()

			ctx := context.Background()
			value := make([]byte, s.size)
			for i := range value {
				value[i] = byte(i % 256)
			}

			b.ResetTimer()
			for i := range b.N {
				key := fmt.Sprintf("bench:datasize:%s:%d", s.name, i)
				if err := cache.Set(ctx, key, value, 5*time.Minute); err != nil {
					b.Fatalf("Set failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkCache_Get_DataSize compares Get performance across data sizes.
// Each Get is a cold Redis read to measure network transfer cost.
func BenchmarkCache_Get_DataSize(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"64B", 64},
		{"1KB", 1_024},
		{"10KB", 10_240},
		{"100KB", 102_400},
	}

	for _, s := range sizes {
		b.Run(s.name, func(b *testing.B) {
			cache, _, cleanup := benchSetupCache(b)
			defer cleanup()

			ctx := context.Background()
			value := make([]byte, s.size)
			for i := range value {
				value[i] = byte(i % 256)
			}

			// Pre-populate a fixed pool of keys
			const poolSize = 500
			for i := 0; i < poolSize; i++ {
				key := fmt.Sprintf("bench:getsize:%s:%d", s.name, i)
				if err := cache.Set(ctx, key, value, 5*time.Minute); err != nil {
					b.Fatalf("pre-Set failed: %v", err)
				}
				cache.local.Delete(key)
			}

			b.ResetTimer()
			for i := range b.N {
				key := fmt.Sprintf("bench:getsize:%s:%d", s.name, i%poolSize)
				cache.local.Delete(key) // force cold read
				_, ok := cache.Get(ctx, key)
				if !ok {
					b.Fatalf("expected cache hit for %s", key)
				}
			}
		})
	}
}

// =========================================================================
// Cache Throughput Benchmark (ops/sec)
// =========================================================================

// BenchmarkCache_Throughput measures sequential Set + Get throughput.
func BenchmarkCache_Throughput(b *testing.B) {
	cache, _, cleanup := benchSetupCache(b)
	defer cleanup()

	ctx := context.Background()
	value := []byte("throughput-test-value")

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:throughput:%d", i)
		if err := cache.Set(ctx, key, value, 5*time.Minute); err != nil {
			b.Fatalf("Set failed: %v", err)
		}
		_, ok := cache.Get(ctx, key)
		if !ok {
			b.Fatal("expected cache hit")
		}
		_ = cache.Invalidate(ctx, key)
	}
}

// =========================================================================
// PubSub Throughput Benchmark
// =========================================================================

// BenchmarkPubSub_Throughput measures publish + handle throughput.
func BenchmarkPubSub_Throughput(b *testing.B) {
	ps, _, _, cleanup := benchSetupPubSub(b)
	defer cleanup()

	ctx := context.Background()
	consumerCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Start a consumer goroutine that reads from Redis directly
	// and calls handleMessage to simulate the subscriber path.
	// Uses a context timeout to prevent goroutine leaks.
	wg.Add(1)
	go func() {
		defer wg.Done()
		pubsub := ps.client.Subscribe(ctx, ChannelInvalidate)
		defer pubsub.Close()

		ch := pubsub.Channel()
		count := 0
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				ps.handleMessage(msg.Payload)
				count++
				if count >= b.N {
					return
				}
			case <-consumerCtx.Done():
				return
			}
		}
	}()

	// Allow subscriber to initialize
	time.Sleep(50 * time.Millisecond)

	b.ResetTimer()
	for i := range b.N {
		key := fmt.Sprintf("bench:pubsub:throughput:%d", i)
		if err := ps.PublishInvalidate(ctx, key); err != nil {
			b.Fatalf("PublishInvalidate failed: %v", err)
		}
	}

	wg.Wait()
}
