// @Author: wangsheng
// @Description:
// @File:  item_service
// @Version: 1.0.0
// @Date: 2022/7/5 10:15 上午
// Copyright 2021 QINIU. All rights reserved

package live

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qbox/livekit/common/auth/qiniumac"
	"github.com/qbox/livekit/utils/rpc"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

const maxItemCount = 100

type IItemService interface {
	AddItems(ctx context.Context, liveId string, items []*model.ItemEntity) error
	DelItems(ctx context.Context, liveId string, items []string) error
	ListItems(ctx context.Context, liveId string, showOffline bool) ([]*model.ItemEntity, error)
	UpdateItemInfo(ctx context.Context, liveId string, item *model.ItemEntity) error
	UpdateItemExtends(ctx context.Context, liveId string, itemId string, extends model.Extends) error

	UpdateItemStatus(ctx context.Context, liveId string, statuses []*model.ItemStatus) error
	UpdateItemOrder(ctx context.Context, liveId string, orders []*model.ItemOrder) error
	UpdateItemOrderSingle(ctx context.Context, liveId string, itemId string, from, to uint) error

	SetDemonstrateItem(ctx context.Context, liveId string, itemId string) error
	DelDemonstrateItem(ctx context.Context, liveId string) error
	GetDemonstrateItem(ctx context.Context, liveId string) (*model.ItemEntity, error)
	GetLiveItem(ctx context.Context, liveId string, itemId string) (*model.ItemEntity, error)
	IsDemonstrateItem(ctx context.Context, liveId, itemId string) (bool, error)

	StartRecordVideo(ctx context.Context, liveId string, itemId string) error
	StopRecordVideo(ctx context.Context, liveId string, demonId int) (*model.ItemDemonstrateRecord, error)
	getRecordVideo(ctx context.Context, demonId int) (*model.ItemDemonstrateRecord, error)
	GetRecordVideo(ctx context.Context, demonId uint) (*model.ItemDemonstrateRecord, error)
	UpdateRecordVideo(ctx context.Context, itemLog *model.ItemDemonstrateRecord) error
	GetListRecordVideo(ctx context.Context, liveId string, itemId string) (*model.ItemDemonstrateRecord, error)
	GetListLiveRecordVideo(ctx context.Context, liveId string) ([]*model.ItemDemonstrateRecord, error)
	GetPreviousItem(ctx context.Context, liveId string) (*int, error)
	DelRecordVideo(ctx context.Context, liveId string, demonItem []uint) error
	saveRecordVideo(ctx context.Context, liveId, itemId string) error
	UpdateItemRecord(ctx context.Context, demonId uint, liveId string, itemId string) error
	DeleteItemRecord(ctx context.Context, demonId uint, liveId string, itemId string) error
}

type ItemService struct {
	Config
}

var itemService IItemService

type Config struct {
	PiliHub   string
	AccessKey string
	SecretKey string
}

func InitService(conf Config) {
	itemService = &ItemService{
		Config: conf,
	}
}

func GetItemService() IItemService {
	return itemService
}

func (s *ItemService) AddItems(ctx context.Context, liveId string, items []*model.ItemEntity) (err error) {
	log := logger.ReqLogger(ctx)

	if len(liveId) == 0 || len(items) == 0 {
		err = api.ErrInvalidArgument
		return
	}

	_, err = GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("get live %s error %s", liveId, err.Error())
		return api.ErrNotFound
	}

	db := mysql.GetLive(log.ReqID())
	tx := db.Begin()
	defer func() {
		if err != nil {
			log.Errorf("add items error %s", err.Error())
			tx.Rollback()
		} else {
			err = tx.Commit().Error
			if err != nil {
				log.Errorf("commit items error %s", err.Error())
			}
		}
	}()

	currentCount, err := s.countItems(tx, liveId)
	if err != nil {
		log.Errorf("count live items error %s", err.Error())
		return
	}
	if currentCount+len(items) > maxItemCount {
		err = api.ErrCodeLiveItemExceed
		log.Errorf("live room %s item exceed with new items %d", liveId, len(items))
		return
	}

	currentOrder, err := s.maxItemOrder(tx, liveId)
	if err != nil {
		log.Errorf("get max order error %s", err.Error())
		return
	}

	itemIds := make([]string, 0, len(items))
	for _, item := range items {
		if !item.IsValid() {
			log.Errorf("invalid item %+v", item)
			err = api.ErrInvalidArgument
			return
		}
		currentOrder++
		item.Order = currentOrder
		item.LiveId = liveId
		item.CreatedAt = timestamp.Now()
		item.UpdatedAt = timestamp.Now()

		itemIds = append(itemIds, item.ItemId)
	}

	err = s.checkItemsExist(tx, liveId, itemIds)
	if err != nil {
		log.Errorf("check item existed error %s", err.Error())
		return
	}

	err = s.batchAddItems(tx, items)
	return
}

