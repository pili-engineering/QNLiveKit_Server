// @Author: wangsheng
// @Description:
// @File:  item_controller
// @Version: 1.0.0
// @Date: 2022/7/1 5:58 下午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/live"
	"github.com/qbox/livekit/module/biz/item/dto"
	"github.com/qbox/livekit/module/biz/item/internal/controller/utils"
	"github.com/qbox/livekit/module/biz/item/internal/impl"
	"github.com/qbox/livekit/utils/logger"
)

var ItemController = &itemController{}

func RegisterRoutes() {
	//itemGroup := group.Group("/item")
	//itemGroup.POST("/add", ItemController.PostItemAdd)
	//itemGroup.POST("/delete", ItemController.PostItemDelete)
	//itemGroup.POST("/status", ItemController.PostItemStatus)
	//itemGroup.POST("/order", ItemController.PostItemOrder)
	//itemGroup.POST("/order/single", ItemController.PostItemOrderSingle)
	//itemGroup.GET("/:liveId", ItemController.GetItems)
	//itemGroup.PUT("/:liveId/:itemId", ItemController.PutItem)
	//itemGroup.PUT("/:liveId/:itemId/extends", ItemController.PutItemExtends)
	//
	//itemGroup.POST("/demonstrate/:liveId/:itemId", ItemController.PostItemDemonstrate)
	//itemGroup.POST("/demonstrate/start/:liveId/:itemId", ItemController.PostStartRecordDemonstrate)
	//itemGroup.DELETE("/demonstrate/:liveId", ItemController.DeleteItemDemonstrate)
	//itemGroup.GET("/demonstrate/:liveId", ItemController.GetItemDemonstrate)
	//
	//itemGroup.GET("/demonstrate/record/:liveId", ItemController.ListLiveRecordVideo)
	//itemGroup.GET("/demonstrate/record/:liveId/:itemId", ItemController.ListrecordVideo)
	//itemGroup.POST("/demonstrate/record/delete", ItemController.DelRecordVideo)

	httpq.ClientHandle(http.MethodPut, "/item/add", ItemController.PostItemAdd)
	httpq.ClientHandle(http.MethodPost, "/item/delete", ItemController.PostItemDelete)
	httpq.ClientHandle(http.MethodPost, "/item/status", ItemController.PostItemStatus)
	httpq.ClientHandle(http.MethodPost, "/item/order", ItemController.PostItemOrder)
	httpq.ClientHandle(http.MethodPost, "/item/order/single", ItemController.PostItemOrderSingle)
	httpq.ClientHandle(http.MethodGet, "/item/:liveId", ItemController.GetItems)
	httpq.ClientHandle(http.MethodPut, "/item/:liveId/:itemId", ItemController.PutItem)
	httpq.ClientHandle(http.MethodPut, "/item/:liveId/:itemId/extends", ItemController.PutItemExtends)

	httpq.ClientHandle(http.MethodPost, "/item/demonstrate/:liveId/:itemId", ItemController.PostItemDemonstrate)
	httpq.ClientHandle(http.MethodPost, "/item/demonstrate/start/:liveId/:itemId", ItemController.PostStartRecordDemonstrate)
	httpq.ClientHandle(http.MethodDelete, "/item/demonstrate/:liveId", ItemController.DeleteItemDemonstrate)
	httpq.ClientHandle(http.MethodGet, "/demonstrate/:liveId", ItemController.GetItemDemonstrate)

	httpq.ClientHandle(http.MethodGet, "/item/demonstrate/record/:liveId", ItemController.ListLiveRecordVideo)
	httpq.ClientHandle(http.MethodGet, "/item/demonstrate/record/:liveId/:itemId", ItemController.ListrecordVideo)
	httpq.ClientHandle(http.MethodPost, "/item/demonstrate/record/delete", ItemController.DelRecordVideo)
}

type itemController struct {
}

type AddItemRequest struct {
	LiveId string         `json:"live_id"`
	Items  []*dto.ItemDto `json:"items"`
}

// PostItemAdd 添加商品
// return nil
func (c *itemController) PostItemAdd(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	request := &AddItemRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(request.LiveId) == 0 || len(request.Items) == 0 {
		log.Errorf("invalid request %+v", request)
		return nil, rest.ErrBadRequest
	}

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	entities := make([]*model.ItemEntity, 0, len(request.Items))
	for _, d := range request.Items {
		entities = append(entities, dto.ItemDtoToEntity(d))
	}

	itemService := impl.GetInstance()
	err := itemService.AddItems(ctx, request.LiveId, entities)
	if err != nil {
		log.Errorf("add items error %s", err.Error())
		return nil, err
	}

	return nil, nil
}

type DelItemRequest struct {
	LiveId string   `json:"live_id"`
	Items  []string `json:"items"`
}

// PostItemDelete 删除商品
// return nil
func (c *itemController) PostItemDelete(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	request := &DelItemRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(request.LiveId) == 0 || len(request.Items) == 0 {
		log.Errorf("invalid request %+v", request)
		return nil, rest.ErrBadRequest
	}

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	err := itemService.DelItems(ctx, request.LiveId, request.Items)
	if err != nil {
		log.Errorf("delete items error %s", err.Error())
		return nil, err
	}

	return nil, nil
}

