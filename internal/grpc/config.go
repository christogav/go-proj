package grpc

import (
	"github.com/spf13/viper"
)

const (
	hostkey = "host"
	portkey = "port"
)

type (
	Config struct {
		Host string
		Port int
	}
)

func newConfig() *Config {
	viper.SetDefault(hostkey, "0.0.0.0")
	viper.SetDefault(portkey, 8888)

	return &Config{
		Host: viper.GetString(hostkey),
		Port: viper.GetInt(portkey),
	}
}
