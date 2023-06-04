package config

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

// FromFile loads config from a file. `config` must be a pointer to a struct.
func FromFile(file string, config interface{}) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return err
	}

	return nil
}

// Load loads config from the default `config.yaml` in a folder matching `path.Base(os.Args[0])` in `./config`.
func Load(config interface{}) error {
	return FromFile(path.Join("config", path.Base(os.Args[0]), "config.yaml"), config)
}
