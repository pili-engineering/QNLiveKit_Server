// @Author: wangsheng
// @Description:
// @File:  live_controller
// @Version: 1.0.0
// @Date: 2022/6/28 11:07 上午
// Copyright 2021 QINIU. All rights reserved

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/live/dto"
	"github.com/qbox/livekit/module/base/live/internal/impl"
	"github.com/qbox/livekit/module/base/live/service"
	"github.com/qbox/livekit/module/base/user"

	"github.com/qbox/livekit/utils/logger"
)

func RegisterLiveRoutes(group *gin.RouterGroup) {
	//userGroup := group.Group("/live")
	//userGroup.POST("", LiveController.PostLiveCreate)
	//userGroup.GET("/:id", LiveController.GetLive)
	//userGroup.POST("/:id/stop", LiveController.PostLiveStop)
	//userGroup.DELETE("/:id", LiveController.DeleteLive)

	httpq.ServerHandle(http.MethodPost, "/live", LiveController.PostLiveCreate)
	httpq.ServerHandle(http.MethodGet, "/live/:id", LiveController.GetLive)
	httpq.ServerHandle(http.MethodPost, "/live/:id/stop", LiveController.PostLiveStop)
	httpq.ServerHandle(http.MethodDelete, "/live/:id", LiveController.DeleteLive)
}

var LiveController = &liveController{}

type liveController struct {
}

// PostLiveCreate 创建直播间
// return dto.LiveInfoDto
func (c *liveController) PostLiveCreate(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	request := &service.CreateLiveRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	anchor, err := user.GetService().FindUser(ctx, request.AnchorId)
	if err != nil {
		log.Errorf("find anchor failed, err: %v", err)
		return nil, err
	}

	liveEntity, err := impl.GetInstance().CreateLive(ctx, request)
	if err != nil {
		log.Errorf("create liveEntity failed, err: %v", err)
		return nil, err
	}

	liveDto := dto.BuildLiveDto(liveEntity, anchor)
	return liveDto, nil
}

// GetLive 查询直播间
// return
func (c *liveController) GetLive(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("id")
	liveEntity, err := impl.GetInstance().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("create liveEntity failed, err: %v", err)
		return nil, err
	}

	anchor, err := user.GetService().FindUser(ctx, liveEntity.AnchorId)
	if err != nil {
		log.Errorf("find anchor failed, err: %v", err)
		return nil, err
	}

	liveDto := dto.BuildLiveDto(liveEntity, anchor)
	return liveDto, nil
}

// PostLiveStop 停止直播间
// return nil
func (c *liveController) PostLiveStop(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("id")
	liveEntity, err := impl.GetInstance().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("find live error %s", err.Error())
		return nil, err
	}

	err = impl.GetInstance().StopLive(ctx, liveId, liveEntity.AnchorId)
	if err != nil {
		log.Errorf("stop live failed, err: %v", err)
		return nil, err
	}

	return nil, nil
}

// DeleteLive 删除直播间
// return nil
func (c *liveController) DeleteLive(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("id")

	liveEntity, err := impl.GetInstance().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("find live error %s", err.Error())
		return nil, err
	}

	err = impl.GetInstance().DeleteLive(ctx, liveId, liveEntity.AnchorId)
	if err != nil {
		log.Errorf("delete live failed, err: %v", err)
		return nil, err
	}
	return nil, nil
}
