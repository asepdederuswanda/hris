package cache

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

// =========================================================================
// Integration Helpers
// =========================================================================



// createCacheForIntegration creates a fully initialized Cache backed by miniredis.
func createCacheForIntegration(t *testing.T) (*Cache, *miniredis.Miniredis) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	cfg := Config{
		RedisAddr:     mr.Addr(),
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    30 * time.Minute,
	}

	cache, err := New(cfg, mustNewLogger(t))
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

// =========================================================================
// Integration Test: Full Lifecycle
// =========================================================================

// TestCacheIntegration_FullLifecycle exercises the complete cache lifecycle:
// New → Set (raw) → Get → SetJSON → Get (JSON) → Invalidate → Get (miss) → Ping → Close.
func TestCacheIntegration_FullLifecycle(t *testing.T) {
	cache, mr := createCacheForIntegration(t)
	ctx := context.Background()

	// ---------- Phase 1: Set raw bytes ----------
	t.Log("Phase 1: Set raw bytes")
	err := cache.Set(ctx, "int:key1", []byte("raw-value"), 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify in Redis
	if !mr.Exists("int:key1") {
		t.Fatal("key should exist in Redis after Set")
	}

	// Verify TTL was applied (default)
	ttl := mr.TTL("int:key1")
	if ttl <= 0 {
		t.Fatal("key should have a TTL in Redis")
	}

	// ---------- Phase 2: Get raw bytes ----------
	t.Log("Phase 2: Get raw bytes")
	val, ok := cache.Get(ctx, "int:key1")
	if !ok {
		t.Fatal("Get should return cache hit")
	}
	if string(val) != "raw-value" {
		t.Fatalf("Get returned %q, want %q", string(val), "raw-value")
	}

	// ---------- Phase 3: Local cache hit ----------
	t.Log("Phase 3: Verify local cache hit")
	// Delete from Redis but it should still be in local cache
	mr.Del("int:key1")
	val, ok = cache.Get(ctx, "int:key1")
	if !ok {
		t.Fatal("Get should hit local cache after Redis deletion")
	}
	if string(val) != "raw-value" {
		t.Fatalf("local cache returned %q, want %q", string(val), "raw-value")
	}

	// ---------- Phase 4: Set JSON ----------
	t.Log("Phase 4: Set JSON structure")
	type profile struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}
	user := profile{UserID: 101, Username: "john.doe", Role: "admin"}

	err = cache.SetJSON(ctx, "int:user:101", user, 0)
	if err != nil {
		t.Fatalf("SetJSON failed: %v", err)
	}

	// ---------- Phase 5: Get and verify JSON ----------
	t.Log("Phase 5: Get and verify JSON")
	raw, ok := cache.Get(ctx, "int:user:101")
	if !ok {
		t.Fatal("Get should return cache hit for JSON key")
	}

	var decoded profile
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if decoded.UserID != 101 || decoded.Username != "john.doe" || decoded.Role != "admin" {
		t.Fatalf("decoded profile mismatch: got %+v", decoded)
	}

	// ---------- Phase 6: Invalidate single key ----------
	t.Log("Phase 6: Invalidate single key")
	err = cache.Invalidate(ctx, "int:key1")
	if err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}

	// Verify gone from Redis
	if mr.Exists("int:key1") {
		t.Fatal("key should be deleted from Redis after Invalidate")
	}

	// Verify cache miss
	_, ok = cache.Get(ctx, "int:key1")
	if ok {
		t.Fatal("Get should return cache miss after Invalidate")
	}

	// ---------- Phase 7: Invalidate non-existent key (should not error) ----------
	t.Log("Phase 7: Invalidate non-existent key")
	err = cache.Invalidate(ctx, "int:nonexistent")
	if err != nil {
		t.Fatalf("Invalidate on non-existent key should not error: %v", err)
	}

	// ---------- Phase 8: InvalidatePrefix ----------
	t.Log("Phase 8: InvalidatePrefix")
	// Set keys with common prefix
	keys := []string{"int:prefix:a", "int:prefix:b", "int:prefix:c"}
	for _, k := range keys {
		if err := cache.Set(ctx, k, []byte(k), 0); err != nil {
			t.Fatalf("Set %s failed: %v", k, err)
		}
	}

	// Verify all exist
	for _, k := range keys {
		if !mr.Exists(k) {
			t.Fatalf("key %s should exist in Redis", k)
		}
	}

	// Invalidate prefix
	err = cache.InvalidatePrefix(ctx, "int:prefix:")
	if err != nil {
		t.Fatalf("InvalidatePrefix failed: %v", err)
	}

	// Verify all gone
	for _, k := range keys {
		if mr.Exists(k) {
			t.Fatalf("key %s should be deleted by InvalidatePrefix", k)
		}
		_, ok := cache.Get(ctx, k)
		if ok {
			t.Fatalf("key %s should be cache miss after InvalidatePrefix", k)
		}
	}

	// ---------- Phase 9: Ping (health check) ----------
	t.Log("Phase 9: Ping health check")
	err = cache.Ping(ctx)
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	t.Log("Integration test: Full lifecycle PASSED")
}

