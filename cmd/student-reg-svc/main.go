package main

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/whittier16/student-reg-svc/internal/app/server"
	"github.com/whittier16/student-reg-svc/internal/pkg/logger"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmot"
	"time"
)

func main() {
	// initialize tracing as early as possible
	opentracing.SetGlobalTracer(apmot.New())
	defer apm.DefaultTracer.Flush(nil)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "run")
	defer span.Finish()

	server, err := server.New()
	if err != nil {
		return err
	}
	err = server.Start(ctx)
	return err
}
