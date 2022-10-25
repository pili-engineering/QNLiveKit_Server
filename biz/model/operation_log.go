package model

import "github.com/qbox/livekit/utils/timestamp"

type OperationLog struct {
	ID         uint                `json:"id" gorm:"primary_key"`
	UserId     string              `json:"user_id"`     // 操作用户UID
	IP         string              `json:"ip" `         //客户端IP地址
	Method     string              `json:"method" `     // 操作方法 PUT|DELETE|POST
	URL        string              `json:"url" `        // 操作路径
	Args       string              `json:"args" `       // 操作内容详情
	StatusCode int                 `json:"status_code"` //返回状态码
	CreatedAt  timestamp.Timestamp `json:"created_at"`  // 日志时间
}

func (OperationLog) TableName() string {
	return "operation_log"
}
