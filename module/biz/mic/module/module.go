package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/module/biz/mic/internal/controller/client"
)

const moduleName = "mic"

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
	return nil
}
