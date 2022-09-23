package cache

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
