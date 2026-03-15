package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request in Response Time Middleware")
		start := time.Now()

		// Create a custom responseWriter to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		// Calculate the duration
		duration := time.Since(start)
		wrappedWriter.Header().Set("X-Response-Time", duration.String()) // doesn't work as of now

		// Log the request details
		fmt.Printf("Method: %s, URL: %s, Status: %d, Duration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.String())
		fmt.Println("Sent response from Response Time Middleware")
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
