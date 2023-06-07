package main

import (
	"context"
	"log"

	"github.com/christogav/go-proj/internal/grpc"
	"github.com/christogav/go-proj/internal/logging"
	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()

	logging.Init()

	app := fx.New(
		ConfigModule,
		logging.Module,
		grpc.Module,

		logging.WithLogger,
	)

	go func() {
		app.Run()
	}()

	<-app.Done()

	if err := app.Stop(ctx); err != nil {
		log.Fatalf("error stopping gracefully")
	}
}
