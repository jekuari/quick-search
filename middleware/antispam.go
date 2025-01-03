package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jekuari/quick-search/cache"
	"github.com/jekuari/quick-search/constants"
	"github.com/jekuari/quick-search/logger"
)

type IpLimit struct {
	Count int       `json:"count"`
	Time  time.Time `json:"time"`
}

func (ipData *IpLimit) Marshal() ([]byte, error) {
	return json.Marshal(ipData)
}

func (ipData *IpLimit) Unmarshal(data []byte) error {
	return json.Unmarshal(data, ipData)
}

func (ipData *IpLimit) Increment() {
	ipData.Count++
}

func (ipData *IpLimit) Reset() {
	ipData.Count = 0
}

func AntiSpam(ctx context.Context, next http.Handler) http.Handler {
	redisClient := cache.GetRedisRateLimitsClient(ctx)
	if redisClient == nil {
		panic("Redis client not found in context")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		ipData := IpLimit{}

		// We should store a json object containing the count and the time
		retrievedData, err := redisClient.Get(ctx, ip).Result()
		if err != nil {
			// create a new IpLimit object
			ipData.Reset()

			// If the key does not exist, we should set it
			ipData.Time = time.Now()
			ipData.Count = 1

			newData, err := ipData.Marshal()
			if err != nil {
				logger.Error("Could not marshal data", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Set the new data in redis
			err = redisClient.Set(ctx, ip, newData, time.Hour).Err()
			if err != nil {
				logger.Error("Could not set value in redis", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Serve request
			next.ServeHTTP(w, r)
			return
		}

		// Unmarshal the data
		err = ipData.Unmarshal([]byte(retrievedData))
		if err != nil {
			logger.Error("Could not unmarshal data from redis", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Reject if the count is over the limit
		if ipData.Count >= constants.SEARCH_HOUR_RATE_LIMIT {
			logger.Error("Too many requests", errors.New("Too many requests"))
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		// The request is within limits, so we can increment the count
		ipData.Increment()

		newData, err := ipData.Marshal()
		if err != nil {
			logger.Error("Could marshal data", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the new data in redis
		err = redisClient.Set(ctx, ip, newData, time.Hour).Err()
		if err != nil {
			logger.Error("Could not set redis data", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Log("Request from: ", ip, " count: ", ipData.Count)
		next.ServeHTTP(w, r)
	})
}
