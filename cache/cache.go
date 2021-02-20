package cache

import (
	"context"
	"encoding/json"
	"log"

	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
	"github.com/inblack67/GQLGenAPI/constants"
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

	err := db.PgConn.Preload("Stories").Find(&dbUsers).Error

	if err != nil {
		log.Println("err populating users in redis = ", err)
		return err
	}

	marshalledUsers, marshallErr := json.Marshal(dbUsers)

	if marshallErr != nil {
		log.Println("marshallErr = ", marshallErr)
		return marshallErr
	}

	setErr := RedisClient.Set(context.Background(), constants.KUsers, marshalledUsers, 0).Err()		// dont expire on your own

	if setErr != nil {
		log.Println("setErr = ", setErr)
		return setErr
	}

	color.HiGreen("users populated in redis")

	return nil
}

// PopulateStories ...
func PopulateStories () error {
	var dbStories []*mymodels.Story

	err := db.PgConn.Preload("User").Find(&dbStories).Error

	if err != nil {
		log.Println("err populating stories in redis = ", err)
		return err
	}

	marshalledStories, marshallErr := json.Marshal(dbStories)

	if marshallErr != nil {
		log.Println("marshallErr = ", marshallErr)
		return marshallErr
	}

	setErr := RedisClient.Set(context.Background(), constants.KStories, marshalledStories, 0).Err()		// dont expire on your own

	if setErr != nil {
		log.Println("setErr = ", setErr)
		return setErr
	}

	color.HiGreen("stories populated in redis")

	return nil
}