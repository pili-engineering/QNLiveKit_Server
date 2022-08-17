package model

import (
	"github.com/qbox/livekit/utils/timestamp"
)

type StatsSingleLiveEntity struct {
	ID        uint   `gorm:"primary_key"`
	LiveId    string `json:"live_id"`
	UserId    string `json:"user_id"` //应用内用户ID
	ItemId    string `json:"item_id"`
	Count     int    `json:"count"`
	Type      int    `json:"type"`
	UpdatedAt timestamp.Timestamp
}

func (StatsSingleLiveEntity) TableName() string {
	return "stats_single_live"
}
