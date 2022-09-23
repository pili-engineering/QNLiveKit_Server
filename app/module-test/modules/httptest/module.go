package httptest

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
)

const moduleName = "httptest"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	instance = &Server{}
	return nil
}

func (m *Module) PreStart() error {
	instance.RegisterApi()
	return nil
}
