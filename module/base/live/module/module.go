package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/module/base/live/internal/controller/client"
	"github.com/qbox/livekit/module/base/live/internal/controller/server"
	"github.com/qbox/livekit/module/base/live/internal/cron"
	"github.com/qbox/livekit/module/base/live/internal/impl"
	"github.com/qbox/livekit/module/base/live/service"
)

const moduleName = "live"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	service.Instance = &impl.Service{}
	return nil
}

func (m *Module) PreStart() error {
	client.RegisterRoutes()
	server.RegisterRoutes()

	cron.RegisterCrons()
	return nil
}
