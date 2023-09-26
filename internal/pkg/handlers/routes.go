package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/whittier16/student-reg-svc/internal/app/config"
	"github.com/whittier16/student-reg-svc/internal/pkg/database/cache"
	"github.com/whittier16/student-reg-svc/internal/pkg/database/db"
	"net/http"
)

// RegisterRoutes register the different handlers' route definitions.
func RegisterRoutes(r *mux.Router, log *logrus.Logger, db *db.MySQL, c *cache.Redis, cfg *config.MainConfig) {
	h := New(log, db, cfg)

	// adding logger middleware
	r = addMiddlewares(r, h)
	r.HandleFunc("/healthz", h.Health())
	r.HandleFunc("/auth", h.Auth())
	r.HandleFunc("/api/register", h.Register()).Methods(http.MethodPost)
	r.HandleFunc("/api/commonstudents", h.GetCommonStudents()).Methods(http.MethodGet)
	r.HandleFunc("/api/suspend", h.Suspend()).Methods(http.MethodPost)
	r.HandleFunc("/api/retrievefornotifications", h.RetrieveNotifications()).Methods(http.MethodPost)
	r.HandleFunc("/api/students", h.CreateStudent()).Methods(http.MethodPost)
	r.HandleFunc("/api/teachers", h.CreateTeacher()).Methods(http.MethodPost)
}

func addMiddlewares(r *mux.Router, h *Service) *mux.Router {
	r.Use(h.LoggerMiddleware())
	r.Use(h.AuthMiddleware())
	return r
}