func (s *ItemService) countItems(db *gorm.DB, liveId string) (int, error) {
	count := 0
	err := db.Model(model.ItemEntity{}).Where("live_id = ?", liveId).Count(&count).Error
	if err != nil {
		return 0, api.ErrDatabase
	}
	return count, nil
}

func (s *ItemService) maxItemOrder(db *gorm.DB, liveId string) (uint, error) {
	var x sql.NullInt64
	err := db.Table(model.ItemEntity{}.TableName()).Select("MAX(`order`)").Where("live_id = ?", liveId).Row().Scan(&x)
	if err != nil {
		return 0, api.ErrDatabase
	}
	return uint(x.Int64), nil
}

func (s *ItemService) checkItemsExist(db *gorm.DB, liveId string, itemIds []string) error {
	items := make([]*model.ItemEntity, 0)
	err := db.Where("live_id = ? and item_id in (?)", liveId, itemIds).Find(&items).Error
	if err != nil {
		return api.ErrDatabase
	}

	if len(items) > 0 {
		return api.ErrAlreadyExist
	}

	return nil
}

//gorm 软删除的坑，做一遍硬删除
func (s *ItemService) deleteOldItems(db *gorm.DB, liveId string, itemIds []string) error {
	err := db.Exec("delete from items where live_id = ? and deleted_at is not null and item_id in ?", liveId, itemIds).Error
	if err != nil {
		return api.ErrDatabase
	} else {
		return nil
	}
}

func (s *ItemService) batchAddItems(db *gorm.DB, items []*model.ItemEntity) error {
	for _, item := range items {
		err := db.Save(item).Error
		if err != nil {
			return api.ErrDatabase
		}
	}

	return nil
}

func (s *ItemService) checkIfCanDelItems(ctx context.Context, liveId string, items []string) ([]string, error) {
	log := logger.ReqLogger(ctx)
	var delItem []string
	var demoItem []string
	for _, i := range items {
		isOnline, err := itemService.IsDemonstrateItem(ctx, liveId, i)
		if err != nil {
			log.Errorf("delete items error %s", err.Error())
			return nil, err
		}
		if isOnline == false {
			delItem = append(delItem, i)
		} else {
			demoItem = append(demoItem, i)
		}
	}
	if len(demoItem) != 0 {
		log.Errorf("some items demonstrating cannot delete %s", demoItem)
	}
	return delItem, nil
}

func (s *ItemService) DelItems(ctx context.Context, liveId string, items []string) error {
	log := logger.ReqLogger(ctx)
	items, err := s.checkIfCanDelItems(ctx, liveId, items)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	db := mysql.GetLive(log.ReqID())
	err = db.Delete(model.ItemEntity{}, "live_id = ? and item_id in (?)", liveId, items).Error
	if err != nil {
		log.Errorf("delete live item error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *ItemService) ListItems(ctx context.Context, liveId string, showOffline bool) ([]*model.ItemEntity, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	query := "live_id = ? "
	if !showOffline {
		query += " and status != 0 "
	}

	items := make([]*model.ItemEntity, 0)
	err := db.Order("`order` desc").Find(&items, query, liveId).Error
	if err != nil {
		log.Errorf("find live items error %s", err.Error())
		return nil, api.ErrDatabase
	}

	return items, nil
}

func (s *ItemService) UpdateItemStatus(ctx context.Context, liveId string, items []*model.ItemStatus) (err error) {
	if len(liveId) == 0 || len(items) == 0 {
		return
	}

	for _, item := range items {
		if item.Status > 2 {
			return api.Error("", api.ErrorCodeInvalidArgument, fmt.Sprintf("invalid item status %d", item.Status))
		}
	}

	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
			if err != nil {
				log.Errorf("commit item status error %s", err.Error())
				err = api.ErrDatabase
			}
		}
	}()

	for _, item := range items {
		updates := map[string]interface{}{
			"status":     item.Status,
			"updated_at": timestamp.Now(),
		}
		err = tx.Model(model.ItemEntity{}).Where("live_id = ? and item_id = ?", liveId, item.ItemId).Update(updates).Error
		if err != nil {
			log.Errorf("update item status error %s", err.Error())
			return
		}
	}

	return
}

