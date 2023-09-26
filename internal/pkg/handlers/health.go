package handlers

import (
	"net/http"
)

func (s *Service) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.response(w, "OK", http.StatusOK)
	}
}
