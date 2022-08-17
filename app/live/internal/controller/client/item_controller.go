// @Author: wangsheng
// @Description:
// @File:  item_controller
// @Version: 1.0.0
// @Date: 2022/7/1 5:58 下午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"net/http"

	"github.com/qbox/livekit/biz/live"
	"github.com/qbox/livekit/common/auth/liveauth"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
)

var ItemController = &itemController{}

func RegisterItemRoutes(group *gin.RouterGroup) {
	itemGroup := group.Group("/item")
	itemGroup.POST("/add", ItemController.PostItemAdd)
	itemGroup.POST("/delete", ItemController.PostItemDelete)
	itemGroup.POST("/status", ItemController.PostItemStatus)
	itemGroup.POST("/order", ItemController.PostItemOrder)
	itemGroup.POST("/order/single", ItemController.PostItemOrderSingle)
	itemGroup.GET("/:liveId", ItemController.GetItems)
	itemGroup.PUT("/:liveId/:itemId", ItemController.PutItem)
	itemGroup.PUT("/:liveId/:itemId/extends", ItemController.PutItemExtends)

	itemGroup.POST("/demonstrate/:liveId/:itemId", ItemController.PostItemDemonstrate)
	itemGroup.POST("/demonstrate/start/:liveId/:itemId", ItemController.PostStartRecordDemonstrate)
	itemGroup.DELETE("/demonstrate/:liveId", ItemController.DeleteItemDemonstrate)
	itemGroup.GET("/demonstrate/:liveId", ItemController.GetItemDemonstrate)

	itemGroup.GET("/demonstrate/record/:liveId", ItemController.ListLiveRecordVideo)
	itemGroup.GET("/demonstrate/record/:liveId/:itemId", ItemController.ListrecordVideo)
	itemGroup.POST("/demonstrate/record/delete", ItemController.DelRecordVideo)
}

type itemController struct {
}

type AddItemRequest struct {
	LiveId string         `json:"live_id"`
	Items  []*dto.ItemDto `json:"items"`
}

