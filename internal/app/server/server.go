package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/whittier16/student-reg-svc/internal/app/config"
	"github.com/whittier16/student-reg-svc/internal/pkg/database/cache"
	"github.com/whittier16/student-reg-svc/internal/pkg/database/db"
	"github.com/whittier16/student-reg-svc/internal/pkg/handlers"
	"github.com/whittier16/student-reg-svc/internal/pkg/logger"
	"go.elastic.co/apm/module/apmgorilla"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Server is our API server
type Server struct {
	logger *logrus.Logger
	router *mux.Router
	cfg    *config.MainConfig
}

// New returns a new instance of Server
func New() (*Server, error) {
	cnf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	// creates a new instance of database
	db, err := db.New(
		cnf.Database.User,
		cnf.Database.Pass,
		cnf.Database.Host,
		cnf.Database.Port,
		cnf.Database.DBName,
	)
	if err != nil {
		return nil, err
	}

	// creates a new instance of redis
	c := cache.NewRedis(cnf.Redis.Host, cnf.Redis.Port, cnf.Redis.Password, cnf.Redis.DB, cnf.Redis.UseTLS)

	// creates a new instance of a logger
	log := logger.NewLogger()

	// creates a new instance of a mux router
	router := mux.NewRouter().StrictSlash(true)
	// instrument the application
	apmgorilla.Instrument(router)
	handlers.RegisterRoutes(router, log, db, c, cnf)

	s := &Server{
		logger: log,
		cfg:    cnf,
		router: router,
	}
	return s, nil
}

// Start starts the API server
func (s *Server) Start(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Start")
	defer span.Finish()

	// setup CORS
	cors := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Cache-Control"},
		AllowCredentials: true,
	}).Handler(s.router)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Server.Address),
		Handler: cors,
	}
	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(stopServer)

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		s.logger.Printf("REST API listening on port %d", s.cfg.Server.Address)
		serverErrors <- server.ListenAndServe()
	}(&wg)

	// blocking run and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting REST API server: %w", err)
	case <-stopServer:
		s.logger.Warn("server received STOP signal")
		// asking listener to shut down
		err := server.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("graceful shutdown did not complete: %w", err)
		}
		wg.Wait()
		s.logger.Info("server was shut down gracefully")
	}
	return nil
}
