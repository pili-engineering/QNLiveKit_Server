package cache

import (
	log "github.com/sirupsen/logrus"

	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
)

const moduleName = "cache"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	if c == nil {
		return Nil
	}

	conf := Config{}
	if err := c.Unmarshal(&conf); err != nil {
		log.Error("Unmarshal config error %v", err)
		return err
	}

	if err := Init(&conf); err != nil {
		return err
	}

	return nil
}