// =========================================================================
// Integration Test: Two-Instance Pub/Sub Invalidation
// =========================================================================

// TestCacheIntegration_TwoInstanceInvalidation verifies that invalidating
// a key on one Cache instance propagates to another instance via Pub/Sub.
func TestCacheIntegration_TwoInstanceInvalidation(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	cfg := Config{
		RedisAddr:     mr.Addr(),
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    30 * time.Minute,
	}

	// Create two independent cache instances on the same Redis
	cache1, err := New(cfg, mustNewLogger(t))
	if err != nil {
		t.Fatalf("failed to create cache1: %v", err)
	}
	defer cache1.Close()

	cache2, err := New(cfg, mustNewLogger(t))
	if err != nil {
		t.Fatalf("failed to create cache2: %v", err)
	}
	defer cache2.Close()

	ctx := context.Background()
	key := "int:shared:key"
	value := []byte("shared-value")

	// Phase 1: Set on cache1, both caches can read from Redis
	t.Log("Phase 1: Set on cache1, read from both")
	err = cache1.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Set on cache1 failed: %v", err)
	}

	// cache2 reads from Redis (cold)
	v2, ok := cache2.Get(ctx, key)
	if !ok {
		t.Fatal("cache2 should read value from Redis")
	}
	if string(v2) != string(value) {
		t.Fatalf("cache2 got %q, want %q", string(v2), string(value))
	}

	// Phase 2: Both caches now have local copies. Remove from Redis
	// to verify local cache works.
	t.Log("Phase 2: Both caches have local copies")
	mr.Del(key)
	v1, ok := cache1.Get(ctx, key)
	if !ok {
		t.Fatal("cache1 should hit local cache after Redis deletion")
	}
	if string(v1) != string(value) {
		t.Fatalf("cache1 local cache: got %q, want %q", string(v1), string(value))
	}
	v2, ok = cache2.Get(ctx, key)
	if !ok {
		t.Fatal("cache2 should hit local cache after Redis deletion")
	}
	if string(v2) != string(value) {
		t.Fatalf("cache2 local cache: got %q, want %q", string(v2), string(value))
	}

	// Phase 3: Invalidate from cache2 — should propagate to cache1 via Pub/Sub
	t.Log("Phase 3: Invalidate from cache2, verify both caches evicted")
	err = cache2.Invalidate(ctx, key)
	if err != nil {
		t.Fatalf("Invalidate from cache2 failed: %v", err)
	}

	// Wait for Pub/Sub propagation
	deadline := time.Now().Add(2 * time.Second)
	for {
		_, ok1 := cache1.Get(ctx, key)
		_, ok2 := cache2.Get(ctx, key)
		if !ok1 && !ok2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for Pub/Sub invalidation to propagate")
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Log("Integration test: Two-instance invalidation PASSED")
}

// =========================================================================
// Integration Test: TTL Expiry
// =========================================================================

