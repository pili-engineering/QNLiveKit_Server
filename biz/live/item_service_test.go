// @Author: wangsheng
// @Description:
// @File:  item_service_test.go
// @Version: 1.0.0
// @Date: 2022/7/6 5:40 下午
// Copyright 2021 QINIU. All rights reserved

package live

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/mysql"
)

const testLiveId = "test_live_1"

func itemSetup() {
	mysql.Init(&mysql.ConfigStructure{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		//Password: "123456",
		Database: "live_test",
		Default:  "live",
	}, &mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		//Password: "123456",
		Database: "live_test",
		Default:  "live",
		ReadOnly: true,
	})

	mysql.GetLive().AutoMigrate(model.ItemEntity{}, model.ItemDemonstrate{}, model.LiveEntity{})

	liveEntity := model.LiveEntity{
		LiveId: testLiveId,
	}
	mysql.GetLive().Save(&liveEntity)

	items := make([]*model.ItemEntity, 0, 50)
	for i := uint(1); i <= 50; i++ {
		item := &model.ItemEntity{
			ItemId:       fmt.Sprintf("item_%d", i),
			Title:        fmt.Sprintf("item_%d", i),
			Tags:         "test,hi",
			Thumbnail:    fmt.Sprintf("http://thumbnail/%d", i),
			Link:         fmt.Sprintf("http://link/%d", i),
			CurrentPrice: "$199",
			OriginPrice:  "$299",
			Status:       0,
			Extends:      map[string]string{"age": "18"},
		}
		if i%2 == 0 {
			item.Status = 0
		} else {
			item.Status = 1
		}
		items = append(items, item)
	}
	itemService := GetItemService()
	itemService.AddItems(context.Background(), testLiveId, items)

}

func itemTearDown() {
	db := mysql.GetLive()
	db.DropTableIfExists(model.ItemEntity{}, model.ItemDemonstrate{}, model.LiveEntity{})
}

func TestItemService_countItems(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	db := mysql.GetLive()

	tests := []struct {
		name    string
		liveId  string
		want    int
		wantErr bool
	}{
		{
			name:    "test_live_1",
			liveId:  "test_live_1",
			want:    50,
			wantErr: false,
		},
		{
			name:    "test_live_2",
			liveId:  "test_live_2",
			want:    0,
			wantErr: false,
		},
	}

	itemService := &ItemService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := itemService.countItems(db, tt.liveId)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestItemService_checkItemsExist(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	db := mysql.GetLive()

	tests := []struct {
		name    string
		liveId  string
		itemIds []string
		wantErr bool
	}{
		{
			name:    "test_live_1",
			liveId:  "test_live_1",
			itemIds: []string{"item_1", "item_51"},
			wantErr: true,
		},
		{
			name:    "test_live_1",
			liveId:  "test_live_1",
			itemIds: []string{"item_51", "item_52"},
			wantErr: false,
		},
		{
			name:    "test_live_2",
			liveId:  "test_live_2",
			itemIds: []string{"item_1", "item_2"},
			wantErr: false,
		},
	}

	itemService := &ItemService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := itemService.checkItemsExist(db, tt.liveId, tt.itemIds); (err != nil) != tt.wantErr {
				t.Errorf("checkItemsExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestItemService_DelItems(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	itemService := GetItemService()
	err := itemService.DelItems(context.Background(), testLiveId, []string{"item_1", "item_2"})
	assert.Nil(t, err)

	items, err := itemService.ListItems(context.Background(), testLiveId, true)
	assert.Nil(t, err)
	assert.Equal(t, 48, len(items))

	err = itemService.DelItems(context.Background(), testLiveId, []string{"item_1", "item_2"})
	assert.Nil(t, err)

	items, err = itemService.ListItems(context.Background(), testLiveId, true)
	assert.Nil(t, err)
	assert.Equal(t, 48, len(items))
}

func TestItemService_ListItems(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	tests := []struct {
		name        string
		liveId      string
		showOffline bool
		want        int
	}{
		{
			name:        "test_live_1",
			liveId:      testLiveId,
			showOffline: true,
			want:        50,
		},
		{
			name:        "test_live_1",
			liveId:      testLiveId,
			showOffline: false,
			want:        25,
		},
		{
			name:        "test_live_2",
			liveId:      "test_live_2",
			showOffline: false,
			want:        0,
		},
		{
			name:        "test_live_2",
			liveId:      "test_live_2",
			showOffline: true,
			want:        0,
		},
	}

	s := GetItemService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListItems(context.Background(), tt.liveId, tt.showOffline)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, len(got))
		})
	}
}

