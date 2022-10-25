package model

type LiveLikeFlush struct {
	Id             uint  `gorm:"primary_key" json:"id"`
	LastUpdateTime int64 `json:"last_update_time"`
}

func (LiveLikeFlush) TableName() string {
	return "live_like_flush"
}
