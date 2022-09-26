package uuid

import (
	"github.com/qbox/livekit/core/application"
	mNode "github.com/qbox/livekit/core/module/node"
)

const moduleName = "module"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) PreStart() error {
	Init(mNode.NodeId())
	return nil
}
