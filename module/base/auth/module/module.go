package module

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/auth/internal/controller/admin"
	"github.com/qbox/livekit/module/base/auth/internal/controller/server"
	"github.com/qbox/livekit/module/base/auth/internal/impl"
)

const moduleName = "auth"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

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

	impl.ConfigService(conf)
	m.SetConfigSuccess()
	return nil
}

func (m *Module) PreStart() error {
	impl.GetInstance().RegisterAuthMiddleware()
	admin.RegisterRoutes()
	server.RegisterRoutes()
	return nil
}

func (m *Module) RequireModules() []string {
	return []string{"httpq", "user", "admin"}
}
