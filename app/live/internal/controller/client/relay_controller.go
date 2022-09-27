// @Author: wangsheng
// @Description:
// @File:  relay_controller
// @Version: 1.0.0
// @Date: 2022/5/26 8:58 上午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"net/http"

	"github.com/qbox/livekit/module/biz/relay"
	"github.com/qbox/livekit/module/fun/rtc"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterRelayRoutes(group *gin.RouterGroup) {
	relayGroup := group.Group("/relay")
	relayGroup.POST("/start", RelayController.PostRelayStart)
	relayGroup.GET("/:id", RelayController.GetRelaySession)
	relayGroup.GET("/:id/token", RelayController.GetRelayToken)
	relayGroup.POST("/:id/stop", RelayController.PostRelayStop)
	relayGroup.POST("/:id/started", RelayController.PostRelayStarted)
	relayGroup.POST("/:id/extends", RelayController.PostRelayExtends)
}

var RelayController = &relayController{}

type relayController struct {
}

type StartRelayRequest struct {
	InitRoomId string        `json:"init_room_id"`
	RecvRoomId string        `json:"recv_room_id"` //目的房间ID
	RecvUserId string        `json:"recv_user_id"` //目的主播用户ID
	Extends    model.Extends `json:"extends"`      //扩展数据
}

func (r *StartRelayRequest) IsValid() bool {
	return r.InitRoomId != "" && r.RecvRoomId != "" && r.RecvUserId != ""
}

type StartRelayResponse struct {
	api.Response
	Data struct {
		RelayId     string `json:"relay_id"`
		RelayStatus int    `json:"relay_status"`
		RelayToken  string `json:"relay_token"`
	} `json:"data"`
}

// 开始跨房PK
// 到这里认为都是协商完成的，不需要进行通知与等待对方确认。
// POST  /client/relay/start
func (*relayController) PostRelayStart(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	uInfo := liveauth.GetUserInfo(ctx)
	if uInfo == nil {
		log.Errorf("user info not exist")
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrNotFound))
		return
	}

	req := StartRelayRequest{}
	ctx.BindJSON(&req)
	if !req.IsValid() {
		log.Errorf("invalid args %+v", req)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if uInfo.UserId == req.RecvUserId || req.InitRoomId == req.RecvRoomId {
		log.Errorf("invalid args %+v", req)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	relayService := relay.GetRelayService()
	relayParams := relay.StartRelayParams{
		InitUserId: uInfo.UserId,
		InitRoomId: req.InitRoomId,
		RecvRoomId: req.RecvRoomId,
		RecvUserId: req.RecvUserId,
		Extends:    req.Extends,
	}
	relaySession, err := relayService.StartRelay(ctx, &relayParams)
	if err != nil {
		log.Errorf("start relay error %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	rtcClient := rtc.GetService()
	token := rtcClient.GetRelayToken(uInfo.UserId, req.RecvRoomId)
	resp := StartRelayResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data: struct {
			RelayId     string `json:"relay_id"`
			RelayStatus int    `json:"relay_status"`
			RelayToken  string `json:"relay_token"`
		}{
			RelayId:     relaySession.SID,
			RelayStatus: relaySession.Status,
			RelayToken:  token,
		},
	}

	ctx.JSON(http.StatusOK, &resp)
}

type GetRelayTokenResponse StartRelayResponse

// 获取跨房token
// 只有跨房的主播，才能获取relay token
// GET /client/relay/{id}/token
func (*relayController) GetRelayToken(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	uInfo := liveauth.GetUserInfo(ctx)

	relayService := relay.GetRelayService()
	relaySession, relayRoom, err := relayService.GetRelayRoom(ctx, uInfo.UserId, id)
	if err != nil {
		log.Errorf("empty relay room error %v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	rtcClient := rtc.GetService()
	token := rtcClient.GetRelayToken(uInfo.UserId, relayRoom)

	resp := GetRelayTokenResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data: struct {
			RelayId     string `json:"relay_id"`
			RelayStatus int    `json:"relay_status"`
			RelayToken  string `json:"relay_token"`
		}{
			RelayId:     relaySession.SID,
			RelayStatus: relaySession.Status,
			RelayToken:  token,
		},
	}

	ctx.JSON(http.StatusOK, &resp)
}

type GetRelaySessionResponse struct {
	api.Response
	Data *model.RelaySession `json:"data"`
}

// 获取跨房会话信息
// GET /client/relay/{id}
func (*relayController) GetRelaySession(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	relayService := relay.GetRelayService()
	relaySession, err := relayService.GetRelaySession(ctx, id)
	if err != nil {
		log.Errorf("get relay session (%s) error %v", id, err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	resp := GetRelaySessionResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     relaySession,
	}

	ctx.JSON(http.StatusOK, &resp)
}

// 停止跨房
// POST /client/relay/{id}/stop
func (*relayController) PostRelayStop(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	uInfo := liveauth.GetUserInfo(ctx)

	relayService := relay.GetRelayService()
	if err := relayService.StopRelay(ctx, uInfo.UserId, id); err != nil {
		log.Errorf("stop relay error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type PostRelayStartedResponse struct {
	api.Response
	Data struct {
		RelayId     string `json:"relay_id"`
		RelayStatus int    `json:"relay_status"`
	} `json:"data"`
}

// 通知服务端，本端跨房已经完成
// POST /client/relay/{id}/started
func (*relayController) PostRelayStarted(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	uInfo := liveauth.GetUserInfo(ctx)
	relayService := relay.GetRelayService()
	if relaySession, err := relayService.ReportRelayStarted(ctx, uInfo.UserId, id); err != nil {
		log.Errorf("report relay started error %v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	} else {
		ctx.JSON(http.StatusOK, &PostRelayStartedResponse{
			Response: api.SuccessResponse(log.ReqID()),
			Data: struct {
				RelayId     string `json:"relay_id"`
				RelayStatus int    `json:"relay_status"`
			}{
				RelayId:     relaySession.SID,
				RelayStatus: relaySession.Status,
			},
		})
		return
	}
}

type PostRelayExtendsRequest struct {
	Extends model.Extends `json:"extends"`
}

// 更新扩展信息
// POST /client/relay/{id}/extends
func (*relayController) PostRelayExtends(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	req := PostRelayExtendsRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	relayService := relay.GetRelayService()

	if err := relayService.UpdateRelayExtends(ctx, id, req.Extends); err != nil {
		log.Errorf("update relay extends error %v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	} else {
		ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
		return
	}
}
