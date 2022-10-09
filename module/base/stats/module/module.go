package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/module/base/stats/internal/controller/client"
	"github.com/qbox/livekit/module/base/stats/internal/cron"
	"github.com/qbox/livekit/module/base/stats/internal/impl"
	"github.com/qbox/livekit/module/base/stats/service"
)

const moduleName = "stats"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	service.Instance = impl.NewServiceImpl()
	return nil
}

func (m *Module) PreStart() error {
	client.RegisterRoutes()

	cron.RegisterCrons()
	return nil
}
