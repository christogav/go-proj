package main

import (
	"context"
	"log"
	"strings"

	"github.com/christogav/go-proj/internal/grpc"
	hellov1 "github.com/christogav/go-proj/internal/hello/v1"
	"github.com/christogav/go-proj/internal/logging"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	EnvPrefix = "TEST_PROJ"
)

func main() {
	ctx := context.Background()

	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/example")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = logging.Init()
	if err != nil {
		panic(err)
	}

	app := fx.New(
		logging.Module,
		grpc.Module,
		hellov1.Module,
	)

	go func() {
		app.Run()
	}()

	<-app.Done()

	if err := app.Stop(ctx); err != nil {
		log.Fatalf("error stopping gracefully")
	}
}
