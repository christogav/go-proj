package main

import (
	"fmt"

	"github.com/christogav/go-proj/internal/config"
	"github.com/christogav/go-proj/internal/grpc"
	"github.com/christogav/go-proj/internal/logging"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

type (
	Config struct {
		App App `yaml:"app"`
	}

	App struct {
		Server grpc.Config    `yaml:"server"`
		Log    logging.Config `yaml:"logging"`
	}
)

var ConfigModule = fx.Module(
	"config",

	fx.Provide(LoadConfig),

	fx.Provide(func(config Config) logging.Config {
		return config.App.Log
	}),

	fx.Provide(func(config Config) grpc.Config {
		return config.App.Server
	}),
)

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := config.Load(&cfg); err != nil {
		return cfg, fmt.Errorf("unable to load configuration: %w", err)
	}

	return cfg, cfg.Validate()
}

func (c *Config) Validate() error {
	// No need to reuse validators as they're one offs and no benefit gained from caching
	validate := validator.New()

	return validate.Struct(c)
}
