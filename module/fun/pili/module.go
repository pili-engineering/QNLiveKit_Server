package pili

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
)

const moduleName = "pili"

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

	InitService(conf)
	return nil
}
