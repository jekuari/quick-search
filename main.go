package main

import (
	"context"
	"net/http"
	"os"

	"github.com/jekuari/quick-search/.git/logger"
	"github.com/jekuari/quick-search/cache"
	constants "github.com/jekuari/quick-search/constants"
	handlers "github.com/jekuari/quick-search/handlers"
	dotenv "github.com/joho/godotenv"
)

var GLOAL_CONTEXT = context.Background()

func main() {
	logger.Log("Server starting...")

	env := os.Getenv("GO_ENV")
	if env != "production" {
		if err := dotenv.Load(".dev.env"); err != nil {
			panic("Error loading .env file")
		}

		logger.Log("Loaded dev environment")
	}

	constants.LoadEnvVariables()

	redisClientSearches := cache.RedisClient(constants.REDIS_DB_SEARCHES)
	redisClientRateLimits := cache.RedisClient(constants.REDIS_DB_RATE_LIMITS)

	ctx := context.WithValue(GLOAL_CONTEXT, constants.REDIS_SEARCHES_CONTEXT_KEY, redisClientSearches)
	ctx = context.WithValue(ctx, constants.REDIS_RATE_LIMITS_CONTEXT_KEY, redisClientRateLimits)

	handlers.Setup(ctx)

	logger.Log("Loaded handlers...")

	logger.Log("Server started on port: ", constants.PORT)

	// set no cache headers for all requests
	// create a blocking channel to keep the server running
	block := make(chan struct{}, 1)
	go func() {
		err := http.ListenAndServe(constants.PORT, nil)
		if err != nil {
			logger.Log("Error starting server: ", err)
		}
	}()
	<-block
}
