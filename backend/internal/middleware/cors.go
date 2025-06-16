package middleware

import (
	"backend/config"
	"net/http"
)

var allowedOrigins = map[string]bool{
	config.FrontendURL:      true,
	"http://localhost:5173": true,
}

func WithCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow frontend origin
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		// Allow credentials like cookies
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// For preflight requests
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}