func TestItemService_UpdateItemStatus(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	s := GetItemService()
	err := s.UpdateItemStatus(context.Background(), "", nil)
	assert.Nil(t, err)

	err = s.UpdateItemStatus(context.Background(), testLiveId, nil)
	assert.Nil(t, err)

	err = s.UpdateItemStatus(context.Background(), testLiveId, []*model.ItemStatus{{
		ItemId: "item_50",
		Status: 2,
	}, {
		ItemId: "item_49",
		Status: 2,
	}})
	assert.Nil(t, err)

	items, err := s.ListItems(context.Background(), testLiveId, false)
	assert.Nil(t, err)
	for _, item := range items {
		if item.ItemId == "item_49" || item.ItemId == "item_50" {
			assert.Equal(t, uint(2), item.Status)
		}
	}

	err = s.UpdateItemStatus(context.Background(), testLiveId, []*model.ItemStatus{{
		ItemId: "item_50",
		Status: 3,
	}})
	assert.NotNil(t, err)
}

func TestItemService_UpdateItemOrder(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	s := &ItemService{}
	orders := []*model.ItemOrder{{"item_1", 3}, {"item_2", 1}, {"item_3", 2}}
	err := s.UpdateItemOrder(context.Background(), testLiveId, orders)
	assert.Nil(t, err)
	item, err := s.GetLiveItem(context.Background(), testLiveId, "item_1")
	assert.NotNil(t, item)
	assert.Equal(t, uint(3), item.Order)

	orders = []*model.ItemOrder{{"item_1", 3}, {"item_2", 1}, {"item_3", 4}}
	err = s.UpdateItemOrder(context.Background(), testLiveId, orders)
	assert.NotNil(t, err)

	orders = []*model.ItemOrder{{"item_1", 3}, {"item_3", 1}, {"item_1", 1}}
	err = s.UpdateItemOrder(context.Background(), testLiveId, orders)
	assert.NotNil(t, err)

	orders = []*model.ItemOrder{{"item_1", 3}, {"item_2", 1}, {"item_3", 1}}
	err = s.UpdateItemOrder(context.Background(), testLiveId, orders)
	assert.NotNil(t, err)
}

func TestItemService_UpdateItemOrderSingle(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	s := &ItemService{}

	err := s.UpdateItemOrderSingle(context.Background(), testLiveId, "item_1", 1, 10)
	assert.Nil(t, err)
	item, _ := s.GetLiveItem(context.Background(), testLiveId, "item_1")
	assert.Equal(t, uint(10), item.Order)
	item, _ = s.GetLiveItem(context.Background(), testLiveId, "item_2")
	assert.Equal(t, uint(1), item.Order)
	item, _ = s.GetLiveItem(context.Background(), testLiveId, "item_10")
	assert.Equal(t, uint(9), item.Order)

	err = s.UpdateItemOrderSingle(context.Background(), testLiveId, "item_50", 50, 41)
	assert.Nil(t, err)
	item, _ = s.GetLiveItem(context.Background(), testLiveId, "item_50")
	assert.Equal(t, uint(41), item.Order)
	item, _ = s.GetLiveItem(context.Background(), testLiveId, "item_41")
	assert.Equal(t, uint(42), item.Order)
	item, _ = s.GetLiveItem(context.Background(), testLiveId, "item_49")
	assert.Equal(t, uint(50), item.Order)

}

func TestItemService_SetDemonstrateItem(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	s := &ItemService{}
	item, err := s.GetDemonstrateItem(context.Background(), testLiveId)
	assert.Nil(t, err)
	assert.Nil(t, item)

	err = s.SetDemonstrateItem(context.Background(), testLiveId, "item_1")
	assert.Nil(t, err)

	item, err = s.GetDemonstrateItem(context.Background(), testLiveId)
	assert.Nil(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "item_1", item.ItemId)

	err = s.SetDemonstrateItem(context.Background(), testLiveId, "item_2")
	assert.Nil(t, err)

	item, err = s.GetDemonstrateItem(context.Background(), testLiveId)
	assert.Nil(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "item_2", item.ItemId)

	err = s.DelDemonstrateItem(context.Background(), testLiveId)
	assert.Nil(t, err)

	item, err = s.GetDemonstrateItem(context.Background(), testLiveId)
	assert.Nil(t, err)
	assert.Nil(t, item)
}

