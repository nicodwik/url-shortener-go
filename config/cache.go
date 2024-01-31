package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisConn *redis.Client
var ctx = context.Background()

func InitCache() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", redisHost, redisPort),
		Password: "",
		DB:       1,
	})

	RedisConn = redis
}

func CacheRemember(key string, timeout int, callback func() interface{}) (string, error) {
	var returnData string
	cachePrefix := "url-shortener-go:redirection:"
	fullCacheKey := cachePrefix + key

	cachedData := RedisConn.Get(ctx, fullCacheKey)

	if len(cachedData.Val()) == 0 {
		value := callback()

		if value != nil {
			bytesValue, err := json.Marshal(value)
			if err != nil {
				return "", err
			}

			err = RedisConn.Set(ctx, fullCacheKey, bytesValue, 15*time.Minute).Err()
			if err != nil {
				return "", err
			}

			returnData = string(bytesValue)
		}
	} else {
		returnData = cachedData.Val()
	}

	return returnData, nil
}

func CacheIncrement(prefix string, key string) error {
	cachePrefix := "url-shortener-go"
	cacheKey := fmt.Sprintf(cachePrefix+":%v:%v", prefix, key)

	// *** INCR method ***
	if err := RedisConn.Incr(ctx, cacheKey).Err(); err != nil {
		return err
	}

	// *** ZADD method ***
	// cacheKey := fmt.Sprintf(cachePrefix+"|%v", prefix)

	// members := []redis.Z{{Score: 1, Member: key}}
	// if err := RedisConn.ZAddArgsIncr(ctx, cacheKey, redis.ZAddArgs{
	// 	Members: members,
	// }).Err(); err != nil {
	// 	return err
	// }

	return nil
}

func CacheGetHitCountTemp() (map[string]interface{}, error) {
	keys := []string{}

	iterator := RedisConn.Scan(ctx, 0, "url-shortener-go:hit-count-temp:*", 0).Iterator()
	for iterator.Next(ctx) {
		keys = append(keys, iterator.Val())
	}

	if err := iterator.Err(); err != nil {
		return nil, err
	}

	urlHitCountInput := make(map[string]interface{})

	if len(keys) > 0 {
		cacheValues := RedisConn.MGet(ctx, keys...)
		if err := cacheValues.Err(); err != nil {
			return nil, err
		}

		for idx, key := range keys {
			splittedString := strings.Split(key, ":")
			shortUrl := splittedString[len(splittedString)-1] // take last word
			urlHitCountInput[shortUrl] = cacheValues.Val()[idx]
		}
	}

	return urlHitCountInput, nil
}

func CacheResetHitCountTemp() error {
	keys := []string{}

	iterator := RedisConn.Scan(ctx, 0, "url-shortener-go:hit-count-temp:*", 0).Iterator()
	for iterator.Next(ctx) {
		keys = append(keys, iterator.Val())
	}

	if err := iterator.Err(); err != nil {
		return err
	}

	for _, key := range keys {
		if err := RedisConn.Del(ctx, key).Err(); err != nil {
			return err
		}
	}

	return nil
}
