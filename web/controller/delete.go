package controller

import "net/http"

// DeleteToken deletes tokens from the pool
func DeleteToken(w http.ResponseWriter, r *http.Request) {
	tokenValue := r.URL.Query().Get("token")
	setNames := []string{"tokens::free", "tokens::blocked", "tokens::withexpiry"}
	for _, setName := range setNames {
		err := GetRedisClient().ZRem(ctx, setName, tokenValue).Err()
		if err != nil {
			_, _ = w.Write([]byte("Token " + tokenValue + " is not available."))
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Token " + tokenValue + " is deleted."))

}
