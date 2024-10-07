package controller

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

const keepAliveInterval = 5 * time.Minute // Token expires if not kept alive within 5 minutes

func KeepTokenAlive(w http.ResponseWriter, r *http.Request) {
	tokenValue := r.URL.Query().Get("token")

	var err error
	client := GetRedisClient()

	if zScoreCmd := client.ZScore(ctx, "tokens::free", tokenValue); zScoreCmd.Err() == nil && zScoreCmd.Val() != 0 {
		err = refreshToken(client, ctx, "tokens::free", tokenValue, zScoreCmd, err)
	} else if zScoreCmd := client.ZScore(ctx, "tokens::withexpiry", tokenValue); zScoreCmd.Err() == nil && zScoreCmd.Val() != 0 {
		err = refreshToken(client, ctx, "tokens::withexpiry", tokenValue, zScoreCmd, err)
	} else {
		err = errors.New("token not found in sets")
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Error in updating expiry of token : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Updated expiry of token: "))

}

// refreshToken is used to refresh the token expiry time
func refreshToken(client *redis.Client, ctx context.Context, key string, token string, zscoreCmd *redis.FloatCmd, err error) error {
	client.ZRem(ctx, key, token)

	originalTtl := zscoreCmd.Val()
	// constant: update ttl keepalive
	additionDuration := time.Duration(5) * time.Minute
	updatedExpiry := time.Unix(int64(originalTtl), 0).Add(additionDuration).Unix()

	err = client.ZAdd(ctx, key, redis.Z{
		Score:  float64(updatedExpiry),
		Member: token,
	}).Err()
	return err
}
