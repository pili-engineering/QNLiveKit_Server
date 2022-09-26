package rtc

import (
	"gopkg.in/validator.v2"
)

type Config struct {
	AppId string `mapstructure:"app_id" validate:"nonzero"` // RTC AppId
}

func (c *Config) Validate() error {
	return validator.NewValidator().Validate(c)
}
