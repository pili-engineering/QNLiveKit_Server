package appinfo

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/module/appinfo/internal/controller/client"
)

const moduleName = "appinfo"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) PreStart() error {
	client.RegisterRoutes()
	return nil
}
