package node

var nodeInfo = Config{
	NodeId: 0,
}

type Config struct {
	NodeId int64 `mapstructure:"node_id"`
}
