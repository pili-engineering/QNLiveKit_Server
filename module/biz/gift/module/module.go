package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	giftConfig "github.com/qbox/livekit/module/biz/gift/config"
	"github.com/qbox/livekit/module/biz/gift/internal/controller/admin"
	"github.com/qbox/livekit/module/biz/gift/internal/controller/client"
	"github.com/qbox/livekit/module/biz/gift/internal/controller/server"
	"github.com/qbox/livekit/module/biz/gift/internal/impl"
)

const moduleName = "gift"

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

	conf := giftConfig.Config{}
	if err := c.Unmarshal(&conf); err != nil {
		return err
	}

	if err := conf.Validate(); err != nil {
		return nil
	}

	impl.ConfigService(conf)

	return nil
}

func (m *Module) PreStart() error {
	if impl.GetInstance() == nil {
		return nil
	}

	admin.RegisterRoute()
	server.RegisterRoutes()
	client.RegisterRoutes()

	return nil
}
