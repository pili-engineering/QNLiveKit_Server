// @Author: wangsheng
// @Description:
// @File:  item_controller
// @Version: 1.0.0
// @Date: 2022/7/1 5:58 下午
// Copyright 2021 QINIU. All rights reserved

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/module/base/live"
	"github.com/qbox/livekit/utils/logger"
)

var ItemController = &itemController{}

func RegisterItemRoutes(group *gin.RouterGroup) {
	itemGroup := group.Group("/item")
	itemGroup.POST("/add", ItemController.PostItemAdd)
	itemGroup.POST("/delete", ItemController.PostItemDelete)
	itemGroup.POST("/status", ItemController.PostItemStatus)
	itemGroup.POST("/order", ItemController.PostItemOrder)
	itemGroup.GET("/:liveId", ItemController.GetItems)

	itemGroup.PUT("/:liveId/:itemId", ItemController.PutItem)
	itemGroup.PUT("/:liveId/:itemId/extends", ItemController.PutItemExtends)
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

	itemService := live.GetItemService()
	err := itemService.UpdateItemOrder(ctx, request.LiveId, request.Items)
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

	itemService := live.GetItemService()
	entities, err := itemService.ListItems(ctx, liveId, true)
	if err != nil {
		log.Errorf("list item error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	response := GetItemsResponse{
		Response: api.SuccessResponse(log.ReqID()),
	}
	if len(entities) > 0 {
		dtos := make([]*dto.ItemDto, 0, len(entities))
		for _, e := range entities {
			dtos = append(dtos, dto.ItemEntityToDto(e))
		}
		response.Data = dtos
	}
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

	itemService := live.GetItemService()
	err := itemService.UpdateItemExtends(ctx, liveId, itemId, extends)
	if err != nil {
		log.Errorf("update item extends error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}
