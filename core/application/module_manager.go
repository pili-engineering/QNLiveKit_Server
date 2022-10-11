package application

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/utils/logger"
)

var moduleManager = &ModuleManager{
	modules: map[string]Module{},
	status:  statusInit,
}

const (
	statusInit     = int32(0)
	statusConfig   = int32(1)
	statusPreStart = int32(2)
	statusStart    = int32(3)
	statusStop     = int32(4)
)

type ModuleManager struct {
	c          *config.Config
	status     int32
	modules    map[string]Module
	moduleLock sync.RWMutex

	startOnce sync.Once
	stopOnce  sync.Once
}

func (m *ModuleManager) Start() error {
	var err error
	m.startOnce.Do(func() {
		err = m.configAllModules()
		if err != nil {
			return
		}

		err = m.preStartAllModules()
		if err != nil {
			return
		}

		err = m.startAllModules()
		if err != nil {
			return
		}
	})

	return err
}

func (m *ModuleManager) configAllModules() error {
	log := logger.New()
	atomic.StoreInt32(&m.status, statusConfig)
	m.moduleLock.RLock()
	defer m.moduleLock.RUnlock()

	for name, module := range m.modules {
		var c *config.Config = nil
		if m.c != nil {
			c = m.c.Sub(name)
		}
		if err := module.Config(c); err != nil {
			log.Errorf("config module %s error %v", name, err)
			return fmt.Errorf("config module %s error %v", name, err)
		}
	}

	for name, module := range m.modules {
		if !module.IsConfigSuccess() {
			continue
		}

		if len(module.RequireModules()) == 0 {
			continue
		}

		for _, rName := range module.RequireModules() {
			rModule := m.modules[rName]
			if rModule == nil || !rModule.IsConfigSuccess() {
				log.Errorf("module %s not config success, required by %s", rName, name)
				return fmt.Errorf("module %s not config success, required by %s", rName, name)
			}
		}
	}

	return nil
}

func (m *ModuleManager) preStartAllModules() error {
	log := logger.New()
	atomic.StoreInt32(&m.status, statusConfig)
	m.moduleLock.RLock()
	defer m.moduleLock.RUnlock()

	for name, module := range m.modules {
		if err := module.PreStart(); err != nil {
			log.Errorf("pre start module %s error %v", name, err)
			return fmt.Errorf("pre start module %s error %v", name, err)
		}
	}

	return nil
}

func (m *ModuleManager) startAllModules() error {
	log := logger.New()
	atomic.StoreInt32(&m.status, statusStart)

	m.moduleLock.RLock()
	defer m.moduleLock.RUnlock()

	for name, module := range m.modules {
		if err := module.Start(); err != nil {
			log.Errorf("start module %s error %v", name, err)
			return fmt.Errorf("start module %s error %v", name, err)
		}
	}

	return nil
}

func (m *ModuleManager) Stop(err error) {
	m.stopOnce.Do(func() {
		atomic.StoreInt32(&m.status, statusStop)
	})
}

func (m *ModuleManager) RegisterModule(name string, module Module) error {
	if s := atomic.LoadInt32(&m.status); s != statusInit {
		return fmt.Errorf("cant register module on status %d", s)
	}

	m.moduleLock.Lock()
	defer m.moduleLock.Unlock()

	if _, exist := m.modules[name]; exist {
		return fmt.Errorf("duplicated module %s", name)
	}

	m.modules[name] = module
	return nil
}
