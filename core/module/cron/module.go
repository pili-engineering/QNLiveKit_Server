package cron

import (
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/config"
)

const ModuleName = "cron"

var _ application.Module = &Module{}

type Module struct {
}

func (m *Module) Config(c *config.Config) error {
	return nil
}

func (m *Module) PreStart() error {
	return nil
}

func (m *Module) Start() error {
	return nil
}

func (m *Module) Stop() error {
	return nil
}
