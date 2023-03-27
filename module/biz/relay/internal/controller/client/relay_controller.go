// @Author: wangsheng
// @Description:
// @File:  relay_controller
// @Version: 1.0.0
// @Date: 2022/5/26 8:58 上午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/module/base/callback"
	"github.com/qbox/livekit/module/store/cache"
	"net/http"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/biz/relay"
	"github.com/qbox/livekit/module/fun/rtc"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterRoutes() {
	//relayGroup := group.Group("/relay")
	//relayGroup.POST("/start", RelayController.PostRelayStart)
	//relayGroup.GET("/:id", RelayController.GetRelaySession)
	//relayGroup.GET("/:id/token", RelayController.GetRelayToken)
	//relayGroup.POST("/:id/stop", RelayController.PostRelayStop)
	//relayGroup.POST("/:id/started", RelayController.PostRelayStarted)
	//relayGroup.POST("/:id/extends", RelayController.PostRelayExtends)

	httpq.ClientHandle(http.MethodPost, "/relay/start", RelayController.PostRelayStart)
	httpq.ClientHandle(http.MethodGet, "/relay/:id", RelayController.GetRelaySession)
	httpq.ClientHandle(http.MethodGet, "/relay/:id/token", RelayController.GetRelayToken)
	httpq.ClientHandle(http.MethodPost, "/relay/:id/stop", RelayController.PostRelayStop)
	httpq.ClientHandle(http.MethodPost, "/relay/:id/started", RelayController.PostRelayStarted)
	httpq.ClientHandle(http.MethodPost, "/relay/:id/extends", RelayController.PostRelayExtends)
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

type StartRelayResult struct {
	RelayId     string `json:"relay_id"`
	RelayStatus int    `json:"relay_status"`
	RelayToken  string `json:"relay_token"`
}

// PostRelayStart 开始跨房PK
// 到这里认为都是协商完成的，不需要进行通知与等待对方确认。
// POST  /client/relay/start
// return *StartRelayResult
func (*relayController) PostRelayStart(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	uInfo := auth.GetUserInfo(ctx)
	req := StartRelayRequest{}
	ctx.BindJSON(&req)
	if !req.IsValid() {
		log.Errorf("invalid args %+v", req)
		return nil, rest.ErrBadRequest
	}

	if uInfo.UserId == req.RecvUserId || req.InitRoomId == req.RecvRoomId {
		log.Errorf("invalid args %+v", req)
		return nil, rest.ErrBadRequest
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
		return nil, err
	}

	rtcClient := rtc.GetService()
	token := rtcClient.GetRelayToken(uInfo.UserId, req.RecvRoomId)
	resp := StartRelayResult{
		RelayId:     relaySession.SID,
		RelayStatus: relaySession.Status,
		RelayToken:  token,
	}
	// 跨房PK开始时需要进行回调
	go func() {
		err = callback.GetCallbackService().Do(ctx, callback.TypePKStarted, relaySession)
		if err != nil {
			log.Errorf("【pk_started】callback failed，errInfo：【%v】", err.Error())
		} else {
			log.Infof("【pk_started】callback success，data【%v】", relaySession)
		}
	}()
	// 将跨房双方信息缓存
	cache.Client.Set(fmt.Sprintf(model.PkIntegral, relaySession.InitUserId), relaySession.SID, 0)
	cache.Client.Set(fmt.Sprintf(model.PkIntegral, relaySession.RecvUserId), relaySession.SID, 0)
	return &resp, nil
}

type GetRelayTokenResult StartRelayResult

// GetRelayToken 获取跨房token
// 只有跨房的主播，才能获取relay token
// GET /client/relay/{id}/token
func (*relayController) GetRelayToken(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		return nil, rest.ErrBadRequest.WithMessage("empty relay session id")
	}
	uInfo := auth.GetUserInfo(ctx)

	relayService := relay.GetRelayService()
	relaySession, relayRoom, err := relayService.GetRelayRoom(ctx, uInfo.UserId, id)
	if err != nil {
		log.Errorf("empty relay room error %v", err)
		return nil, err
	}

	rtcClient := rtc.GetService()
	token := rtcClient.GetRelayToken(uInfo.UserId, relayRoom)

	resp := GetRelayTokenResult{
		RelayId:     relaySession.SID,
		RelayStatus: relaySession.Status,
		RelayToken:  token,
	}
	return &resp, nil
}

