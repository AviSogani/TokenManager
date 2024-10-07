package controller

import (
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

// GenerateToken generates unique tokens in the pool
func GenerateToken(w http.ResponseWriter, _ *http.Request) {
	token, _ := uuid.NewV7()

	expireTtl := time.Now().UTC().Add(keepAliveInterval).Unix()
	GetRedisClient().ZAdd(ctx, "tokens::free", redis.Z{
		Score:  float64(expireTtl),
		Member: token.String(),
	})

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("New token generated : " + token.String()))

}
