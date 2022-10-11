package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/module/biz/item/internal/controller/client"
	"github.com/qbox/livekit/module/biz/item/internal/controller/server"
	"github.com/qbox/livekit/module/biz/item/internal/impl"
	"github.com/qbox/livekit/module/biz/item/service"
)

const moduleName = "item"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	service.Instance = &impl.ItemService{}
	m.SetConfigSuccess()
	return nil
}

func (m *Module) PreStart() error {
	client.RegisterRoutes()
	server.RegisterRoutes()
	return nil
}

func (m *Module) RequireModules() []string {
	return []string{"live"}
}