func (s *ItemService) UpdateItemOrder(ctx context.Context, liveId string, orders []*model.ItemOrder) (err error) {
	log := logger.ReqLogger(ctx)

	orderMap := make(map[uint]bool)
	itemMap := make(map[string]bool)
	for _, order := range orders {
		if _, ok := itemMap[order.ItemId]; ok {
			return api.Error(log.ReqID(), api.ErrorCodeInvalidArgument, fmt.Sprintf("duplicated item id %s", order.ItemId))
		}

		if order.Order <= 0 {
			return api.Error(log.ReqID(), api.ErrorCodeInvalidArgument, fmt.Sprintf("invalid item order %d", order.Order))
		}

		if _, ok := orderMap[order.Order]; ok {
			return api.Error(log.ReqID(), api.ErrorCodeInvalidArgument, fmt.Sprintf("duplicated item order %d", order.Order))
		}
		orderMap[order.Order] = true
		itemMap[order.ItemId] = true
	}

	//检查 order 是不是有冲突
	itemEntities, err := s.ListItems(ctx, liveId, true)
	if err != nil {
		return
	}

	if err = s.checkOrderConflict(itemEntities, orders); err != nil {
		return
	}

	db := mysql.GetLive(log.ReqID())
	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
			if err != nil {
				log.Errorf("item order commit error %s", err.Error())
			}
		}
	}()

	for _, o := range orders {
		updates := map[string]interface{}{
			"order":      o.Order,
			"updated_at": timestamp.Now(),
		}
		err = tx.Model(model.ItemEntity{}).Where("live_id = ? and item_id = ?", liveId, o.ItemId).Update(updates).Error
		if err != nil {
			log.Errorf("update item status error %s", err.Error())
			return
		}
	}

	return nil
}

func (s *ItemService) checkOrderConflict(items []*model.ItemEntity, orders []*model.ItemOrder) error {
	checkItems := make([]*model.ItemEntity, 0, len(items))
	for _, item := range items {
		for _, o := range orders {
			if o.ItemId == item.ItemId {
				goto outer
			}
		}

		checkItems = append(checkItems, item)
	outer:
	}

	for _, o := range orders {
		for _, item := range checkItems {
			if o.Order == item.Order {
				return api.Error("", api.ErrorCodeInvalidArgument, fmt.Sprintf("order %d conflict with %s", item.Order, item.ItemId))
			}
		}
	}

	return nil
}

func (s *ItemService) UpdateItemOrderSingle(ctx context.Context, liveId string, itemId string, from, to uint) (err error) {
	log := logger.ReqLogger(ctx)
	if from == to {
		err = api.Error("", api.ErrorCodeInvalidArgument, "from equal to")
		return
	}

	itemEntities, err := s.ListItems(ctx, liveId, true)
	if err != nil {
		log.Errorf("list item error %s", err.Error())
		return
	}

	affectedItems, err := calcAffectedItems(itemEntities, itemId, from, to)
	if err != nil {
		log.Errorf("calc affected items error %s", err.Error())
		return
	}
	db := mysql.GetLive(log.ReqID())
	tx := db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
			if err != nil {
				log.Errorf("item order commit error %s", err.Error())
			}
		}
	}()

	for _, item := range affectedItems {
		updates := map[string]interface{}{
			"order":      item.Order,
			"updated_at": timestamp.Now(),
		}
		err = tx.Model(model.ItemEntity{}).Where("live_id = ? and item_id = ?", liveId, item.ItemId).Update(updates).Error
		if err != nil {
			log.Errorf("update item status error %s", err.Error())
			return
		}
	}

	return
}

func calcAffectedItems(items []*model.ItemEntity, itemId string, from, to uint) ([]*model.ItemEntity, error) {
	min, max := sortUint(from, to)
	affectedItems := make([]*model.ItemEntity, 0)
	for _, item := range items {
		if item.ItemId == itemId && item.Order != from {
			return nil, api.Error("", api.ErrorCodeInvalidArgument, "from not match")
		}

		if item.Order >= min && item.Order <= max {
			affectedItems = append(affectedItems, item)
		}
	}

	if from > to { //order 从大到小，中间的order 分别 +1
		for _, item := range affectedItems {
			if item.ItemId == itemId {
				item.Order = to
			} else {
				item.Order = item.Order + 1
			}
		}
	} else { //order 从小到大，中间的order 分别 -1
		for _, item := range affectedItems {
			if item.ItemId == itemId {
				item.Order = to
			} else {
				item.Order = item.Order - 1
			}
		}
	}

	return affectedItems, nil
}

