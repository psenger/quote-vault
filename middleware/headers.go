package middleware

import (
	"net/http"
	"time"
)

// SecurityHeaders adds common security headers to responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// API information headers
		w.Header().Set("X-API-Version", "v1")
		w.Header().Set("X-Powered-By", "Quote Vault API")
		
		// Caching headers for API responses
		if r.Method == "GET" {
			w.Header().Set("Cache-Control", "public, max-age=300")
			w.Header().Set("ETag", generateETag(r.URL.Path))
		} else {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		}
		
		next.ServeHTTP(w, r)
	})
}

// ResponseTime adds response time header
func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		w.Header().Set("X-Response-Time", duration.String())
	})
}

// generateETag creates a simple ETag based on the path
func generateETag(path string) string {
	return `"` + path + `-` + time.Now().Format("20060102") + `"`
}