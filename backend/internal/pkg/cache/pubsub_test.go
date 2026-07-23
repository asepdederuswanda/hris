package cache

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// setupTestPubSub creates a miniredis server and returns a ready-to-use PubSub
// along with a callback recorder to track invalidation events.
type invalidationRecord struct {
	keys []string
}

type pubSubTestHarness struct {
	mr           *miniredis.Miniredis
	client       *redis.Client
	pubsub       *PubSub
	records      []invalidationRecord
	recordsMu    sync.Mutex
	onInvalidate func(keys ...string)
	logger       *zap.Logger
}

func setupTestPubSub(t *testing.T) *pubSubTestHarness {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	logger := mustNewPubSubLogger(t)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
		DB:   0,
	})

	h := &pubSubTestHarness{
		mr:      mr,
		client:  client,
		logger:  logger,
		records: make([]invalidationRecord, 0),
	}

	h.onInvalidate = func(keys ...string) {
		h.recordsMu.Lock()
		defer h.recordsMu.Unlock()
		h.records = append(h.records, invalidationRecord{keys: keys})
	}

	ps, err := NewPubSub(client, logger, h.onInvalidate)
	if err != nil {
		mr.Close()
		client.Close()
		t.Fatalf("NewPubSub failed: %v", err)
	}
	h.pubsub = ps

	t.Cleanup(func() {
		ps.Close()
		client.Close()
		mr.Close()
	})

	return h
}

func mustNewPubSubLogger(t *testing.T) *zap.Logger {
	t.Helper()
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	return logger
}

// getRecords returns a copy of all recorded invalidation events.
func (h *pubSubTestHarness) getRecords() []invalidationRecord {
	h.recordsMu.Lock()
	defer h.recordsMu.Unlock()
	result := make([]invalidationRecord, len(h.records))
	copy(result, h.records)
	return result
}

