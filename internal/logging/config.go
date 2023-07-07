package logging

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const (
	levelKey = "log.level"
	devKey   = "log.dev"
)

type (
	// Config defines a common definition of logging options.
	Config struct {
		Level string
		Dev   bool
	}
)

func newConfig() *Config {
	viper.SetDefault(levelKey, zerolog.TraceLevel)
	viper.SetDefault(devKey, false)

	return &Config{
		Level: viper.GetString(levelKey),
		Dev:   viper.GetBool(devKey),
	}
}
