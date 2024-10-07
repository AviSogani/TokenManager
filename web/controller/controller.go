package controller

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"log"
)

var ctx = context.Background()

var redisClient *redis.Client

// initializeCron initiates all the crons in the app
func initializeCron() {
	// Create a new cron job runner
	c := cron.New()

	// Schedule the cleanUp function to run every 5 minutes
	_, err := c.AddFunc("@every 5m", cleanUp)
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	// Schedule the releaseBlockedTokens function to run every minute
	_, err1 := c.AddFunc("@every 60s", releaseBlockedTokens)
	if err1 != nil {
		log.Fatalf("Error adding cron job: %v", err1)
	}

	// Start the cron scheduler
	c.Start()
}

// initializeRedisClient initializes the Redis client
func initializeRedisClient() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379", // Redis server address
	})

	return redisClient
}

// GetRedisClient Singleton client
func GetRedisClient() *redis.Client {
	if redisClient == nil {
		redisClient = initializeRedisClient()
	}
	return redisClient
}

func Init() {
	initializeRedisClient()
	initializeCron()
}