// TestCacheIntegration_TTLExpiry verifies that keys expire correctly
// both in Redis (via miniredis FastForward) and in local cache.
func TestCacheIntegration_TTLExpiry(t *testing.T) {
	cache, mr := createCacheForIntegration(t)
	ctx := context.Background()

	shortTTL := 100 * time.Millisecond

	// Set with short TTL
	err := cache.Set(ctx, "int:ttl:key", []byte("ttl-value"), shortTTL)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Immediate read should succeed
	val, ok := cache.Get(ctx, "int:ttl:key")
	if !ok {
		t.Fatal("immediate Get should succeed")
	}
	if string(val) != "ttl-value" {
		t.Fatalf("got %q, want %q", string(val), "ttl-value")
	}

	// Advance miniredis clock past TTL
	mr.FastForward(shortTTL + 50*time.Millisecond)

	// Redis key should be expired
	if mr.Exists("int:ttl:key") {
		t.Fatal("Redis key should have expired")
	}

	// Wait for real wall clock to advance past local cache expiry
	time.Sleep(shortTTL + 50*time.Millisecond)

	// Now both Redis AND local cache should be expired
	_, ok = cache.Get(ctx, "int:ttl:key")
	if ok {
		t.Fatal("expected cache miss after TTL expiry (local + Redis)")
	}
}

// =========================================================================
// Integration Test: Concurrent Workflow
// =========================================================================

// TestCacheIntegration_ConcurrentWorkflow runs a mixed workload
// with multiple goroutines performing Set/Get/Invalidate operations.
func TestCacheIntegration_ConcurrentWorkflow(t *testing.T) {
	cache, _ := createCacheForIntegration(t)
	ctx := context.Background()

	var wg sync.WaitGroup
	n := 20

	// Concurrent writes
	t.Log("Spawning concurrent workers...")
	for i := range n {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			key := "int:concurrent:key:" + string(rune('a'+id%26))
			rawKey := "int:concurrent:raw:" + string(rune('a'+id%26))

			// Set
			if err := cache.Set(ctx, key, []byte{byte(id)}, 0); err != nil {
				t.Errorf("worker %d: Set failed: %v", id, err)
				return
			}

			// Get (should hit)
			val, ok := cache.Get(ctx, key)
			if !ok {
				t.Errorf("worker %d: Get miss after Set", id)
				return
			}
			if len(val) != 1 || val[0] != byte(id) {
				t.Errorf("worker %d: Get returned wrong value: %v", id, val)
				return
			}

			// SetJSON
			type data struct {
				ID int `json:"id"`
			}
			if err := cache.SetJSON(ctx, rawKey, data{ID: id}, 0); err != nil {
				t.Errorf("worker %d: SetJSON failed: %v", id, err)
				return
			}

			// Get JSON
			raw, ok := cache.Get(ctx, rawKey)
			if !ok {
				t.Errorf("worker %d: Get miss for JSON key", id)
				return
			}
			var decoded data
			if err := json.Unmarshal(raw, &decoded); err != nil {
				t.Errorf("worker %d: JSON unmarshal failed: %v", id, err)
				return
			}
			if decoded.ID != id {
				t.Errorf("worker %d: JSON decoded ID=%d, want %d", id, decoded.ID, id)
				return
			}

			// Invalidate
			if err := cache.Invalidate(ctx, key); err != nil {
				t.Errorf("worker %d: Invalidate failed: %v", id, err)
				return
			}

			// Verify invalidated
			_, ok = cache.Get(ctx, key)
			if ok {
				t.Errorf("worker %d: cache hit after Invalidate", id)
			}
		}(i)
	}
	wg.Wait()
	t.Log("Integration test: Concurrent workflow PASSED")
}

// =========================================================================
// Integration Test: Cache Miss Scenarios
// =========================================================================

