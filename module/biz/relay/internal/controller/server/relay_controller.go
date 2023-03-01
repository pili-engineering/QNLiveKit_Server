// Copyright 2023 QINIU. All rights reserved
// @Description: 服务端跨房相关接口
// @Version: 1.0.0
// @Date: 2023/03/01 11:08
// @Author: fengyuan-liang@foxmail.com

package server

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/biz/relay"
	"github.com/qbox/livekit/utils/logger"
	"net/http"
)

func RegisterRoutes() {
	httpq.ServerHandle(http.MethodPost, "/relay/:id/extends", RelayController.PostRelayExtends)
}

var RelayController = &relayController{}

type relayController struct {
}

type PostRelayExtendsRequest struct {
	Extends model.Extends `json:"extends"`
}

func (c relayController) PostRelayExtends(ctx *gin.Context) (interface{}, error) {
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
