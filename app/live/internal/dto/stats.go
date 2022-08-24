package dto

import (
	"github.com/qbox/livekit/biz/model"
)

type SingleLiveInfo struct {
	LiveId string `json:"live_id"`
	UserId string `json:"user_id"`
	BizId  string `json:"biz_id"`
	Count  int    `json:"count"`
	Type   int    `json:"type"`
}

func StatsSLDtoToEntity(d *SingleLiveInfo) *model.StatsSingleLiveEntity {
	if d == nil {
		return nil
	}

	return &model.StatsSingleLiveEntity{
		LiveId: d.LiveId,
		BizId:  d.BizId,
		UserId: d.UserId,
		Type:   d.Type,
		Count:  d.Count,
	}
}
