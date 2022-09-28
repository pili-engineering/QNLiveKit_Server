package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/auth/internal/impl"
	"github.com/qbox/livekit/module/base/auth/service"
)

const moduleName = "auth"

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	conf := auth.Config{}

	if c != nil {
		if err := c.Unmarshal(&conf); err != nil {
			return err
		}
	}

	service.Instance = impl.NewService(conf)
	return nil
}

func (m *Module) PreStart() error {
	service.Instance.RegisterAuthMiddleware()
	return nil
}
