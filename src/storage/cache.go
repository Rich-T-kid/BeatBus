package storage

import (
	"BeatBus/internal"
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

type messageQueue struct {
	client *redis.Client
	logger *log.Logger
	mu     sync.RWMutex
}

var cfg = internal.GetConfig()

func newRedisClient(redisURI string) *redis.Client {
	const port = "6379"
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
	return client
}
func NewMessageQueue(l *log.Logger) *messageQueue {
	client := newRedisClient(cfg.RedisURI)
	return &messageQueue{
		client: client,
		logger: l,
	}
}