func sortUint(a, b uint) (uint, uint) {
	if a > b {
		return b, a
	} else {
		return a, b
	}
}

func (s *ItemService) UpdateItemInfo(ctx context.Context, liveId string, item *model.ItemEntity) error {
	log := logger.ReqLogger(ctx)
	old, err := s.GetLiveItem(ctx, liveId, item.ItemId)
	if err != nil {
		log.Errorf("get live item error %s", err.Error())
		return err
	}

	if old == nil {
		log.Errorf("live item (%s, %s) not exist", liveId, item.ItemId)
		return api.ErrNotFound
	}

	updates := s.itemToUpdates(old, item)
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = timestamp.Now()

	db := mysql.GetLive(log.ReqID())
	if err = db.Model(old).Update(updates).Error; err != nil {
		log.Errorf("update item error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *ItemService) itemToUpdates(originItem, updateItem *model.ItemEntity) map[string]interface{} {
	updates := make(map[string]interface{})

	if len(updateItem.Title) > 0 && updateItem.Title != originItem.Title {
		updates["title"] = updateItem.Title
	}
	if len(updateItem.Tags) > 0 && updateItem.Tags != originItem.Tags {
		updates["tags"] = updateItem.Tags
	}
	if len(updateItem.Thumbnail) > 0 && updateItem.Thumbnail != originItem.Thumbnail {
		updates["thumbnail"] = updateItem.Thumbnail
	}
	if len(updateItem.Link) > 0 && updateItem.Link != originItem.Link {
		updates["link"] = updateItem.Link
	}
	if len(updateItem.CurrentPrice) > 0 && updateItem.CurrentPrice != originItem.CurrentPrice {
		updates["current_price"] = updateItem.CurrentPrice
	}
	if len(updateItem.OriginPrice) > 0 && updateItem.OriginPrice != originItem.OriginPrice {
		updates["origin_price"] = updateItem.OriginPrice
	}
	if updateItem.RecordId > 0 && updateItem.RecordId != originItem.RecordId {
		updates["record_id"] = updateItem.RecordId
	}
	if updateItem.RecordId == 0 && updateItem.RecordId != originItem.RecordId {
		updates["record_id"] = nil
	}
	if len(updateItem.Extends) > 0 {
		updates["extends"] = model.CombineExtends(originItem.Extends, updateItem.Extends)
	}

	return updates
}

func (s *ItemService) UpdateItemExtends(ctx context.Context, liveId string, itemId string, extends model.Extends) error {
	log := logger.ReqLogger(ctx)
	old, err := s.GetLiveItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("get live item error %s", err.Error())
		return err
	}

	if old == nil {
		log.Errorf("live item (%s, %s) not exist", liveId, itemId)
		return api.ErrNotFound
	}

	updates := make(map[string]interface{})
	updates["extends"] = model.CombineExtends(old.Extends, extends)
	updates["updated_at"] = timestamp.Now()

	db := mysql.GetLive(log.ReqID())
	if err = db.Model(old).Update(updates).Error; err != nil {
		log.Errorf("update item error %s", err.Error())
		return api.ErrDatabase
	}

	return nil
}

func (s *ItemService) StartRecordVideo(ctx context.Context, liveId string, itemId string) error {
	log := logger.ReqLogger(ctx)

	// 停止当前直播间的上一个商品的录制讲解
	p, err := s.GetPreviousItem(ctx, liveId)
	if err != nil {
		log.Errorf("record previous item error %s", err.Error())
	}
	if err == nil && p != nil {
		_, err = itemService.StopRecordVideo(ctx, liveId, *p)
	}

	// status
	err = s.saveRecordVideo(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("save item demonstrate log error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *ItemService) saveRecordVideo(ctx context.Context, liveId, itemId string) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	return db.Exec("insert into item_demonstrate_log(live_id, item_id, start,status) values(?, ?, ? ,? ) ",
		liveId, itemId, timestamp.Now(), model.RecordStatusDefault).Error
}

func (s *ItemService) StopRecordVideo(ctx context.Context, liveId string, demonId int) (demonstrateLog *model.ItemDemonstrateRecord, err error) {
	info, err := GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Info("get Live_entities table error %s", err.Error())
		return nil, err
	}
	//info.PushUrl = "rtmp://pili-publish.qnsdk.com/sdk-live/qn_live_kit-1556829451990339584"
	split := strings.Split(info.PushUrl, "/")
	demonstrateLog, err = s.getRecordVideo(ctx, demonId)
	if err != nil {
		log.Info("get DemonstrateLog table error %s", err.Error())
		return nil, err
	}
	if demonstrateLog == nil {
		log.Info("donnot find  record video in the date  error ")
		return nil, api.ErrNotFound
	}
	demonstrateLog.End = timestamp.Now()
	reqValue := &model.StreamsDemonstrateReq{
		Fname: demonstrateLog.LiveId + demonstrateLog.ItemId,
		Start: demonstrateLog.Start.Unix(),
		End:   demonstrateLog.End.Unix(),
	}
	encodedStreamTitle := base64.StdEncoding.EncodeToString([]byte(split[len(split)-1]))
	streamResp, err := s.postDemonstrateStreams(ctx, reqValue, encodedStreamTitle)
	if err != nil {
		log.Info("POST DemonstrateLog  error %s", err.Error())
		demonstrateLog.Status = model.RecordStatusFail
		s.UpdateRecordVideo(ctx, demonstrateLog)
		return nil, err
	}
	demonstrateLog.Fname = streamResp.Fname
	demonstrateLog.Status = model.RecordStatusSuccess
	s.UpdateRecordVideo(ctx, demonstrateLog)
	return demonstrateLog, nil
}

func (s *ItemService) UpdateItemRecord(ctx context.Context, demonId uint, liveId string, itemId string) error {
	item, err := s.GetLiveItem(ctx, liveId, itemId)
	if err != nil {
		return err
	}
	item.RecordId = demonId
	err = s.UpdateItemInfo(ctx, liveId, item)
	if err != nil {
		return err
	}
	return nil
}

func (s *ItemService) DeleteItemRecord(ctx context.Context, demonId uint, liveId string, itemId string) error {
	item, err := s.GetLiveItem(ctx, liveId, itemId)
	if err != nil {
		return err
	}
	item.RecordId = 0
	err = s.UpdateItemInfo(ctx, liveId, item)
	if err != nil {
		return err
	}
	return nil
}

func (s *ItemService) postDemonstrateStreams(ctx context.Context, reqValue *model.StreamsDemonstrateReq, encodedStreamTitle string) (*model.StreamsDemonstrateResponse, error) {
	url := "https://pili.qiniuapi.com" + "/v2/hubs/" + s.PiliHub + "/streams/" + encodedStreamTitle + "/saveas"
	mac := &qiniumac.Mac{
		AccessKey: s.AccessKey,
		SecretKey: []byte(s.SecretKey),
	}
	c := &http.Client{
		Transport: qiniumac.NewTransport(mac, nil),
	}
	client := rpc.Client{
		Client: c,
	}

	resp := &model.StreamsDemonstrateResponse{}
	err := client.CallWithJSON(logger.ReqLogger(ctx), resp, url, reqValue)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ItemService) GetRecordVideo(ctx context.Context, demonId uint) (*model.ItemDemonstrateRecord, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	demonstrateLog := model.ItemDemonstrateRecord{}
	result := db.First(&demonstrateLog, "id = ? ", demonId)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		} else {
			return nil, api.ErrDatabase
		}
	}
	return &demonstrateLog, nil
}

