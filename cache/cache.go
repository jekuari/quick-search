package cache

import (
	"context"

	"github.com/jekuari/quick-search/constants"
	"github.com/jekuari/quick-search/logger"
	"github.com/redis/go-redis/v9"
)

func RedisClient(db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     constants.REDIS_URL,
		Password: "",
		DB:       0,
	})

	return rdb
}

// get redis client from context
func GetRedisSearchesClient(ctx context.Context) *redis.Client {
	redisClient := ctx.Value(constants.REDIS_SEARCHES_CONTEXT_KEY)
	if redisClient == nil {
		logger.Log("redisClient is nil")
		return nil
	}

	return redisClient.(*redis.Client)
}

func GetRedisRateLimitsClient(ctx context.Context) *redis.Client {
	redisClient := ctx.Value(constants.REDIS_RATE_LIMITS_CONTEXT_KEY)
	if redisClient == nil {
		logger.Log("redisClient is nil")
		return nil
	}

	return redisClient.(*redis.Client)
}
