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

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
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
	//itemGroup.GET("/:liveId", ItemController.GetItems)
	//
	//itemGroup.PUT("/:liveId/:itemId", ItemController.PutItem)
	//itemGroup.PUT("/:liveId/:itemId/extends", ItemController.PutItemExtends)

	httpq.ServerHandle(http.MethodPost, "/item/add", ItemController.PostItemAdd)
	httpq.ServerHandle(http.MethodPost, "/item/delete", ItemController.PostItemDelete)
	httpq.ServerHandle(http.MethodPost, "/item/status", ItemController.PostItemStatus)
	httpq.ServerHandle(http.MethodPost, "/item/order", ItemController.PostItemOrder)
	httpq.ServerHandle(http.MethodGet, "/item/:liveId", ItemController.GetItems)

	httpq.ServerHandle(http.MethodPut, "/item/:liveId/:itemId", ItemController.PutItem)
	httpq.ServerHandle(http.MethodPut, "/item/:liveId/:itemId/extends", ItemController.PutItemExtends)
}

type itemController struct {
}

type AddItemRequest struct {
	LiveId string         `json:"live_id"`
	Items  []*dto.ItemDto `json:"items"`
}

// PostItemAdd 在直播间添加商品
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

// PostItemDelete 删除直播间商品
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

// PostItemStatus 批量修改商品状态
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

// PostItemOrder 修改商品的顺序
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

	itemService := impl.GetInstance()
	err := itemService.UpdateItemOrder(ctx, request.LiveId, request.Items)
	if err != nil {
		log.Errorf("update item order error %s", err.Error())
		return nil, err
	}

	return nil, nil
}

// GetItems 获取商品列表
// return []*dto.ItemDto
func (c *itemController) GetItems(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")

	itemService := impl.GetInstance()
	entities, err := itemService.ListItems(ctx, liveId, true)
	if err != nil {
		log.Errorf("list item error %s", err.Error())
		return nil, err
	}

	dtos := make([]*dto.ItemDto, 0, len(entities))
	if len(entities) > 0 {
		for _, e := range entities {

			dtos = append(dtos, utils.ItemEntityToDto(e))
		}
	}
	return dtos, nil
}

// PutItem 更新商品信息
// return nil
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

	itemService := impl.GetInstance()
	itemEntity := dto.ItemDtoToEntity(&itemDto)
	err := itemService.UpdateItemInfo(ctx, liveId, itemEntity)
	if err != nil {
		log.Errorf("update item info error %s", err.Error())
		return nil, err
	}

	return nil, nil
}

// PutItemExtends 更新商品的扩展信息
// return nil
func (c *itemController) PutItemExtends(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	itemId := ctx.Param("itemId")

	extends := make(map[string]string)
	if err := ctx.BindJSON(&extends); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	itemService := impl.GetInstance()
	err := itemService.UpdateItemExtends(ctx, liveId, itemId, extends)
	if err != nil {
		log.Errorf("update item extends error %s", err.Error())
		return nil, err
	}

	return nil, nil
}
