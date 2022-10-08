package config

import (
	"gopkg.in/validator.v2"
)

type Config struct {
	GiftAddr string `mapstructure:"gift_addr" validate:"nonzero"`
}

func (c *Config) Validate() error {
	return validator.NewValidator().Validate(c)
}
