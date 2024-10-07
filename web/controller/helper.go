package controller

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func cleanUp() {
	currentTimeString := fmt.Sprintf("%v", time.Now().UTC().Unix())
	GetRedisClient().ZRemRangeByScore(context.Background(), "tokens::free", "-inf", currentTimeString).Err()
}

// unblocker unblocks a token by taking in the key
func unblocker(ctx context.Context, key string, client *redis.Client) {
	originalKeyTtl := client.ZScore(ctx, "tokens::withexpiry", key).Val()

	client.ZRem(ctx, "tokens::blocked", key)
	client.ZRem(ctx, "tokens::withexpiry", key)
	client.ZAdd(ctx, "tokens::free", redis.Z{
		Score:  originalKeyTtl,
		Member: key,
	})
}

// Unblock all tokens that can be unblocked
func releaseBlockedTokens() {
	currTimeString := fmt.Sprintf("%v", time.Now().UTC().Unix())
	keyList, _ := GetRedisClient().ZRangeByScore(context.Background(), "tokens::blocked", &redis.ZRangeBy{
		Min: "-inf",
		Max: currTimeString,
	}).Result()

	for _, key := range keyList {
		unblocker(ctx, key, GetRedisClient())
	}
}