type UpdateItemStatusRequest struct {
	LiveId string              `json:"live_id"`
	Items  []*model.ItemStatus `json:"items"`
}

// PostItemStatus 更新商品状态
// return nil
func (c *itemController) PostItemStatus(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	request := &UpdateItemStatusRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(request.LiveId) == 0 || len(request.Items) == 0 {
		log.Errorf("invalid request %+v", request)
		return nil, rest.ErrBadRequest
	}

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	err := itemService.UpdateItemStatus(ctx, request.LiveId, request.Items)
	if err != nil {
		log.Errorf("update item status error %s", err.Error())
		return nil, err
	}

	return nil, nil
}

type UpdateItemOrderRequest struct {
	LiveId string             `json:"live_id"`
	Items  []*model.ItemOrder `json:"items"`
}

// PostItemOrder 修改商品顺序
// return nil
func (c *itemController) PostItemOrder(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	request := &UpdateItemOrderRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(request.LiveId) == 0 || len(request.Items) == 0 {
		log.Errorf("invalid request %+v", request)
		return nil, rest.ErrBadRequest
	}

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	err := itemService.UpdateItemOrder(ctx, request.LiveId, request.Items)
	if err != nil {
		log.Errorf("update item order error %s", err.Error())
		return nil, err
	}

	return nil, nil
}

type UpdateItemOrderSingleRequest struct {
	LiveId string `json:"live_id"`
	ItemId string `json:"item_id"`
	From   uint   `json:"from"`
	To     uint   `json:"to"`
}

// PostItemOrderSingle 更改单个善品的顺序
// return nil
func (c *itemController) PostItemOrderSingle(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	request := &UpdateItemOrderSingleRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(request.LiveId) == 0 {
		log.Errorf("invalid request %+v", request)
		return nil, rest.ErrBadRequest
	}

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	err := itemService.UpdateItemOrderSingle(ctx, request.LiveId, request.ItemId, request.From, request.To)
	if err != nil {
		log.Errorf("update item order error %s", err.Error())
		return nil, err
	}

	return nil, nil
}

// GetItems 查询商品信息
// return []*dto.ItemDto
func (c *itemController) GetItems(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")

	userInfo := auth.GetUserInfo(ctx)
	liveEntity, err := live.GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("get live info error %s", err.Error())
		return nil, err
	}

	showOffline := userInfo.UserId == liveEntity.AnchorId
	itemService := impl.GetInstance()
	entities, err := itemService.ListItems(ctx, liveId, showOffline)
	if err != nil {
		log.Errorf("list item error %s", err.Error())
		return nil, err
	}

	demonItem, err := itemService.GetDemonstrateItem(ctx, liveId)
	if err != nil {
		log.Errorf("get demonstrate item error %s", err.Error())
	}

	dtos := make([]*dto.ItemDto, 0, len(entities))
	if demonItem != nil {
		dtos = append(dtos, utils.ItemEntityToDto(demonItem))
	}

	if len(entities) > 0 {
		for _, e := range entities {
			if demonItem == nil || e.ItemId != demonItem.ItemId {
				dtos = append(dtos, utils.ItemEntityToDto(e))
			}
		}
	}

	return dtos, nil
}

// PutItem 更新商品信息
func (c *itemController) PutItem(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	itemDto := dto.ItemDto{}
	if err := ctx.BindJSON(&itemDto); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	itemDto.ItemId = itemId

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	itemEntity := dto.ItemDtoToEntity(&itemDto)
	err := itemService.UpdateItemInfo(ctx, liveId, itemEntity)
	if err != nil {
		log.Errorf("update item info error %s", err.Error())
		return nil, err
	}

	return nil, nil
}

// PutItemExtends 更新商品扩展信息
func (c *itemController) PutItemExtends(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	extends := make(map[string]string)
	if err := ctx.BindJSON(&extends); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	err := itemService.UpdateItemExtends(ctx, liveId, itemId, extends)
	if err != nil {
		log.Errorf("update item extends error %s", err.Error())
		return nil, err
	}

	return nil, nil
}

// PostItemDemonstrate 讲解商品
func (c *itemController) PostItemDemonstrate(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	item, err := itemService.GetLiveItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("get live item error %s", err.Error())
		return nil, err
	}
	if item.Status == model.ItemStatusOffline {
		log.Errorf("item offline, cannot demonstrate ")
		return nil, rest.ErrBadRequest.WithMessage("item offline, cannot demonstrate")
	}

	err = itemService.SetDemonstrateItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("set demonstrate item error %s", err.Error())
		return nil, err
	}
	return nil, nil
}

