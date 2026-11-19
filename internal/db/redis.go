package db

import (
"context"
"log"

"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis(ctx context.Context, redisURL string) error {
opt, err := redis.ParseURL(redisURL)
if err != nil {
return err
}

RedisClient = redis.NewClient(opt)
if err := RedisClient.Ping(ctx).Err(); err != nil {
return err
}

log.Println("Connected to Redis")
return nil
}