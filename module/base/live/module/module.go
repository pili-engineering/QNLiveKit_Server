package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
)

const moduleName = "live"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	return nil
}

func (m *Module) PreStart() error {
	client.RegisterRoutes()
	server.RegisterRoutes()

	return nil
}
