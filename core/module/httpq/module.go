package httpq

import (
	"fmt"

	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
)

const ModuleName = "httpq"

func init() {
	application.RegisterModule(ModuleName, &Module{})
}

var _ application.Module = &Module{}

type Module struct {
}

func (m *Module) Config(c *config.Config) error {
	conf := Config{}
	if err := c.Unmarshal(&conf); err != nil {
		return fmt.Errorf("unmarshal config error %v", err)
	}

	if err := conf.Validate(); err != nil {
		return fmt.Errorf("validate config error %v", err)
	}

	instance = newServer(&conf)

	return nil
}

func (m *Module) PreStart() error {
	return nil
}

func (m *Module) Start() error {
	return nil
}

func (m *Module) Stop() error {
	return nil
}
