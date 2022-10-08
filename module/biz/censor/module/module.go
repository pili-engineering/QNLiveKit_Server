package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	censorConfig "github.com/qbox/livekit/module/biz/censor/config"
	"github.com/qbox/livekit/module/biz/censor/internal/controller/admin"
	"github.com/qbox/livekit/module/biz/censor/internal/impl"
)

const moduleName = "censor"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	if c == nil {
		return nil
	}

	conf := censorConfig.Config{}
	if err := c.Unmarshal(&conf); err != nil {
		return nil
	}

	if err := conf.Validate(); err != nil {
		return err
	}

	impl.ConfigCensorService(&conf)

	return nil
}

func (m *Module) PreStart() error {
	if impl.GetInstance() == nil {
		return nil
	}

	admin.RegisterRoutes()
	return impl.GetInstance().PreStart()
}
