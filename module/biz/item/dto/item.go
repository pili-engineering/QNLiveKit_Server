// @Author: wangsheng
// @Description:
// @File:  item
// @Version: 1.0.0
// @Date: 2022/7/1 5:59 下午
// Copyright 2021 QINIU. All rights reserved

package dto

import (
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/utils/timestamp"
)

type ItemDto struct {
	LiveId       string        `json:"live_id"`       //直播间ID
	ItemId       string        `json:"item_id"`       //商品id
	Order        uint          `json:"order"`         //商品排序
	Title        string        `json:"title"`         //商品标题
	Tags         string        `json:"tags"`          //商品标签
	Thumbnail    string        `json:"thumbnail"`     //商品缩略图
	Link         string        `json:"link"`          //商品链接
	CurrentPrice string        `json:"current_price"` //商品当前售价
	OriginPrice  string        `json:"origin_price"`  //商品原始售价(划线价)
	Status       uint          `json:"status"`        //商品状态
	Record       *RecordDto    `json:"record"`        //商品讲解回放
	Extends      model.Extends `json:"extends"`       //扩展属性
}

func ItemDtoToEntity(d *ItemDto) *model.ItemEntity {
	if d == nil {
		return nil
	}

	e := &model.ItemEntity{
		LiveId:       d.LiveId,
		ItemId:       d.ItemId,
		Title:        d.Title,
		Tags:         d.Tags,
		Thumbnail:    d.Thumbnail,
		Link:         d.Link,
		CurrentPrice: d.CurrentPrice,
		OriginPrice:  d.OriginPrice,
		Status:       d.Status,
		Extends:      d.Extends,
	}
	if d.Record != nil {
		e.RecordId = d.Record.ID
	}
	return e
}

type RecordDto struct {
	ID        uint                `json:"id"`
	RecordUrl string              `json:"record_url"`
	Start     timestamp.Timestamp `json:"start"`
	End       timestamp.Timestamp `json:"end"`
	Status    uint                `json:"status"`
	LiveId    string              `json:"live_id"`
	ItemId    string              `json:"item_id"`
}
