package uuid

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// Init 初始化 snowflake node
func Init(nodeID int64) {
	var err error
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		panic(err)
	}
}

// Gen 生成 uuid
func Gen() string {
	return node.Generate().String()
}

// GetTimeFromUUID 从 UUID 解析时间, 注意，如果传入错误的 ID, 可能会解析失败并返回空的时间
func GetTimeFromUUID(id string) time.Time {
	flakeID, err := snowflake.ParseString(id)
	if err == nil {
		tInt64 := flakeID.Time()
		t := time.Unix(tInt64/1000, 0)
		return t
	}

	return time.Time{}
}
