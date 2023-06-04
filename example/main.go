package main

import (
	"context"
	"log"

	"github.com/christogav/go-proj/internal/logging"
	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()

	app := fx.New(
		ConfigModule,
		logging.Module,
	)

	go func() {
		app.Run()
	}()

	<-app.Done()

	if err := app.Stop(ctx); err != nil {
		log.Fatalf("error stopping gracefully")
	}
}
