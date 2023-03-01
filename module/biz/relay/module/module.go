package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/module/biz/relay/internal/controller/client"
	"github.com/qbox/livekit/module/biz/relay/internal/controller/server"
)

const moduleName = "relay"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) PreStart() error {
	client.RegisterRoutes()
	server.RegisterRoutes()
	return nil
}

func (m *Module) RequireModules() []string {
	return []string{"mysql", "rtc", "live"}
}
