package account

import (
	"gopkg.in/validator.v2"
)

var defaultConfig = &Config{}

type Config struct {
	AccessKey string `mapstructure:"access_key" validate:"nonzero"`
	SecretKey string `mapstructure:"secret_key" validate:"nonzero"`
}

func (c *Config) Validate() error {
	return validator.NewValidator().Validate(c)
}

func AccessKey() string {
	return defaultConfig.AccessKey
}

func SecretKey() string {
	return defaultConfig.SecretKey
}
