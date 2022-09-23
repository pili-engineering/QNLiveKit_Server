package cache

import (
	"github.com/qbox/livekit/core/application"
)

const moduleName = "cache"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}
