// @Author: wangsheng
// @Description:
// @File:  app_controller
// @Version: 1.0.0
// @Date: 2022/6/17 3:55 下午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"net/http"

	"github.com/qbox/livekit/app/live/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterAppRoutes(group *gin.RouterGroup) {
	appGroup := group.Group("/app")

	appGroup.GET("/config", AppController.GetConfig)
}

type appController struct {
}

var AppController = &appController{}

type AppConfig struct {
	IMAppID string `json:"im_app_id"`
}
type GetConfigResponse struct {
	api.Response
	Data AppConfig `json:"data"`
}

func (*appController) GetConfig(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	appConfig := AppConfig{
		IMAppID: config.AppConfig.ImConfig.AppId,
	}
	resp := GetConfigResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     appConfig,
	}

	ctx.JSON(http.StatusOK, &resp)
}
