package application

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qbox/livekit/core/config"
)

func TestModuleManager_RegisterModule(t *testing.T) {
	mm := &ModuleManager{
		modules: map[string]Module{},
		status:  statusInit,
	}

	err := mm.RegisterModule("a", &EmptyModule{})
	assert.Nil(t, err)

	err = mm.RegisterModule("b", &EmptyModule{})
	assert.Nil(t, err)

	err = mm.RegisterModule("a", &EmptyModule{})
	assert.NotNil(t, err)

	mm.Start()
	err = mm.RegisterModule("c", &EmptyModule{})
	assert.NotNil(t, err)
}

func TestModuleManager_StartSuccess(t *testing.T) {
	mm := &ModuleManager{
		modules: map[string]Module{},
		status:  statusInit,
	}

	moduleManager.RegisterModule("a", &EmptyModule{})
	moduleManager.RegisterModule("b", &EmptyModule{})
	moduleManager.RegisterModule("c", &EmptyModule{})

	err := mm.Start()
	assert.Nil(t, err)
}

func TestModuleManager_StartConfigError(t *testing.T) {
	mm := &ModuleManager{
		modules: map[string]Module{},
		status:  statusInit,
	}

	moduleManager.RegisterModule("a", &EmptyModule{})
	moduleManager.RegisterModule("b", &EmptyModule{})
	moduleManager.RegisterModule("c", &ConfigErrorModule{})

	err := mm.Start()
	assert.Nil(t, err)
}

func TestModuleManager_StartStartError(t *testing.T) {
	mm := &ModuleManager{
		modules: map[string]Module{},
		status:  statusInit,
	}

	moduleManager.RegisterModule("a", &EmptyModule{})
	moduleManager.RegisterModule("b", &EmptyModule{})
	moduleManager.RegisterModule("c", &StartErrorModule{})

	err := mm.Start()
	assert.Nil(t, err)
}

type ConfigErrorModule struct {
	EmptyModule
}

func (m *ConfigErrorModule) Config(c *config.Config) error {
	return fmt.Errorf("config error")
}

type StartErrorModule struct {
	EmptyModule
}

func (m *StartErrorModule) Start() error {
	return fmt.Errorf("start error")
}
