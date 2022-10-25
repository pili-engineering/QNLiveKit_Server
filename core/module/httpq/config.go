package httpq

import (
	"gopkg.in/validator.v2"
)

type Config struct {
	Addr string `mapstructure:"addr" validate:"nonzero"`
}

func (c *Config) Validate() error {
	validate := validator.NewValidator()
	return validate.Validate(c)
}
