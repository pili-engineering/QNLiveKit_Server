package dto

import (
	"github.com/qbox/livekit/biz/model"
)

type SingleLiveInfo struct {
	LiveId string `json:"live_id"`
	UserId string `json:"user_id"`
	ItemId string `json:"item_id"`
	Count  int    `json:"count"`
	Type   int    `json:"type"`
}

func StatsSLDtoToEntity(d *SingleLiveInfo) *model.StatsSingleLiveEntity {
	if d == nil {
		return nil
	}

	return &model.StatsSingleLiveEntity{
		LiveId: d.LiveId,
		ItemId: d.ItemId,
		UserId: d.UserId,
		Type:   d.Type,
		Count:  d.Count,
	}
}
