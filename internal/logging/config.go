package logging

type (
	// Config defines a common definition of logging options.
	Config struct {
		Level string `yaml:"level"`
		Local bool   `yaml:"local"`
	}
)