func TestItemService_getLiveItem(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	tests := []struct {
		name   string
		liveId string
		itemId string
		want   string
	}{
		{
			name:   "get",
			liveId: testLiveId,
			itemId: "item_1",
			want:   "item_1",
		},
		{
			name:   "get 2",
			liveId: testLiveId,
			itemId: "item_50",
			want:   "item_50",
		},
		{
			name:   "item id",
			liveId: testLiveId,
			itemId: "item_51",
			want:   "",
		},
		{
			name:   "live id",
			liveId: "live_1",
			itemId: "item_1",
			want:   "",
		},
	}

	s := &ItemService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetLiveItem(context.Background(), tt.liveId, tt.itemId)
			assert.Nil(t, err)
			if tt.want == "" {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, tt.itemId)
			}
		})
	}
}

func TestItemService_UpdateItemInfo(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	tests := []struct {
		name    string
		liveId  string
		item    *model.ItemEntity
		updated bool
		wantErr bool
	}{
		{
			name:    "not exit",
			liveId:  "test_live_no",
			item:    &model.ItemEntity{ItemId: "item_1", Tags: "new", Extends: map[string]string{"name": "item_1"}},
			updated: false,
			wantErr: true,
		},
		{
			name:    "not exit",
			liveId:  testLiveId,
			item:    &model.ItemEntity{ItemId: "item_100", Tags: "new", Extends: map[string]string{"name": "item_1"}},
			updated: false,
			wantErr: true,
		},
		{
			name:    "not updated",
			liveId:  testLiveId,
			item:    &model.ItemEntity{ItemId: "item_1"},
			updated: false,
			wantErr: false,
		},
		{
			name:   "updated",
			liveId: testLiveId,
			item: &model.ItemEntity{
				ItemId:       "item_2",
				Title:        "new title",
				Tags:         "new tag",
				Thumbnail:    "new thumbnail",
				Link:         "",
				CurrentPrice: "new p1",
				OriginPrice:  "new p2",
				Extends:      map[string]string{"name": "item_1"},
			},
			updated: true,
			wantErr: false,
		},
	}

	s := &ItemService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdateItemInfo(context.Background(), tt.liveId, tt.item)
			assert.Equal(t, tt.wantErr, err != nil)

			if tt.updated {
				cur, _ := s.GetLiveItem(context.Background(), tt.liveId, tt.item.ItemId)
				assert.Equal(t, tt.item.Title, cur.Title)
				assert.Equal(t, tt.item.Tags, cur.Tags)
				assert.Equal(t, tt.item.Thumbnail, cur.Thumbnail)
				assert.Equal(t, tt.item.CurrentPrice, cur.CurrentPrice)
				assert.Equal(t, tt.item.OriginPrice, cur.OriginPrice)
				assert.Equal(t, tt.item.Extends["name"], cur.Extends["name"])
			}
		})
	}
}

func TestItemService_UpdateItemExtends(t *testing.T) {
	itemSetup()
	defer itemTearDown()

	tests := []struct {
		name    string
		liveId  string
		itemId  string
		extends model.Extends
		updated bool
		wantErr bool
	}{
		{
			name:    "not exit",
			liveId:  "test_live_no",
			itemId:  "item_1",
			extends: map[string]string{"name": "item_1"},
			updated: false,
			wantErr: true,
		},
		{
			name:    "not exit",
			liveId:  testLiveId,
			itemId:  "item_100",
			extends: map[string]string{"name": "item_1"},
			updated: false,
			wantErr: true,
		},
		{
			name:    "not updated",
			liveId:  testLiveId,
			itemId:  "item_1",
			extends: map[string]string{},
			updated: false,
			wantErr: false,
		},
		{
			name:    "updated",
			liveId:  testLiveId,
			itemId:  "item_1",
			extends: map[string]string{"name": "item_1"},
			updated: true,
			wantErr: false,
		},
	}
	s := &ItemService{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdateItemExtends(context.Background(), tt.liveId, tt.itemId, tt.extends)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.updated {
				cur, _ := s.GetLiveItem(context.Background(), tt.liveId, tt.itemId)
				assert.Equal(t, tt.extends["name"], cur.Extends["name"])
			}
		})
	}
}
