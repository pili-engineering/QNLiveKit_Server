package rtc

import (
	"gopkg.in/validator.v2"
)

type Config struct {
	AppId     string `mapstructure:"app_id" validate:"nonzero"`     // RTC AppId
	AccessKey string `mapstructure:"access_key" validate:"nonzero"` // AK
	SecretKey string `mapstructure:"secret_key" validate:"nonzero"` // SK
}

func (c *Config) Validate() error {
	return validator.NewValidator().Validate(c)
}
