package trace

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/module/account"
)

const moduleName = "trace"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) PreStart() error {
	instance.client = newRpcClient(account.AccessKey(), account.SecretKey())
	return nil
}
