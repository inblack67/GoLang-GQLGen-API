package cache

import (
	"context"
	"log"

	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
	"github.com/inblack67/GQLGenAPI/db"
	"github.com/inblack67/GQLGenAPI/mymodels"
)

// RedisClient ...
var (
	RedisClient *redis.Client
)

// StartRedis ...
func StartRedis() (context.Context) {
	RedisClient = redis.NewClient(&redis.Options{})

	color.Blue("Redis is here")

	ctx := context.Background()

	RedisClient.FlushAll(ctx)

	return  ctx
}

// PopulateUsers ...
func PopulateUsers () error {
	var dbUsers []*mymodels.User

	err := db.PgConn.Find(&dbUsers).Error

	if err != nil {
		log.Println("err populating users in redis = ", err)
		return err
	}

	color.HiGreen("users populated in redis")

	return nil
}

// PopulateStories ...
func PopulateStories () error {
	var dbStories []*mymodels.Story

	err := db.PgConn.Find(&dbStories).Error

	if err != nil {
		log.Println("err populating stories in redis = ", err)
		return err
	}

	color.HiGreen("stories populated in redis")

	return nil
}