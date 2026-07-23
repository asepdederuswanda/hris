package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// ChannelRedis adalah channel Redis Pub/Sub untuk invalidasi cache.
const ChannelInvalidate = "hris:cache:invalidate"

// invalidationMessage adalah payload yang dikirim via Pub/Sub.
type invalidationMessage struct {
	Keys []string `json:"keys"`
}

// PubSub mengelola Redis Pub/Sub untuk distributed cache invalidation.
// Saat satu instance meng-invalidate cache, semua instance lain akan
// menerima notifikasi dan ikut meng-evict local cache mereka.
type PubSub struct {
	client     *redis.Client
	logger     *zap.Logger
	onInvalidate func(keys ...string)
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// NewPubSub membuat PubSub baru dan memulai subscriber goroutine.
// Parameter onInvalidate akan dipanggil setiap kali ada pesan invalidasi
// dari instance lain.
func NewPubSub(client *redis.Client, logger *zap.Logger, onInvalidate func(keys ...string)) (*PubSub, error) {
	ctx, cancel := context.WithCancel(context.Background())

	ps := &PubSub{
		client:       client,
		logger:       logger,
		onInvalidate: onInvalidate,
		ctx:          ctx,
		cancel:       cancel,
	}

	// Start subscriber in background
	if err := ps.startSubscriber(); err != nil {
		cancel()
		return nil, err
	}

	logger.Info("Pub/Sub cache invalidation started",
		zap.String("channel", ChannelInvalidate),
	)

	return ps, nil
}

// startSubscriber memulai goroutine yang subscribe ke channel invalidasi.
func (ps *PubSub) startSubscriber() error {
	pubsub := ps.client.Subscribe(ps.ctx, ChannelInvalidate)

	// Verifikasi subscription
	_, err := pubsub.Receive(ps.ctx)
	if err != nil {
		return fmt.Errorf("cache/pubsub: failed to subscribe to %s: %w", ChannelInvalidate, err)
	}

	ps.wg.Add(1)
	go func() {
		defer ps.wg.Done()
		ps.listen(pubsub)
	}()

	return nil
}

// listen menerima pesan dari channel dan memanggil callback.
func (ps *PubSub) listen(pubsub *redis.PubSub) {
	ch := pubsub.Channel()
	for {
		select {
		case <-ps.ctx.Done():
			pubsub.Close()
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			ps.handleMessage(msg.Payload)
		}
	}
}

// handleMessage memproses pesan invalidasi dari instance lain.
func (ps *PubSub) handleMessage(payload string) {
	var msg invalidationMessage
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		ps.logger.Warn("Cache Pub/Sub: failed to unmarshal invalidation message",
			zap.String("payload", payload),
			zap.Error(err),
		)
		return
	}

	if len(msg.Keys) > 0 {
		ps.onInvalidate(msg.Keys...)
	}
}

// PublishInvalidate mengirim pesan invalidasi ke semua instance.
// Instance lain akan menerima notifikasi dan meng-evict local cache mereka.
func (ps *PubSub) PublishInvalidate(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	msg := invalidationMessage{Keys: keys}
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("cache/pubsub: failed to marshal invalidation message: %w", err)
	}

	if err := ps.client.Publish(ctx, ChannelInvalidate, data).Err(); err != nil {
		return fmt.Errorf("cache/pubsub: failed to publish invalidation: %w", err)
	}

	ps.logger.Debug("Published cache invalidation",
		zap.Int("keys", len(keys)),
	)
	return nil
}

// Close menutup subscriber dan melepaskan resources.
func (ps *PubSub) Close() {
	ps.cancel()
	ps.wg.Wait()
}