func (c *itemController) PostItemAdd(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	request := &AddItemRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(request.LiveId) == 0 || len(request.Items) == 0 {
		log.Errorf("invalid request %+v", request)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	entities := make([]*model.ItemEntity, 0, len(request.Items))
	for _, d := range request.Items {
		entities = append(entities, dto.ItemDtoToEntity(d))
	}

	itemService := live.GetItemService()
	err := itemService.AddItems(ctx, request.LiveId, entities)
	if err != nil {
		log.Errorf("add items error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type DelItemRequest struct {
	LiveId string   `json:"live_id"`
	Items  []string `json:"items"`
}

func (c *itemController) PostItemDelete(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	request := &DelItemRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(request.LiveId) == 0 || len(request.Items) == 0 {
		log.Errorf("invalid request %+v", request)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	err := itemService.DelItems(ctx, request.LiveId, request.Items)
	if err != nil {
		log.Errorf("delete items error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type UpdateItemStatusRequest struct {
	LiveId string              `json:"live_id"`
	Items  []*model.ItemStatus `json:"items"`
}

func (c *itemController) PostItemStatus(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	request := &UpdateItemStatusRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(request.LiveId) == 0 || len(request.Items) == 0 {
		log.Errorf("invalid request %+v", request)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	err := itemService.UpdateItemStatus(ctx, request.LiveId, request.Items)
	if err != nil {
		log.Errorf("update item status error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type UpdateItemOrderRequest struct {
	LiveId string             `json:"live_id"`
	Items  []*model.ItemOrder `json:"items"`
}

func (c *itemController) PostItemOrder(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	request := &UpdateItemOrderRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(request.LiveId) == 0 || len(request.Items) == 0 {
		log.Errorf("invalid request %+v", request)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	err := itemService.UpdateItemOrder(ctx, request.LiveId, request.Items)
	if err != nil {
		log.Errorf("update item order error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type UpdateItemOrderSingleRequest struct {
	LiveId string `json:"live_id"`
	ItemId string `json:"item_id"`
	From   uint   `json:"from"`
	To     uint   `json:"to"`
}

func (c *itemController) PostItemOrderSingle(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	request := &UpdateItemOrderSingleRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(request.LiveId) == 0 {
		log.Errorf("invalid request %+v", request)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, request.LiveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	err := itemService.UpdateItemOrderSingle(ctx, request.LiveId, request.ItemId, request.From, request.To)
	if err != nil {
		log.Errorf("update item order error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type GetItemsResponse struct {
	api.Response
	Data []*dto.ItemDto `json:"data"`
}

func (c *itemController) GetItems(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")

	userInfo := liveauth.GetUserInfo(ctx)
	liveEntity, err := live.GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("get live info error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	showOffline := userInfo.UserId == liveEntity.AnchorId
	itemService := live.GetItemService()
	entities, err := itemService.ListItems(ctx, liveId, showOffline)
	if err != nil {
		log.Errorf("list item error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	demonItem, err := itemService.GetDemonstrateItem(ctx, liveId)
	if err != nil {
		log.Errorf("get demonstrate item error %s", err.Error())
	}

	response := GetItemsResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     make([]*dto.ItemDto, 0),
	}
	dtos := make([]*dto.ItemDto, 0, len(entities))
	if demonItem != nil {
		dtos = append(dtos, dto.ItemEntityToDto(demonItem))
	}

	if len(entities) > 0 {
		for _, e := range entities {
			if demonItem == nil || e.ItemId != demonItem.ItemId {
				dtos = append(dtos, dto.ItemEntityToDto(e))
			}
		}
	}
	response.Data = dtos

	ctx.JSON(http.StatusOK, response)
}

func (c *itemController) PutItem(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	itemDto := dto.ItemDto{}
	if err := ctx.BindJSON(&itemDto); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	itemDto.ItemId = itemId

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	itemEntity := dto.ItemDtoToEntity(&itemDto)
	err := itemService.UpdateItemInfo(ctx, liveId, itemEntity)
	if err != nil {
		log.Errorf("update item info error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

func (c *itemController) PutItemExtends(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	extends := make(map[string]string)
	if err := ctx.BindJSON(&extends); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	err := itemService.UpdateItemExtends(ctx, liveId, itemId, extends)
	if err != nil {
		log.Errorf("update item extends error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

func (c *itemController) PostItemDemonstrate(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	item, err := itemService.GetLiveItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("get live item error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	if item.Status == model.ItemStatusOffline {
		log.Errorf("item offline, cannot demonstrate ")
		ctx.AbortWithStatusJSON(http.StatusOK, api.Error(log.ReqID(), 501, "item offline, cannot demonstrate "))
		return
	}
	err = itemService.SetDemonstrateItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("set demonstrate item error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

func (c *itemController) PostStartRecordDemonstrate(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	item, err := itemService.GetLiveItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("get live item error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	if item.Status == model.ItemStatusOffline {
		log.Errorf("item offline, cannot demonstrate ")
		ctx.AbortWithStatusJSON(http.StatusOK, api.Error(log.ReqID(), 501, "item offline, cannot demonstrate "))
		return
	}
	err = itemService.SetDemonstrateItem(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("set demonstrate item error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	err = itemService.StartRecordVideo(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("set demonstrate  Log error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

func (c *itemController) ListLiveRecordVideo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	demonstrateLog, err := itemService.GetListLiveRecordVideo(ctx, liveId)
	if err != nil {
		log.Errorf("ListrecordVideo error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	var data []*dto.RecordDto
	for _, v := range demonstrateLog {
		data = append(data, dto.RecordEntityToDto(v))
	}
	response := &ListLiveDemonstrateLog{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     data,
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *itemController) ListrecordVideo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	demonstrateLog, err := itemService.GetListRecordVideo(ctx, liveId, itemId)
	if err != nil {
		log.Errorf("ListrecordVideo error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	response := &ListDemonstrateLog{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     dto.RecordEntityToDto(demonstrateLog),
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *itemController) DelRecordVideo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	request := &DelDemonItemRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	if len(request.DemonItems) == 0 {
		log.Errorf("invalid request %+v", request)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	itemService := live.GetItemService()
	err := itemService.DelRecordVideo(ctx, request.LiveId, request.DemonItems)
	if err != nil {
		log.Errorf("delete demonstrate error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type DelDemonItemRequest struct {
	LiveId     string `json:"live_id"`
	DemonItems []int  `json:"demonstrate_item"`
}

type StopItemDemonstrateLogResponse struct {
	api.Response
	Data *dto.RecordDto
}

type ListDemonstrateLog struct {
	api.Response
	Data *dto.RecordDto
}
type ListLiveDemonstrateLog struct {
	api.Response
	Data []*dto.RecordDto
}

func (c *itemController) DeleteItemDemonstrate(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")

	userInfo := liveauth.GetUserInfo(ctx)
	if err := live.GetService().CheckLiveAnchor(ctx, liveId, userInfo.UserId); err != nil {
		log.Errorf("check live anchor error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	itemService := live.GetItemService()
	err := itemService.DelDemonstrateItem(ctx, liveId)
	if err != nil {
		log.Errorf("delete demonstrate item error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	demonId, err := itemService.GetPreviousItem(ctx, liveId)
	if err != nil {
		log.Errorf("delete demonstrate item error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	if demonId == nil {
		ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
	}
	demonstrateLog, err := itemService.StopRecordVideo(ctx, liveId, *demonId)
	if err != nil {
		log.Errorf("record and stop demonstrate log error %+v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	demonstrateLog.Fname = "pili-playback.qnsdk.com/" + demonstrateLog.Fname
	response := &StopItemDemonstrateLogResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     dto.RecordEntityToDto(demonstrateLog),
	}
	ctx.JSON(http.StatusOK, response)
}

type GetItemDemonstrateResponse struct {
	api.Response
	Data *dto.ItemDto `json:"data"`
}

func (c *itemController) GetItemDemonstrate(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")

	itemService := live.GetItemService()
	itemEntity, err := itemService.GetDemonstrateItem(ctx, liveId)
	if err != nil {
		log.Errorf("get demonstrate item error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	response := GetItemDemonstrateResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     dto.ItemEntityToDto(itemEntity),
	}
	ctx.JSON(http.StatusOK, response)
}