// TestCacheIntegration_MissScenarios verifies various cache miss conditions.
func TestCacheIntegration_MissScenarios(t *testing.T) {
	cache, _ := createCacheForIntegration(t)
	ctx := context.Background()

	// Scenario 1: Non-existent key
	t.Log("Scenario 1: Non-existent key")
	val, ok := cache.Get(ctx, "int:miss:never-set")
	if ok {
		t.Fatal("expected cache miss for non-existent key")
	}
	if val != nil {
		t.Fatalf("expected nil value, got %v", val)
	}

	// Scenario 2: Key set then invalidated
	t.Log("Scenario 2: Key set then invalidated")
	err := cache.Set(ctx, "int:miss:invalidated", []byte("data"), 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	err = cache.Invalidate(ctx, "int:miss:invalidated")
	if err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}
	_, ok = cache.Get(ctx, "int:miss:invalidated")
	if ok {
		t.Fatal("expected cache miss for invalidated key")
	}

	// Scenario 3: Empty value (valid key, but empty data)
	t.Log("Scenario 3: Empty value")
	err = cache.Set(ctx, "int:miss:empty", []byte{}, 0)
	if err != nil {
		t.Fatalf("Set empty value failed: %v", err)
	}
	val, ok = cache.Get(ctx, "int:miss:empty")
	if !ok {
		t.Fatal("expected cache hit for empty value")
	}
	if len(val) != 0 {
		t.Fatalf("expected empty slice, got %v", val)
	}

	// Scenario 4: Multiple invalidations of same key (idempotent)
	t.Log("Scenario 4: Idempotent invalidations")
	err = cache.Set(ctx, "int:miss:idempotent", []byte("data"), 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := cache.Invalidate(ctx, "int:miss:idempotent"); err != nil {
			t.Fatalf("Invalidate iteration %d failed: %v", i, err)
		}
	}
	_, ok = cache.Get(ctx, "int:miss:idempotent")
	if ok {
		t.Fatal("expected cache miss after multiple invalidations")
	}
}

// =========================================================================
// Integration Test: Large Payload
// =========================================================================

// TestCacheIntegration_LargePayload verifies cache handles large data correctly.
func TestCacheIntegration_LargePayload(t *testing.T) {
	cache, _ := createCacheForIntegration(t)
	ctx := context.Background()

	// Create a 50KB payload
	payload := make([]byte, 50*1024)
	for i := range payload {
		payload[i] = byte(i % 256)
	}

	err := cache.Set(ctx, "int:large:50kb", payload, 0)
	if err != nil {
		t.Fatalf("Set large payload failed: %v", err)
	}

	val, ok := cache.Get(ctx, "int:large:50kb")
	if !ok {
		t.Fatal("expected cache hit for large payload")
	}
	if len(val) != len(payload) {
		t.Fatalf("payload size mismatch: got %d, want %d", len(val), len(payload))
	}
	for i := range payload {
		if val[i] != payload[i] {
			t.Fatalf("payload mismatch at byte %d: got %d, want %d", i, val[i], payload[i])
		}
	}
}

// =========================================================================
// Integration Test: End-to-End with Multiple Data Types
// =========================================================================

// TestCacheIntegration_MultipleDataTypes exercises the cache with
// various data shapes to ensure JSON marshal/unmarshal works correctly.
func TestCacheIntegration_MultipleDataTypes(t *testing.T) {
	cache, _ := createCacheForIntegration(t)
	ctx := context.Background()

	// Define various data structures
	type address struct {
		Street  string `json:"street"`
		City    string `json:"city"`
		ZipCode string `json:"zip_code"`
	}

	type department struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Active    bool   `json:"active"`
	}

	type employee struct {
		ID         int        `json:"id"`
		Name       string     `json:"name"`
		Email      string     `json:"email"`
		Address    address    `json:"address"`
		Departments []department `json:"departments"`
	}

	emp := employee{
		ID:    1001,
		Name:  "Jane Smith",
		Email: "jane@company.com",
		Address: address{
			Street:  "123 Main St",
			City:    "Jakarta",
			ZipCode: "12345",
		},
		Departments: []department{
			{ID: 1, Name: "Engineering", Active: true},
			{ID: 2, Name: "Product", Active: true},
		},
	}

	// Set
	err := cache.SetJSON(ctx, "int:complex:employee:1001", emp, 0)
	if err != nil {
		t.Fatalf("SetJSON failed: %v", err)
	}

	// Get and verify
	raw, ok := cache.Get(ctx, "int:complex:employee:1001")
	if !ok {
		t.Fatal("expected cache hit for complex object")
	}

	var decoded employee
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.ID != emp.ID || decoded.Name != emp.Name || decoded.Email != emp.Email {
		t.Fatalf("employee mismatch:\n  got:  %+v\n  want: %+v", decoded, emp)
	}
	if decoded.Address.Street != emp.Address.Street || decoded.Address.City != emp.Address.City {
		t.Fatalf("address mismatch:\n  got:  %+v\n  want: %+v", decoded.Address, emp.Address)
	}
	if len(decoded.Departments) != len(emp.Departments) {
		t.Fatalf("departments count: got %d, want %d", len(decoded.Departments), len(emp.Departments))
	}
	if decoded.Departments[0].Name != "Engineering" || !decoded.Departments[0].Active {
		t.Fatal("department 0 mismatch")
	}
	if decoded.Departments[1].Name != "Product" || !decoded.Departments[1].Active {
		t.Fatal("department 1 mismatch")
	}
}

