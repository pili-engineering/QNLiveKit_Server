package cache

import (
	"fmt"
)

const (
	maxRetries = 3
	idleConns  = 8
	poolSize   = 100

	// TypeCluster 集群模式
	TypeCluster = "cluster"
	// TypeNode 单机模式
	TypeNode = "node"
)

type Config struct {
	Type     string   `mapstructure:"type"`
	Addr     string   `mapstructure:"addr"`
	Addrs    []string `mapstructure:"addrs"`
	Password string   `mapstructure:"password"`
}

func (c *Config) Validate() error {
	switch c.Type {
	case TypeCluster:
		if len(c.Addrs) == 0 {
			return fmt.Errorf("type %s with empty addrs", c.Type)
		}
	case TypeNode:
		if len(c.Addr) == 0 {
			return fmt.Errorf("type %s with empty addr", c.Type)
		}
	default:
		return fmt.Errorf("invalid type %s", c.Type)
	}

	return nil
}
