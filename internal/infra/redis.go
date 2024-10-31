package infra

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Redis struct {
	client *redis.Client
}

func NewRedis(config *redis.Options) *Redis {
	client := redis.NewClient(config)

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)

	return &Redis{client: client}
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}