// =========================================================================
// Integration Test: Cross-Instance Prefix Invalidation
// =========================================================================

// TestCacheIntegration_CrossInstancePrefixInvalidation verifies that
// prefix-based invalidation on one cache instance propagates to others.
func TestCacheIntegration_CrossInstancePrefixInvalidation(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	cfg := Config{
		RedisAddr:     mr.Addr(),
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    30 * time.Minute,
	}

	cacheA, err := New(cfg, mustNewLogger(t))
	if err != nil {
		t.Fatalf("failed to create cacheA: %v", err)
	}
	defer cacheA.Close()

	cacheB, err := New(cfg, mustNewLogger(t))
	if err != nil {
		t.Fatalf("failed to create cacheB: %v", err)
	}
	defer cacheB.Close()

	ctx := context.Background()

	// Set keys on cacheA
	keys := []string{"int:shared:cfg:db", "int:shared:cfg:cache", "int:shared:other"}
	for _, k := range keys {
		if err := cacheA.Set(ctx, k, []byte(k), 0); err != nil {
			t.Fatalf("Set %s failed: %v", k, err)
		}
		// Warm up cacheB local cache
		cacheB.Get(ctx, k)
	}

	// Verify all keys cached locally in both instances (delete from Redis first
	// to prove local cache is serving the data, not Redis)
	for _, k := range keys {
		mr.Del(k)
	}
	for _, k := range keys {
		_, ok := cacheA.Get(ctx, k)
		if !ok {
			t.Fatalf("cacheA local cache miss for %s", k)
		}
		_, ok = cacheB.Get(ctx, k)
		if !ok {
			t.Fatalf("cacheB local cache miss for %s", k)
		}
	}

	// Re-set keys in Redis (InvalidatePrefix needs keys to exist in Redis
	// to find them via SCAN + delete them). Also re-warm cacheB local cache.
	for _, k := range keys {
		if err := cacheA.Set(ctx, k, []byte(k), 0); err != nil {
			t.Fatalf("re-Set %s failed: %v", k, err)
		}
		cacheB.Get(ctx, k)
	}

	// Invalidate prefix from cacheB — should propagate to cacheA via Pub/Sub
	err = cacheB.InvalidatePrefix(ctx, "int:shared:cfg:")
	if err != nil {
		t.Fatalf("InvalidatePrefix from cacheB failed: %v", err)
	}

	// Wait for Pub/Sub propagation
	deadline := time.Now().Add(2 * time.Second)
	for {
		_, aOk := cacheA.Get(ctx, "int:shared:cfg:db")
		_, bOk := cacheB.Get(ctx, "int:shared:cfg:db")
		if !aOk && !bOk {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for cross-instance prefix invalidation")
		}
		time.Sleep(10 * time.Millisecond)
	}

	// The 'other' key should still be cached locally
	_, okA := cacheA.Get(ctx, "int:shared:other")
	_, okB := cacheB.Get(ctx, "int:shared:other")
	if !okA || !okB {
		t.Fatal("non-prefix key should still be cached in both instances")
	}
}
