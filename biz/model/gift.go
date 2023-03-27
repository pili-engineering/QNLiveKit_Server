package model

import "github.com/qbox/livekit/utils/timestamp"

type GiftEntity struct {
	ID            uint                 `gorm:"primary_key"`
	Type          int                  `json:"type"`
	GiftId        int                  `json:"gift_id"`
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

type LiveGift struct {
	ID        uint                `gorm:"primary_key" json:"id"`
	BizId     string              `json:"biz_id"`
	UserId    string              `json:"user_id"`
	GiftId    int                 `json:"gift_id"`
	Amount    int                 `json:"amount"`
	Status    int                 `json:"status"`
	LiveId    string              `json:"live_id"`
	AnchorId  string              `json:"anchor_id"`
	CreatedAt timestamp.Timestamp `json:"created_at"`
	UpdatedAt timestamp.Timestamp `json:"updated_at"`
}

func (e LiveGift) TableName() string {
	return "live_gift"
}

const (
	SendGiftStatusWait    = iota //刚创建
	SendGiftStatusSuccess        //发送礼物成功
	SendGiftStatusFailure        //发送礼物失败
)

// PkIntegral redis gift相关key
const (
	PkIntegral = "pkIntegral:%v" // 送礼积分key前缀，仅在直播中维护，直播结束会清理
)
