package application

import (
	"github.com/qbox/livekit/core/config"
)

type Module interface {
	// Config 完成模块的配置工作
	// 具体的模块，应该根据配置信息，完成实例的创建工作
	// Config 阶段，模块不应该使用到其他的Module
	Config(c *config.Config) error

	// PreStart 预启动阶段，可以执行：
	// 载入数据等操作
	// 向其他module 注册自己信息
	PreStart() error

	// Start 对于一些需要工作Loop 的module，在这里启动
	Start() error

	// Stop 停止module 比如
	// 停止所有 WorkLoop
	// 释放资源
	// io/cache flush
	Stop() error

	// IsConfigSuccess 模块是否已经完成配置
	IsConfigSuccess() bool

	// RequireModules 返回本模块依赖的其他模块列表
	RequireModules() []string
}

var _ Module = &EmptyModule{}

type EmptyModule struct {
	configSuccess bool
}

func (m *EmptyModule) Config(c *config.Config) error {
	m.configSuccess = true
	return nil
}

func (m *EmptyModule) PreStart() error {
	return nil
}

func (m *EmptyModule) Start() error {
	return nil
}

func (m *EmptyModule) Stop() error {
	return nil
}

func (m *EmptyModule) IsConfigSuccess() bool {
	return m.configSuccess
}

func (m *EmptyModule) RequireModules() []string {
	return nil
}

func (m *EmptyModule) SetConfigSuccess() {
	m.configSuccess = true
}
