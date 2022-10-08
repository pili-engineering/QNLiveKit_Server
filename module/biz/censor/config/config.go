package config

import (
	"gopkg.in/validator.v2"
)

type Config struct {
	Callback string `mapstructure:"callback" validate:"nonzero"`
	Bucket   string `mapstructure:"bucket" validate:"nonzero"`
	Addr     string `mapstructure:"addr" validate:"nonzero"`
}

func (c *Config) Validate() error {
	return validator.NewValidator().Validate(c)
}
