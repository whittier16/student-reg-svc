package main

import (
	"context"
	"github.com/whittier16/student-reg-svc/internal/app/server"
	"github.com/whittier16/student-reg-svc/internal/pkg/logger"
	"time"
)

func main() {
	logger := logger.NewLogger()
	logger.Println("Start")
	start := time.Now()
	if err := run(context.Background()); err != nil {
		logger.Fatalf("%+v", err)
	}
	logger.Println("Since", time.Since(start))
	logger.Println("Done")
}

func run(ctx context.Context) error {
	server, err := server.New()
	if err != nil {
		return err
	}
	err = server.Start(ctx)
	return err
}
