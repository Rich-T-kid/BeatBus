package storage

import (
	"BeatBus/internal"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type messageQueue struct {
	client *redis.Client
	logger *log.Logger
	mu     sync.RWMutex
}

var (
	cfg                = internal.GetConfig()
	ErrKeyDoesNotExist = fmt.Errorf("key does not exist")
	rdsClient          *redis.Client // singleton instance
)

func newRedisClient(redisURI string) *redis.Client {
	const port = "6379"
	if rdsClient != nil {
		return rdsClient
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisURI + ":" + port,
		Password: "", // No password set
		DB:       0,  // Use default DB

	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		fmt.Printf("Failed to ping redis connection: redisURI = %s, port : %s\n", redisURI, port)
		panic(err)
	}
	rdsClient = client
	return client
}
func NewMessageQueue(l *log.Logger) *messageQueue {
	client := newRedisClient(cfg.RedisURI)
	return &messageQueue{
		client: client,
		logger: l,
	}
}
func (mq *messageQueue) EnsureKeyExists(key string) error {
	ctx := context.Background()
	val, err := mq.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if val == 0 {
		return ErrKeyDoesNotExist
	}
	return nil
}
func (mq *messageQueue) SetKeyWithExpiry(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	mq.mu.Lock()
	defer mq.mu.Unlock()
	err := mq.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		mq.logger.Println("Failed to set key in redis:", err)
		return err
	}
	return nil
}
func (mq *messageQueue) Incr(key string) error {
	ctx := context.Background()
	mq.mu.Lock()
	defer mq.mu.Unlock()
	err := mq.client.Incr(ctx, key).Err()
	if err != nil {
		mq.logger.Println("Failed to increment key in redis:", err)
		return err
	}
	return nil
}
func (mq *messageQueue) UpdateChannel(channel string, message interface{}) error {
	ctx := context.Background()
	err := mq.client.Publish(ctx, channel, message).Err()
	if err != nil {
		mq.logger.Println("Failed to publish message to channel:", err)
		return err
	}
	mq.logger.Printf("Published message to channel %s: %v\n", channel, message)
	return nil
}
func (mq *messageQueue) SubscribeChannel(channel string) *redis.PubSub {
	ctx := context.Background()
	pubsub := mq.client.Subscribe(ctx, channel)
	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(ctx)
	if err != nil {
		mq.logger.Println("Failed to subscribe to channel:", err)
		return nil
	}
	return pubsub
}
