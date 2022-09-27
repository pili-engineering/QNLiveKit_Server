package module

import (
	"github.com/qbox/livekit/core/application"
)

const moduleName = "item"

type Module struct {
	application.EmptyModule
}

func (m *Module) PreStart() error {
	return nil
}
