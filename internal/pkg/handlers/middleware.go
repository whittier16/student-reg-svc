package handlers

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// LoggerMiddleware logs the incoming HTTP request and response. Enable it only for debug purpose disable it on production.
func (h *Handler) LoggerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" {
				// Call the next handler don't log if it is internal request from health check of Kubernetes
				next.ServeHTTP(w, r)
				return
			}

			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			requestBody, err := h.readRequestBody(r)
			if err != nil {
				h.response(w, err, 0)
				return
			}
			h.restoreRequestBody(r, requestBody)

			logMessage := fmt.Sprintf("path:%s, method: %s, requestBody: %v", r.URL.EscapedPath(), r.Method, string(requestBody))

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)

			logMessage = fmt.Sprintf("%s, responseStatus: %d, responseBody: %s", logMessage, wrapped.Status(), string(wrapped.Body()))
			h.logger.Infof("%s, duration: %v", logMessage, time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}

// CORSMiddleware wraps the logic for collecting metrics
func (h *Handler) CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Authorization, Cache-Control")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// AuthMiddleware authorizes requests using valid JWT token from the header
func (h *Handler) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" ||
				r.URL.Path == "/auth" {
				// Call the next handler don't log if it is internal request from health check and auth
				next.ServeHTTP(w, r)
				return
			}

			if r.Header["Token"] != nil {
				requestBody, err := h.readRequestBody(r)
				if err != nil {
					h.response(w, err, http.StatusInternalServerError)
					return
				}
				h.restoreRequestBody(r, requestBody)

				token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("there was an error")
					}
					return []byte(h.cfg.JWT.Secret), nil
				})

				if err != nil {
					h.respondWithError(w, err.Error(), http.StatusUnauthorized)
				}

				if token.Valid {
					wrapped := wrapResponseWriter(w)
					next.ServeHTTP(wrapped, r)
				}
			} else {
				h.respondWithError(w, "Not Authorized", http.StatusUnauthorized)
			}
		}
		return http.HandlerFunc(fn)
	}
}
