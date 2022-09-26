package prome

import (
	"context"

	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
)

const moduleName = "prome"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	if c == nil {
		return nil
	}

	conf := Config{}
	if err := c.Unmarshal(&conf); err != nil {
		return err
	}

	if err := validateConfig(conf); err != nil {
		return err
	}

	instance = newService(&conf)

	return nil
}

func (m *Module) Start() error {
	if instance != nil {
		return instance.Start(context.Background())
	} else {
		return nil
	}
}
