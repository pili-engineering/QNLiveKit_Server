package utils

import (
	"context"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/module/biz/item/dto"
	"github.com/qbox/livekit/module/biz/item/internal/impl"
	"github.com/qbox/livekit/module/fun/pili"
)

func ItemEntityToDto(e *model.ItemEntity) *dto.ItemDto {
	if e == nil {
		return nil
	}

	itemDto := &dto.ItemDto{
		LiveId:       e.LiveId,
		ItemId:       e.ItemId,
		Order:        e.Order,
		Title:        e.Title,
		Tags:         e.Tags,
		Thumbnail:    e.Thumbnail,
		Link:         e.Link,
		CurrentPrice: e.CurrentPrice,
		OriginPrice:  e.OriginPrice,
		Status:       e.Status,
		Extends:      e.Extends,
	}

	if e.RecordId > 0 {
		record, _ := impl.GetInstance().GetRecordVideo(context.Background(), e.RecordId)
		if record != nil {
			itemDto.Record = RecordEntityToDto(record)
		}
	}

	return itemDto
}

func RecordEntityToDto(e *model.ItemDemonstrateRecord) *dto.RecordDto {
	if e == nil {
		return nil
	}

	r := &dto.RecordDto{
		ID:        e.ID,
		RecordUrl: e.Fname,
		Start:     e.Start,
		End:       e.End,
		LiveId:    e.LiveId,
		Status:    e.Status,
		ItemId:    e.ItemId,
	}

	if e.Status == 0 {
		//r.RecordUrl = config.AppConfig.RtcConfig.RtcPlayBackUrl + "/" + e.Fname
		r.RecordUrl = pili.GetService().PlaybackURL(e.Fname)
	}
	return r
}
