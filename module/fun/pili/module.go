package pili

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/core/module/trace"
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
	m.SetConfigSuccess()
	return nil
}

func (m *Module) PreStart() error {
	if service == nil {
		return nil
	}
	trace.SetPiliHub(service.PiliHub())

	return service.setupAccount()
}

func (m *Module) RequireModules() []string {
	return []string{"account", "trace", "appinfo"}
}