func (s *ItemService) getRecordVideo(ctx context.Context, demonId int) (*model.ItemDemonstrateRecord, error) {
	t := uint(demonId)
	return s.GetRecordVideo(ctx, t)
}

func (s *ItemService) GetListRecordVideo(ctx context.Context, liveId string, itemId string) (*model.ItemDemonstrateRecord, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	demonstrateLog := model.ItemDemonstrateRecord{}
	result := db.Last(&demonstrateLog, "live_id = ? and item_id = ? and status = 0", liveId, itemId)
	if result.Error != nil {
		log.Errorf("GetList DemosrateLog %v", result.Error)
		return nil, result.Error
	}
	return &demonstrateLog, nil
}

func (s *ItemService) GetListLiveRecordVideo(ctx context.Context, liveId string) ([]*model.ItemDemonstrateRecord, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	demonstrateLogs := make([]*model.ItemDemonstrateRecord, 0)
	result := db.Find(&demonstrateLogs, "live_id = ?", liveId)
	if result.Error != nil {
		log.Errorf("GetList DemosrateLog %v", result.Error)
		return nil, result.Error
	}
	return demonstrateLogs, nil
}

func (s *ItemService) UpdateRecordVideo(ctx context.Context, itemLog *model.ItemDemonstrateRecord) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	updates := map[string]interface{}{
		"fname":  itemLog.Fname,
		"status": itemLog.Status,
		"end":    itemLog.End,
	}
	var err error
	if itemLog.ID != 0 {
		err = db.Model(model.ItemDemonstrateRecord{}).Where("id = ? ", itemLog.ID).Update(updates).Error
	} else {
		err = db.Model(model.ItemDemonstrateRecord{}).Where("live_id = ? and item_id = ?", itemLog.LiveId, itemLog.ItemId).Update(updates).Error
	}
	if err != nil {
		log.Errorf("update itemLog error %s", err.Error())
		return err
	}
	return nil
}

