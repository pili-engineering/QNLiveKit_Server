package admin

import (
	"github.com/qbox/livekit/core/application"
)

const moduleName = "admin"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) RequreModules() []string {
	return []string{"mysql"}
}
