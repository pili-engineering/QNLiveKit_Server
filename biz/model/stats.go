package model

import (
	"github.com/qbox/livekit/utils/timestamp"
)

type StatsSingleLiveEntity struct {
	ID        uint   `gorm:"primary_key"`
	LiveId    string `json:"live_id"`
	UserId    string `json:"user_id"` //应用内用户ID
	BizId     string `json:"biz_id"`
	Count     int    `json:"count"`
	Type      int    `json:"type"`
	UpdatedAt timestamp.Timestamp
}

func (StatsSingleLiveEntity) TableName() string {
	return "stats_single_live"
}

func (e *StatsSingleLiveEntity) SaveSql() (string, []interface{}) {
	sql := "INSERT INTO stats_single_live(`live_id`, `user_id`, `type`, `biz_id`, `count`, `updated_at`) values(?, ?, ?, ?, ?, ?) " +
		" ON DUPLICATE KEY UPDATE count = ?, updated_at = ?"
	params := []interface{}{
		e.LiveId, e.UserId, e.Type, e.BizId, e.Count, e.UpdatedAt,
		e.Count, e.UpdatedAt,
	}
	return sql, params
}

const (
	StatsTypeLive    = 1 //观看直播
	StatsTypeItem    = 2 //查看商品
	StatsTypeComment = 3 //评论
	StatsTypeLike    = 4 //点赞
)

const (
	StatsTypeDescLive    = "Live"
	StatsTypeDescItem    = "Item"
	StatsTypeDescComment = "Comment"
	StatsTypeDescLike    = "Like"
)

var StatsTypeDescription map[int]string = map[int]string{
	StatsTypeLive:    StatsTypeDescLive,
	StatsTypeItem:    StatsTypeDescItem,
	StatsTypeComment: StatsTypeDescComment,
	StatsTypeLike:    StatsTypeDescLike,
}
