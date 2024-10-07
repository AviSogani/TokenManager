package controller

import "net/http"

// Health checks if the app is healthy
func Health(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("App is healthy."))
}
