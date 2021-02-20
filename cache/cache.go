package cache

import (
	"context"

	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
)

// RedisClient ...
var (
	RedisClient *redis.Client
)

// StartRedis ...
func StartRedis() (context.Context){
	RedisClient = redis.NewClient(&redis.Options{})

	color.Blue("Redis is here")

	ctx := context.Background()

	RedisClient.FlushAll(ctx)

	return  ctx
}

// SET ...
func SET(key string, value string) error{
	ctx := StartRedis()
	err := RedisClient.Set(ctx, key, value, 0).Err()
	return err
}

// GET ...
func GET(key string) (string, error){
	ctx := StartRedis()
	val , err := RedisClient.Get(ctx, key).Result()
	return val, err
}