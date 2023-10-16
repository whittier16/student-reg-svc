package handlers

import (
	"net/http"
)

// Health checks the up status of the server
func (h *Handler) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.response(w, "OK", http.StatusOK)
	}
}
