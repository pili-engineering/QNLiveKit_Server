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

	"github.com/qbox/livekit/module/base/live/dto"
	"github.com/qbox/livekit/module/base/live/service"

	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterLiveRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/live")
	userGroup.POST("", LiveController.PostLiveCreate)
	userGroup.GET("/:id", LiveController.GetLive)
	userGroup.POST("/:id/stop", LiveController.PostLiveStop)
	userGroup.DELETE("/:id", LiveController.DeleteLive)
}

var LiveController = &liveController{}

type liveController struct {
}

type LiveResponse struct {
	api.Response
	Data dto.LiveInfoDto `json:"data"`
}

func (c *liveController) PostLiveCreate(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	request := &service.CreateLiveRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.Response{
			RequestId: log.ReqID(),
			Code:      api.ErrorCodeInvalidArgument,
			Message:   "invalid request",
		})
		return
	}

	anchor, err := service.GetService().FindUser(ctx, request.AnchorId)
	if err != nil {
		log.Errorf("find anchor failed, err: %v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	liveEntity, err := service.GetService().CreateLive(ctx, request)
	if err != nil {
		log.Errorf("create liveEntity failed, err: %v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "create live failed",
			RequestId: log.ReqID(),
		})
		return
	}

	response := &LiveResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.LiveId = liveEntity.LiveId
	response.Data.Title = liveEntity.Title
	response.Data.Notice = liveEntity.Notice
	response.Data.CoverUrl = liveEntity.CoverUrl
	response.Data.Extends = liveEntity.Extends
	response.Data.AnchorInfo.UserId = anchor.UserId
	response.Data.AnchorInfo.ImUserid = anchor.ImUserid
	response.Data.AnchorInfo.Nick = anchor.Nick
	response.Data.AnchorInfo.Avatar = anchor.Avatar
	response.Data.AnchorInfo.Extends = anchor.Extends
	response.Data.RoomToken = ""
	response.Data.PkId = liveEntity.PkId
	response.Data.OnlineCount = liveEntity.OnlineCount
	response.Data.StartTime = liveEntity.StartAt.Unix()
	response.Data.EndTime = liveEntity.EndAt.Unix()
	response.Data.ChatId = liveEntity.ChatId
	response.Data.PushUrl = liveEntity.PushUrl
	response.Data.HlsUrl = liveEntity.HlsPlayUrl
	response.Data.RtmpUrl = liveEntity.RtmpPlayUrl
	response.Data.FlvUrl = liveEntity.FlvPlayUrl
	response.Data.Pv = 0
	response.Data.Uv = 0
	response.Data.TotalCount = 0
	response.Data.TotalMics = 0
	response.Data.LiveStatus = liveEntity.Status

	ctx.JSON(http.StatusOK, response)
}

func (c *liveController) GetLive(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	liveId := ctx.Param("id")
	liveEntity, err := service.GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("create liveEntity failed, err: %v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "create live failed",
			RequestId: log.ReqID(),
		})
		return
	}

	anchor, err := service.GetService().FindUser(ctx, liveEntity.AnchorId)
	if err != nil {
		log.Errorf("find anchor failed, err: %v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	response := &LiveResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.LiveId = liveEntity.LiveId
	response.Data.Title = liveEntity.Title
	response.Data.Notice = liveEntity.Notice
	response.Data.CoverUrl = liveEntity.CoverUrl
	response.Data.Extends = liveEntity.Extends
	response.Data.AnchorInfo.UserId = anchor.UserId
	response.Data.AnchorInfo.ImUserid = anchor.ImUserid
	response.Data.AnchorInfo.Nick = anchor.Nick
	response.Data.AnchorInfo.Avatar = anchor.Avatar
	response.Data.AnchorInfo.Extends = anchor.Extends
	response.Data.RoomToken = ""
	response.Data.PkId = liveEntity.PkId
	response.Data.OnlineCount = liveEntity.OnlineCount
	response.Data.StartTime = liveEntity.StartAt.Unix()
	response.Data.EndTime = liveEntity.EndAt.Unix()
	response.Data.ChatId = liveEntity.ChatId
	response.Data.PushUrl = liveEntity.PushUrl
	response.Data.HlsUrl = liveEntity.HlsPlayUrl
	response.Data.RtmpUrl = liveEntity.RtmpPlayUrl
	response.Data.FlvUrl = liveEntity.FlvPlayUrl
	response.Data.Pv = 0
	response.Data.Uv = 0
	response.Data.TotalCount = 0
	response.Data.TotalMics = 0
	response.Data.LiveStatus = liveEntity.Status
	response.Data.StopReason = liveEntity.StopReason
	response.Data.StopUserId = liveEntity.StopUserId
	if liveEntity.StopAt != nil {
		response.Data.StopTime = liveEntity.StopAt.Unix()
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *liveController) PostLiveStop(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("id")
	liveEntity, err := service.GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("find live error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	err = service.GetService().StopLive(ctx, liveId, liveEntity.AnchorId)
	if err != nil {
		log.Errorf("stop live failed, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.Response{
		Code:      200,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

func (c *liveController) DeleteLive(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("id")

	liveEntity, err := service.GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("find live error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	err = service.GetService().DeleteLive(ctx, liveId, liveEntity.AnchorId)
	if err != nil {
		log.Errorf("delete live failed, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "delete live failed",
			RequestId: log.ReqID(),
		})
		return
	}
	ctx.JSON(http.StatusOK, api.Response{
		Code:      200,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}