// PostStartRecordDemonstrate 开始录制讲解商品
func (c *itemController) PostStartRecordDemonstrate(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	item, err := itemService.GetLiveItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("get live item error %s", err.Error())
		return nil, err
	}
	if item.Status == model.ItemStatusOffline {
		log.Errorf("item offline, cannot demonstrate ")
		return nil, rest.ErrBadRequest.WithMessage("item offline, cannot demonstrate ")
	}

	err = itemService.SetDemonstrateItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("set demonstrate item error %s", err.Error())
		return nil, err
	}

	err = itemService.StartRecordVideo(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("set demonstrate  Log error %s", err.Error())
		return nil, err
	}
	demonItem, err := itemService.GetPreviousItem(ctx, liveId)
	if err != nil {
		log.Errorf("set demonstrate  Log,GetPreviousItem  error %s", err.Error())
		return nil, err
	}

	err = itemService.UpdateItemRecord(ctx, uint(*demonItem), liveId, itemId)
	if err != nil {
		log.Errorf("Demonstrate Record  donnot save to item %s", err.Error())
		return nil, err
	}

	return nil, nil
}

// ListLiveRecordVideo 查看商品讲解记录
// return []*dto.RecordDto
func (c *itemController) ListLiveRecordVideo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	demonstrateLog, err := itemService.GetListLiveRecordVideo(ctx, liveId)
	if err != nil {
		log.Errorf("ListrecordVideo error %s", err.Error())
		return nil, err
	}
	var data []*dto.RecordDto
	for _, v := range demonstrateLog {
		data = append(data, utils.RecordEntityToDto(v))
	}
	return data, nil
}

// ListrecordVideo 查看讲解记录
// return *dto.RecordDto
func (c *itemController) ListrecordVideo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	demonstrateLog, err := itemService.GetListRecordVideo(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("ListrecordVideo error %s", err.Error())
		return nil, err
	}

	return utils.RecordEntityToDto(demonstrateLog), nil
}

type DelDemonItemResult struct {
	LiveId            string `json:"live_id"`
	FailureDemonItems []uint `json:"failure_demon_items,omitempty"`
}

// DelRecordVideo 删除讲解记录
// return *DelDemonItemResult
func (c *itemController) DelRecordVideo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	request := &DelDemonItemRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	if len(request.DemonItems) == 0 {
		log.Errorf("invalid request %+v", request)
		return nil, rest.ErrBadRequest.WithMessage("empty demonstrate_item")
	}
	itemService := impl.GetInstance()

	d := &DelDemonItemResult{
		LiveId: request.LiveId,
	}
	for _, v := range request.DemonItems {
		video, err := itemService.GetRecordVideo(ctx, v)
		if video.LiveId != request.LiveId {
			log.Errorf("demonstrate_record_id donnot equal to request live_id  %+v", request)
			return nil, rest.ErrBadRequest.WithMessage(fmt.Sprintf("video %d not int live %s", video.ID, request.LiveId))
		}
		if err != nil {
			d.FailureDemonItems = append(d.FailureDemonItems, v)
			log.Errorf("delete demonstrate error %s", err.Error())
		} else if video == nil {
			continue
		} else {
			err = itemService.DeleteItemRecord(ctx, v, video.LiveId, video.ItemId)
			if err != nil {
				d.FailureDemonItems = append(d.FailureDemonItems, v)
				log.Errorf("delete demonstrate error %s", err.Error())
			} else {
				err := itemService.DelRecordVideo(ctx, request.LiveId, []uint{v})
				if err != nil {
					d.FailureDemonItems = append(d.FailureDemonItems, v)
					log.Errorf("delete demonstrate error %s", err.Error())
				}
			}
		}
	}

	return d, nil
}

type DelDemonItemRequest struct {
	LiveId     string `json:"live_id"`
	DemonItems []uint `json:"demonstrate_item"`
}

// DeleteItemDemonstrate 删除商品简介记录
// return *dto.RecordDto
func (c *itemController) DeleteItemDemonstrate(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")

	userInfo := auth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		return nil, err
	}

	itemService := impl.GetInstance()
	err := itemService.DelDemonstrateItem(ctx, liveId)
	if err != nil {
		log.Errorf("delete demonstrate item error %s", err.Error())
		return nil, err
	}

	demonId, err := itemService.GetPreviousItem(ctx, liveId)
	if err != nil {
		log.Errorf("delete demonstrate item error %s", err.Error())
		return nil, err
	}
	if demonId == nil {
		return nil, nil
	}

	demonstrateLog, err := itemService.StopRecordVideo(ctx, liveId, *demonId)
	if err != nil {
		log.Errorf("record and stop demonstrate log error %+v", err)
		return nil, err
	}
	return utils.RecordEntityToDto(demonstrateLog), nil
}

// GetItemDemonstrate 获取上平讲解记录
// return *dto.ItemDto
func (c *itemController) GetItemDemonstrate(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")

	itemService := impl.GetInstance()
	itemEntity, err := itemService.GetDemonstrateItem(ctx, liveId)
	if err != nil {
		log.Errorf("get demonstrate item error %s", err.Error())
		return nil, err
	}

	return utils.ItemEntityToDto(itemEntity), nil
}
