package controller

import "net/http"

func UnblockToken(w http.ResponseWriter, r *http.Request) {
	tokenValue := r.URL.Query().Get("token")
	unblocker(ctx, tokenValue, GetRedisClient())

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Token unblocked."))
}
