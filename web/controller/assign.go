package controller

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

const tokenLifetime = 60 * time.Second // Token auto-release after 60 seconds

func addToBlockedTokens(ctx context.Context, key string, client *redis.Client) error {
	availableKeyTtl, _ := client.ZScore(ctx, "tokens::free", key).Result()

	err := client.ZRem(ctx, "tokens::free", key).Err()
	if err != nil {
		return err
	}

	client.ZAdd(ctx, "tokens::withexpiry", redis.Z{
		Score:  availableKeyTtl,
		Member: key,
	})

	blockedKeyTtl := time.Now().UTC().Add(tokenLifetime).Unix()
	client.ZAdd(ctx, "tokens::blocked", redis.Z{
		Score:  float64(blockedKeyTtl),
		Member: key,
	})

	return nil
}

// AssignToken assigns unique token from the pool
func AssignToken(w http.ResponseWriter, _ *http.Request) {
	randomToken, err := redisClient.ZRandMember(ctx, "tokens::free", 1).Result()
	if err == nil {
		if len(randomToken) == 0 || randomToken[0] == "" {
			_, _ = w.Write([]byte("No token available in set"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Println("Token assigned: ", randomToken)

		err = addToBlockedTokens(ctx, randomToken[0], redisClient)
		if err != nil {
			_, _ = w.Write([]byte("No token available in set"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jsonErr := json.NewEncoder(w).Encode(map[string]string{"token": randomToken[0]})
		if jsonErr != nil {
			_, _ = w.Write([]byte("No token available in set"))
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}
	} else {
		log.Println("Error", err)
	}
}
