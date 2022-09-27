package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/module/base/user/internal/controller/client"
	"github.com/qbox/livekit/module/base/user/internal/controller/server"
)

const moduleName = "user"

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