// GetRelaySession 获取跨房会话信息
// GET /client/relay/{id}
// return  *model.RelaySession
func (*relayController) GetRelaySession(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		return nil, rest.ErrBadRequest.WithMessage("empty relay session id")
	}

	relayService := relay.GetRelayService()
	relaySession, err := relayService.GetRelaySession(ctx, id)
	if err != nil {
		log.Errorf("get relay session (%s) error %v", id, err)
		return nil, err
	}

	return relaySession, nil
}

// PostRelayStop 停止跨房
// POST /client/relay/{id}/stop
// return nil
func (*relayController) PostRelayStop(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		return nil, rest.ErrBadRequest.WithMessage("empty relay session id")
	}
	uInfo := auth.GetUserInfo(ctx)
	relayService := relay.GetRelayService()
	// 跨房结束进行回调
	// 获取跨房信息
	session, err := relayService.GetRelaySession(ctx, id)
	if err == nil {
		err = callback.GetCallbackService().Do(ctx, callback.TypePKStopped, session)
		if err != nil {
			log.Errorf("【pk_stopped】callback failed，data：【%v】errInfo：【%v】", session, err.Error())
		} else {
			log.Infof("【pk_stopped】callback success，data【%v】", session)
		}
	} else {
		log.Errorf("【pk_stopped】callback failed，session is empty, data：【%v】errInfo：【%v】", session, err.Error())
	}
	if err = relayService.StopRelay(ctx, uInfo.UserId, id); err != nil {
		log.Errorf("stop relay error %v", err)
		return nil, err
	}
	// 跨房结束删除送礼的积分记录值
	if pkId, err := cache.Client.Get(fmt.Sprintf(model.PkIntegral, session.InitUserId)); err == nil {
		cache.Client.Del(fmt.Sprintf(model.PkIntegral, pkId))
	}
	cache.Client.Del(fmt.Sprintf(model.PkIntegral, session.InitUserId))
	cache.Client.Del(fmt.Sprintf(model.PkIntegral, session.RecvUserId))
	return nil, nil
}

type PostRelayStartedResult struct {
	RelayId     string `json:"relay_id"`
	RelayStatus int    `json:"relay_status"`
}

// PostRelayStarted 通知服务端，本端跨房已经完成
// POST /client/relay/{id}/started
// return *PostRelayStartedResult
func (*relayController) PostRelayStarted(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		return nil, rest.ErrBadRequest.WithMessage("empty relay session id")
	}

	uInfo := auth.GetUserInfo(ctx)
	relayService := relay.GetRelayService()
	if relaySession, err := relayService.ReportRelayStarted(ctx, uInfo.UserId, id); err != nil {
		log.Errorf("report relay started error %v", err)
		return nil, err
	} else {
		return &PostRelayStartedResult{
			RelayId:     relaySession.SID,
			RelayStatus: relaySession.Status,
		}, nil
	}
}

type PostRelayExtendsRequest struct {
	Extends model.Extends `json:"extends"`
}

// PostRelayExtends 更新扩展信息
// POST /client/relay/{id}/extends
// return nil
func (*relayController) PostRelayExtends(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	id, ok := ctx.Params.Get("id")
	if !ok {
		log.Errorf("empty relay session id")
		return nil, rest.ErrBadRequest.WithMessage("empty relay session id")
	}

	req := PostRelayExtendsRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	relayService := relay.GetRelayService()
	if err := relayService.UpdateRelayExtends(ctx, id, req.Extends); err != nil {
		log.Errorf("update relay extends error %v", err)
		return nil, err
	} else {
		return nil, nil
	}
}
