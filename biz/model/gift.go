package model

import "github.com/qbox/livekit/utils/timestamp"

type GiftEntity struct {
	ID            uint                 `gorm:"primary_key"`
	Type          int                  `json:"type"`
	Name          string               `json:"name"`
	Amount        int                  `json:"amount"`
	Img           string               `json:"img"`
	AnimationType int                  `json:"animation_type"`
	AnimationImg  string               `json:"animation_img"`
	Order         int                  `json:"order"`
	CreatedAt     timestamp.Timestamp  `json:"created_at"`
	UpdatedAt     timestamp.Timestamp  `json:"updated_at"`
	DeletedAt     *timestamp.Timestamp `json:"deleted_at"`
	Extends       Extends              `json:"extends" gorm:"type:varchar(512)"`
}

func (e GiftEntity) TableName() string {
	return "gift_config"
}
