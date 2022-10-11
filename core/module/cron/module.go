package cron

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
)

const ModuleName = "cron"

func init() {
	application.RegisterModule(ModuleName, &Module{})
}

var _ application.Module = &Module{}

type Module struct {
	application.EmptyModule
}

func (m *Module) Config(c *config.Config) error {
	singleTaskNode := int64(0)
	if c != nil {
		conf := &Config{}
		if err := c.Unmarshal(&conf); err != nil {
			return err
		}
		singleTaskNode = conf.SingleTaskNode
	}

	instance = newService(singleTaskNode)
	m.SetConfigSuccess()
	return nil
}

func (m *Module) Start() error {
	instance.StartCron()
	return nil
}

func (m *Module) Stop() error {
	return instance.StopCron()
}
