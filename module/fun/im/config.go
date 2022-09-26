package im

import (
	"gopkg.in/validator.v2"
)

type Config struct {
	AppId    string `mapstructure:"app_id" validate:"nonzero"`
	Endpoint string `mapstructure:"endpoint" validate:"nonzero"`
	Token    string `mapstructure:"token" validate:"nonzero"` //管理员AccessToken
}

func (c *Config) Validate() error {
	return validator.NewValidator().Validate(c)
}
