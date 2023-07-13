package config

import (
	validate "github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

// Config...
type Config struct {
	PORT       string `validate:"required" envconfig:"PORT" default:"8080"`
	PARENTDIR  string `validate:"required" envconfig:"PARENT_DIR"`
	KERNEL     string `validate:"required" envconfig:"KERNEL" default:"vmlinux.bin"`
	InitrdPath string `validate:"required" envconfig:"INITRD_PATH" default:"initrd.cpio"`
	FCBIN      string `validate:"required" envconfig:"FC_BINARY" default:"firecracker"`
	LogLevel   string `validate:"required" envconfig:"LOG_LEVEL" default:"debug"`
	DBNAME     string `validate:"required" envconfig:"DB_NAME"`
}

// Load configuration loads conf variables
func Load(prefix string) (*Config, error) {

	c := new(Config)

	if err := envconfig.Process(prefix, c); err != nil {
		return nil, err
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// validate configuration
func (c *Config) validate() error {

	var validator = validate.New()

	if err := validator.Struct(c); err != nil {
		return err
	}
	return nil
}