func (s *ItemService) DelRecordVideo(ctx context.Context, liveId string, demonItem []uint) error {
	log := logger.ReqLogger(ctx)
	if len(demonItem) == 0 {
		return nil
	}
	db := mysql.GetLive(log.ReqID())
	err := db.Delete(model.ItemDemonstrateRecord{}, "live_id = ? and id in (?) ", liveId, demonItem).Error
	if err != nil {
		log.Errorf("delete demonstrate Log error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *ItemService) SetDemonstrateItem(ctx context.Context, liveId string, itemId string) error {
	log := logger.ReqLogger(ctx)

	// 停止当前直播间的上一个商品的录制讲解
	p, err := s.GetPreviousItem(ctx, liveId)
	if err != nil {
		log.Errorf("record previous item error %s", err.Error())
	}
	if err == nil && p != nil {
		_, err = itemService.StopRecordVideo(ctx, liveId, *p)
	}

	err = s.saveDemonstrateItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("save item demonstrate error %s", err.Error())
		return api.ErrDatabase
	}

	return nil
}

func (s *ItemService) DelDemonstrateItem(ctx context.Context, liveId string) error {
	log := logger.ReqLogger(ctx)
	err := s.saveDemonstrateItem(ctx, liveId, "")
	if err != nil {
		log.Errorf("save item demonstrate error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *ItemService) saveDemonstrateItem(ctx context.Context, liveId, itemId string) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	return db.Exec("insert into item_demonstrate(live_id, item_id, updated_at) values(?, ?, ?) "+
		"ON DUPLICATE KEY UPDATE item_id = ?, updated_at = ?",
		liveId, itemId, timestamp.Now(),
		itemId, timestamp.Now()).Error
}

func (s *ItemService) IsDemonstrateItem(ctx context.Context, liveId, itemId string) (bool, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	d := model.ItemDemonstrate{}
	result := db.First(&d, "live_id= ? and item_id = ? ", liveId, itemId)
	if result.Error != nil {
		if result.RecordNotFound() {
			return false, nil
		} else {
			return false, api.ErrDatabase
		}
	}
	return true, nil
}

func (s *ItemService) GetDemonstrateItem(ctx context.Context, liveId string) (*model.ItemEntity, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	d := model.ItemDemonstrate{}
	result := db.First(&d, "live_id = ?", liveId)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		} else {
			return nil, api.ErrDatabase
		}
	}

	if len(d.ItemId) == 0 {
		return nil, nil
	}

	return s.GetLiveItem(ctx, liveId, d.ItemId)
}
func (s *ItemService) GetPreviousItem(ctx context.Context, liveId string) (*int, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	d := model.ItemDemonstrateRecord{}
	result := db.Last(&d, "live_id = ? and end is null", liveId)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		} else {
			return nil, api.ErrDatabase
		}
	}

	id := int(d.ID)
	return &id, nil
}

func (s *ItemService) GetLiveItem(ctx context.Context, liveId string, itemId string) (*model.ItemEntity, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())

	itemEntity := model.ItemEntity{}
	result := db.First(&itemEntity, "live_id = ? and item_id = ?", liveId, itemId)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, nil
		} else {
			return nil, api.ErrDatabase
		}
	}

	return &itemEntity, nil
}