// waitForRecords waits until at least n invalidation events are recorded,
// or a timeout elapses.
func (h *pubSubTestHarness) waitForRecords(t *testing.T, n int, timeout time.Duration) []invalidationRecord {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		records := h.getRecords()
		if len(records) >= n {
			return records
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %d invalidation records, got %d", n, len(records))
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

// =========================================================================
// Test: NewPubSub and Close
// =========================================================================

func TestPubSub_NewAndClose(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewPubSubLogger(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	callback := func(keys ...string) {}

	ps, err := NewPubSub(client, logger, callback)
	if err != nil {
		t.Fatalf("NewPubSub failed: %v", err)
	}

	// Close should not panic or error (no return value)
	ps.Close()
}

func TestPubSub_NewPubSub_NilCallback(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewPubSubLogger(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	// Creating with nil callback should not panic at creation.
	// Note: handleMessage() calls ps.onInvalidate(keys...) directly,
	// so receiving a message with a nil callback WILL panic.
	// This test only verifies creation and close — no messages flow.
	ps, err := NewPubSub(client, logger, nil)
	if err != nil {
		t.Fatalf("NewPubSub with nil callback failed: %v", err)
	}
	ps.Close()
}

// =========================================================================
// Test: PublishInvalidate
// =========================================================================

func TestPubSub_PublishInvalidate_SingleKey(t *testing.T) {
	h := setupTestPubSub(t)
	ctx := context.Background()

	err := h.pubsub.PublishInvalidate(ctx, "test:key1")
	if err != nil {
		t.Fatalf("PublishInvalidate failed: %v", err)
	}

	records := h.waitForRecords(t, 1, time.Second)
	if len(records[0].keys) != 1 || records[0].keys[0] != "test:key1" {
		t.Fatalf("expected [test:key1], got %v", records[0].keys)
	}
}

func TestPubSub_PublishInvalidate_MultipleKeys(t *testing.T) {
	h := setupTestPubSub(t)
	ctx := context.Background()

	keys := []string{"key:a", "key:b", "key:c"}
	err := h.pubsub.PublishInvalidate(ctx, keys...)
	if err != nil {
		t.Fatalf("PublishInvalidate failed: %v", err)
	}

	records := h.waitForRecords(t, 1, time.Second)
	if len(records[0].keys) != 3 {
		t.Fatalf("expected 3 keys, got %d: %v", len(records[0].keys), records[0].keys)
	}

	expected := map[string]bool{"key:a": true, "key:b": true, "key:c": true}
	for _, k := range records[0].keys {
		if !expected[k] {
			t.Fatalf("unexpected key %q in invalidation", k)
		}
		delete(expected, k)
	}
}

func TestPubSub_PublishInvalidate_EmptyKeys(t *testing.T) {
	h := setupTestPubSub(t)
	ctx := context.Background()

	// Publishing empty keys should not error and should not trigger callback
	err := h.pubsub.PublishInvalidate(ctx)
	if err != nil {
		t.Fatalf("PublishInvalidate with empty keys should not error: %v", err)
	}

	// Give a moment to ensure no callback is triggered
	time.Sleep(50 * time.Millisecond)
	records := h.getRecords()
	if len(records) != 0 {
		t.Fatalf("expected no invalidation records for empty keys, got %d", len(records))
	}
}

func TestPubSub_PublishInvalidate_SameKeyMultipleTimes(t *testing.T) {
	h := setupTestPubSub(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		err := h.pubsub.PublishInvalidate(ctx, "repeated:key")
		if err != nil {
			t.Fatalf("PublishInvalidate attempt %d failed: %v", i, err)
		}
	}

	records := h.waitForRecords(t, 3, time.Second)
	if len(records) != 3 {
		t.Fatalf("expected 3 invalidation records, got %d", len(records))
	}
	for i, r := range records {
		if len(r.keys) != 1 || r.keys[0] != "repeated:key" {
			t.Fatalf("record %d: expected [repeated:key], got %v", i, r.keys)
		}
	}
}

// =========================================================================
// Test: Cross-Instance Message Propagation
// =========================================================================

func TestPubSub_CrossInstance_Propagation(t *testing.T) {
	// Create two PubSub instances sharing the same Redis to verify
	// messages propagate correctly.
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewPubSubLogger(t)

	client1 := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client1.Close()

	client2 := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client2.Close()

	var mu1, mu2 sync.Mutex
	var received1, received2 []string

	ps1, err := NewPubSub(client1, logger, func(keys ...string) {
		mu1.Lock()
		defer mu1.Unlock()
		received1 = append(received1, keys...)
	})
	if err != nil {
		t.Fatalf("NewPubSub ps1 failed: %v", err)
	}
	defer ps1.Close()

	ps2, err := NewPubSub(client2, logger, func(keys ...string) {
		mu2.Lock()
		defer mu2.Unlock()
		received2 = append(received2, keys...)
	})
	if err != nil {
		t.Fatalf("NewPubSub ps2 failed: %v", err)
	}
	defer ps2.Close()

	ctx := context.Background()

	// Publish from ps1, both ps1 and ps2 should receive
	err = ps1.PublishInvalidate(ctx, "from:ps1")
	if err != nil {
		t.Fatalf("PublishInvalidate from ps1 failed: %v", err)
	}

	// Wait for both to receive
	deadline := time.Now().Add(time.Second)
	for {
		mu1.Lock()
		len1 := len(received1)
		mu1.Unlock()
		mu2.Lock()
		len2 := len(received2)
		mu2.Unlock()

		if len1 >= 1 && len2 >= 1 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for propagation: ps1=%d, ps2=%d", len1, len2)
		}
		time.Sleep(10 * time.Millisecond)
	}

	mu1.Lock()
	if len(received1) != 1 || received1[0] != "from:ps1" {
		t.Fatalf("ps1 received %v, expected [from:ps1]", received1)
	}
	mu1.Unlock()

	mu2.Lock()
	if len(received2) != 1 || received2[0] != "from:ps1" {
		t.Fatalf("ps2 received %v, expected [from:ps1]", received2)
	}
	mu2.Unlock()

	// Publish from ps2, both should receive
	err = ps2.PublishInvalidate(ctx, "from:ps2", "extra:key")
	if err != nil {
		t.Fatalf("PublishInvalidate from ps2 failed: %v", err)
	}

	deadline = time.Now().Add(time.Second)
	for {
		mu1.Lock()
		len1 := len(received1)
		mu1.Unlock()
		mu2.Lock()
		len2 := len(received2)
		mu2.Unlock()

		if len1 >= 2 && len2 >= 2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for second propagation: ps1=%d, ps2=%d", len1, len2)
		}
		time.Sleep(10 * time.Millisecond)
	}

	mu1.Lock()
	if len(received1) != 3 {
		t.Fatalf("ps1 expected 3 total keys (from:ps1 + from:ps2 + extra:key), got %d: %v", len(received1), received1)
	}
	mu1.Unlock()
}

// =========================================================================
// Test: handleMessage (direct unit test)
// =========================================================================

func TestPubSub_HandleMessage_ValidJSON(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewPubSubLogger(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	var capturedKeys []string
	ps, err := NewPubSub(client, logger, func(keys ...string) {
		capturedKeys = append(capturedKeys, keys...)
	})
	if err != nil {
		t.Fatalf("NewPubSub failed: %v", err)
	}
	defer ps.Close()

	// Directly call handleMessage (private, but accessible within same package)
	payload := `{"keys":["direct:key1","direct:key2"]}`
	ps.handleMessage(payload)

	if len(capturedKeys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(capturedKeys), capturedKeys)
	}
	if capturedKeys[0] != "direct:key1" || capturedKeys[1] != "direct:key2" {
		t.Fatalf("expected [direct:key1 direct:key2], got %v", capturedKeys)
	}
}

func TestPubSub_HandleMessage_InvalidJSON(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewPubSubLogger(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	var callbackCalled atomic.Bool
	ps, err := NewPubSub(client, logger, func(keys ...string) {
		callbackCalled.Store(true)
	})
	if err != nil {
		t.Fatalf("NewPubSub failed: %v", err)
	}
	defer ps.Close()

	// Invalid JSON should not trigger callback
	ps.handleMessage("{invalid json}")

	if callbackCalled.Load() {
		t.Fatal("callback should not be called for invalid JSON")
	}
}

func TestPubSub_HandleMessage_EmptyKeys(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewPubSubLogger(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	var callbackCalled atomic.Bool
	ps, err := NewPubSub(client, logger, func(keys ...string) {
		callbackCalled.Store(true)
	})
	if err != nil {
		t.Fatalf("NewPubSub failed: %v", err)
	}
	defer ps.Close()

	// Empty keys array should not trigger callback (len(msg.Keys) > 0 check)
	ps.handleMessage(`{"keys":[]}`)

	if callbackCalled.Load() {
		t.Fatal("callback should not be called for empty keys array")
	}
}

func TestPubSub_HandleMessage_MalformedButClose(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewPubSubLogger(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	var callbackCalled atomic.Bool
	ps, err := NewPubSub(client, logger, func(keys ...string) {
		callbackCalled.Store(true)
	})
	if err != nil {
		t.Fatalf("NewPubSub failed: %v", err)
	}
	defer ps.Close()

	// Edge case: extra fields should be ignored (lenient JSON parsing)
	ps.handleMessage(`{"keys":["valid:key"],"extra":"ignored"}`)

	if !callbackCalled.Load() {
		t.Fatal("callback should be called for valid keys even with extra fields")
	}
}

// =========================================================================
// Test: ChannelConstant
// =========================================================================

func TestPubSub_ChannelConstant(t *testing.T) {
	if ChannelInvalidate != "hris:cache:invalidate" {
		t.Fatalf("expected ChannelInvalidate to be 'hris:cache:invalidate', got %q", ChannelInvalidate)
	}
}

// =========================================================================
// Test: Message JSON Structure
// =========================================================================

func TestPubSub_MessageJSON_Structure(t *testing.T) {
	// Verify that the message structure marshals/unmarshals correctly
	msg := invalidationMessage{Keys: []string{"a", "b"}}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded invalidationMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if len(decoded.Keys) != 2 || decoded.Keys[0] != "a" || decoded.Keys[1] != "b" {
		t.Fatalf("roundtrip failed: got %v", decoded.Keys)
	}
}

func TestPubSub_MessageJSON_Empty(t *testing.T) {
	msg := invalidationMessage{}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal empty message failed: %v", err)
	}

	var decoded invalidationMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal empty message failed: %v", err)
	}

	if len(decoded.Keys) != 0 {
		t.Fatalf("expected empty keys, got %v", decoded.Keys)
	}
}

// =========================================================================
// Test: Concurrent Publish
// =========================================================================

func TestPubSub_ConcurrentPublish(t *testing.T) {
	h := setupTestPubSub(t)
	ctx := context.Background()
	var wg sync.WaitGroup
	n := 30

	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "concurrent:key:" + string(rune('a'+i%26))
			err := h.pubsub.PublishInvalidate(ctx, key)
			if err != nil {
				t.Errorf("concurrent PublishInvalidate failed: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// Wait for all messages to be received
	records := h.waitForRecords(t, n, 2*time.Second)
	if len(records) != n {
		t.Fatalf("expected %d invalidation events, got %d", n, len(records))
	}
}

// =========================================================================
// Test: Subscribe Failure Handling
// =========================================================================

func TestPubSub_NewPubSub_InvalidAddress(t *testing.T) {
	logger := mustNewPubSubLogger(t)
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:1", // Port 1 is unlikely to have a Redis server
	})

	// This should fail because Redis is not reachable
	_, err := NewPubSub(client, logger, func(keys ...string) {})
	if err == nil {
		client.Close()
		t.Fatal("expected error when subscribing to unreachable Redis")
	}
	client.Close()
}

// =========================================================================
// Test: Close Idempotency
// =========================================================================

func TestPubSub_Close_Idempotent(t *testing.T) {
	h := setupTestPubSub(t)

	// First close
	h.pubsub.Close()

	// Second close should not panic (context already cancelled)
	h.pubsub.Close()
}

// =========================================================================
// Test: Callback Not Called After Close
// =========================================================================

func TestPubSub_CallbackNotCalledAfterClose(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	logger := mustNewPubSubLogger(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	var callbackCount atomic.Int32
	ps, err := NewPubSub(client, logger, func(keys ...string) {
		callbackCount.Add(1)
	})
	if err != nil {
		t.Fatalf("NewPubSub failed: %v", err)
	}

	// Close before publishing
	ps.Close()

	// Try to publish after close — the context is cancelled, so the
	// subscriber goroutine has exited and messages won't be processed.
	ctx := context.Background()
	_ = ps.PublishInvalidate(ctx, "after:close")

	// Give time for any stray processing
	time.Sleep(50 * time.Millisecond)

	if count := callbackCount.Load(); count > 0 {
		t.Fatalf("expected 0 callback invocations after close, got %d", count)
	}
}

// =========================================================================
// Test: Publish After Close (Graceful Handling)
// =========================================================================

func TestPubSub_PublishAfterClose(t *testing.T) {
	h := setupTestPubSub(t)
	ctx := context.Background()

	// Close first
	h.pubsub.Close()

	// Publishing after close should still succeed (Publish uses its own context)
	// The message goes to Redis but there's no subscriber to receive it.
	err := h.pubsub.PublishInvalidate(ctx, "after:close")
	if err != nil {
		t.Fatalf("PublishInvalidate after close should not error: %v", err)
	}
}
