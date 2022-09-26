package account

import (
	"fmt"

	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
)

const moduleName = "account"

func init() {
	application.RegisterModule(moduleName, &Module{})
}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	if c == nil {
		return fmt.Errorf("empty account info")
	}

	if err := c.Unmarshal(&defaultConfig); err != nil {
		return err
	}

	return defaultConfig.Validate()
}
