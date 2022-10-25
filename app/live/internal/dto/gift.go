package dto

import "github.com/qbox/livekit/biz/model"

type GiftConfigDto struct {
	GiftId        int           `json:"gift_id"`
	Type          int           `json:"type"`
	Name          string        `json:"name"`
	Amount        int           `json:"amount"`
	Img           string        `json:"img"`
	AnimationType int           `json:"animation_type"`
	AnimationImg  string        `json:"animation_img"`
	Order         int           `json:"order"`
	Extends       model.Extends `json:"extends"`
}

func GiftDtoToEntity(d *GiftConfigDto) *model.GiftEntity {
	return &model.GiftEntity{
		GiftId:        d.GiftId,
		Type:          d.Type,
		Name:          d.Name,
		Amount:        d.Amount,
		Img:           d.Img,
		AnimationType: d.AnimationType,
		AnimationImg:  d.AnimationImg,
		Order:         d.Order,
		Extends:       d.Extends,
	}
}

func GiftEntityToDto(d *model.GiftEntity) *GiftConfigDto {
	return &GiftConfigDto{
		GiftId:        d.GiftId,
		Type:          d.Type,
		Name:          d.Name,
		Amount:        d.Amount,
		Img:           d.Img,
		AnimationType: d.AnimationType,
		AnimationImg:  d.AnimationImg,
		Order:         d.Order,
		Extends:       d.Extends,
	}
}
