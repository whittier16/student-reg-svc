package handlers

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging. This type will implement http.ResponseWriter.
type responseWriter struct {
	http.ResponseWriter
	status      int
	body        []byte
	wroteHeader bool
	wroteBody   bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteBody {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	if rw.wroteBody {
		return 0, nil
	}
	i, err := rw.ResponseWriter.Write(body)
	if err != nil {
		return 0, err
	}
	rw.body = body
	return i, err
}

func (rw *responseWriter) Body() []byte {
	return rw.body
}

// LoggerMiddleware logs the incoming HTTP request and response. Enable it only for debug purpose disable it on production.
func (s *Service) LoggerMiddleware() func(http.Handler) http.Handler {
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

			requestBody, err := s.readRequestBody(r)
			if err != nil {
				s.response(w, err, 0)
				return
			}
			s.restoreRequestBody(r, requestBody)

			logMessage := fmt.Sprintf("path:%s, method: %s, requestBody: %v", r.URL.EscapedPath(), r.Method, string(requestBody))

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)

			logMessage = fmt.Sprintf("%s, responseStatus: %d, responseBody: %s", logMessage, wrapped.Status(), string(wrapped.Body()))
			s.logger.Infof("%s, duration: %v", logMessage, time.Since(start))
		}
		return http.HandlerFunc(fn)
	}
}

func (s *Service) AuthMiddlware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" ||
				r.URL.Path == "/auth" {
				// Call the next handler don't log if it is internal request from health check of Kubernetes
				next.ServeHTTP(w, r)
				return
			}

			if r.Header["Token"] != nil {
				requestBody, err := s.readRequestBody(r)
				if err != nil {
					s.response(w, err, http.StatusInternalServerError)
					return
				}
				s.restoreRequestBody(r, requestBody)

				token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("there was an error")
					}
					return sampleSecretKey, nil
				})

				if err != nil {
					s.respondWithError(w, err.Error(), http.StatusInternalServerError)
				}

				if token.Valid {
					wrapped := wrapResponseWriter(w)
					next.ServeHTTP(wrapped, r)
				}
			} else {
				s.respondWithError(w, "Not Authorized", http.StatusUnauthorized)
			}
		}
		return http.HandlerFunc(fn)
	}
}
