package application

import (
	"github.com/qbox/livekit/core/config"
)

type Module interface {
	Config(c *config.Config) error
	Start() error
	Stop() error
}

var _ Module = &EmptyModule{}

type EmptyModule struct {
}

func (m *EmptyModule) Config(c *config.Config) error {
	return nil
}

func (m *EmptyModule) Start() error {
	return nil
}

func (m *EmptyModule) Stop() error {
	return nil
}
