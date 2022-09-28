// @Author: wangsheng
// @Description:
// @File:  item
// @Version: 1.0.0
// @Date: 2022/7/5 10:19 上午
// Copyright 2021 QINIU. All rights reserved

package model

import (
	"github.com/qbox/livekit/utils/timestamp"
)

const ItemStatusOffline = 0 //下架
const ItemStatusOnline = 1  //上架
const ItemStatusLocked = 2  //锁定，不能购买

type ItemEntity struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt timestamp.Timestamp
	UpdatedAt timestamp.Timestamp
	DeletedAt *timestamp.Timestamp `sql:"index"`

	LiveId       string  `json:"live_id"`
	ItemId       string  `json:"item_id"`
	Order        uint    `json:"order"`
	Title        string  `json:"title"`
	Tags         string  `json:"tags"`
	Thumbnail    string  `json:"thumbnail"`
	Link         string  `json:"link"`
	CurrentPrice string  `json:"current_price"`
	OriginPrice  string  `json:"origin_price"`
	Status       uint    `json:"status"`
	RecordId     uint    `json:"record_id"`
	Extends      Extends `json:"extends" gorm:"type:varchar(512)"`
}

func (e ItemEntity) TableName() string {
	return "items"
}

func (e *ItemEntity) IsValid() bool {
	if len(e.ItemId) == 0 || len(e.Title) == 0 || len(e.CurrentPrice) == 0 {
		return false
	}

	if e.Status > 2 {
		return false
	}

	return true
}

type ItemStatus struct {
	ItemId string `json:"item_id"`
	Status uint   `json:"status"`
}

type ItemOrder struct {
	ItemId string `json:"item_id"`
	Order  uint   `json:"order"`
}

type ItemDemonstrate struct {
	ID        uint   `gorm:"primary_key"`
	LiveId    string `json:"live_id" gorm:"unique_index"`
	ItemId    string `json:"item_id"`
	UpdatedAt timestamp.Timestamp
}

func (ItemDemonstrate) TableName() string {
	return "item_demonstrate"
}

type ItemDemonstrateRecord struct {
	ID           uint `gorm:"primary_key"`
	Start        timestamp.Timestamp
	End          timestamp.Timestamp
	Status       uint   `json:"status"`
	Format       uint   `json:"format"`
	LiveId       string `json:"live_id"`
	ItemId       string `json:"item_id"`
	ExpireDays   uint   `json:"expireDays"`
	Fname        string `json:"fname"`
	PersistentID string `json:"persistentID"`
}

func (ItemDemonstrateRecord) TableName() string {
	return "item_demonstrate_log"
}

//状态码0成功，1等待处理，2正在处理，3处理失败。
const (
	RecordStatusSuccess = iota
	RecordStatusWait
	RecordStatusProcessing
	RecordStatusFail
	RecordStatusDefault
)
