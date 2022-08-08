// @Author: wangsheng
// @Description:
// @File:  item
// @Version: 1.0.0
// @Date: 2022/7/1 5:59 下午
// Copyright 2021 QINIU. All rights reserved

package dto

import (
	"github.com/qbox/livekit/biz/model"
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
	Extends      model.Extends `json:"extends"`       //扩展属性
}

func ItemDtoToEntity(d *ItemDto) *model.ItemEntity {
	if d == nil {
		return nil
	}

	return &model.ItemEntity{
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
}

func ItemEntityToDto(e *model.ItemEntity) *ItemDto {
	if e == nil {
		return nil
	}

	return &ItemDto{
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
}
